package repository

import (
	"agent-balam/models"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// FamilyRepo handles family table queries.
type FamilyRepo struct {
	db *sqlx.DB
}

// NewFamilyRepo creates a new FamilyRepo.
func NewFamilyRepo(db *sqlx.DB) *FamilyRepo {
	return &FamilyRepo{db: db}
}

// List returns paginated families with member and policy counts.
func (r *FamilyRepo) List(search string, page, limit int) ([]models.FamilyListItem, int, error) {
	offset := (page - 1) * limit
	args := []interface{}{}
	where := ""
	if search != "" {
		where = "WHERE f.headname ILIKE $1 OR f.mobileno ILIKE $1 OR f.familycode ILIKE $1"
		args = append(args, "%"+search+"%")
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM family f %s`, where)
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limitArg := len(args) + 1
	offsetArg := len(args) + 2
	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		SELECT f.familycode, f.headname, f.mobileno, f.address, f.pincode, f._id,
		       COUNT(DISTINCT c._id)       AS member_count,
		       COUNT(DISTINCT p.policy_no) AS policy_count
		FROM family f
		LEFT JOIN client c ON c.familycode = f.familycode
		LEFT JOIN policies p ON p.familycode = f.familycode
		%s
		GROUP BY f._id, f.familycode, f.headname, f.mobileno, f.address, f.pincode
		ORDER BY f.headname
		LIMIT $%d OFFSET $%d`, where, limitArg, offsetArg)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.FamilyListItem
	for rows.Next() {
		var item models.FamilyListItem
		if err := rows.Scan(&item.FamilyCode, &item.HeadName, &item.Mobile, &item.Address,
			&item.Pincode, &item.ID, &item.MemberCount, &item.PolicyCount); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.FamilyListItem{}
	}
	return items, total, rows.Err()
}

// FindByCode returns a single family by its familycode.
func (r *FamilyRepo) FindByCode(familyCode string) (*models.Family, error) {
	var f models.Family
	if err := r.db.Get(&f, `SELECT _id, familycode, headname, address, email, mobileno, pincode, religion, designation, last_update
	                          FROM family WHERE familycode=$1`, familyCode); err != nil {
		return nil, fmt.Errorf("family not found: %w", err)
	}
	return &f, nil
}

// Create inserts a new family, auto-generating a code if none provided.
func (r *FamilyRepo) Create(req models.CreateFamilyRequest) (*models.Family, error) {
	code := strings.ToUpper(req.FamilyCode)
	if code == "" {
		var seq int
		if err := r.db.QueryRow(`SELECT COALESCE(MAX(_id),0)+1 FROM family`).Scan(&seq); err != nil {
			return nil, err
		}
		code = fmt.Sprintf("F%04d", seq)
	}

	var id int
	err := r.db.QueryRow(`SELECT COALESCE(MAX(_id),0)+1 FROM family`).Scan(&id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	_, err = r.db.Exec(`INSERT INTO family (_id, familycode, headname, address, email, mobileno, pincode, religion, designation, last_update)
	                     VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		id, code, req.HeadName, req.Address, req.Email, req.Mobile,
		req.Pincode, req.Religion, req.Designation, now)
	if err != nil {
		return nil, fmt.Errorf("create family: %w", err)
	}
	return r.FindByCode(code)
}

// Update updates editable fields on a family.
func (r *FamilyRepo) Update(familyCode string, req models.UpdateFamilyRequest) error {
	_, err := r.db.Exec(`UPDATE family SET headname=$1, address=$2, email=$3, mobileno=$4, pincode=$5, religion=$6, designation=$7, last_update=$8
	                      WHERE familycode=$9`,
		req.HeadName, req.Address, req.Email, req.Mobile,
		req.Pincode, req.Religion, req.Designation, time.Now(), familyCode)
	return err
}

// CodeExists returns true if a familycode is taken.
func (r *FamilyRepo) CodeExists(code string) (bool, error) {
	var n int
	err := r.db.QueryRow(`SELECT COUNT(1) FROM family WHERE familycode=$1`, code).Scan(&n)
	return n > 0, err
}
