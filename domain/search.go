package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
)

// SearchService handles cross-entity global search.
type SearchService struct {
	familyRepo *repository.FamilyRepo
	clientRepo *repository.ClientRepo
	policyRepo *repository.PolicyRepo
}

// NewSearchService creates a new SearchService.
func NewSearchService(fr *repository.FamilyRepo, cr *repository.ClientRepo, pr *repository.PolicyRepo) *SearchService {
	return &SearchService{familyRepo: fr, clientRepo: cr, policyRepo: pr}
}

// Search looks up q across families (head name/code), clients (name/mobile), and
// policies (policy number) in one call, so the frontend can offer a single search box.
func (s *SearchService) Search(q string) (*models.GlobalSearchResponse, error) {
	families, _, err := s.familyRepo.List(q, 1, 10)
	if err != nil {
		return nil, err
	}

	clients, err := s.clientRepo.Search(q)
	if err != nil {
		return nil, err
	}
	if len(clients) > 10 {
		clients = clients[:10]
	}

	policies, err := s.policyRepo.SearchByNo(q)
	if err != nil {
		return nil, err
	}

	return &models.GlobalSearchResponse{
		Families: families,
		Clients:  clients,
		Policies: policies,
	}, nil
}
