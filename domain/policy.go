package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"fmt"
	"time"
)

// PolicyService handles policy business logic.
type PolicyService struct {
	policyRepo *repository.PolicyRepo
	planRepo   *repository.PlanRepo
	fupRepo    *repository.FUPRepo
	loanRepo   *repository.LoanRepo
	sbRepo     *repository.SBRepo
}

// NewPolicyService creates a new PolicyService.
func NewPolicyService(pr *repository.PolicyRepo, planR *repository.PlanRepo, fr *repository.FUPRepo, lr *repository.LoanRepo, sbR *repository.SBRepo) *PolicyService {
	return &PolicyService{policyRepo: pr, planRepo: planR, fupRepo: fr, loanRepo: lr, sbRepo: sbR}
}

// List returns paginated policies with filters.
func (s *PolicyService) List(status, familyCode string, dueThisMonth bool, lapsingIn, page, limit int) ([]models.PolicyListItem, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.policyRepo.List(repository.ListFilter{
		Status:       status,
		FamilyCode:   familyCode,
		DueThisMonth: dueThisMonth,
		LapsingIn:    lapsingIn,
		Page:         page,
		Limit:        limit,
	})
}

// GetByNo returns a full policy detail.
func (s *PolicyService) GetByNo(policyNo int) (*models.PolicyDetail, error) {
	policy, err := s.policyRepo.FindByNo(policyNo)
	if err != nil {
		return nil, errors.New("policy_not_found")
	}

	history, err := s.fupRepo.History(policyNo)
	if err != nil {
		return nil, err
	}

	loans, err := s.loanRepo.List(policyNo)
	if err != nil {
		return nil, err
	}

	sbs, err := s.sbRepo.List("", false)
	if err != nil {
		return nil, err
	}
	var policySBs []models.SB
	for _, sb := range sbs {
		if sb.PolicyNo == policyNo {
			policySBs = append(policySBs, sb)
		}
	}
	if policySBs == nil {
		policySBs = []models.SB{}
	}

	return &models.PolicyDetail{
		Policy:     *policy,
		FUPHistory: history,
		Loans:      loans,
		SBRecords:  policySBs,
	}, nil
}

// Create creates a new policy, validating no duplicate policy number.
func (s *PolicyService) Create(req models.CreatePolicyRequest) error {
	exists, err := s.policyRepo.Exists(req.PolicyNo)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("policy_no_conflict")
	}

	plan, err := s.planRepo.FindByNo(req.PlanNo)
	if err != nil {
		return fmt.Errorf("plan_not_found")
	}

	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		return errors.New("invalid_issue_date")
	}
	matDate, err := time.Parse("2006-01-02", req.MatDate)
	if err != nil {
		return errors.New("invalid_mat_date")
	}
	nextPremium, err := time.Parse("2006-01-02", req.NextPremium)
	if err != nil {
		return errors.New("invalid_next_premium")
	}

	return s.policyRepo.Create(req, plan.PlanName, issueDate, matDate, nextPremium)
}

// Update applies selective updates to a policy.
func (s *PolicyService) Update(policyNo int, req models.UpdatePolicyRequest) (*models.Policy, error) {
	if _, err := s.policyRepo.FindByNo(policyNo); err != nil {
		return nil, errors.New("policy_not_found")
	}

	var nextPremium *time.Time
	if req.NextPremium != "" {
		t, err := time.Parse("2006-01-02", req.NextPremium)
		if err != nil {
			return nil, errors.New("invalid_next_premium")
		}
		nextPremium = &t
	}

	if err := s.policyRepo.Update(policyNo, req, nextPremium); err != nil {
		return nil, err
	}
	return s.policyRepo.FindByNo(policyNo)
}
