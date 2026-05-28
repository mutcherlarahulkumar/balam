package repository

import (
	"agent-balam/models"
	"time"

	"github.com/jmoiron/sqlx"
)

// LoanRepo handles loan table queries.
type LoanRepo struct {
	db *sqlx.DB
}

// NewLoanRepo creates a new LoanRepo.
func NewLoanRepo(db *sqlx.DB) *LoanRepo {
	return &LoanRepo{db: db}
}

// List returns all loans, optionally filtered by policyNo.
func (r *LoanRepo) List(policyNo int) ([]models.Loan, error) {
	query := `SELECT _id, policy_no, ldate, loan, intduedate, loaninterest, COALESCE(details,'') AS details FROM loan`
	args := []interface{}{}
	if policyNo > 0 {
		query += ` WHERE policy_no=$1`
		args = append(args, policyNo)
	}
	query += ` ORDER BY ldate DESC`
	var items []models.Loan
	err := r.db.Select(&items, query, args...)
	if items == nil {
		items = []models.Loan{}
	}
	return items, err
}

// Create inserts a loan record.
func (r *LoanRepo) Create(req models.CreateLoanRequest, loanDate, intDueDate time.Time) error {
	_, err := r.db.Exec(`INSERT INTO loan (policy_no, ldate, loan, intduedate, loaninterest) VALUES ($1,$2,$3,$4,$5)`,
		req.PolicyNo, loanDate, int(req.LoanAmount), intDueDate, int(req.LoanInterest))
	return err
}
