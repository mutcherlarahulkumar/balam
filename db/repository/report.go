package repository

import (
	"agent-balam/models"

	"github.com/jmoiron/sqlx"
)

// ReportRepo handles report cache table queries.
type ReportRepo struct {
	db *sqlx.DB
}

// NewReportRepo creates a new ReportRepo.
func NewReportRepo(db *sqlx.DB) *ReportRepo {
	return &ReportRepo{db: db}
}

// Cashflow returns cached cashflow rows for a family.
func (r *ReportRepo) Cashflow(familyCode string) ([]models.ReportCashflow, error) {
	var items []models.ReportCashflow
	err := r.db.Select(&items, `SELECT familycode, perscode, policy_no, TO_CHAR(due_date,'YYYY-MM-DD'), name, age, sb_amount, bonus, total, type
	                              FROM rpt_cashflow WHERE familycode=$1 ORDER BY due_date`, familyCode)
	if items == nil {
		items = []models.ReportCashflow{}
	}
	return items, err
}

// CashInOut returns cash-in vs cash-out per year for all policies.
func (r *ReportRepo) CashInOut() ([]models.ReportCashInOut, error) {
	var items []models.ReportCashInOut
	err := r.db.Select(&items, `SELECT year, cashin_mat, cashin_sb, cashout, netamt FROM rpt_cashinout ORDER BY year`)
	if items == nil {
		items = []models.ReportCashInOut{}
	}
	return items, err
}

// CurrentStatus returns the current status snapshot for a family.
func (r *ReportRepo) CurrentStatus(familyCode string) ([]models.ReportCurrentStatus, error) {
	var items []models.ReportCurrentStatus
	err := r.db.Select(&items, `SELECT familycode, perscode, policy_no, due_status, premium_paid, vested_bonus, surr_value, loan_taken, COALESCE(loan_date,''), COALESCE(loan_amt,0)
	                              FROM rpt_currstatus WHERE familycode=$1`, familyCode)
	if items == nil {
		items = []models.ReportCurrentStatus{}
	}
	return items, err
}

// Calendar returns the 12-month premium calendar for a family.
func (r *ReportRepo) Calendar(familyCode string) ([]models.ReportCalendar, error) {
	var items []models.ReportCalendar
	err := r.db.Select(&items, `SELECT familycode, perscode, policy_no, jan, feb, mar, apr, may, jun, jul, aug, sep, oct, nov, dec
	                              FROM rpt_prcalendar WHERE familycode=$1`, familyCode)
	if items == nil {
		items = []models.ReportCalendar{}
	}
	return items, err
}

// RefreshCashflow rebuilds the rpt_cashflow cache for a family from policies + sb data.
func (r *ReportRepo) RefreshCashflow(familyCode string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	_, err = tx.Exec(`DELETE FROM rpt_cashflow WHERE familycode=$1`, familyCode)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO rpt_cashflow (familycode, perscode, policy_no, due_date, name, age, sb_amount, bonus, total, type, plan, term, ppt)
		SELECT p.familycode, p.perscode, CAST(p.policy_no AS VARCHAR), s.sb_duedate,
		       COALESCE(c.name,''), COALESCE(p.age,0),
		       s.sb_amount, 0, s.sb_amount, 'SB',
		       p.plan, p.term, p.ppt
		FROM sb s
		JOIN policies p ON p.policy_no = s.policy_no
		LEFT JOIN client c ON c.familycode=p.familycode AND c.perscode=p.perscode
		WHERE p.familycode=$1`, familyCode)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO rpt_cashflow (familycode, perscode, policy_no, due_date, name, age, sb_amount, bonus, total, type, plan, term, ppt)
		SELECT p.familycode, p.perscode, CAST(p.policy_no AS VARCHAR), p.mat_date,
		       COALESCE(c.name,''), COALESCE(p.age,0),
		       p.sum_assured, 0, p.sum_assured, 'MAT',
		       p.plan, p.term, p.ppt
		FROM policies p
		LEFT JOIN client c ON c.familycode=p.familycode AND c.perscode=p.perscode
		WHERE p.familycode=$1 AND p.mat_date IS NOT NULL`, familyCode)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// RefreshCalendar rebuilds the premium calendar cache for a family.
func (r *ReportRepo) RefreshCalendar(familyCode string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	_, err = tx.Exec(`DELETE FROM rpt_prcalendar WHERE familycode=$1`, familyCode)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO rpt_prcalendar (familycode, perscode, policy_no, jan, feb, mar, apr, may, jun, jul, aug, sep, oct, nov, dec)
		SELECT p.familycode, p.perscode, CAST(p.policy_no AS VARCHAR),
		       CASE WHEN EXTRACT(MONTH FROM p.next_premium) = 1  OR p.payment_mode IN ('M','H','Q') THEN p.premium ELSE 0 END,
		       CASE WHEN EXTRACT(MONTH FROM p.next_premium) = 2  OR p.payment_mode IN ('M','H','Q') THEN p.premium ELSE 0 END,
		       CASE WHEN EXTRACT(MONTH FROM p.next_premium) = 3  OR p.payment_mode IN ('M','Q') THEN p.premium ELSE 0 END,
		       CASE WHEN EXTRACT(MONTH FROM p.next_premium) = 4  OR p.payment_mode = 'M' THEN p.premium ELSE 0 END,
		       CASE WHEN EXTRACT(MONTH FROM p.next_premium) = 5  OR p.payment_mode = 'M' THEN p.premium ELSE 0 END,
		       CASE WHEN EXTRACT(MONTH FROM p.next_premium) = 6  OR p.payment_mode IN ('M','H') THEN p.premium ELSE 0 END,
		       CASE WHEN EXTRACT(MONTH FROM p.next_premium) = 7  OR p.payment_mode = 'M' THEN p.premium ELSE 0 END,
		       CASE WHEN EXTRACT(MONTH FROM p.next_premium) = 8  OR p.payment_mode = 'M' THEN p.premium ELSE 0 END,
		       CASE WHEN EXTRACT(MONTH FROM p.next_premium) = 9  OR p.payment_mode IN ('M','Q') THEN p.premium ELSE 0 END,
		       CASE WHEN EXTRACT(MONTH FROM p.next_premium) = 10 OR p.payment_mode = 'M' THEN p.premium ELSE 0 END,
		       CASE WHEN EXTRACT(MONTH FROM p.next_premium) = 11 OR p.payment_mode = 'M' THEN p.premium ELSE 0 END,
		       CASE WHEN EXTRACT(MONTH FROM p.next_premium) = 12 OR p.payment_mode IN ('M','H','Q') THEN p.premium ELSE 0 END
		FROM policies p
		WHERE p.familycode=$1 AND p.status='IF'`, familyCode)
	if err != nil {
		return err
	}

	return tx.Commit()
}
