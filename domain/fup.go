package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"time"
)

// FUPService handles FUP business logic.
type FUPService struct {
	fupRepo    *repository.FUPRepo
	policyRepo *repository.PolicyRepo
}

// NewFUPService creates a new FUPService.
func NewFUPService(fr *repository.FUPRepo, pr *repository.PolicyRepo) *FUPService {
	return &FUPService{fupRepo: fr, policyRepo: pr}
}

// DuePolicies returns policies with due/overdue premiums.
func (s *FUPService) DuePolicies(month string, overdueDays int) ([]models.FUPDueItem, error) {
	return s.fupRepo.DuePolicies(month, overdueDays)
}

// UpdateFUP validates the old FUP and applies the update transactionally.
func (s *FUPService) UpdateFUP(req models.UpdateFUPRequest, agentName string) error {
	policy, err := s.policyRepo.FindByNo(req.PolicyNo)
	if err != nil {
		return errors.New("policy_not_found")
	}

	oldFUP, err := time.Parse("2006-01-02", req.OldFUP)
	if err != nil {
		return errors.New("invalid_old_fup: use YYYY-MM-DD")
	}
	newFUP, err := time.Parse("2006-01-02", req.NewFUP)
	if err != nil {
		return errors.New("invalid_new_fup: use YYYY-MM-DD")
	}

	if policy.NextPremium != nil && !policy.NextPremium.Truncate(24*time.Hour).Equal(oldFUP.Truncate(24*time.Hour)) {
		return errors.New("old_fup_mismatch")
	}

	return s.fupRepo.UpdateFUP(req.PolicyNo, oldFUP, newFUP, agentName)
}

// History returns the full FUP audit trail for a policy.
func (s *FUPService) History(policyNo int) ([]models.FUPHistory, error) {
	if _, err := s.policyRepo.FindByNo(policyNo); err != nil {
		return nil, errors.New("policy_not_found")
	}
	return s.fupRepo.History(policyNo)
}
