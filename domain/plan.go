package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"strconv"
	"strings"
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
// TODO(LIC-API — Phase 2): Sync plan catalogue from GET /api/v2/plans/list on agent login
// to pick up new plans as LIC launches them.
func (s *PlanService) List(planType, search string) ([]models.PlanResponse, error) {
	plans, err := s.planRepo.List(planType, search)
	if err != nil {
		return nil, err
	}

	resp := make([]models.PlanResponse, 0, len(plans))
	for _, p := range plans {
		pr := models.PlanResponse{
			PlanNo:     p.PlanNo,
			PlanName:   p.PlanName,
			PlanType:   p.PlanType,
			TermPPT:    p.TermPPT == "Y" || p.TermPPT == "1" || p.TermPPT == "true",
			SBSchedule: parseSBSchedule(p.SBYear, p.SBBenefit),
			Stax:       p.Stax,
			LapsDays:   p.LapsDays,
			GSTRates: models.GSTInfo{
				FirstYear: 0,
				Renewal:   0,
				Note:      "GST 0% effective 22 Sep 2025 for all individual policies",
			},
		}
		resp = append(resp, pr)
	}
	return resp, nil
}

// parseSBSchedule converts the pipe-delimited sb_year and sb_benefit strings into a
// structured slice. Example: sb_year="05|10|15", sb_benefit="20|20|20" → [{year:5, benefitPct:20}, ...]
func parseSBSchedule(sbYear, sbBenefit string) []models.SBScheduleEntry {
	if sbYear == "" || sbBenefit == "" {
		return []models.SBScheduleEntry{}
	}
	years := strings.Split(sbYear, "|")
	benefits := strings.Split(sbBenefit, "|")
	entries := make([]models.SBScheduleEntry, 0, len(years))
	for i, y := range years {
		yr, err := strconv.Atoi(strings.TrimSpace(y))
		if err != nil {
			continue
		}
		pct := 0.0
		if i < len(benefits) {
			pct, _ = strconv.ParseFloat(strings.TrimSpace(benefits[i]), 64)
		}
		entries = append(entries, models.SBScheduleEntry{Year: yr, BenefitPct: pct})
	}
	return entries
}
