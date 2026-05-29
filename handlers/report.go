package handlers

import (
	"agent-balam/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReportHandler handles report endpoints.
type ReportHandler struct {
	svc *domain.ReportService
}

// NewReportHandler creates a new ReportHandler.
func NewReportHandler(svc *domain.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// Cashflow handles GET /reports/cashflow.
func (h *ReportHandler) Cashflow(c *gin.Context) {
	familyCode := c.Query("familyCode")
	if familyCode == "" {
		respondError(c, http.StatusBadRequest, "missing_family_code", "familyCode query parameter is required")
		return
	}
	data, err := h.svc.Cashflow(familyCode)
	if err != nil {
		h.handleErr(c, err)
		return
	}
	respond(c, http.StatusOK, gin.H{"data": data})
}

// CashInOut handles GET /reports/cashinout.
func (h *ReportHandler) CashInOut(c *gin.Context) {
	data, err := h.svc.CashInOut()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to generate report")
		return
	}
	respond(c, http.StatusOK, gin.H{"data": data})
}

// CurrentStatus handles GET /reports/status.
func (h *ReportHandler) CurrentStatus(c *gin.Context) {
	familyCode := c.Query("familyCode")
	if familyCode == "" {
		respondError(c, http.StatusBadRequest, "missing_family_code", "familyCode query parameter is required")
		return
	}
	data, err := h.svc.CurrentStatus(familyCode)
	if err != nil {
		h.handleErr(c, err)
		return
	}
	respond(c, http.StatusOK, gin.H{"data": data})
}

// Calendar handles GET /reports/calendar.
func (h *ReportHandler) Calendar(c *gin.Context) {
	familyCode := c.Query("familyCode")
	if familyCode == "" {
		respondError(c, http.StatusBadRequest, "missing_family_code", "familyCode query parameter is required")
		return
	}
	data, err := h.svc.Calendar(familyCode)
	if err != nil {
		h.handleErr(c, err)
		return
	}
	respond(c, http.StatusOK, gin.H{"data": data})
}

// Refresh handles POST /reports/refresh.
func (h *ReportHandler) Refresh(c *gin.Context) {
	var body struct {
		FamilyCode string `json:"familyCode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondValidation(c, err)
		return
	}
	if err := h.svc.Refresh(body.FamilyCode); err != nil {
		h.handleErr(c, err)
		return
	}
	respond(c, http.StatusOK, gin.H{"message": "Reports refreshed successfully"})
}

func (h *ReportHandler) handleErr(c *gin.Context, err error) {
	if err.Error() == "family_not_found" {
		respondError(c, http.StatusNotFound, "family_not_found", "Family not found")
		return
	}
	respondError(c, http.StatusInternalServerError, "server_error", "Failed to generate report")
}
