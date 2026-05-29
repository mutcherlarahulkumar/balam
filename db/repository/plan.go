package repository

import (
	"agent-balam/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// PlanRepo handles plans table queries.
type PlanRepo struct {
	db *sqlx.DB
}

// NewPlanRepo creates a new PlanRepo.
func NewPlanRepo(db *sqlx.DB) *PlanRepo {
	return &PlanRepo{db: db}
}

const planCols = `
	COALESCE(plan_no,'')     AS plan_no,
	COALESCE(plan_name,'')   AS plan_name,
	COALESCE(plan_type,'')   AS plan_type,
	COALESCE(term_ppt,'')    AS term_ppt,
	COALESCE(sb_year,'')     AS sb_year,
	COALESCE(sb_benefit,'')  AS sb_benefit,
	COALESCE(stax,'')        AS stax,
	COALESCE(lapsdays,0)     AS lapsdays`

// List returns all plans, optionally filtered by planType or search term.
func (r *PlanRepo) List(planType, search string) ([]models.Plan, error) {
	query := `SELECT ` + planCols + ` FROM plans WHERE 1=1`
	args := []interface{}{}
	n := 1
	if planType != "" {
		query += fmt.Sprintf(` AND plan_type ILIKE $%d`, n)
		args = append(args, "%"+planType+"%")
		n++
	}
	if search != "" {
		query += fmt.Sprintf(` AND (plan_name ILIKE $%d OR plan_no ILIKE $%d)`, n, n)
		args = append(args, "%"+search+"%")
	}
	query += ` ORDER BY plan_no`

	var plans []models.Plan
	err := r.db.Select(&plans, query, args...)
	if plans == nil {
		plans = []models.Plan{}
	}
	return plans, err
}

// FindByNo returns a single plan by plan number.
func (r *PlanRepo) FindByNo(planNo string) (*models.Plan, error) {
	var p models.Plan
	if err := r.db.Get(&p, `SELECT `+planCols+` FROM plans WHERE plan_no=$1`, planNo); err != nil {
		return nil, err
	}
	return &p, nil
}
