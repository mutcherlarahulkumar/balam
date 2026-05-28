package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
)

// AgentService handles agent profile business logic.
type AgentService struct {
	repo *repository.AgentRepo
}

// NewAgentService creates a new AgentService.
func NewAgentService(repo *repository.AgentRepo) *AgentService {
	return &AgentService{repo: repo}
}

// Profile returns the public profile for an agent.
func (s *AgentService) Profile(agentID int) (*models.AgentPublic, error) {
	agent, err := s.repo.FindByID(agentID)
	if err != nil {
		return nil, errors.New("agent_not_found")
	}
	return &models.AgentPublic{
		ID:          agent.ID,
		Name:        agent.Name,
		AgentCode:   agent.Login,
		Branch:      agent.Branch,
		Club:        agent.Club,
		LicenceNo:   agent.LicenceNo,
		AgSince:     agent.AgSince,
		RenewalDate: agent.RenewalDate,
		PAN:         agent.PAN,
		Mobile:      agent.Mobile,
		Email:       agent.Email,
		Photo:       agent.Photo,
		Slogan:      agent.Slogan,
		NewPortal:   agent.NewPortal,
	}, nil
}

// UpdateProfile updates editable profile fields, preserving existing values where the request field is empty.
func (s *AgentService) UpdateProfile(agentID int, req models.UpdateProfileRequest) (*models.AgentPublic, error) {
	existing, err := s.repo.FindByID(agentID)
	if err != nil {
		return nil, errors.New("agent_not_found")
	}

	name := coalesce(req.Name, existing.Name)
	mobile := coalesce(req.Mobile, existing.Mobile)
	email := coalesce(req.Email, existing.Email)
	photo := coalesce(req.Photo, existing.Photo)
	slogan := coalesce(req.Slogan, existing.Slogan)
	address := coalesce(req.Address, existing.Address)

	if err := s.repo.UpdateProfile(agentID, name, mobile, email, photo, slogan, address); err != nil {
		return nil, err
	}
	return s.Profile(agentID)
}

func coalesce(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
