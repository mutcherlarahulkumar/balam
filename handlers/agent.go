package handlers

import (
	"agent-balam/domain"
	"agent-balam/middlewares"
	"agent-balam/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AgentHandler handles agent profile endpoints.
type AgentHandler struct {
	svc *domain.AgentService
}

// NewAgentHandler creates a new AgentHandler.
func NewAgentHandler(svc *domain.AgentService) *AgentHandler {
	return &AgentHandler{svc: svc}
}

// GetProfile handles GET /agent/profile.
func (h *AgentHandler) GetProfile(c *gin.Context) {
	profile, err := h.svc.Profile(middlewares.AgentID(c))
	if err != nil {
		respondError(c, http.StatusNotFound, "agent_not_found", "Agent not found")
		return
	}
	respond(c, http.StatusOK, profile)
}

// UpdateProfile handles PUT /agent/profile.
func (h *AgentHandler) UpdateProfile(c *gin.Context) {
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	profile, err := h.svc.UpdateProfile(middlewares.AgentID(c), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to update profile")
		return
	}
	respond(c, http.StatusOK, profile)
}
