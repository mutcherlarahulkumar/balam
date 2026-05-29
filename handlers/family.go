package handlers

import (
	"agent-balam/domain"
	"agent-balam/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FamilyHandler handles family endpoints.
type FamilyHandler struct {
	svc *domain.FamilyService
}

// NewFamilyHandler creates a new FamilyHandler.
func NewFamilyHandler(svc *domain.FamilyService) *FamilyHandler {
	return &FamilyHandler{svc: svc}
}

// List handles GET /families.
func (h *FamilyHandler) List(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	items, total, err := h.svc.List(search, page, limit)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to list families")
		return
	}
	respond(c, http.StatusOK, models.PaginatedResponse{Data: items, Total: total, Page: page, Limit: limit})
}

// Get handles GET /families/:familyCode.
func (h *FamilyHandler) Get(c *gin.Context) {
	detail, err := h.svc.GetByCode(c.Param("familyCode"))
	if err != nil {
		respondError(c, http.StatusNotFound, "family_not_found", "Family not found")
		return
	}
	respond(c, http.StatusOK, detail)
}

// Create handles POST /families.
func (h *FamilyHandler) Create(c *gin.Context) {
	var req models.CreateFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidation(c, err)
		return
	}

	family, err := h.svc.Create(req)
	if err != nil {
		switch err.Error() {
		case "family_code_taken":
			respondError(c, http.StatusConflict, "family_code_taken", "A family with this code already exists")
		default:
			respondError(c, http.StatusInternalServerError, "server_error", "Failed to create family")
		}
		return
	}
	respond(c, http.StatusCreated, family)
}

// Update handles PUT /families/:familyCode.
func (h *FamilyHandler) Update(c *gin.Context) {
	var req models.UpdateFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidation(c, err)
		return
	}

	family, err := h.svc.Update(c.Param("familyCode"), req)
	if err != nil {
		switch err.Error() {
		case "family_not_found":
			respondError(c, http.StatusNotFound, "family_not_found", "Family not found")
		default:
			respondError(c, http.StatusInternalServerError, "server_error", "Failed to update family")
		}
		return
	}
	respond(c, http.StatusOK, family)
}
