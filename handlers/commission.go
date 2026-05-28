package handlers

import (
	"agent-balam/domain"
	"agent-balam/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommissionHandler handles commission endpoints.
type CommissionHandler struct {
	svc *domain.CommissionService
}

// NewCommissionHandler creates a new CommissionHandler.
func NewCommissionHandler(svc *domain.CommissionService) *CommissionHandler {
	return &CommissionHandler{svc: svc}
}

// List handles GET /commission.
func (h *CommissionHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Query("year"), c.Query("month"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to fetch commissions")
		return
	}
	respond(c, http.StatusOK, gin.H{"data": items})
}

// Create handles POST /commission.
func (h *CommissionHandler) Create(c *gin.Context) {
	var req models.CreateCommissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.svc.Create(req); err != nil {
		switch err.Error() {
		case "policy_not_found":
			respondError(c, http.StatusNotFound, "policy_not_found", "Policy not found")
		default:
			respondError(c, http.StatusBadRequest, "create_failed", err.Error())
		}
		return
	}
	respond(c, http.StatusCreated, gin.H{"message": "Commission record added"})
}

// Summary handles GET /commission/summary.
func (h *CommissionHandler) Summary(c *gin.Context) {
	summary, err := h.svc.Summary()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to fetch summary")
		return
	}
	respond(c, http.StatusOK, summary)
}

// Calculate handles GET /commission/calculate.
func (h *CommissionHandler) Calculate(c *gin.Context) {
	policyNo, err := strconv.Atoi(c.Query("policyNo"))
	if err != nil || policyNo <= 0 {
		respondError(c, http.StatusBadRequest, "invalid_policy_no", "policyNo must be a positive integer")
		return
	}
	year, _ := strconv.Atoi(c.DefaultQuery("year", "1"))

	estimate, err := h.svc.Calculate(policyNo, year)
	if err != nil {
		switch err.Error() {
		case "policy_not_found":
			respondError(c, http.StatusNotFound, "policy_not_found", "Policy not found")
		default:
			respondError(c, http.StatusInternalServerError, "server_error", "Calculation failed")
		}
		return
	}
	respond(c, http.StatusOK, estimate)
}
