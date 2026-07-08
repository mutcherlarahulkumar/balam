package handlers

import (
	"agent-balam/domain"
	"agent-balam/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PolicyHandler handles policy endpoints.
type PolicyHandler struct {
	svc *domain.PolicyService
}

// NewPolicyHandler creates a new PolicyHandler.
func NewPolicyHandler(svc *domain.PolicyService) *PolicyHandler {
	return &PolicyHandler{svc: svc}
}

// List handles GET /policies.
func (h *PolicyHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	lapsingIn, _ := strconv.Atoi(c.Query("lapsingIn"))
	dueThisMonth := c.Query("dueThisMonth") == "true"

	items, total, err := h.svc.List(
		c.Query("status"), c.Query("familyCode"),
		dueThisMonth, lapsingIn, page, limit,
	)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to list policies")
		return
	}
	respond(c, http.StatusOK, models.PaginatedResponse{Data: items, Total: total, Page: page, Limit: limit})
}

// Get handles GET /policies/:policyNo.
func (h *PolicyHandler) Get(c *gin.Context) {
	policyNo, err := strconv.Atoi(c.Param("policyNo"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_policy_no", "Policy number must be an integer")
		return
	}

	detail, err := h.svc.GetByNo(policyNo)
	if err != nil {
		respondError(c, http.StatusNotFound, "policy_not_found", "Policy not found")
		return
	}
	respond(c, http.StatusOK, detail)
}

// Create handles POST /policies.
func (h *PolicyHandler) Create(c *gin.Context) {
	var req models.CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidation(c, err)
		return
	}

	if err := h.svc.Create(req); err != nil {
		switch err.Error() {
		case "policy_no_conflict":
			respondError(c, http.StatusConflict, "policy_no_conflict", "A policy with this number already exists")
		case "plan_not_found":
			respondError(c, http.StatusBadRequest, "plan_not_found", "Plan not found")
		case "mat_date_before_issue_date":
			respondError(c, http.StatusUnprocessableEntity, "mat_date_before_issue_date", "Maturity date must be after the issue date")
		case "next_premium_before_issue_date":
			respondError(c, http.StatusUnprocessableEntity, "next_premium_before_issue_date", "Next premium date cannot be earlier than the issue date")
		default:
			respondError(c, http.StatusBadRequest, "create_failed", err.Error())
		}
		return
	}
	respond(c, http.StatusCreated, gin.H{"message": "Policy created successfully", "policyNo": req.PolicyNo})
}

// Update handles PUT /policies/:policyNo.
func (h *PolicyHandler) Update(c *gin.Context) {
	policyNo, err := strconv.Atoi(c.Param("policyNo"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_policy_no", "Policy number must be an integer")
		return
	}

	var req models.UpdatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidation(c, err)
		return
	}

	policy, err := h.svc.Update(policyNo, req)
	if err != nil {
		switch err.Error() {
		case "policy_not_found":
			respondError(c, http.StatusNotFound, "policy_not_found", "Policy not found")
		default:
			respondError(c, http.StatusBadRequest, "update_failed", err.Error())
		}
		return
	}
	respond(c, http.StatusOK, policy)
}
