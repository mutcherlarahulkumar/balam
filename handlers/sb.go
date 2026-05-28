package handlers

import (
	"agent-balam/domain"
	"agent-balam/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SBHandler handles survival benefit endpoints.
type SBHandler struct {
	svc *domain.SBService
}

// NewSBHandler creates a new SBHandler.
func NewSBHandler(svc *domain.SBService) *SBHandler {
	return &SBHandler{svc: svc}
}

// List handles GET /sb.
func (h *SBHandler) List(c *gin.Context) {
	unpaidOnly := c.Query("unpaidOnly") == "true"
	items, err := h.svc.List(c.Query("year"), unpaidOnly)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to fetch SB records")
		return
	}
	respond(c, http.StatusOK, gin.H{"data": items})
}

// Create handles POST /sb.
func (h *SBHandler) Create(c *gin.Context) {
	var req models.CreateSBRequest
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
	respond(c, http.StatusCreated, gin.H{"message": "SB record added"})
}

// MarkPaid handles PUT /sb/:id/mark-paid.
func (h *SBHandler) MarkPaid(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_id", "SB ID must be an integer")
		return
	}

	var req models.MarkSBPaidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.svc.MarkPaid(id, req); err != nil {
		switch err.Error() {
		case "sb_not_found":
			respondError(c, http.StatusNotFound, "sb_not_found", "SB record not found")
		case "sb_already_paid":
			respondError(c, http.StatusConflict, "sb_already_paid", "SB has already been marked as paid")
		default:
			respondError(c, http.StatusBadRequest, "mark_paid_failed", err.Error())
		}
		return
	}
	respond(c, http.StatusOK, gin.H{"message": "SB marked as paid"})
}
