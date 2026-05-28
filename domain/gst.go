package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"fmt"
	"strconv"
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
// GST Council 54th meeting, effective 22 Sep 2025.
var gstCutoff = time.Date(2025, 9, 22, 0, 0, 0, 0, time.UTC)

// Calculate returns the GST breakdown for a given policy and premium year.
// Uses policy.issue_date (not today) to determine the correct GST regime:
//   - issue_date >= 22 Sep 2025 → 0%
//   - issue_date < 22 Sep 2025  → historical rate from plans.stax
//
// TODO: Store GST regime change history in DB (a gst_regime table with effective_date + rates)
// so premium receipts for old policies are historically accurate across future regime changes.
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

	// Determine reference date: use issue_date when available, otherwise today.
	referenceDate := time.Now()
	if policy.IssueDate != nil {
		referenceDate = *policy.IssueDate
	}

	var gstRate float64
	if !referenceDate.Before(gstCutoff) {
		// Post-cutoff: 0% GST for all individual policies
		gstRate = 0
	} else {
		// Pre-cutoff: parse plans.stax column ("firstYear/renewal") or fall back to plan-type heuristic
		gstRate = gstRateFromStax(plan.Stax, plan.PlanType, premiumYear)
	}

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

// gstRateFromStax parses the plans.stax field ("firstYear/renewal", e.g. "4.5/2.25") and
// returns the appropriate rate for the given premium year.
// Falls back to plan-type string matching when stax is absent or malformed.
func gstRateFromStax(stax, planType string, premiumYear int) float64 {
	stax = strings.TrimSpace(stax)
	if stax != "" {
		parts := strings.SplitN(stax, "/", 2)
		if len(parts) == 2 {
			firstYear, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			renewal, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				if premiumYear == 1 {
					return firstYear
				}
				return renewal
			}
		}
		// Single value (e.g. "18" for Term)
		if rate, err := strconv.ParseFloat(stax, 64); err == nil {
			return rate
		}
	}

	// Fallback: derive from plan type name
	return gstRateByPlanType(planType, premiumYear)
}

// gstRateByPlanType is the fallback when stax column is absent.
func gstRateByPlanType(planType string, premiumYear int) float64 {
	pt := strings.ToLower(planType)
	switch {
	case strings.Contains(pt, "term") || strings.Contains(pt, "health"):
		return 18
	case strings.Contains(pt, "annuity"):
		return 1.8
	default: // Endowment, ULIP, Money-Back, whole-life, children, etc.
		if premiumYear == 1 {
			return 4.5
		}
		return 2.25
	}
}
