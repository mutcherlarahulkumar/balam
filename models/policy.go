package models

import "time"

// Policy represents a row in the policies table.
type Policy struct {
	ID             int        `db:"_id"          json:"id"`
	PolicyNo       int        `db:"policy_no"    json:"policyNo"`
	FamilyCode     string     `db:"familycode"   json:"familyCode"`
	PersCode       string     `db:"perscode"     json:"persCode"`
	IssueDate      *time.Time `db:"issue_date"   json:"issueDate"`
	MatDate        *time.Time `db:"mat_date"     json:"matDate"`
	PaymentMode    string     `db:"payment_mode" json:"paymentMode"`
	Premium        float64    `db:"premium"      json:"premium"`
	SumAssured     float64    `db:"sum_assured"  json:"sumAssured"`
	Plan           int        `db:"plan"         json:"plan"`
	PlanName       string     `db:"plan_name"    json:"planName"`
	Term           int        `db:"term"         json:"term"`
	PPT            int        `db:"ppt"          json:"ppt"`
	PremiumEndDate *time.Time `db:"premium_end_date" json:"premiumEndDate"`
	NextPremium    *time.Time `db:"next_premium" json:"nextPremium"`
	MatAmount      float64    `db:"mat_amount"   json:"matAmount"`
	Nominee        string     `db:"nominee"      json:"nominee"`
	Relation       string     `db:"relation"     json:"relation"`
	AgCode         string     `db:"agcode"       json:"agCode"`
	StatCD         string     `db:"stat_cd"      json:"statCd"`
	Status         string     `db:"status"       json:"status"`
	FUPStatus      string     `db:"fupstatus"    json:"fupStatus"`
	Age            int        `db:"age"          json:"age"`
	LastPaid       *time.Time `db:"lastpaid"     json:"lastPaid"`
	NEFT           string     `db:"neft"         json:"neft"`
	LastFUP        *time.Time `db:"lastfup"      json:"lastFup"`
	Branch         string     `db:"branch"       json:"branch"`
	DAB            int        `db:"dab"          json:"dab"`
	TermRider      int        `db:"termrider"    json:"termRider"`
}

// PolicyListItem is the compact view returned in GET /policies listing.
type PolicyListItem struct {
	ID             int        `json:"id"`
	PolicyNo       int        `json:"policyNo"`
	FamilyCode     string     `json:"familyCode"`
	ClientName     string     `json:"clientName"`
	PlanNo         string     `json:"planNo"`
	PlanName       string     `json:"planName"`
	Term           int        `json:"term"`
	PPT            int        `json:"ppt"`
	Premium        float64    `json:"premium"`
	SumAssured     float64    `json:"sumAssured"`
	PaymentMode    string     `json:"paymentMode"`
	PremiumEndDate *time.Time `json:"premiumEndDate"`
	NextPremium    *time.Time `json:"nextPremium"`
	MatDate        *time.Time `json:"matDate"`
	Status         string     `json:"status"`
	FUPStatus      string     `json:"fupStatus"`
	DaysUntilLapse int        `json:"daysUntilLapse"`
}

// CreatePolicyRequest holds data for POST /policies.
type CreatePolicyRequest struct {
	PolicyNo    int     `json:"policyNo"    binding:"required,gt=0"`
	FamilyCode  string  `json:"familyCode"  binding:"required"`
	PersCode    string  `json:"persCode"    binding:"required"`
	PlanNo      string  `json:"planNo"      binding:"required"`
	IssueDate   string  `json:"issueDate"   binding:"required"`
	MatDate     string  `json:"matDate"     binding:"required"`
	Term        int     `json:"term"        binding:"required,min=1"`
	PPT         int     `json:"ppt"         binding:"required,min=1"`
	SumAssured  float64 `json:"sumAssured"  binding:"required,gt=0"`
	Premium     float64 `json:"premium"     binding:"required,gt=0"`
	PaymentMode string  `json:"paymentMode" binding:"required,oneof=Y H Q M S"`
	NextPremium string  `json:"nextPremium" binding:"required"`
	Nominee     string  `json:"nominee"     binding:"required,max=80"`
	Relation    string  `json:"relation"    binding:"required,max=20"`
	Branch      string  `json:"branch"      binding:"omitempty,max=15"`
	NEFT        string  `json:"neft"        binding:"omitempty,oneof=YES NO"`
	DAB         int     `json:"dab"         binding:"omitempty,min=0"`
	TermRider   int     `json:"termRider"   binding:"omitempty,min=0"`
}

// UpdatePolicyRequest holds data for PUT /policies/:policyNo.
type UpdatePolicyRequest struct {
	Status      string `json:"status"      binding:"omitempty,oneof=IF LA PU SU MA CL EX"`
	Nominee     string `json:"nominee"     binding:"omitempty,max=80"`
	Relation    string `json:"relation"    binding:"omitempty,max=20"`
	NEFT        string `json:"neft"        binding:"omitempty,oneof=YES NO"`
	NextPremium string `json:"nextPremium" binding:"omitempty"`
	FUPStatus   string `json:"fupStatus"   binding:"omitempty,oneof=DUE PAID OVERDUE LAPSED"`
}

// PolicyDetail extends Policy with full history.
type PolicyDetail struct {
	Policy
	PlanDetails *PlanResponse `json:"planDetails"`
	FUPHistory  []FUPHistory  `json:"fupHistory"`
	Loans       []Loan        `json:"loans"`
	SBRecords   []SB          `json:"sbRecords"`
}
