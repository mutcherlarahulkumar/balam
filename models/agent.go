// Package models contains all domain structs, request/response types, and DB row types.
package models

import "time"

// Agent represents a row in the agent table.
type Agent struct {
	ID          int        `db:"_id"        json:"id"`
	Name        string     `db:"name"        json:"name"`
	Address     string     `db:"address"     json:"address"`
	Mobile      string     `db:"mobile"      json:"mobile"`
	Email       string     `db:"email"       json:"email"`
	Login       string     `db:"login"       json:"login"`
	Password    string     `db:"password"    json:"-"`
	Branch      string     `db:"branch"      json:"branch"`
	Club        string     `db:"club"        json:"club"`
	LicenceNo   string     `db:"licence_no"  json:"licenceNo"`
	AgSince     *time.Time `db:"ag_since"    json:"agSince"`
	RenewalDate *time.Time `db:"ren_dt"      json:"renewalDate"`
	PAN         string     `db:"pan"         json:"pan"`
	Photo       string     `db:"photo"       json:"photo"`
	Slogan      string     `db:"slogan"      json:"slogan"`
	NewPortal   bool       `db:"newportal"   json:"newPortal"`
	UserKey     string     `db:"userkey"     json:"-"`
	AuthToken   string     `db:"authtoken"   json:"-"`
}

// AgentPublic is the safe public representation returned in API responses.
type AgentPublic struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	AgentCode   string     `json:"agentCode"`
	Branch      string     `json:"branch"`
	Club        string     `json:"club"`
	LicenceNo   string     `json:"licenceNo"`
	AgSince     *time.Time `json:"agSince"`
	RenewalDate *time.Time `json:"renewalDate"`
	PAN         string     `json:"pan"`
	Mobile      string     `json:"mobile"`
	Email       string     `json:"email"`
	Photo       string     `json:"photo"`
	Slogan      string     `json:"slogan"`
	NewPortal   bool       `json:"newPortal"`
}

// LoginRequest holds credentials for POST /auth/login.
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password"   binding:"required,min=6"`
}

// LoginResponse is the body returned on successful login.
type LoginResponse struct {
	Token     string      `json:"token"`
	Agent     AgentPublic `json:"agent"`
	ExpiresAt time.Time   `json:"expiresAt"`
}

// RegisterRequest holds data for POST /auth/register.
type RegisterRequest struct {
	Name      string `json:"name"      binding:"required,min=2,max=80"`
	Email     string `json:"email"     binding:"required,email"`
	AgentCode string `json:"agentCode" binding:"required"`
	Password  string `json:"password"  binding:"required,min=8"`
	Branch    string `json:"branch"    binding:"required"`
	Mobile    string `json:"mobile"    binding:"required,min=10,max=15"`
	LicenceNo string `json:"licenceNo" binding:"required"`
}

// RefreshResponse is returned by POST /auth/refresh.
type RefreshResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// ChangePasswordRequest holds data for POST /auth/change-password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword"     binding:"required,min=8"`
}

// UpdateProfileRequest holds editable profile fields for PUT /agent/profile.
type UpdateProfileRequest struct {
	Name    string `json:"name"    binding:"omitempty,min=2,max=80"`
	Mobile  string `json:"mobile"  binding:"omitempty,min=10,max=15"`
	Email   string `json:"email"   binding:"omitempty,email"`
	Photo   string `json:"photo"   binding:"omitempty"`
	Slogan  string `json:"slogan"  binding:"omitempty,max=200"`
	Address string `json:"address" binding:"omitempty,max=200"`
}
