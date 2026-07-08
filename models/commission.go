package models

import "time"

// Commission represents a row in the commission table.
type Commission struct {
	ID         int        `db:"_id"        json:"id"`
	PolicyNo   int        `db:"policy_no"  json:"policyNo"`
	BillDate   *time.Time `db:"bill_date"  json:"billDate"`
	FirstComm  float64    `db:"first_comm" json:"firstComm"`
	SecondComm float64    `db:"second_comm" json:"secondComm"`
	ThirdComm  float64    `db:"third_comm" json:"thirdComm"`
	BonusComm  float64    `db:"bonus_comm" json:"bonusComm"`
	SingleComm float64    `db:"single_comm" json:"singleComm"`
	SubComm    float64    `db:"sub_comm"   json:"subComm"`
	PayDate    *time.Time `db:"paydate"    json:"payDate"`
}

// CreateCommissionRequest holds data for POST /commission.
type CreateCommissionRequest struct {
	PolicyNo   int     `json:"policyNo"   binding:"required,gt=0"`
	BillDate   string  `json:"billDate"   binding:"required"`
	FirstComm  float64 `json:"firstComm"  binding:"omitempty,min=0"`
	SecondComm float64 `json:"secondComm" binding:"omitempty,min=0"`
	ThirdComm  float64 `json:"thirdComm"  binding:"omitempty,min=0"`
	BonusComm  float64 `json:"bonusComm"  binding:"omitempty,min=0"`
	SingleComm float64 `json:"singleComm" binding:"omitempty,min=0"`
	PayDate    string  `json:"payDate"    binding:"omitempty"`
}

// CommissionSummary is the response body for GET /commission/summary.
type CommissionSummary struct {
	CurrentMonth  MonthlyCommission  `json:"currentMonth"`
	Yearly        []YearlyCommission `json:"yearly"`
}

// MonthlyCommission holds aggregated commission for one month.
type MonthlyCommission struct {
	Month           string  `json:"month"`
	TotalCommission float64 `json:"totalCommission"`
	PoliciesBilled  int     `json:"policiesBilled"`
}

// YearlyCommission holds aggregated commission for one year.
type YearlyCommission struct {
	Year      int     `json:"year"`
	FirstYear float64 `json:"firstYear"`
	Renewal   float64 `json:"renewal"`
	Bonus     float64 `json:"bonus"`
	Gross     float64 `json:"gross"`
}

// CommissionEstimate is returned by GET /commission/calculate.
type CommissionEstimate struct {
	PolicyNo            int     `json:"policyNo"`
	Premium             float64 `json:"premium"`
	PlanNo              string  `json:"planNo"`
	Year                int     `json:"year"`
	BaseCommissionPct   float64 `json:"baseCommissionPct"`
	BonusCommissionPct  float64 `json:"bonusCommissionPct"`
	TotalPct            float64 `json:"totalPct"`
	EstimatedCommission float64 `json:"estimatedCommission"`
	Note                string  `json:"note"`
}
