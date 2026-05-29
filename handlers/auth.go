package handlers

import (
	"agent-balam/domain"
	"agent-balam/middlewares"
	"agent-balam/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authSvc  *domain.AuthService
	agentSvc *domain.AgentService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authSvc *domain.AuthService, agentSvc *domain.AgentService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, agentSvc: agentSvc}
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidation(c, err)
		return
	}

	resp, err := h.authSvc.Login(req)
	if err != nil {
		switch err.Error() {
		case "account_locked":
			respondError(c, http.StatusTooManyRequests, "account_locked",
				"Too many failed login attempts. Account locked for 15 minutes.")
		case "account_terminated":
			respondError(c, http.StatusForbidden, "account_terminated",
				"This agent account has been terminated.")
		default:
			respondError(c, http.StatusUnauthorized, "invalid_credentials",
				"Email/agent code or password is incorrect")
		}
		return
	}
	respond(c, http.StatusOK, resp)
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidation(c, err)
		return
	}

	id, err := h.authSvc.Register(req)
	if err != nil {
		switch err.Error() {
		case "email_taken":
			respondError(c, http.StatusConflict, "email_taken", "An account with this email already exists")
		case "agent_code_taken":
			respondError(c, http.StatusConflict, "agent_code_taken", "This agent code is already registered")
		default:
			respondError(c, http.StatusInternalServerError, "server_error", "Failed to create account")
		}
		return
	}
	respond(c, http.StatusCreated, gin.H{"id": id, "message": "Account created successfully"})
}

// Refresh handles POST /auth/refresh.
func (h *AuthHandler) Refresh(c *gin.Context) {
	agentID := middlewares.AgentID(c)
	resp, err := h.authSvc.Refresh(agentID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to refresh token")
		return
	}
	respond(c, http.StatusOK, resp)
}

// ChangePassword handles POST /auth/change-password.
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidation(c, err)
		return
	}

	agentID := middlewares.AgentID(c)
	if err := h.authSvc.ChangePassword(agentID, req); err != nil {
		switch err.Error() {
		case "wrong_current_password":
			respondError(c, http.StatusUnauthorized, "wrong_current_password", "Current password is incorrect")
		default:
			respondError(c, http.StatusInternalServerError, "server_error", "Failed to change password")
		}
		return
	}
	respond(c, http.StatusOK, gin.H{"message": "Password changed successfully"})
}
