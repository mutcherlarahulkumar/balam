package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"time"
)

// LoanService handles loan business logic.
type LoanService struct {
	loanRepo   *repository.LoanRepo
	policyRepo *repository.PolicyRepo
}

// NewLoanService creates a new LoanService.
func NewLoanService(lr *repository.LoanRepo, pr *repository.PolicyRepo) *LoanService {
	return &LoanService{loanRepo: lr, policyRepo: pr}
}

// List returns all loans.
func (s *LoanService) List() ([]models.Loan, error) {
	return s.loanRepo.List(0)
}

// Create validates the policy and inserts a loan record.
func (s *LoanService) Create(req models.CreateLoanRequest) error {
	if _, err := s.policyRepo.FindByNo(req.PolicyNo); err != nil {
		return errors.New("policy_not_found")
	}
	loanDate, err := time.Parse("2006-01-02", req.LoanDate)
	if err != nil {
		return errors.New("invalid_loan_date")
	}
	intDueDate, err := time.Parse("2006-01-02", req.IntDueDate)
	if err != nil {
		return errors.New("invalid_interest_due_date")
	}
	if intDueDate.Before(loanDate) {
		return errors.New("interest_due_before_loan_date")
	}
	return s.loanRepo.Create(req, loanDate, intDueDate)
}
