package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"fmt"
	"strings"
	"time"
)

// GSTService handles GST calculation business logic.
type GSTService struct {
	policyRepo *repository.PolicyRepo
	planRepo   *repository.PlanRepo
}

// NewGSTService creates a new GSTService.
func NewGSTService(pr *repository.PolicyRepo, planR *repository.PlanRepo) *GSTService {
	return &GSTService{policyRepo: pr, planRepo: planR}
}

// gstCutoff is the date GST became 0% for all individual life insurance policies.
var gstCutoff = time.Date(2025, 9, 22, 0, 0, 0, 0, time.UTC)

// Calculate returns the GST breakdown for a given policy and premium year.
func (s *GSTService) Calculate(policyNo, premiumYear int) (*models.GSTCalculateResponse, error) {
	policy, err := s.policyRepo.FindByNo(policyNo)
	if err != nil {
		return nil, errors.New("policy_not_found")
	}

	planNo := fmt.Sprintf("%d", policy.Plan)
	plan, err := s.planRepo.FindByNo(planNo)
	if err != nil {
		return nil, errors.New("plan_not_found")
	}

	gstRate := gstRate(plan.PlanType, premiumYear, time.Now())
	gstAmount := policy.Premium * gstRate / 100

	return &models.GSTCalculateResponse{
		PolicyNo:       policyNo,
		PlanNo:         planNo,
		PlanType:       plan.PlanType,
		BasePremium:    policy.Premium,
		PremiumYear:    premiumYear,
		GSTRate:        gstRate,
		GSTAmount:      gstAmount,
		TotalPremium:   policy.Premium + gstAmount,
		Regulation:     "GST 0% on all individual life insurance policies effective 22 Sep 2025 (GST Council 54th meeting)",
		HistoricalNote: "Pre-Sep 2025: Endowment 4.5% Y1 / 2.25% renewal. Term/Health 18%. Annuity 1.8% single premium.",
	}, nil
}

// gstRate returns the applicable GST percentage given plan type, year, and the effective date.
func gstRate(planType string, year int, effective time.Time) float64 {
	if !effective.Before(gstCutoff) {
		return 0
	}
	pt := strings.ToLower(planType)
	switch {
	case strings.Contains(pt, "term") || strings.Contains(pt, "health"):
		return 18
	case strings.Contains(pt, "annuity"):
		return 1.8
	default:
		if year == 1 {
			return 4.5
		}
		return 2.25
	}
}
