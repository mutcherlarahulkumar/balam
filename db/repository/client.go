package repository

import (
	"agent-balam/models"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// ClientRepo handles client table queries.
type ClientRepo struct {
	db *sqlx.DB
}

// NewClientRepo creates a new ClientRepo.
func NewClientRepo(db *sqlx.DB) *ClientRepo {
	return &ClientRepo{db: db}
}

const clientCols = `_id,
	COALESCE(familycode,'')  AS familycode,
	COALESCE(perscode,'')    AS perscode,
	COALESCE(name,'')        AS name,
	COALESCE(mobileno,'')    AS mobileno,
	dob_r,
	COALESCE(sex,'')         AS sex,
	COALESCE(address,'')     AS address,
	COALESCE(email,'')       AS email,
	COALESCE(occupation,'')  AS occupation,
	COALESCE(age,0)          AS age,
	COALESCE(client_type,'') AS client_type,
	COALESCE(comment,'')     AS comment,
	COALESCE(source,'')      AS source,
	COALESCE(city,'')        AS city,
	COALESCE(state,'')       AS state`

// List returns paginated clients filtered by optional familyCode and clientType.
func (r *ClientRepo) List(search, familyCode, clientType string, page, limit int) ([]models.Client, int, error) {
	offset := (page - 1) * limit
	conditions := []string{"1=1"}
	args := []interface{}{}
	n := 1

	if search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR mobileno ILIKE $%d)", n, n))
		args = append(args, "%"+search+"%")
		n++
	}
	if familyCode != "" {
		conditions = append(conditions, fmt.Sprintf("familycode = $%d", n))
		args = append(args, familyCode)
		n++
	}
	if clientType != "" {
		conditions = append(conditions, fmt.Sprintf("client_type = $%d", n))
		args = append(args, clientType)
		n++
	}

	where := "WHERE " + joinConditions(conditions)
	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM client %s`, where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	rows, err := r.db.Queryx(fmt.Sprintf(
		`SELECT %s FROM client %s ORDER BY name LIMIT $%d OFFSET $%d`,
		clientCols, where, n, n+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		if err := rows.StructScan(&c); err != nil {
			return nil, 0, err
		}
		clients = append(clients, c)
	}
	if clients == nil {
		clients = []models.Client{}
	}
	return clients, total, rows.Err()
}

// FindByID returns a client by primary key.
func (r *ClientRepo) FindByID(id int) (*models.Client, error) {
	var c models.Client
	if err := r.db.Get(&c, fmt.Sprintf(`SELECT %s FROM client WHERE _id=$1`, clientCols), id); err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}
	return &c, nil
}

// FindByFamilyAndPers returns a client by familyCode + persCode (used to compute age on policy creation).
func (r *ClientRepo) FindByFamilyAndPers(familyCode, persCode string) (*models.Client, error) {
	var c models.Client
	if err := r.db.Get(&c, fmt.Sprintf(`SELECT %s FROM client WHERE familycode=$1 AND perscode=$2 LIMIT 1`, clientCols), familyCode, persCode); err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}
	return &c, nil
}

// Search searches by name, mobile, or policy_no.
func (r *ClientRepo) Search(q string) ([]models.Client, error) {
	rows, err := r.db.Queryx(fmt.Sprintf(`
		SELECT %s FROM client c
		WHERE c.name ILIKE $1 OR c.mobileno ILIKE $1
		   OR EXISTS (
		       SELECT 1 FROM policies p
		       WHERE p.familycode = c.familycode AND p.perscode = c.perscode
		         AND CAST(p.policy_no AS TEXT) ILIKE $1
		   )
		LIMIT 50`, clientCols), "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clients []models.Client
	for rows.Next() {
		var c models.Client
		if err := rows.StructScan(&c); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	if clients == nil {
		clients = []models.Client{}
	}
	return clients, rows.Err()
}

// Create inserts a new client.
func (r *ClientRepo) Create(req models.CreateClientRequest, dob *time.Time, age int) (*models.Client, error) {
	var id int
	if err := r.db.QueryRow(`SELECT COALESCE(MAX(_id),0)+1 FROM client`).Scan(&id); err != nil {
		return nil, err
	}
	_, err := r.db.Exec(`INSERT INTO client (_id, familycode, perscode, name, mobileno, dob_r, sex, address, email, occupation, age, client_type)
	                      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		id, req.FamilyCode, req.PersCode, req.Name, req.Mobile, dob, req.Sex,
		req.Address, req.Email, req.Occupation, age, req.ClientType)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	return r.FindByID(id)
}

// Update updates client fields.
func (r *ClientRepo) Update(id int, req models.UpdateClientRequest, dob *time.Time, age int) error {
	_, err := r.db.Exec(`UPDATE client SET name=$1, mobileno=$2, email=$3, sex=$4, address=$5, occupation=$6, client_type=$7, dob_r=$8, age=$9
	                      WHERE _id=$10`,
		req.Name, req.Mobile, req.Email, req.Sex, req.Address, req.Occupation, req.ClientType, dob, age, id)
	return err
}

// BankDetails returns all bank records for a client, decrypting sensitive fields on read.
func (r *ClientRepo) BankDetails(clientID int) ([]models.BankDetail, error) {
	rows, err := r.db.Query(`SELECT _id, clientid, bankname, accountnumber, ifsecode, micrcode, familycode, perscode,
	                                 COALESCE(aadhar,''), COALESCE(pan,''), COALESCE(ckyc,'')
	                          FROM bankdetails WHERE clientid=$1`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.BankDetail
	for rows.Next() {
		var b models.BankDetail
		var encAadhar, encPAN, encCKYC string
		if err := rows.Scan(&b.ID, &b.ClientID, &b.BankName, &b.AccountNumber, &b.IFSECode, &b.MICRCode,
			&b.FamilyCode, &b.PersCode, &encAadhar, &encPAN, &encCKYC); err != nil {
			return nil, err
		}
		b.Aadhar, _ = decryptField(encAadhar)
		b.PAN, _ = decryptField(encPAN)
		b.CKYC, _ = decryptField(encCKYC)
		items = append(items, b)
	}
	if items == nil {
		items = []models.BankDetail{}
	}
	return items, rows.Err()
}

// Documents returns all documents for a client.
func (r *ClientRepo) Documents(clientID int) ([]models.Document, error) {
	var items []models.Document
	err := r.db.Select(&items, `SELECT _id, clientid, policy_no, title, photo, profile FROM documents WHERE clientid=$1`, clientID)
	if items == nil {
		items = []models.Document{}
	}
	return items, err
}

// CreateBankDetail inserts a bank record with sensitive fields encrypted at rest.
func (r *ClientRepo) CreateBankDetail(clientID int, req models.CreateBankDetailRequest, familyCode, persCode string) (*models.BankDetail, error) {
	encAadhar, err := encryptField(req.Aadhar)
	if err != nil {
		return nil, err
	}
	encPAN, err := encryptField(req.PAN)
	if err != nil {
		return nil, err
	}
	encCKYC, err := encryptField(req.CKYC)
	if err != nil {
		return nil, err
	}
	var id int
	err = r.db.QueryRow(`INSERT INTO bankdetails (clientid, bankname, accountnumber, ifsecode, micrcode, familycode, perscode, aadhar, pan, ckyc)
	                      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING _id`,
		clientID, req.BankName, req.AccountNumber, req.IFSECode, req.MICRCode,
		familyCode, persCode, encAadhar, encPAN, encCKYC).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("create bank detail: %w", err)
	}
	details, err := r.BankDetails(clientID)
	if err != nil {
		return nil, err
	}
	for _, d := range details {
		if d.ID == id {
			return &d, nil
		}
	}
	return nil, fmt.Errorf("bank detail not found after insert")
}

// encryptField encrypts a sensitive field value using the domain crypto package.
// Import-cycle avoided by duplicating the env-key-based AES call here via the package.
func encryptField(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	// Delegate to domain crypto — but since we're in repository we cannot import domain.
	// Instead we call the shared package-level helper registered at startup.
	if globalEncryptor == nil {
		return value, nil // fallback: store plain if no encryptor configured (test mode)
	}
	return globalEncryptor(value)
}

func decryptField(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	if globalDecryptor == nil {
		return value, nil
	}
	return globalDecryptor(value)
}

// globalEncryptor and globalDecryptor are set by InitCrypto from main at startup.
var (
	globalEncryptor func(string) (string, error)
	globalDecryptor func(string) (string, error)
)

// InitCrypto wires the encryption/decryption functions from the domain layer into the repository.
// Called once from main before serving requests.
func InitCrypto(enc, dec func(string) (string, error)) {
	globalEncryptor = enc
	globalDecryptor = dec
}

func joinConditions(conds []string) string {
	result := ""
	for i, c := range conds {
		if i > 0 {
			result += " AND "
		}
		result += c
	}
	return result
}
