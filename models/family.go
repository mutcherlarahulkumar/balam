package models

import "time"

// Family represents a row in the family table.
type Family struct {
	ID          int        `db:"_id"         json:"id"`
	FamilyCode  string     `db:"familycode"  json:"familyCode"`
	HeadName    string     `db:"headname"    json:"headName"`
	Address     string     `db:"address"     json:"address"`
	Email       string     `db:"email"       json:"email"`
	Mobile      string     `db:"mobileno"    json:"mobile"`
	Pincode     string     `db:"pincode"     json:"pincode"`
	Religion    string     `db:"religion"    json:"religion"`
	Designation string     `db:"designation" json:"designation"`
	LastUpdate  *time.Time `db:"last_update" json:"lastUpdate"`
}

// FamilyListItem is a compact view used in GET /families listing.
type FamilyListItem struct {
	ID          int    `json:"id"`
	FamilyCode  string `json:"familyCode"`
	HeadName    string `json:"headName"`
	Mobile      string `json:"mobile"`
	Address     string `json:"address"`
	Pincode     string `json:"pincode"`
	MemberCount int    `json:"memberCount"`
	PolicyCount int    `json:"policyCount"`
}

// CreateFamilyRequest holds data for POST /families.
type CreateFamilyRequest struct {
	FamilyCode  string `json:"familyCode"  binding:"omitempty,max=15"`
	HeadName    string `json:"headName"    binding:"required,min=2,max=80"`
	Address     string `json:"address"     binding:"omitempty,max=250"`
	Mobile      string `json:"mobile"      binding:"omitempty,min=10,max=15"`
	Email       string `json:"email"       binding:"omitempty,email"`
	Pincode     string `json:"pincode"     binding:"omitempty,max=10"`
	Religion    string `json:"religion"    binding:"omitempty,max=20"`
	Designation string `json:"designation" binding:"omitempty,max=80"`
}

// UpdateFamilyRequest holds data for PUT /families/:familyCode.
type UpdateFamilyRequest struct {
	HeadName    string `json:"headName"    binding:"omitempty,min=2,max=80"`
	Address     string `json:"address"     binding:"omitempty,max=250"`
	Mobile      string `json:"mobile"      binding:"omitempty,min=10,max=15"`
	Email       string `json:"email"       binding:"omitempty,email"`
	Pincode     string `json:"pincode"     binding:"omitempty,max=10"`
	Religion    string `json:"religion"    binding:"omitempty,max=20"`
	Designation string `json:"designation" binding:"omitempty,max=80"`
}

// FamilyDetail extends Family with members and policies for GET /families/:familyCode.
type FamilyDetail struct {
	Family
	Members  []Client   `json:"members"`
	Policies []PolicyListItem `json:"policies"`
}
