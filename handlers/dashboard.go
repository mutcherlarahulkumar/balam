package handlers

import (
	"agent-balam/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles the aggregated agent home-screen endpoint.
type DashboardHandler struct {
	svc *domain.DashboardService
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(svc *domain.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// Summary handles GET /dashboard.
func (h *DashboardHandler) Summary(c *gin.Context) {
	data, err := h.svc.Summary()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to load dashboard")
		return
	}
	respond(c, http.StatusOK, data)
}
