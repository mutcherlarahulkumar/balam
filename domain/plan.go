package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
)

// PlanService handles plan catalogue business logic.
type PlanService struct {
	planRepo *repository.PlanRepo
}

// NewPlanService creates a new PlanService.
func NewPlanService(repo *repository.PlanRepo) *PlanService {
	return &PlanService{planRepo: repo}
}

// List returns the plan catalogue as API response objects.
func (s *PlanService) List(planType, search string) ([]models.PlanResponse, error) {
	plans, err := s.planRepo.List(planType, search)
	if err != nil {
		return nil, err
	}

	resp := make([]models.PlanResponse, 0, len(plans))
	for _, p := range plans {
		pr := models.PlanResponse{
			PlanNo:   p.PlanNo,
			PlanName: p.PlanName,
			PlanType: p.PlanType,
			TermPPT:  p.TermPPT == "Y" || p.TermPPT == "1" || p.TermPPT == "true",
			Stax:     p.Stax,
			LapsDays: p.LapsDays,
			GSTRates: models.GSTInfo{
				FirstYear: 0,
				Renewal:   0,
				Note:      "GST 0% effective 22 Sep 2025 for all individual policies",
			},
		}
		if p.SBYear != "" {
			pr.SBYears = &p.SBYear
		}
		if p.SBBenefit != "" {
			pr.SBBenefits = &p.SBBenefit
		}
		resp = append(resp, pr)
	}
	return resp, nil
}
