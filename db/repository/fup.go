package repository

import (
	"agent-balam/models"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// FUPRepo handles fup-related table queries.
type FUPRepo struct {
	db *sqlx.DB
}

// NewFUPRepo creates a new FUPRepo.
func NewFUPRepo(db *sqlx.DB) *FUPRepo {
	return &FUPRepo{db: db}
}

// DuePolicies returns policies with premiums due/overdue, optionally filtered by month (YYYY-MM) and overdue days.
func (r *FUPRepo) DuePolicies(month string, overdueDays int) ([]models.FUPDueItem, error) {
	query := `
		SELECT p.policy_no, COALESCE(c.name,'') AS client_name,
		       COALESCE(c.mobileno,'') AS mobile,
		       COALESCE(p.plan_name,'') AS plan_name,
		       p.premium, p.next_premium, p.payment_mode,
		       GREATEST(0, EXTRACT(DAY FROM CURRENT_DATE - p.next_premium)::int) AS days_overdue,
		       p.next_premium + COALESCE(pl.lapsdays,180) * INTERVAL '1 day' AS lapse_date,
		       GREATEST(0, EXTRACT(DAY FROM (p.next_premium + COALESCE(pl.lapsdays,180) * INTERVAL '1 day') - CURRENT_DATE)::int) AS days_until_lapse
		FROM policies p
		LEFT JOIN client c ON c.familycode=p.familycode AND c.perscode=p.perscode
		LEFT JOIN plans pl ON pl.plan_no=CAST(p.plan AS TEXT)
		WHERE p.next_premium <= CURRENT_DATE
		  AND (p.status IS NULL OR p.status NOT IN ('SU','MA','PU'))`

	args := []interface{}{}
	n := 1

	if month != "" {
		query += ` AND TO_CHAR(p.next_premium, 'YYYY-MM') = $` + fmt.Sprintf("%d", n)
		args = append(args, month)
		n++
	}
	if overdueDays > 0 {
		query += ` AND EXTRACT(DAY FROM CURRENT_DATE - p.next_premium) >= $` + fmt.Sprintf("%d", n)
		args = append(args, overdueDays)
		n++
	}
	query += ` ORDER BY p.next_premium`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.FUPDueItem
	for rows.Next() {
		var item models.FUPDueItem
		if err := rows.Scan(&item.PolicyNo, &item.ClientName, &item.Mobile, &item.PlanName,
			&item.Premium, &item.NextPremium, &item.PaymentMode, &item.DaysOverdue,
			&item.LapseDate, &item.DaysUntilLapse); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.FUPDueItem{}
	}
	return items, rows.Err()
}

// UpdateFUP updates next_premium and fupstatus on policies, appends to both history tables.
func (r *FUPRepo) UpdateFUP(policyNo int, oldFUP, newFUP time.Time, updatedBy string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	now := time.Now()

	_, err = tx.Exec(`INSERT INTO fupupdate (policy_no, oldfup, newfup, name, dateupdated) VALUES ($1,$2,$3,$4,$5)`,
		policyNo, oldFUP, newFUP, updatedBy, now)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO fuphistory (policy_no, oldfup, newfup, name, dateupdated) VALUES ($1,$2,$3,$4,$5)`,
		policyNo, oldFUP, newFUP, updatedBy, now)
	if err != nil {
		return err
	}

	fupStatus := "DUE"
	if newFUP.After(time.Now()) {
		fupStatus = "PAID"
	}
	_, err = tx.Exec(`UPDATE policies SET next_premium=$1, fupstatus=$2, lastfup=$3 WHERE policy_no=$4`,
		newFUP, fupStatus, now, policyNo)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// History returns all FUP changes for a policy.
func (r *FUPRepo) History(policyNo int) ([]models.FUPHistory, error) {
	var items []models.FUPHistory
	err := r.db.Select(&items, `SELECT _id, policy_no, oldfup, newfup, name, dateupdated
	                              FROM fuphistory WHERE policy_no=$1 ORDER BY dateupdated DESC`, policyNo)
	if items == nil {
		items = []models.FUPHistory{}
	}
	return items, err
}
