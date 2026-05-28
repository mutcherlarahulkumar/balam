package models

import "time"

// Lead represents a row in the leaddata table.
type Lead struct {
	ID         int        `db:"_id"        json:"id"`
	LeadID     string     `db:"lead_id"    json:"leadId"`
	SearchTerm string     `db:"search_term" json:"searchTerm"`
	Name       string     `db:"name"       json:"name"`
	Mobile     string     `db:"mobile"     json:"mobile"`
	Address    string     `db:"address"    json:"address"`
	Status     int        `db:"status"     json:"status"`
	DateAdded  *time.Time `db:"date_added" json:"dateAdded"`
}

// CreateLeadRequest holds data for POST /leads.
type CreateLeadRequest struct {
	Name       string `json:"name"       binding:"required,min=2,max=150"`
	Mobile     string `json:"mobile"     binding:"required,min=10,max=15"`
	Address    string `json:"address"    binding:"omitempty,max=500"`
	SearchTerm string `json:"searchTerm" binding:"omitempty,max=100"`
}

// Activity represents a row in the prospect_activity table.
type Activity struct {
	ID           int        `db:"_id"           json:"id"`
	ClientID     int        `db:"clientid"      json:"clientId"`
	ActivityDate *time.Time `db:"activity_date" json:"activityDate"`
	ActivityTime *time.Time `db:"activity_time" json:"activityTime"`
	ActivityType string     `db:"activity_type" json:"activityType"`
	Details      string     `db:"details"       json:"details"`
	ReminderDate *time.Time `db:"reminder_date" json:"reminderDate"`
	ReminderTime *time.Time `db:"reminder_time" json:"reminderTime"`
	Status       string     `db:"status"        json:"status"`
}

// CreateActivityRequest holds data for POST /activities.
type CreateActivityRequest struct {
	ClientID     int    `json:"clientId"      binding:"required"`
	ActivityType string `json:"activityType"  binding:"required,oneof=CALL MEETING DEMO EMAIL PROPOSAL MEDICAL REMINDER"`
	ActivityDate string `json:"activityDate"  binding:"required"`
	ActivityTime string `json:"activityTime"  binding:"omitempty"`
	Details      string `json:"details"       binding:"omitempty,max=1000"`
	ReminderDate string `json:"reminderDate"  binding:"omitempty"`
	ReminderTime string `json:"reminderTime"  binding:"omitempty"`
	Status       string `json:"status"        binding:"omitempty,oneof=PENDING DONE CANCELLED"`
}
