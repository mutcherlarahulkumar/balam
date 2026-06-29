package repository

import (
	"agent-balam/models"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// PolicyRepo handles policies table queries.
type PolicyRepo struct {
	db *sqlx.DB
}

// NewPolicyRepo creates a new PolicyRepo.
func NewPolicyRepo(db *sqlx.DB) *PolicyRepo {
	return &PolicyRepo{db: db}
}

// ListFilter holds filter parameters for policy listing.
type ListFilter struct {
	Status       string
	FamilyCode   string
	DueThisMonth bool
	LapsingIn    int
	Page         int
	Limit        int
}

// List returns paginated policies with computed fields joined from client.
func (r *PolicyRepo) List(f ListFilter) ([]models.PolicyListItem, int, error) {
	conditions := []string{"1=1"}
	args := []interface{}{}
	n := 1

	if f.Status != "" {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", n))
		args = append(args, f.Status)
		n++
	}
	if f.FamilyCode != "" {
		conditions = append(conditions, fmt.Sprintf("p.familycode = $%d", n))
		args = append(args, f.FamilyCode)
		n++
	}
	if f.DueThisMonth {
		conditions = append(conditions, fmt.Sprintf(
			"DATE_TRUNC('month', p.next_premium) = DATE_TRUNC('month', CURRENT_DATE)"))
	}
	if f.LapsingIn > 0 {
		conditions = append(conditions, fmt.Sprintf(
			"p.next_premium + (pl.lapsdays * INTERVAL '1 day') - CURRENT_DATE <= $%d * INTERVAL '1 day'", n))
		args = append(args, f.LapsingIn)
		n++
	}

	where := "WHERE " + joinConditions(conditions)
	var total int
	if err := r.db.QueryRow(fmt.Sprintf(
		`SELECT COUNT(*) FROM policies p LEFT JOIN plans pl ON pl.plan_no = CAST(p.plan AS TEXT) %s`, where),
		args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Limit
	args = append(args, f.Limit, offset)
	query := fmt.Sprintf(`
		SELECT p._id, p.policy_no, COALESCE(p.familycode,'') AS familycode,
		       COALESCE(c.name, '') AS client_name,
		       COALESCE(CAST(p.plan AS TEXT),'') AS plan_no,
		       COALESCE(NULLIF(p.plan_name,''), pl.plan_name, '') AS plan_name,
		       COALESCE(p.term,0) AS term, COALESCE(p.ppt,0) AS ppt, COALESCE(p.premium,0) AS premium, COALESCE(p.sum_assured,0) AS sum_assured,
		       COALESCE(p.payment_mode,'') AS payment_mode,
		       (p.issue_date + (p.ppt || ' years')::interval)::date AS premium_end_date,
		       p.next_premium, p.mat_date,
		       COALESCE(p.status,'') AS status,
		       COALESCE(p.fupstatus,'') AS fupstatus,
		       GREATEST(0, EXTRACT(DAY FROM (p.next_premium + COALESCE(pl.lapsdays,180) * INTERVAL '1 day' - CURRENT_TIMESTAMP))::int) AS days_until_lapse
		FROM policies p
		LEFT JOIN client c ON c.familycode = p.familycode AND c.perscode = p.perscode
		LEFT JOIN plans pl ON pl.plan_no = CAST(p.plan AS TEXT)
		%s
		ORDER BY p.next_premium ASC NULLS LAST
		LIMIT $%d OFFSET $%d`, where, n, n+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.PolicyListItem
	for rows.Next() {
		var item models.PolicyListItem
		if err := rows.Scan(&item.ID, &item.PolicyNo, &item.FamilyCode, &item.ClientName,
			&item.PlanNo, &item.PlanName, &item.Term, &item.PPT, &item.Premium,
			&item.SumAssured, &item.PaymentMode, &item.PremiumEndDate, &item.NextPremium, &item.MatDate,
			&item.Status, &item.FUPStatus, &item.DaysUntilLapse); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.PolicyListItem{}
	}
	return items, total, rows.Err()
}

// SearchByNo returns a light view of policies whose policy number contains q (max 10).
func (r *PolicyRepo) SearchByNo(q string) ([]models.PolicyListItem, error) {
	rows, err := r.db.Query(`
		SELECT p._id, p.policy_no, COALESCE(p.familycode,'') AS familycode,
		       COALESCE(c.name,'') AS client_name,
		       COALESCE(CAST(p.plan AS TEXT),'') AS plan_no,
		       COALESCE(NULLIF(p.plan_name,''), pl.plan_name, '') AS plan_name,
		       COALESCE(p.term,0) AS term, COALESCE(p.ppt,0) AS ppt, COALESCE(p.premium,0) AS premium, COALESCE(p.sum_assured,0) AS sum_assured,
		       COALESCE(p.payment_mode,'') AS payment_mode,
		       (p.issue_date + (p.ppt || ' years')::interval)::date AS premium_end_date,
		       p.next_premium, p.mat_date,
		       COALESCE(p.status,'') AS status,
		       COALESCE(p.fupstatus,'') AS fupstatus,
		       GREATEST(0, EXTRACT(DAY FROM (p.next_premium + COALESCE(pl.lapsdays,180) * INTERVAL '1 day' - CURRENT_TIMESTAMP))::int) AS days_until_lapse
		FROM policies p
		LEFT JOIN client c ON c.familycode = p.familycode AND c.perscode = p.perscode
		LEFT JOIN plans pl ON pl.plan_no = CAST(p.plan AS TEXT)
		WHERE CAST(p.policy_no AS TEXT) ILIKE $1
		ORDER BY p.next_premium ASC NULLS LAST
		LIMIT 10`, "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.PolicyListItem
	for rows.Next() {
		var item models.PolicyListItem
		if err := rows.Scan(&item.ID, &item.PolicyNo, &item.FamilyCode, &item.ClientName,
			&item.PlanNo, &item.PlanName, &item.Term, &item.PPT, &item.Premium,
			&item.SumAssured, &item.PaymentMode, &item.PremiumEndDate, &item.NextPremium, &item.MatDate,
			&item.Status, &item.FUPStatus, &item.DaysUntilLapse); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.PolicyListItem{}
	}
	return items, rows.Err()
}

// FindByNo returns a full policy row.
func (r *PolicyRepo) FindByNo(policyNo int) (*models.Policy, error) {
	var p models.Policy
	query := `SELECT p._id, p.policy_no,
	                 COALESCE(p.familycode,'')   AS familycode,
	                 COALESCE(p.perscode,'')     AS perscode,
	                 p.issue_date, p.mat_date,
	                 COALESCE(p.payment_mode,'')   AS payment_mode,
	                 COALESCE(p.premium,0)       AS premium,
	                 COALESCE(p.sum_assured,0)   AS sum_assured,
	                 COALESCE(p.plan,0)          AS plan,
	                 COALESCE(NULLIF(p.plan_name,''), pl.plan_name, '') AS plan_name,
	                 COALESCE(p.term,0)          AS term,
	                 COALESCE(p.ppt,0)           AS ppt,
	                 p.next_premium,
	                 COALESCE(p.mat_amount,0)    AS mat_amount,
	                 COALESCE(p.nominee,'')      AS nominee,
	                 COALESCE(p.relation,'')     AS relation,
	                 COALESCE(p.agcode,'')       AS agcode,
	                 COALESCE(p.stat_cd,'')      AS stat_cd,
	                 COALESCE(p.status,'')       AS status,
	                 COALESCE(p.fupstatus,'')    AS fupstatus,
	                 COALESCE(p.age,0)           AS age,
	                 p.lastpaid,
	                 COALESCE(p.neft,'')         AS neft,
	                 p.lastfup,
	                 COALESCE(p.branch,'')       AS branch,
	                 COALESCE(p.dab,0)           AS dab,
	                 COALESCE(p.termrider,0)     AS termrider,
	                 (p.issue_date + (COALESCE(p.ppt,0) || ' years')::interval)::date AS premium_end_date
	          FROM policies p
	          LEFT JOIN plans pl ON pl.plan_no = CAST(p.plan AS TEXT)
	          WHERE p.policy_no=$1`
	if err := r.db.Get(&p, query, policyNo); err != nil {
		return nil, fmt.Errorf("policy not found: %w", err)
	}
	return &p, nil
}

// Create inserts a new policy. age is computed from client DOB at issue date.
func (r *PolicyRepo) Create(req models.CreatePolicyRequest, planName string, issueDate, matDate, nextPremium time.Time, age int) error {
	var id int
	if err := r.db.QueryRow(`SELECT COALESCE(MAX(_id),0)+1 FROM policies`).Scan(&id); err != nil {
		return err
	}
	_, err := r.db.Exec(`
		INSERT INTO policies (_id, policy_no, familycode, perscode, plan, plan_name, issue_date, mat_date,
		                      term, ppt, sum_assured, premium, payment_mode, next_premium, nominee, relation,
		                      branch, neft, dab, termrider, status, fupstatus, age)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,'IF','DUE',$21)`,
		id, req.PolicyNo, req.FamilyCode, req.PersCode, req.PlanNo, planName,
		issueDate, matDate, req.Term, req.PPT, req.SumAssured, req.Premium,
		req.PaymentMode, nextPremium, req.Nominee, req.Relation,
		req.Branch, req.NEFT, req.DAB, req.TermRider, age)
	return err
}

// Update applies selective updates to a policy.
func (r *PolicyRepo) Update(policyNo int, req models.UpdatePolicyRequest, nextPremium *time.Time) error {
	sets := []string{}
	args := []interface{}{}
	n := 1

	if req.Status != "" {
		sets = append(sets, fmt.Sprintf("status=$%d", n))
		args = append(args, req.Status)
		n++
	}
	if req.Nominee != "" {
		sets = append(sets, fmt.Sprintf("nominee=$%d", n))
		args = append(args, req.Nominee)
		n++
	}
	if req.Relation != "" {
		sets = append(sets, fmt.Sprintf("relation=$%d", n))
		args = append(args, req.Relation)
		n++
	}
	if req.NEFT != "" {
		sets = append(sets, fmt.Sprintf("neft=$%d", n))
		args = append(args, req.NEFT)
		n++
	}
	if nextPremium != nil {
		sets = append(sets, fmt.Sprintf("next_premium=$%d", n))
		args = append(args, *nextPremium)
		n++
	}
	if req.FUPStatus != "" {
		sets = append(sets, fmt.Sprintf("fupstatus=$%d", n))
		args = append(args, req.FUPStatus)
		n++
	}

	if len(sets) == 0 {
		return nil
	}

	args = append(args, policyNo)
	_, err := r.db.Exec(
		fmt.Sprintf("UPDATE policies SET %s WHERE policy_no=$%d", joinSets(sets), n),
		args...)
	return err
}

// PoliciesByFamilyCode returns a light view of policies for a family.
func (r *PolicyRepo) PoliciesByFamilyCode(familyCode string) ([]models.PolicyListItem, error) {
	rows, err := r.db.Query(`
		SELECT p._id, p.policy_no, COALESCE(p.familycode,'') AS familycode,
		       COALESCE(c.name,'') AS client_name,
		       COALESCE(CAST(p.plan AS TEXT),''), COALESCE(NULLIF(p.plan_name,''), pl.plan_name, ''),
		       COALESCE(p.term,0) AS term, COALESCE(p.ppt,0) AS ppt, COALESCE(p.premium,0) AS premium, COALESCE(p.sum_assured,0) AS sum_assured,
		       COALESCE(p.payment_mode,'') AS payment_mode,
		       (p.issue_date + (p.ppt || ' years')::interval)::date,
		       p.next_premium, p.mat_date,
		       COALESCE(p.status,''), COALESCE(p.fupstatus,''),
		       GREATEST(0, EXTRACT(DAY FROM (p.next_premium + COALESCE(pl.lapsdays,180) * INTERVAL '1 day' - CURRENT_TIMESTAMP))::int)
		FROM policies p
		LEFT JOIN client c ON c.familycode=p.familycode AND c.perscode=p.perscode
		LEFT JOIN plans pl ON pl.plan_no=CAST(p.plan AS TEXT)
		WHERE p.familycode=$1 ORDER BY p.next_premium`, familyCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.PolicyListItem
	for rows.Next() {
		var item models.PolicyListItem
		if err := rows.Scan(&item.ID, &item.PolicyNo, &item.FamilyCode, &item.ClientName,
			&item.PlanNo, &item.PlanName, &item.Term, &item.PPT, &item.Premium,
			&item.SumAssured, &item.PaymentMode, &item.PremiumEndDate, &item.NextPremium, &item.MatDate,
			&item.Status, &item.FUPStatus, &item.DaysUntilLapse); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.PolicyListItem{}
	}
	return items, rows.Err()
}

// Exists returns true if a policy number already exists.
func (r *PolicyRepo) Exists(policyNo int) (bool, error) {
	var n int
	err := r.db.QueryRow(`SELECT COUNT(1) FROM policies WHERE policy_no=$1`, policyNo).Scan(&n)
	return n > 0, err
}

func joinSets(sets []string) string {
	result := ""
	for i, s := range sets {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}
