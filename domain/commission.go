package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"fmt"
	"time"
)

// noBonusPlans lists plan numbers where bonus commission is not payable.
// Spec: plans 939, 940, 951, 860 and all single-premium policies (paymentMode=S).
var noBonusPlans = map[string]bool{"939": true, "940": true, "951": true, "860": true}

// CommissionService handles commission business logic.
type CommissionService struct {
	commRepo   *repository.CommissionRepo
	policyRepo *repository.PolicyRepo
}

// NewCommissionService creates a new CommissionService.
func NewCommissionService(cr *repository.CommissionRepo, pr *repository.PolicyRepo) *CommissionService {
	return &CommissionService{commRepo: cr, policyRepo: pr}
}

// List returns commission records filtered by year and month.
func (s *CommissionService) List(year, month string) ([]models.Commission, error) {
	return s.commRepo.List(year, month)
}

// Create inserts a new commission record, validating the policy exists.
// TODO(LIC-API — Phase 2): Sync actual paid commission from GET /api/v2/commission/list.
// Manual entry is only for MVP estimation.
func (s *CommissionService) Create(req models.CreateCommissionRequest) error {
	if _, err := s.policyRepo.FindByNo(req.PolicyNo); err != nil {
		return errors.New("policy_not_found")
	}
	billDate, err := time.Parse("2006-01-02", req.BillDate)
	if err != nil {
		return errors.New("invalid_bill_date")
	}
	var payDate *time.Time
	if req.PayDate != "" {
		t, err := time.Parse("2006-01-02", req.PayDate)
		if err != nil {
			return errors.New("invalid_pay_date")
		}
		payDate = &t
	}
	return s.commRepo.Create(req, &billDate, payDate)
}

// Summary returns monthly and yearly commission aggregations.
func (s *CommissionService) Summary() (*models.CommissionSummary, error) {
	return s.commRepo.Summary()
}

// Calculate estimates commission for a given policy and year.
// Bonus commission is excluded for plans 939, 940, 951, 860 and single-premium policies.
// TODO: Build a commission_rates lookup table mapping planNo + term + year → exact commission %.
// Currently approximate. Exact rates require LIC's internal commission chart per plan.
func (s *CommissionService) Calculate(policyNo, year int) (*models.CommissionEstimate, error) {
	policy, err := s.policyRepo.FindByNo(policyNo)
	if err != nil {
		return nil, errors.New("policy_not_found")
	}

	planNo := fmt.Sprintf("%d", policy.Plan)
	isSinglePremium := policy.PaymentMode == "S"
	bonusEligible := !isSinglePremium && !noBonusPlans[planNo]

	basePct, bonusPct := commissionRates(year, isSinglePremium, bonusEligible)
	estimated := policy.Premium * (basePct + bonusPct) / 100

	note := "Estimate only. Actual amount set by LIC billing."
	if !bonusEligible {
		note += " Bonus commission not applicable for this plan/payment mode."
	}

	return &models.CommissionEstimate{
		PolicyNo:            policyNo,
		Premium:             policy.Premium,
		PlanNo:              planNo,
		Year:                year,
		BaseCommissionPct:   basePct,
		BonusCommissionPct:  bonusPct,
		TotalPct:            basePct + bonusPct,
		EstimatedCommission: estimated,
		Note:                note,
	}, nil
}

// commissionRates returns base and bonus commission percentages by policy year.
// Rates per spec (effective Oct 2024):
//
//	Year 1: 20% base + 8% bonus (= 28% total with bonus)
//	Year 2-3: 7.5% base
//	Year 4+: 5% base
//	Single premium: 2% flat, no bonus
//
// Bonus commission = up to 40% of first-year base commission (40% × 20% = 8%).
func commissionRates(year int, isSinglePremium, bonusEligible bool) (base, bonus float64) {
	if isSinglePremium {
		return 2, 0
	}
	switch {
	case year == 1:
		b := 0.0
		if bonusEligible {
			b = 8 // 40% of 20% base
		}
		return 20, b
	case year <= 3:
		return 7.5, 0
	default:
		return 5, 0
	}
}
