package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"fmt"
	"time"
)

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
func (s *CommissionService) Calculate(policyNo, year int) (*models.CommissionEstimate, error) {
	policy, err := s.policyRepo.FindByNo(policyNo)
	if err != nil {
		return nil, errors.New("policy_not_found")
	}

	basePct, bonusPct := commissionRates(year)
	estimated := policy.Premium * (basePct + bonusPct) / 100

	return &models.CommissionEstimate{
		PolicyNo:            policyNo,
		Premium:             policy.Premium,
		PlanNo:              fmt.Sprintf("%d", policy.Plan),
		Year:                year,
		BaseCommissionPct:   basePct,
		BonusCommissionPct:  bonusPct,
		TotalPct:            basePct + bonusPct,
		EstimatedCommission: estimated,
		Note:                "Estimate only. Actual amount set by LIC billing.",
	}, nil
}

// commissionRates returns base and bonus commission percentages by policy year.
func commissionRates(year int) (base, bonus float64) {
	switch {
	case year == 1:
		return 20, 8
	case year <= 3:
		return 7.5, 0
	default:
		return 5, 0
	}
}
