// Package repository provides database access methods for all domain entities.
package repository

import (
	"agent-balam/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// AgentRepo handles agent table queries.
type AgentRepo struct {
	db *sqlx.DB
}

// NewAgentRepo creates a new AgentRepo.
func NewAgentRepo(db *sqlx.DB) *AgentRepo {
	return &AgentRepo{db: db}
}

const agentSelectCols = `_id, name, address, mobile, email, login, password, branch, club,
	                 licence_no, ag_since, ren_dt, pan, photo, slogan, newportal, terminated, userkey, authtoken`

// FindByIdentifier looks up an agent by email or agentCode (login field).
func (r *AgentRepo) FindByIdentifier(identifier string) (*models.Agent, error) {
	var agent models.Agent
	query := `SELECT ` + agentSelectCols + `
	          FROM agent
	          WHERE email = $1 OR login = $1
	          ORDER BY _id DESC
	          LIMIT 1`
	if err := r.db.Get(&agent, query, identifier); err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}
	return &agent, nil
}

// FindByID returns an agent by primary key.
func (r *AgentRepo) FindByID(id int) (*models.Agent, error) {
	var agent models.Agent
	query := `SELECT ` + agentSelectCols + ` FROM agent WHERE _id = $1`
	if err := r.db.Get(&agent, query, id); err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}
	return &agent, nil
}

// RecordLoginAttempt inserts an attempt row for rate-limit tracking.
func (r *AgentRepo) RecordLoginAttempt(identifier string, succeeded bool) error {
	_, err := r.db.Exec(`INSERT INTO login_attempts (identifier, attempted_at, succeeded) VALUES ($1, NOW(), $2)`,
		identifier, succeeded)
	return err
}

// RecentFailedAttempts returns the count of failed login attempts in the last 15 minutes.
func (r *AgentRepo) RecentFailedAttempts(identifier string) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM login_attempts
	                       WHERE identifier=$1 AND succeeded=false AND attempted_at > NOW() - INTERVAL '15 minutes'`,
		identifier).Scan(&count)
	return count, err
}

// ClearLoginAttempts removes all attempts for an identifier (called on successful login).
func (r *AgentRepo) ClearLoginAttempts(identifier string) error {
	_, err := r.db.Exec(`DELETE FROM login_attempts WHERE identifier=$1`, identifier)
	return err
}

// Create inserts a new agent and returns the inserted ID.
func (r *AgentRepo) Create(name, email, login, hashedPassword, branch, mobile, licenceNo string) (int, error) {
	var id int
	query := `INSERT INTO agent (name, email, login, password, branch, mobile, licence_no)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING _id`
	err := r.db.QueryRow(query, name, email, login, hashedPassword, branch, mobile, licenceNo).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create agent: %w", err)
	}
	return id, nil
}

// UpdateProfile updates editable profile fields.
func (r *AgentRepo) UpdateProfile(id int, name, mobile, email, photo, slogan, address string) error {
	query := `UPDATE agent SET name=$1, mobile=$2, email=$3, photo=$4, slogan=$5, address=$6
	          WHERE _id=$7`
	_, err := r.db.Exec(query, name, mobile, email, photo, slogan, address, id)
	return err
}

// UpdatePassword updates the hashed password for an agent.
func (r *AgentRepo) UpdatePassword(id int, hashedPassword string) error {
	_, err := r.db.Exec(`UPDATE agent SET password=$1 WHERE _id=$2`, hashedPassword, id)
	return err
}

// EmailExists returns true if an email is already taken.
func (r *AgentRepo) EmailExists(email string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(1) FROM agent WHERE email=$1`, email).Scan(&count)
	return count > 0, err
}

// AgentCodeExists returns true if the login/agentCode is already taken.
func (r *AgentRepo) AgentCodeExists(agentCode string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(1) FROM agent WHERE login=$1`, agentCode).Scan(&count)
	return count > 0, err
}
