package models

import "time"

// Client represents a row in the client table.
type Client struct {
	ID          int        `db:"_id"          json:"id"`
	FamilyCode  string     `db:"familycode"   json:"familyCode"`
	PersCode    string     `db:"perscode"     json:"persCode"`
	Name        string     `db:"name"         json:"name"`
	Mobile      string     `db:"mobileno"     json:"mobile"`
	DOB         *time.Time `db:"dob_r"        json:"dob"`
	Sex         string     `db:"sex"          json:"sex"`
	Address     string     `db:"address"      json:"address"`
	Email       string     `db:"email"        json:"email"`
	Occupation  string     `db:"occupation"   json:"occupation"`
	Age         int        `db:"age"          json:"age"`
	ClientType  string     `db:"client_type"  json:"clientType"`
	Comment     string     `db:"comment"      json:"comment"`
	Source      string     `db:"source"       json:"source"`
	City        string     `db:"city"         json:"city"`
	State       string     `db:"state"        json:"state"`
}

// CreateClientRequest holds data for POST /clients.
type CreateClientRequest struct {
	FamilyCode string `json:"familyCode" binding:"required,max=15"`
	PersCode   string `json:"persCode"   binding:"required,max=10"`
	Name       string `json:"name"       binding:"required,min=2,max=80"`
	DOB        string `json:"dob"        binding:"omitempty"`
	Sex        string `json:"sex"        binding:"omitempty,oneof=M F O"`
	Mobile     string `json:"mobile"     binding:"omitempty,min=10,max=15"`
	Email      string `json:"email"      binding:"omitempty,email"`
	Occupation string `json:"occupation" binding:"omitempty,max=50"`
	ClientType string `json:"clientType" binding:"omitempty,oneof=C P N"`
	Address    string `json:"address"    binding:"omitempty,max=255"`
}

// UpdateClientRequest holds data for PUT /clients/:id.
type UpdateClientRequest struct {
	Name       string `json:"name"       binding:"omitempty,min=2,max=80"`
	DOB        string `json:"dob"        binding:"omitempty"`
	Sex        string `json:"sex"        binding:"omitempty,oneof=M F O"`
	Mobile     string `json:"mobile"     binding:"omitempty,min=10,max=15"`
	Email      string `json:"email"      binding:"omitempty,email"`
	Occupation string `json:"occupation" binding:"omitempty,max=50"`
	ClientType string `json:"clientType" binding:"omitempty,oneof=C P N"`
	Address    string `json:"address"    binding:"omitempty,max=255"`
}

// ClientDetail extends Client with linked data.
type ClientDetail struct {
	Client
	Policies    []PolicyListItem `json:"policies"`
	BankDetails []BankDetail     `json:"bankDetails"`
	Documents   []Document       `json:"documents"`
}

// BankDetail represents a row in the bankdetails table.
// Aadhar, PAN, and CKYC are stored encrypted at rest (AES-256-GCM via domain.Encrypt).
type BankDetail struct {
	ID            int    `db:"_id"           json:"id"`
	ClientID      int    `db:"clientid"      json:"clientId"`
	BankName      string `db:"bankname"      json:"bankName"`
	AccountNumber string `db:"accountnumber" json:"accountNumber"`
	IFSECode      string `db:"ifsecode"      json:"ifseCode"`
	MICRCode      string `db:"micrcode"      json:"micrCode"`
	FamilyCode    string `db:"familycode"    json:"familyCode"`
	PersCode      string `db:"perscode"      json:"persCode"`
	// Aadhar is AES-256-GCM encrypted before storage. Value returned by API is decrypted.
	Aadhar        string `db:"aadhar"        json:"aadhar"`
	// PAN is AES-256-GCM encrypted before storage. Value returned by API is decrypted.
	PAN           string `db:"pan"           json:"pan"`
	// CKYC is AES-256-GCM encrypted before storage.
	CKYC          string `db:"ckyc"          json:"ckyc"`
}

// CreateBankDetailRequest holds data for POST /clients/:id/bank.
type CreateBankDetailRequest struct {
	BankName      string `json:"bankName"      binding:"required,max=200"`
	AccountNumber string `json:"accountNumber" binding:"required,max=50"`
	IFSECode      string `json:"ifseCode"      binding:"omitempty,max=20"`
	MICRCode      string `json:"micrCode"      binding:"omitempty,max=30"`
	Aadhar        string `json:"aadhar"        binding:"omitempty,max=15"`
	PAN           string `json:"pan"           binding:"omitempty,max=10"`
	CKYC          string `json:"ckyc"          binding:"omitempty,max=50"`
}

// Document represents a row in the documents table.
type Document struct {
	ID       int    `db:"_id"       json:"id"`
	ClientID int    `db:"clientid"  json:"clientId"`
	PolicyNo int    `db:"policy_no" json:"policyNo"`
	Title    string `db:"title"     json:"title"`
	Photo    string `db:"photo"     json:"photo"`
	Profile  bool   `db:"profile"   json:"profile"`
}
