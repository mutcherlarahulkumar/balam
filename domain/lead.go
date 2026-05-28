package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"time"
)

// LeadService handles leads and activity business logic.
type LeadService struct {
	leadRepo   *repository.LeadRepo
	clientRepo *repository.ClientRepo
}

// NewLeadService creates a new LeadService.
func NewLeadService(lr *repository.LeadRepo, cr *repository.ClientRepo) *LeadService {
	return &LeadService{leadRepo: lr, clientRepo: cr}
}

// ListLeads returns all leads.
func (s *LeadService) ListLeads() ([]models.Lead, error) {
	return s.leadRepo.ListLeads()
}

// CreateLead inserts a new lead.
func (s *LeadService) CreateLead(req models.CreateLeadRequest) (*models.Lead, error) {
	return s.leadRepo.CreateLead(req)
}

// ListActivities returns activities, optionally scoped to a client.
func (s *LeadService) ListActivities(clientID int) ([]models.Activity, error) {
	return s.leadRepo.ListActivities(clientID)
}

// TodayActivities returns today's pending reminders.
func (s *LeadService) TodayActivities() ([]models.Activity, error) {
	return s.leadRepo.TodayActivities()
}

// CreateActivity validates the client exists and inserts an activity.
func (s *LeadService) CreateActivity(req models.CreateActivityRequest) (*models.Activity, error) {
	if _, err := s.clientRepo.FindByID(req.ClientID); err != nil {
		return nil, errors.New("client_not_found")
	}

	actDate, err := time.Parse("2006-01-02", req.ActivityDate)
	if err != nil {
		return nil, errors.New("invalid_activity_date")
	}

	var actTime, remDate, remTime *time.Time
	if req.ActivityTime != "" {
		t, err := time.Parse("15:04:05", req.ActivityTime)
		if err != nil {
			return nil, errors.New("invalid_activity_time: use HH:MM:SS")
		}
		actTime = &t
	}
	if req.ReminderDate != "" {
		t, err := time.Parse("2006-01-02", req.ReminderDate)
		if err != nil {
			return nil, errors.New("invalid_reminder_date")
		}
		remDate = &t
	}
	if req.ReminderTime != "" {
		t, err := time.Parse("15:04:05", req.ReminderTime)
		if err != nil {
			return nil, errors.New("invalid_reminder_time: use HH:MM:SS")
		}
		remTime = &t
	}

	return s.leadRepo.CreateActivity(req, &actDate, actTime, remDate, remTime)
}
