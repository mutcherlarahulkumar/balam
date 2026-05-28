package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"time"
)

// SBService handles survival benefit business logic.
type SBService struct {
	sbRepo     *repository.SBRepo
	policyRepo *repository.PolicyRepo
}

// NewSBService creates a new SBService.
func NewSBService(sr *repository.SBRepo, pr *repository.PolicyRepo) *SBService {
	return &SBService{sbRepo: sr, policyRepo: pr}
}

// List returns SB records.
func (s *SBService) List(year string, unpaidOnly bool) ([]models.SB, error) {
	return s.sbRepo.List(year, unpaidOnly)
}

// Create validates the policy and inserts an SB record.
func (s *SBService) Create(req models.CreateSBRequest) error {
	if _, err := s.policyRepo.FindByNo(req.PolicyNo); err != nil {
		return errors.New("policy_not_found")
	}
	dueDate, err := time.Parse("2006-01-02", req.SBDueDate)
	if err != nil {
		return errors.New("invalid_sb_due_date")
	}
	return s.sbRepo.Create(req, dueDate)
}

// MarkPaid marks an SB as paid.
func (s *SBService) MarkPaid(id int, req models.MarkSBPaidRequest) error {
	sb, err := s.sbRepo.FindByID(id)
	if err != nil {
		return errors.New("sb_not_found")
	}
	if sb.SBPayDate != nil {
		return errors.New("sb_already_paid")
	}
	paidDate, err := time.Parse("2006-01-02", req.PaidDate)
	if err != nil {
		return errors.New("invalid_paid_date")
	}
	return s.sbRepo.MarkPaid(id, paidDate, req.ChequeNo)
}
