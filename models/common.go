package models

// PaginatedResponse wraps a data slice with pagination metadata.
type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

// ErrorResponse is the standard error body returned on all error codes.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// GSTCalculateResponse is returned by GET /gst/calculate.
type GSTCalculateResponse struct {
	PolicyNo        int     `json:"policyNo"`
	PlanNo          string  `json:"planNo"`
	PlanType        string  `json:"planType"`
	BasePremium     float64 `json:"basePremium"`
	PremiumYear     int     `json:"premiumYear"`
	GSTRate         float64 `json:"gstRate"`
	GSTAmount       float64 `json:"gstAmount"`
	TotalPremium    float64 `json:"totalPremium"`
	Regulation      string  `json:"regulation"`
	HistoricalNote  string  `json:"historicalNote"`
}

// ReportCashflow is a row from the rpt_cashflow report.
type ReportCashflow struct {
	FamilyCode string  `db:"familycode" json:"familyCode"`
	PersCode   string  `db:"perscode"   json:"persCode"`
	PolicyNo   string  `db:"policy_no"  json:"policyNo"`
	DueDate    string  `db:"due_date"   json:"dueDate"`
	Name       string  `db:"name"       json:"name"`
	Age        int     `db:"age"        json:"age"`
	SBAmount   float64 `db:"sb_amount"  json:"sbAmount"`
	Bonus      float64 `db:"bonus"      json:"bonus"`
	Total      float64 `db:"total"      json:"total"`
	Type       string  `db:"type"       json:"type"`
}

// ReportCashInOut is a row from rpt_cashinout.
type ReportCashInOut struct {
	Year      int     `db:"year"       json:"year"`
	CashInMat float64 `db:"cashin_mat" json:"cashInMat"`
	CashInSB  float64 `db:"cashin_sb"  json:"cashInSB"`
	CashOut   float64 `db:"cashout"    json:"cashOut"`
	NetAmt    float64 `db:"netamt"     json:"netAmt"`
}

// ReportCurrentStatus is a row from rpt_currstatus.
type ReportCurrentStatus struct {
	FamilyCode   string  `db:"familycode"    json:"familyCode"`
	PersCode     string  `db:"perscode"      json:"persCode"`
	PolicyNo     string  `db:"policy_no"     json:"policyNo"`
	DueStatus    string  `db:"due_status"    json:"dueStatus"`
	PremiumPaid  float64 `db:"premium_paid"  json:"premiumPaid"`
	VestedBonus  float64 `db:"vested_bonus"  json:"vestedBonus"`
	SurrValue    float64 `db:"surr_value"    json:"surrValue"`
	LoanTaken    float64 `db:"loan_taken"    json:"loanTaken"`
	LoanDate     string  `db:"loan_date"     json:"loanDate"`
	LoanAmt      float64 `db:"loan_amt"      json:"loanAmt"`
}

// ReportCalendar is a row from rpt_prcalendar.
type ReportCalendar struct {
	FamilyCode string  `db:"familycode" json:"familyCode"`
	PersCode   string  `db:"perscode"   json:"persCode"`
	PolicyNo   string  `db:"policy_no"  json:"policyNo"`
	Jan        float64 `db:"jan"        json:"jan"`
	Feb        float64 `db:"feb"        json:"feb"`
	Mar        float64 `db:"mar"        json:"mar"`
	Apr        float64 `db:"apr"        json:"apr"`
	May        float64 `db:"may"        json:"may"`
	Jun        float64 `db:"jun"        json:"jun"`
	Jul        float64 `db:"jul"        json:"jul"`
	Aug        float64 `db:"aug"        json:"aug"`
	Sep        float64 `db:"sep"        json:"sep"`
	Oct        float64 `db:"oct"        json:"oct"`
	Nov        float64 `db:"nov"        json:"nov"`
	Dec        float64 `db:"dec"        json:"dec"`
}
