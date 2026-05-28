package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"fmt"
)

// FamilyService handles family business logic.
type FamilyService struct {
	familyRepo *repository.FamilyRepo
	clientRepo *repository.ClientRepo
	policyRepo *repository.PolicyRepo
}

// NewFamilyService creates a new FamilyService.
func NewFamilyService(fr *repository.FamilyRepo, cr *repository.ClientRepo, pr *repository.PolicyRepo) *FamilyService {
	return &FamilyService{familyRepo: fr, clientRepo: cr, policyRepo: pr}
}

// List returns paginated families with counts.
func (s *FamilyService) List(search string, page, limit int) ([]models.FamilyListItem, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.familyRepo.List(search, page, limit)
}

// GetByCode returns a family with members and policies.
func (s *FamilyService) GetByCode(familyCode string) (*models.FamilyDetail, error) {
	family, err := s.familyRepo.FindByCode(familyCode)
	if err != nil {
		return nil, errors.New("family_not_found")
	}

	clients, _, err := s.clientRepo.List("", familyCode, "", 1, 100)
	if err != nil {
		return nil, err
	}

	policies, err := s.policyRepo.PoliciesByFamilyCode(familyCode)
	if err != nil {
		return nil, err
	}

	return &models.FamilyDetail{
		Family:   *family,
		Members:  clients,
		Policies: policies,
	}, nil
}

// Create creates a new family, rejecting duplicate codes.
func (s *FamilyService) Create(req models.CreateFamilyRequest) (*models.Family, error) {
	if req.FamilyCode != "" {
		exists, err := s.familyRepo.CodeExists(req.FamilyCode)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("family_code_taken")
		}
	}
	return s.familyRepo.Create(req)
}

// Update updates a family, returning 404 if it doesn't exist.
func (s *FamilyService) Update(familyCode string, req models.UpdateFamilyRequest) (*models.Family, error) {
	_, err := s.familyRepo.FindByCode(familyCode)
	if err != nil {
		return nil, errors.New("family_not_found")
	}
	if err := s.familyRepo.Update(familyCode, req); err != nil {
		return nil, err
	}
	return s.familyRepo.FindByCode(familyCode)
}
