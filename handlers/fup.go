package handlers

import (
	"agent-balam/domain"
	"agent-balam/middlewares"
	"agent-balam/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FUPHandler handles FUP endpoints.
type FUPHandler struct {
	svc       *domain.FUPService
	agentSvc  *domain.AgentService
}

// NewFUPHandler creates a new FUPHandler.
func NewFUPHandler(svc *domain.FUPService, agentSvc *domain.AgentService) *FUPHandler {
	return &FUPHandler{svc: svc, agentSvc: agentSvc}
}

// Due handles GET /fup/due.
// Query params: year (YYYY), month (1-12), overdueDays (int).
func (h *FUPHandler) Due(c *gin.Context) {
	overdueDays, _ := strconv.Atoi(c.Query("overdueDays"))
	items, err := h.svc.DuePolicies(c.Query("year"), c.Query("month"), overdueDays)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to fetch due policies")
		return
	}
	respond(c, http.StatusOK, gin.H{"data": items})
}

// Update handles POST /fup/update.
func (h *FUPHandler) Update(c *gin.Context) {
	var req models.UpdateFUPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidation(c, err)
		return
	}

	agentID := middlewares.AgentID(c)
	profile, err := h.agentSvc.Profile(agentID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to fetch agent")
		return
	}

	if err := h.svc.UpdateFUP(req, profile.Name); err != nil {
		switch err.Error() {
		case "policy_not_found":
			respondError(c, http.StatusNotFound, "policy_not_found", "Policy not found")
		case "new_fup_not_after_old_fup":
			respondError(c, http.StatusUnprocessableEntity, "new_fup_not_after_old_fup", "New FUP date must be later than the old FUP date")
		case "old_fup_mismatch":
			respondError(c, http.StatusUnprocessableEntity, "fup_mismatch", "Provided oldFup does not match the current next premium on the policy")
		default:
			respondError(c, http.StatusBadRequest, "fup_update_failed", err.Error())
		}
		return
	}
	respond(c, http.StatusOK, gin.H{"message": "FUP updated successfully"})
}

// MultipleDues handles GET /fup/multipledue/:policyNo — returns instalment arrears.
func (h *FUPHandler) MultipleDues(c *gin.Context) {
	policyNo, err := strconv.Atoi(c.Param("policyNo"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_policy_no", "Policy number must be an integer")
		return
	}

	dues, err := h.svc.MultipleDues(policyNo)
	if err != nil {
		switch err.Error() {
		case "policy_not_found":
			respondError(c, http.StatusNotFound, "policy_not_found", "Policy not found")
		default:
			respondError(c, http.StatusInternalServerError, "server_error", "Failed to fetch multiple dues")
		}
		return
	}
	respond(c, http.StatusOK, gin.H{"policyNo": policyNo, "dues": dues})
}

// History handles GET /fup/history/:policyNo.
func (h *FUPHandler) History(c *gin.Context) {
	policyNo, err := strconv.Atoi(c.Param("policyNo"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_policy_no", "Policy number must be an integer")
		return
	}

	history, err := h.svc.History(policyNo)
	if err != nil {
		switch err.Error() {
		case "policy_not_found":
			respondError(c, http.StatusNotFound, "policy_not_found", "Policy not found")
		default:
			respondError(c, http.StatusInternalServerError, "server_error", "Failed to fetch FUP history")
		}
		return
	}
	respond(c, http.StatusOK, gin.H{"policyNo": policyNo, "history": history})
}
