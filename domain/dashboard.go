package domain

import "agent-balam/models"

const dashboardPreviewLimit = 5

// DashboardService aggregates data from other services for the agent home screen.
type DashboardService struct {
	fupSvc  *FUPService
	leadSvc *LeadService
	commSvc *CommissionService
	sbSvc   *SBService
}

// NewDashboardService creates a new DashboardService.
func NewDashboardService(fupSvc *FUPService, leadSvc *LeadService, commSvc *CommissionService, sbSvc *SBService) *DashboardService {
	return &DashboardService{fupSvc: fupSvc, leadSvc: leadSvc, commSvc: commSvc, sbSvc: sbSvc}
}

// Summary returns a single aggregated payload combining: due/overdue premiums,
// today's scheduled activities, this month's commission total, unpaid survival
// benefits, and recent leads — everything an agent needs on login in one call.
func (s *DashboardService) Summary() (*models.DashboardResponse, error) {
	due, err := s.fupSvc.DuePolicies("", "", 0)
	if err != nil {
		return nil, err
	}

	activities, err := s.leadSvc.TodayActivities()
	if err != nil {
		return nil, err
	}

	commSummary, err := s.commSvc.Summary()
	if err != nil {
		return nil, err
	}

	unpaidSB, err := s.sbSvc.List("", "", true)
	if err != nil {
		return nil, err
	}

	leads, err := s.leadSvc.ListLeads()
	if err != nil {
		return nil, err
	}

	resp := &models.DashboardResponse{
		DuePremiums:         models.DashboardFUP{Total: len(due), Preview: truncateFUP(due, dashboardPreviewLimit)},
		TodayActivities:     activities,
		CommissionThisMonth: commSummary.CurrentMonth,
		UnpaidSB:            models.DashboardSB{Total: len(unpaidSB), Preview: truncateSB(unpaidSB, dashboardPreviewLimit)},
		Leads:               models.DashboardLeads{Total: len(leads), Preview: truncateLeads(leads, dashboardPreviewLimit)},
	}
	return resp, nil
}

func truncateFUP(items []models.FUPDueItem, n int) []models.FUPDueItem {
	if len(items) <= n {
		return items
	}
	return items[:n]
}

func truncateSB(items []models.SB, n int) []models.SB {
	if len(items) <= n {
		return items
	}
	return items[:n]
}

func truncateLeads(items []models.Lead, n int) []models.Lead {
	if len(items) <= n {
		return items
	}
	return items[:n]
}
