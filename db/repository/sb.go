package repository

import (
	"agent-balam/models"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

// SBRepo handles survival benefits table queries.
type SBRepo struct {
	db *sqlx.DB
}

// NewSBRepo creates a new SBRepo.
func NewSBRepo(db *sqlx.DB) *SBRepo {
	return &SBRepo{db: db}
}

// List returns SB records, optionally filtered by year, month, and unpaid status.
func (r *SBRepo) List(year, month string, unpaidOnly bool) ([]models.SB, error) {
	query := `SELECT _id, policy_no, sb_duedate, sb_amount, i_no, sb_paydate, ch_no, COALESCE(details,'') FROM sb WHERE 1=1`
	args := []interface{}{}
	n := 1
	if y, err := strconv.Atoi(year); err == nil && y > 0 {
		query += fmt.Sprintf(` AND EXTRACT(YEAR FROM sb_duedate) = $%d`, n)
		args = append(args, y)
		n++
	}
	if m, err := strconv.Atoi(month); err == nil && m >= 1 && m <= 12 {
		query += fmt.Sprintf(` AND EXTRACT(MONTH FROM sb_duedate) = $%d`, n)
		args = append(args, m)
		n++
	}
	if unpaidOnly {
		query += ` AND sb_paydate IS NULL`
	}
	query += ` ORDER BY sb_duedate`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.SB
	for rows.Next() {
		var item models.SB
		if err := rows.Scan(&item.ID, &item.PolicyNo, &item.SBDueDate, &item.SBAmount,
			&item.InstalNo, &item.SBPayDate, &item.ChequeNo, &item.Details); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.SB{}
	}
	return items, rows.Err()
}

// ListByPolicy returns all SB records for a specific policy number.
func (r *SBRepo) ListByPolicy(policyNo int) ([]models.SB, error) {
	rows, err := r.db.Query(`SELECT _id, policy_no, sb_duedate, sb_amount, i_no, sb_paydate, ch_no, COALESCE(details,'') FROM sb WHERE policy_no=$1 ORDER BY sb_duedate`, policyNo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.SB
	for rows.Next() {
		var item models.SB
		if err := rows.Scan(&item.ID, &item.PolicyNo, &item.SBDueDate, &item.SBAmount, &item.InstalNo, &item.SBPayDate, &item.ChequeNo, &item.Details); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.SB{}
	}
	return items, rows.Err()
}

// Create inserts a new SB record.
func (r *SBRepo) Create(req models.CreateSBRequest, dueDate time.Time) error {
	_, err := r.db.Exec(`INSERT INTO sb (policy_no, sb_duedate, sb_amount, i_no) VALUES ($1,$2,$3,$4)`,
		req.PolicyNo, dueDate, req.SBAmount, req.InstalmentNo)
	return err
}

// MarkPaid sets the paid date and cheque number on an SB record.
func (r *SBRepo) MarkPaid(id int, paidDate time.Time, chequeNo string) error {
	_, err := r.db.Exec(`UPDATE sb SET sb_paydate=$1, ch_no=$2 WHERE _id=$3`, paidDate, chequeNo, id)
	return err
}

// FindByID returns a single SB record.
func (r *SBRepo) FindByID(id int) (*models.SB, error) {
	var item models.SB
	if err := r.db.QueryRow(`SELECT _id, policy_no, sb_duedate, sb_amount, i_no, sb_paydate, ch_no, COALESCE(details,'') FROM sb WHERE _id=$1`, id).
		Scan(&item.ID, &item.PolicyNo, &item.SBDueDate, &item.SBAmount, &item.InstalNo, &item.SBPayDate, &item.ChequeNo, &item.Details); err != nil {
		return nil, err
	}
	return &item, nil
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
