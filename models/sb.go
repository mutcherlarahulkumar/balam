package models

import "time"

// SB represents a row in the sb (survival benefit) table.
type SB struct {
	ID        int        `db:"_id"        json:"id"`
	PolicyNo  int        `db:"policy_no"  json:"policyNo"`
	SBDueDate *time.Time `db:"sb_duedate" json:"sbDueDate"`
	SBAmount  float64    `db:"sb_amount"  json:"sbAmount"`
	InstalNo  int        `db:"i_no"       json:"instalmentNo"`
	SBPayDate *time.Time `db:"sb_paydate" json:"sbPayDate"`
	ChequeNo  string     `db:"ch_no"      json:"chequeNo"`
	Details   string     `db:"details"    json:"details"`
}

// CreateSBRequest holds data for POST /sb.
type CreateSBRequest struct {
	PolicyNo    int     `json:"policyNo"     binding:"required,gt=0"`
	SBDueDate   string  `json:"sbDueDate"    binding:"required"`
	SBAmount    float64 `json:"sbAmount"     binding:"required,gt=0"`
	InstalmentNo int    `json:"instalmentNo" binding:"required,min=1"`
}

// MarkSBPaidRequest holds data for PUT /sb/:id/mark-paid.
type MarkSBPaidRequest struct {
	PaidDate string `json:"paidDate" binding:"required"`
	ChequeNo string `json:"chequeNo" binding:"omitempty,max=10"`
}
