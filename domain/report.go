package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
)

// ReportService handles report generation business logic.
type ReportService struct {
	reportRepo *repository.ReportRepo
	familyRepo *repository.FamilyRepo
}

// NewReportService creates a new ReportService.
func NewReportService(rr *repository.ReportRepo, fr *repository.FamilyRepo) *ReportService {
	return &ReportService{reportRepo: rr, familyRepo: fr}
}

// Cashflow returns the cashflow report for a family.
func (s *ReportService) Cashflow(familyCode string) ([]models.ReportCashflow, error) {
	if err := s.assertFamilyExists(familyCode); err != nil {
		return nil, err
	}
	return s.reportRepo.Cashflow(familyCode)
}

// CashInOut returns the cash-in vs cash-out report.
func (s *ReportService) CashInOut() ([]models.ReportCashInOut, error) {
	return s.reportRepo.CashInOut()
}

// CurrentStatus returns the current policy status snapshot for a family.
func (s *ReportService) CurrentStatus(familyCode string) ([]models.ReportCurrentStatus, error) {
	if err := s.assertFamilyExists(familyCode); err != nil {
		return nil, err
	}
	return s.reportRepo.CurrentStatus(familyCode)
}

// Calendar returns the 12-month premium calendar for a family.
func (s *ReportService) Calendar(familyCode string) ([]models.ReportCalendar, error) {
	if err := s.assertFamilyExists(familyCode); err != nil {
		return nil, err
	}
	return s.reportRepo.Calendar(familyCode)
}

// Refresh rebuilds all report caches for a family.
func (s *ReportService) Refresh(familyCode string) error {
	if err := s.assertFamilyExists(familyCode); err != nil {
		return err
	}
	if err := s.reportRepo.RefreshCashflow(familyCode); err != nil {
		return err
	}
	return s.reportRepo.RefreshCalendar(familyCode)
}

func (s *ReportService) assertFamilyExists(familyCode string) error {
	if _, err := s.familyRepo.FindByCode(familyCode); err != nil {
		return errors.New("family_not_found")
	}
	return nil
}
