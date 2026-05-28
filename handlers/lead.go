package handlers

import (
	"agent-balam/domain"
	"agent-balam/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LeadHandler handles lead and activity endpoints.
type LeadHandler struct {
	svc *domain.LeadService
}

// NewLeadHandler creates a new LeadHandler.
func NewLeadHandler(svc *domain.LeadService) *LeadHandler {
	return &LeadHandler{svc: svc}
}

// ListLeads handles GET /leads.
func (h *LeadHandler) ListLeads(c *gin.Context) {
	items, err := h.svc.ListLeads()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to fetch leads")
		return
	}
	respond(c, http.StatusOK, gin.H{"data": items})
}

// CreateLead handles POST /leads.
func (h *LeadHandler) CreateLead(c *gin.Context) {
	var req models.CreateLeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}
	lead, err := h.svc.CreateLead(req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to create lead")
		return
	}
	respond(c, http.StatusCreated, lead)
}

// ListActivities handles GET /activities.
func (h *LeadHandler) ListActivities(c *gin.Context) {
	clientID, _ := strconv.Atoi(c.Query("clientId"))
	items, err := h.svc.ListActivities(clientID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to fetch activities")
		return
	}
	respond(c, http.StatusOK, gin.H{"data": items})
}

// TodayActivities handles GET /activities/today.
func (h *LeadHandler) TodayActivities(c *gin.Context) {
	items, err := h.svc.TodayActivities()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to fetch today's activities")
		return
	}
	respond(c, http.StatusOK, gin.H{"data": items})
}

// CreateActivity handles POST /activities.
func (h *LeadHandler) CreateActivity(c *gin.Context) {
	var req models.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}
	activity, err := h.svc.CreateActivity(req)
	if err != nil {
		switch err.Error() {
		case "client_not_found":
			respondError(c, http.StatusNotFound, "client_not_found", "Client not found")
		default:
			respondError(c, http.StatusBadRequest, "create_failed", err.Error())
		}
		return
	}
	respond(c, http.StatusCreated, activity)
}
