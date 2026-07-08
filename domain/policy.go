package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// PolicyService handles policy business logic.
type PolicyService struct {
	policyRepo *repository.PolicyRepo
	planRepo   *repository.PlanRepo
	fupRepo    *repository.FUPRepo
	loanRepo   *repository.LoanRepo
	sbRepo     *repository.SBRepo
	clientRepo *repository.ClientRepo
}

// NewPolicyService creates a new PolicyService.
func NewPolicyService(
	pr *repository.PolicyRepo,
	planR *repository.PlanRepo,
	fr *repository.FUPRepo,
	lr *repository.LoanRepo,
	sbR *repository.SBRepo,
	cr *repository.ClientRepo,
) *PolicyService {
	return &PolicyService{
		policyRepo: pr,
		planRepo:   planR,
		fupRepo:    fr,
		loanRepo:   lr,
		sbRepo:     sbR,
		clientRepo: cr,
	}
}

// List returns paginated policies with filters.
// TODO(LIC-API — Phase 2): Sync all policies from GET /api/v2/policies/list on login.
// New policies sold via ANANDA portal arrive via POST /api/v2/newbusiness/update.
// velead_id links to ANANDA campaign lead; custid links to LIC customer portal ID.
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

	sbs, err := s.sbRepo.ListByPolicy(policyNo)
	if err != nil {
		return nil, err
	}

	var planDetails *models.PlanResponse
	if plan, err := s.planRepo.FindByNo(strconv.Itoa(policy.Plan)); err == nil {
		pr := buildPlanResponse(*plan)
		planDetails = &pr
	}

	return &models.PolicyDetail{
		Policy:      *policy,
		PlanDetails: planDetails,
		FUPHistory:  history,
		Loans:       loans,
		SBRecords:   sbs,
	}, nil
}

// Create creates a new policy, validating no duplicate policy number.
// Computes `age` from the linked client's date of birth at issue date.
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
	if !matDate.After(issueDate) {
		return errors.New("mat_date_before_issue_date")
	}
	if nextPremium.Before(issueDate) {
		return errors.New("next_premium_before_issue_date")
	}

	// Compute age from client DOB at issue date
	age := 0
	client, err := s.clientRepo.FindByFamilyAndPers(req.FamilyCode, req.PersCode)
	if err == nil && client.DOB != nil {
		age = ageAt(*client.DOB, issueDate)
	}

	return s.policyRepo.Create(req, plan.PlanName, issueDate, matDate, nextPremium, age)
}

// Update applies selective updates to a policy.
// TODO(LIC-API — Phase 2): Status (stat_cd, status) should be refreshed from
// GET /api/v2/policies/status?policyNo=:policyNo on each sync.
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

// ageAt returns the age in whole years of someone born on dob, measured at referenceDate.
func ageAt(dob, referenceDate time.Time) int {
	years := referenceDate.Year() - dob.Year()
	// Haven't had their birthday yet this year
	if referenceDate.YearDay() < dob.YearDay() {
		years--
	}
	if years < 0 {
		return 0
	}
	return years
}
