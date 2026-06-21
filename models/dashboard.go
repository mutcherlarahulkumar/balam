package models

// DashboardResponse aggregates the data an agent needs on their home screen into a
// single call, avoiding multiple round trips to /fup/due, /activities/today,
// /commission/summary, /sb and /leads on every app open.
type DashboardResponse struct {
	DuePremiums         DashboardFUP      `json:"duePremiums"`
	TodayActivities     []Activity        `json:"todayActivities"`
	CommissionThisMonth MonthlyCommission `json:"commissionThisMonth"`
	UnpaidSB            DashboardSB       `json:"unpaidSB"`
	Leads               DashboardLeads    `json:"leads"`
	LapsingPolicies     DashboardPolicies `json:"lapsingPolicies"`
}

// DashboardPolicies summarizes policies lapsing soon: total count plus a short preview.
type DashboardPolicies struct {
	Total   int              `json:"total"`
	Preview []PolicyListItem `json:"preview"`
}

// DashboardFUP summarizes due/overdue premiums: total count plus a short preview.
type DashboardFUP struct {
	Total   int          `json:"total"`
	Preview []FUPDueItem `json:"preview"`
}

// DashboardSB summarizes unpaid survival benefits: total count plus a short preview.
type DashboardSB struct {
	Total   int  `json:"total"`
	Preview []SB `json:"preview"`
}

// DashboardLeads summarizes leads: total count plus a short preview (most recent first).
type DashboardLeads struct {
	Total   int    `json:"total"`
	Preview []Lead `json:"preview"`
}
