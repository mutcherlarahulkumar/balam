package repository

import (
	"agent-balam/models"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

// CommissionRepo handles commission table queries.
type CommissionRepo struct {
	db *sqlx.DB
}

// NewCommissionRepo creates a new CommissionRepo.
func NewCommissionRepo(db *sqlx.DB) *CommissionRepo {
	return &CommissionRepo{db: db}
}

// List returns commission records filtered by year and month.
func (r *CommissionRepo) List(year, month string) ([]models.Commission, error) {
	query := `SELECT _id, policy_no, bill_date, first_comm, second_comm, third_comm, bonus_comm, single_comm, sub_comm, paydate
	          FROM commission WHERE 1=1`
	args := []interface{}{}
	n := 1
	if y, err := strconv.Atoi(year); err == nil && y > 0 {
		query += fmt.Sprintf(` AND EXTRACT(YEAR FROM bill_date) = $%d`, n)
		args = append(args, y)
		n++
	}
	if m, err := strconv.Atoi(month); err == nil && m >= 1 && m <= 12 {
		query += fmt.Sprintf(` AND EXTRACT(MONTH FROM bill_date) = $%d`, n)
		args = append(args, m)
		n++
	}
	query += ` ORDER BY bill_date DESC`

	var items []models.Commission
	err := r.db.Select(&items, query, args...)
	if items == nil {
		items = []models.Commission{}
	}
	return items, err
}

// Create inserts a new commission record.
func (r *CommissionRepo) Create(req models.CreateCommissionRequest, billDate, payDate *time.Time) error {
	_, err := r.db.Exec(`INSERT INTO commission (policy_no, bill_date, first_comm, second_comm, third_comm, bonus_comm, single_comm, paydate)
	                      VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		req.PolicyNo, billDate, req.FirstComm, req.SecondComm, req.ThirdComm, req.BonusComm, req.SingleComm, payDate)
	return err
}

// Summary returns monthly and yearly aggregations.
func (r *CommissionRepo) Summary() (*models.CommissionSummary, error) {
	var monthly models.MonthlyCommission
	err := r.db.QueryRow(`
		SELECT TO_CHAR(DATE_TRUNC('month', CURRENT_DATE),'YYYY-MM'),
		       COALESCE(SUM(first_comm+second_comm+third_comm+bonus_comm+single_comm),0),
		       COUNT(DISTINCT policy_no)
		FROM commission
		WHERE DATE_TRUNC('month', bill_date) = DATE_TRUNC('month', CURRENT_DATE)`).
		Scan(&monthly.Month, &monthly.TotalCommission, &monthly.PoliciesBilled)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT EXTRACT(YEAR FROM bill_date)::int,
		       COALESCE(SUM(first_comm),0),
		       COALESCE(SUM(second_comm+third_comm),0),
		       COALESCE(SUM(bonus_comm),0),
		       COALESCE(SUM(first_comm+second_comm+third_comm+bonus_comm+single_comm),0)
		FROM commission
		GROUP BY EXTRACT(YEAR FROM bill_date)
		ORDER BY 1 DESC LIMIT 5`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var yearly []models.YearlyCommission
	for rows.Next() {
		var y models.YearlyCommission
		if err := rows.Scan(&y.Year, &y.FirstYear, &y.Renewal, &y.Bonus, &y.Gross); err != nil {
			return nil, err
		}
		yearly = append(yearly, y)
	}
	if yearly == nil {
		yearly = []models.YearlyCommission{}
	}
	return &models.CommissionSummary{CurrentMonth: monthly, Yearly: yearly}, rows.Err()
}
