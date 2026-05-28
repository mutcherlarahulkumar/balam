package models

import "time"

// Loan represents a row in the loan table.
type Loan struct {
	ID           int        `db:"_id"          json:"id"`
	PolicyNo     int        `db:"policy_no"    json:"policyNo"`
	LoanDate     *time.Time `db:"ldate"        json:"loanDate"`
	LoanAmount   int        `db:"loan"         json:"loanAmount"`
	IntDueDate   *time.Time `db:"intduedate"   json:"interestDueDate"`
	LoanInterest int        `db:"loaninterest" json:"loanInterest"`
	Details      string     `db:"details"      json:"details"`
}

// CreateLoanRequest holds data for POST /loans.
type CreateLoanRequest struct {
	PolicyNo     int     `json:"policyNo"         binding:"required"`
	LoanDate     string  `json:"loanDate"         binding:"required"`
	LoanAmount   float64 `json:"loanAmount"       binding:"required,gt=0"`
	IntDueDate   string  `json:"interestDueDate"  binding:"required"`
	LoanInterest float64 `json:"loanInterest"     binding:"omitempty,min=0"`
}
