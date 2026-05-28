package handlers

import (
	"agent-balam/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PlanHandler handles plan catalogue endpoints.
type PlanHandler struct {
	svc *domain.PlanService
}

// NewPlanHandler creates a new PlanHandler.
func NewPlanHandler(svc *domain.PlanService) *PlanHandler {
	return &PlanHandler{svc: svc}
}

// List handles GET /plans.
func (h *PlanHandler) List(c *gin.Context) {
	plans, err := h.svc.List(c.Query("type"), c.Query("search"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to fetch plans")
		return
	}
	respond(c, http.StatusOK, gin.H{"data": plans})
}
