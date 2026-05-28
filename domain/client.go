package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"time"
)

// ClientService handles client business logic.
type ClientService struct {
	clientRepo *repository.ClientRepo
	familyRepo *repository.FamilyRepo
	policyRepo *repository.PolicyRepo
}

// NewClientService creates a new ClientService.
func NewClientService(cr *repository.ClientRepo, fr *repository.FamilyRepo, pr *repository.PolicyRepo) *ClientService {
	return &ClientService{clientRepo: cr, familyRepo: fr, policyRepo: pr}
}

// List returns paginated clients.
func (s *ClientService) List(search, familyCode, clientType string, page, limit int) ([]models.Client, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.clientRepo.List(search, familyCode, clientType, page, limit)
}

// GetByID returns a client with linked data.
func (s *ClientService) GetByID(id int) (*models.ClientDetail, error) {
	client, err := s.clientRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("client_not_found")
	}

	policies, err := s.policyRepo.PoliciesByFamilyCode(client.FamilyCode)
	if err != nil {
		return nil, err
	}

	bank, err := s.clientRepo.BankDetails(id)
	if err != nil {
		return nil, err
	}

	docs, err := s.clientRepo.Documents(id)
	if err != nil {
		return nil, err
	}

	return &models.ClientDetail{
		Client:      *client,
		Policies:    policies,
		BankDetails: bank,
		Documents:   docs,
	}, nil
}

// Search searches clients by name, mobile, or policy number.
func (s *ClientService) Search(q string) ([]models.Client, error) {
	if q == "" {
		return []models.Client{}, nil
	}
	return s.clientRepo.Search(q)
}

// Create creates a new client, validating that the family exists.
func (s *ClientService) Create(req models.CreateClientRequest) (*models.Client, error) {
	if _, err := s.familyRepo.FindByCode(req.FamilyCode); err != nil {
		return nil, errors.New("family_not_found")
	}

	var dob *time.Time
	age := 0
	if req.DOB != "" {
		parsed, err := time.Parse("2006-01-02", req.DOB)
		if err != nil {
			return nil, errors.New("invalid_dob_format: use YYYY-MM-DD")
		}
		dob = &parsed
		age = calculateAge(parsed)
	}

	if req.ClientType == "" {
		req.ClientType = "N"
	}
	return s.clientRepo.Create(req, dob, age)
}

// Update updates a client.
func (s *ClientService) Update(id int, req models.UpdateClientRequest) (*models.Client, error) {
	existing, err := s.clientRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("client_not_found")
	}

	var dob *time.Time
	age := existing.Age
	if req.DOB != "" {
		parsed, err := time.Parse("2006-01-02", req.DOB)
		if err != nil {
			return nil, errors.New("invalid_dob_format: use YYYY-MM-DD")
		}
		dob = &parsed
		age = calculateAge(parsed)
	}

	if err := s.clientRepo.Update(id, req, dob, age); err != nil {
		return nil, err
	}
	return s.clientRepo.FindByID(id)
}

func calculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		years--
	}
	return years
}
