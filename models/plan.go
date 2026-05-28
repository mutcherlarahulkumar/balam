package models

// Plan represents a row in the plans table.
type Plan struct {
	PlanNo   string `db:"plan_no"   json:"planNo"`
	PlanName string `db:"plan_name" json:"planName"`
	PlanType string `db:"plan_type" json:"planType"`
	TermPPT  string `db:"term_ppt"  json:"termPpt"`
	SBYear   string `db:"sb_year"   json:"sbYears"`
	SBBenefit string `db:"sb_benefit" json:"sbBenefits"`
	Stax     string `db:"stax"      json:"stax"`
	LapsDays int    `db:"lapsdays"  json:"lapsDays"`
}

// PlanResponse is the API response shape for a plan.
type PlanResponse struct {
	PlanNo    string   `json:"planNo"`
	PlanName  string   `json:"planName"`
	PlanType  string   `json:"planType"`
	TermPPT   bool     `json:"termPpt"`
	SBYears   *string  `json:"sbYears"`
	SBBenefits *string `json:"sbBenefits"`
	Stax      string   `json:"stax"`
	LapsDays  int      `json:"lapsDays"`
	GSTRates  GSTInfo  `json:"gstRates"`
}

// GSTInfo holds the current GST rate metadata returned with plans.
type GSTInfo struct {
	FirstYear float64 `json:"firstYear"`
	Renewal   float64 `json:"renewal"`
	Note      string  `json:"note"`
}
