package repository

import (
	"agent-balam/models"
	"time"

	"github.com/jmoiron/sqlx"
)

// LeadRepo handles leaddata and prospect_activity table queries.
type LeadRepo struct {
	db *sqlx.DB
}

// NewLeadRepo creates a new LeadRepo.
func NewLeadRepo(db *sqlx.DB) *LeadRepo {
	return &LeadRepo{db: db}
}

// ListLeads returns all leads.
func (r *LeadRepo) ListLeads() ([]models.Lead, error) {
	var items []models.Lead
	err := r.db.Select(&items, `SELECT _id, COALESCE(lead_id,'') AS lead_id, COALESCE(search_term,'') AS search_term, name, mobile, COALESCE(address,'') AS address, COALESCE(status,0) AS status, date_added FROM leaddata ORDER BY date_added DESC`)
	if items == nil {
		items = []models.Lead{}
	}
	return items, err
}

// CreateLead inserts a new lead.
func (r *LeadRepo) CreateLead(req models.CreateLeadRequest) (*models.Lead, error) {
	now := time.Now()
	var id int
	err := r.db.QueryRow(`INSERT INTO leaddata (name, mobile, address, search_term, status, date_added) VALUES ($1,$2,$3,$4,0,$5) RETURNING _id`,
		req.Name, req.Mobile, req.Address, req.SearchTerm, now).Scan(&id)
	if err != nil {
		return nil, err
	}
	var item models.Lead
	err = r.db.QueryRow(`SELECT _id, COALESCE(lead_id,''), COALESCE(search_term,''), name, mobile, COALESCE(address,''), COALESCE(status,0), date_added FROM leaddata WHERE _id=$1`, id).
		Scan(&item.ID, &item.LeadID, &item.SearchTerm, &item.Name, &item.Mobile, &item.Address, &item.Status, &item.DateAdded)
	return &item, err
}

// ListActivities returns all activities, optionally for a client.
func (r *LeadRepo) ListActivities(clientID int) ([]models.Activity, error) {
	query := `SELECT _id, clientid, activity_date, activity_time, activity_type, COALESCE(details,''), reminder_date, reminder_time, COALESCE(status,'')
	          FROM prospect_activity`
	args := []interface{}{}
	if clientID > 0 {
		query += ` WHERE clientid=$1`
		args = append(args, clientID)
	}
	query += ` ORDER BY activity_date DESC`
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Activity
	for rows.Next() {
		var a models.Activity
		if err := rows.Scan(&a.ID, &a.ClientID, &a.ActivityDate, &a.ActivityTime,
			&a.ActivityType, &a.Details, &a.ReminderDate, &a.ReminderTime, &a.Status); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	if items == nil {
		items = []models.Activity{}
	}
	return items, rows.Err()
}

// TodayActivities returns activities where reminder_date is today and status is PENDING.
func (r *LeadRepo) TodayActivities() ([]models.Activity, error) {
	rows, err := r.db.Query(`SELECT _id, clientid, activity_date, activity_time, activity_type, COALESCE(details,''), reminder_date, reminder_time, COALESCE(status,'')
	                          FROM prospect_activity
	                          WHERE DATE(reminder_date) = CURRENT_DATE AND status='PENDING'
	                          ORDER BY reminder_time`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Activity
	for rows.Next() {
		var a models.Activity
		if err := rows.Scan(&a.ID, &a.ClientID, &a.ActivityDate, &a.ActivityTime,
			&a.ActivityType, &a.Details, &a.ReminderDate, &a.ReminderTime, &a.Status); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	if items == nil {
		items = []models.Activity{}
	}
	return items, rows.Err()
}

// CreateActivity inserts a new prospect activity.
func (r *LeadRepo) CreateActivity(req models.CreateActivityRequest, actDate, actTime, remDate, remTime *time.Time) (*models.Activity, error) {
	status := req.Status
	if status == "" {
		status = "PENDING"
	}
	var id int
	err := r.db.QueryRow(`INSERT INTO prospect_activity (clientid, activity_date, activity_time, activity_type, details, reminder_date, reminder_time, status)
	                       VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING _id`,
		req.ClientID, actDate, actTime, req.ActivityType, req.Details, remDate, remTime, status).Scan(&id)
	if err != nil {
		return nil, err
	}
	var a models.Activity
	err = r.db.QueryRow(`SELECT _id, clientid, activity_date, activity_time, activity_type, COALESCE(details,''), reminder_date, reminder_time, COALESCE(status,'')
	                      FROM prospect_activity WHERE _id=$1`, id).
		Scan(&a.ID, &a.ClientID, &a.ActivityDate, &a.ActivityTime, &a.ActivityType, &a.Details, &a.ReminderDate, &a.ReminderTime, &a.Status)
	return &a, err
}
