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
func (s *FUPService) DuePolicies(year, month string, overdueDays int) ([]models.FUPDueItem, error) {
	return s.fupRepo.DuePolicies(year, month, overdueDays)
}

// UpdateFUP validates the old FUP and applies the update transactionally.
// TODO(LIC-API — Phase 2): Sync queue fupupdate to be drained by pushing each row to
// POST /api/v2/fup/update on LIC server. Delete from queue on success.
// LIC server is the authority on FUP — conflicts resolved in LIC's favour.
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

	if !newFUP.After(oldFUP) {
		return errors.New("new_fup_not_after_old_fup")
	}
	if policy.NextPremium != nil && !policy.NextPremium.Truncate(24*time.Hour).Equal(oldFUP.Truncate(24*time.Hour)) {
		return errors.New("old_fup_mismatch")
	}

	fupStatus := ComputeFUPStatus(newFUP, policy.PaymentMode, 180)
	return s.fupRepo.UpdateFUP(req.PolicyNo, oldFUP, newFUP, agentName, fupStatus)
}

// History returns the full FUP audit trail for a policy.
func (s *FUPService) History(policyNo int) ([]models.FUPHistory, error) {
	if _, err := s.policyRepo.FindByNo(policyNo); err != nil {
		return nil, errors.New("policy_not_found")
	}
	return s.fupRepo.History(policyNo)
}

// MultipleDues returns all instalment arrears for a policy from the multipledue table.
func (s *FUPService) MultipleDues(policyNo int) ([]models.MultipleDue, error) {
	if _, err := s.policyRepo.FindByNo(policyNo); err != nil {
		return nil, errors.New("policy_not_found")
	}
	return s.fupRepo.MultipleDues(policyNo)
}

// ComputeFUPStatus derives the fupstatus string from the new FUP date and lapse days.
// Status ladder:
//
//	PAID    → newFUP is in the future (premium not yet due)
//	DUE     → premium is due today or overdue but within grace period
//	OVERDUE → grace period exceeded but policy not yet lapsed
//	LAPSED  → past lapse date
func ComputeFUPStatus(newFUP time.Time, paymentMode string, defaultLapsDays int) string {
	now := time.Now().Truncate(24 * time.Hour)
	fup := newFUP.Truncate(24 * time.Hour)

	if fup.After(now) {
		return "PAID"
	}

	lapsDays := defaultLapsDays
	// Term / Health plans lapse after 15 days — callers should pass plan-specific lapsdays.
	lapseDate := fup.AddDate(0, 0, lapsDays)

	switch {
	case now.After(lapseDate):
		return "LAPSED"
	case now.After(fup):
		return "OVERDUE"
	default:
		// now == fup (today is the due date)
		return "DUE"
	}
}

// NextFUPDate advances a FUP date by one instalment period based on payment mode.
// Follows the rule: M=+1 month, Q=+3 months, H=+6 months, Y=+12 months.
func NextFUPDate(current time.Time, paymentMode string) time.Time {
	switch paymentMode {
	case "M":
		return current.AddDate(0, 1, 0)
	case "Q":
		return current.AddDate(0, 3, 0)
	case "H":
		return current.AddDate(0, 6, 0)
	default: // Y and anything else
		return current.AddDate(1, 0, 0)
	}
}
