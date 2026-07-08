package handlers

import (
	"agent-balam/domain"
	"agent-balam/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoanHandler handles loan endpoints.
type LoanHandler struct {
	svc *domain.LoanService
}

// NewLoanHandler creates a new LoanHandler.
func NewLoanHandler(svc *domain.LoanService) *LoanHandler {
	return &LoanHandler{svc: svc}
}

// List handles GET /loans.
func (h *LoanHandler) List(c *gin.Context) {
	items, err := h.svc.List()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to fetch loans")
		return
	}
	respond(c, http.StatusOK, gin.H{"data": items})
}

// Create handles POST /loans.
func (h *LoanHandler) Create(c *gin.Context) {
	var req models.CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidation(c, err)
		return
	}

	if err := h.svc.Create(req); err != nil {
		switch err.Error() {
		case "policy_not_found":
			respondError(c, http.StatusNotFound, "policy_not_found", "Policy not found")
		case "interest_due_before_loan_date":
			respondError(c, http.StatusUnprocessableEntity, "interest_due_before_loan_date", "Interest due date cannot be earlier than the loan date")
		default:
			respondError(c, http.StatusBadRequest, "create_failed", err.Error())
		}
		return
	}
	respond(c, http.StatusCreated, gin.H{"message": "Loan record added"})
}
