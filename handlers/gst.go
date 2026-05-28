package handlers

import (
	"agent-balam/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GSTHandler handles GST calculator endpoints.
type GSTHandler struct {
	svc *domain.GSTService
}

// NewGSTHandler creates a new GSTHandler.
func NewGSTHandler(svc *domain.GSTService) *GSTHandler {
	return &GSTHandler{svc: svc}
}

// Calculate handles GET /gst/calculate.
func (h *GSTHandler) Calculate(c *gin.Context) {
	policyNo, err := strconv.Atoi(c.Query("policyNo"))
	if err != nil || policyNo <= 0 {
		respondError(c, http.StatusBadRequest, "invalid_policy_no", "policyNo must be a positive integer")
		return
	}
	premiumYear, _ := strconv.Atoi(c.DefaultQuery("premiumYear", "1"))

	result, err := h.svc.Calculate(policyNo, premiumYear)
	if err != nil {
		switch err.Error() {
		case "policy_not_found":
			respondError(c, http.StatusNotFound, "policy_not_found", "Policy not found")
		case "plan_not_found":
			respondError(c, http.StatusNotFound, "plan_not_found", "Plan not found for this policy")
		default:
			respondError(c, http.StatusInternalServerError, "server_error", "GST calculation failed")
		}
		return
	}
	respond(c, http.StatusOK, result)
}
