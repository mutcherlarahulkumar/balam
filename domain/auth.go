// Package domain contains all business logic for the Balam LIC Agent CRM.
package domain

import (
	"agent-balam/db/repository"
	"agent-balam/models"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const maxFailedAttempts = 5

// AuthService handles authentication business logic.
type AuthService struct {
	agentRepo *repository.AgentRepo
}

// NewAuthService creates a new AuthService.
func NewAuthService(repo *repository.AgentRepo) *AuthService {
	return &AuthService{agentRepo: repo}
}

// Login validates credentials and returns a signed JWT.
// Returns "account_locked" if the identifier has >= 5 failed attempts in the last 15 minutes.
// Returns "account_terminated" if the agent record has terminated=true.
// TODO(LIC-API — Phase 2): On first login with agentCode, call POST https://licindia.in/api/v2/auth/login
// to get userkey + authtoken. Store in agent.userkey and agent.authtoken. Use machine_id (device
// fingerprint hash) for device binding. Support MPIN for offline re-auth without network.
func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	failedCount, err := s.agentRepo.RecentFailedAttempts(req.Identifier)
	if err == nil && failedCount >= maxFailedAttempts {
		log.Printf("[auth] login blocked: identifier=%q locked after %d failed attempts", req.Identifier, failedCount)
		return nil, errors.New("account_locked")
	}

	agent, err := s.agentRepo.FindByIdentifier(req.Identifier)
	if err != nil {
		log.Printf("[auth] login failed: identifier=%q agent not found: %v", req.Identifier, err)
		_ = s.agentRepo.RecordLoginAttempt(req.Identifier, false)
		return nil, errors.New("invalid_credentials")
	}

	log.Printf("[auth] agent found: id=%d login=%q email=%q passwordLen=%d", agent.ID, agent.Login, agent.Email, len(agent.Password))

	if agent.Terminated != nil && *agent.Terminated {
		log.Printf("[auth] login blocked: identifier=%q account terminated", req.Identifier)
		return nil, errors.New("account_terminated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(agent.Password), []byte(req.Password)); err != nil {
		log.Printf("[auth] login failed: identifier=%q bcrypt mismatch: %v", req.Identifier, err)
		_ = s.agentRepo.RecordLoginAttempt(req.Identifier, false)
		return nil, errors.New("invalid_credentials")
	}

	log.Printf("[auth] login success: identifier=%q agentID=%d", req.Identifier, agent.ID)

	// Clear failed attempts on success
	_ = s.agentRepo.ClearLoginAttempts(req.Identifier)
	_ = s.agentRepo.RecordLoginAttempt(req.Identifier, true)

	token, expiresAt, err := generateToken(agent.ID)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		Agent: models.AgentPublic{
			ID:          agent.ID,
			Name:        agent.Name,
			AgentCode:   agent.Login,
			Branch:      agent.Branch,
			Club:        agent.Club,
			LicenceNo:   agent.LicenceNo,
			AgSince:     agent.AgSince,
			RenewalDate: agent.RenewalDate,
			PAN:         agent.PAN,
			Mobile:      agent.Mobile,
			Email:       agent.Email,
			Photo:       agent.Photo,
			Slogan:      agent.Slogan,
			NewPortal:   agent.NewPortal,
		},
		ExpiresAt: expiresAt,
	}, nil
}

// Register creates a new agent account.
func (s *AuthService) Register(req models.RegisterRequest) (int, error) {
	emailExists, err := s.agentRepo.EmailExists(req.Email)
	if err != nil {
		return 0, err
	}
	if emailExists {
		return 0, fmt.Errorf("email_taken")
	}

	codeExists, err := s.agentRepo.AgentCodeExists(req.AgentCode)
	if err != nil {
		return 0, err
	}
	if codeExists {
		return 0, fmt.Errorf("agent_code_taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	id, err := s.agentRepo.Create(req.Name, req.Email, req.AgentCode, string(hashed), req.Branch, req.Mobile, req.LicenceNo)
	if err != nil {
		return 0, err
	}
	// Clear any prior failed attempts so a freshly registered account is never locked out.
	_ = s.agentRepo.ClearLoginAttempts(req.AgentCode)
	_ = s.agentRepo.ClearLoginAttempts(req.Email)
	return id, nil
}

// Refresh issues a new token for an already-authenticated agent.
func (s *AuthService) Refresh(agentID int) (*models.RefreshResponse, error) {
	token, expiresAt, err := generateToken(agentID)
	if err != nil {
		return nil, err
	}
	return &models.RefreshResponse{Token: token, ExpiresAt: expiresAt}, nil
}

// ChangePassword validates the old password and sets a new one.
func (s *AuthService) ChangePassword(agentID int, req models.ChangePasswordRequest) error {
	agent, err := s.agentRepo.FindByID(agentID)
	if err != nil {
		return errors.New("agent_not_found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(agent.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("wrong_current_password")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.agentRepo.UpdatePassword(agentID, string(hashed))
}

// ValidateToken parses and validates a JWT, returning the agent ID on success.
func ValidateToken(tokenStr string) (int, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid_token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid_claims")
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("invalid_subject")
	}
	return int(sub), nil
}

func generateToken(agentID int) (string, time.Time, error) {
	secret := os.Getenv("JWT_SECRET")
	hours := 24
	if h, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS")); err == nil && h > 0 {
		hours = h
	}
	expiresAt := time.Now().Add(time.Duration(hours) * time.Hour)

	claims := jwt.MapClaims{
		"sub": agentID,
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	return signed, expiresAt, err
}
