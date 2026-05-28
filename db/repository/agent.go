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

// FindByIdentifier looks up an agent by email or agentCode (login field).
func (r *AgentRepo) FindByIdentifier(identifier string) (*models.Agent, error) {
	var agent models.Agent
	query := `SELECT _id, name, address, mobile, email, login, password, branch, club,
	                 licence_no, ag_since, ren_dt, pan, photo, slogan, newportal, userkey, authtoken
	          FROM agent
	          WHERE email = $1 OR login = $1
	          LIMIT 1`
	if err := r.db.Get(&agent, query, identifier); err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}
	return &agent, nil
}

// FindByID returns an agent by primary key.
func (r *AgentRepo) FindByID(id int) (*models.Agent, error) {
	var agent models.Agent
	query := `SELECT _id, name, address, mobile, email, login, password, branch, club,
	                 licence_no, ag_since, ren_dt, pan, photo, slogan, newportal, userkey, authtoken
	          FROM agent WHERE _id = $1`
	if err := r.db.Get(&agent, query, id); err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}
	return &agent, nil
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
