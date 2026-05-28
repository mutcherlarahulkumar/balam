package models

import "time"

// FUPHistory represents a row in the fuphistory table.
type FUPHistory struct {
	ID          int        `db:"_id"         json:"id"`
	PolicyNo    int        `db:"policy_no"   json:"policyNo"`
	OldFUP      *time.Time `db:"oldfup"      json:"oldFup"`
	NewFUP      *time.Time `db:"newfup"      json:"newFup"`
	UpdatedBy   string     `db:"name"        json:"updatedBy"`
	UpdatedAt   *time.Time `db:"dateupdated" json:"updatedAt"`
}

// FUPDueItem is one row returned by GET /fup/due.
type FUPDueItem struct {
	PolicyNo       int        `json:"policyNo"`
	ClientName     string     `json:"clientName"`
	Mobile         string     `json:"mobile"`
	PlanName       string     `json:"planName"`
	Premium        float64    `json:"premium"`
	NextPremium    *time.Time `json:"nextPremium"`
	DaysOverdue    int        `json:"daysOverdue"`
	PaymentMode    string     `json:"paymentMode"`
	LapseDate      *time.Time `json:"lapseDate"`
	DaysUntilLapse int        `json:"daysUntilLapse"`
}

// UpdateFUPRequest holds data for POST /fup/update.
type UpdateFUPRequest struct {
	PolicyNo int    `json:"policyNo" binding:"required"`
	OldFUP   string `json:"oldFup"   binding:"required"`
	NewFUP   string `json:"newFup"   binding:"required"`
	Reason   string `json:"reason"   binding:"omitempty,max=500"`
}
