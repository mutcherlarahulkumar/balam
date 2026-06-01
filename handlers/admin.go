package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// AdminHandler handles admin-only operations like bulk data import.
type AdminHandler struct {
	db *sqlx.DB
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(db *sqlx.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// allowedSeedPrefixes are statement types allowed in the full seed file (includes TRUNCATE).
var allowedSeedPrefixes = []string{
	"begin", "commit", "truncate", "insert", "select setval(", "--", "\n", "\r",
}

// allowedImportPrefixes are statement types for user-supplied import files (no TRUNCATE — server handles clearing).
var allowedImportPrefixes = []string{
	"insert", "select setval(", "--", "\n", "\r",
}

// resetTablesSQL clears all user data and resets sequences so imports start fresh.
const resetTablesSQL = `TRUNCATE TABLE
	login_attempts, rpt_prcalendar, rpt_currstatus, rpt_cashinout,
	rpt_cashflow, prospect_activity, leaddata, documents, bankdetails,
	fupupdate, fuphistory, multipledue, sb, commission,
	loan, policies, client, family, plans, agent
RESTART IDENTITY CASCADE;`

func isAllowedSQL(sql string, prefixes []string) bool {
	for _, line := range strings.Split(sql, "\n") {
		trimmed := strings.ToLower(strings.TrimSpace(line))
		if trimmed == "" {
			continue
		}
		allowed := false
		for _, pfx := range prefixes {
			if strings.HasPrefix(trimmed, pfx) {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}
	return true
}

func readUploadedSQL(c *gin.Context) (content string, filename string, ok bool) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		respondError(c, http.StatusBadRequest, "missing_file", "Attach a .sql file in the 'file' field.")
		return "", "", false
	}
	defer file.Close()

	if header.Size > 50<<20 {
		respondError(c, http.StatusRequestEntityTooLarge, "file_too_large", "SQL file must be under 50 MB.")
		return "", "", false
	}

	raw, err := io.ReadAll(file)
	if err != nil {
		slog.Error("import: read file", "err", err)
		respondError(c, http.StatusInternalServerError, "server_error", "Could not read uploaded file.")
		return "", "", false
	}
	return string(raw), header.Filename, true
}

// ImportSQL handles POST /v1/admin/import-sql.
// Accepts a multipart .sql file produced by the Balam migrate script.
// The file may contain TRUNCATE + INSERT + SELECT setval() statements.
func (h *AdminHandler) ImportSQL(c *gin.Context) {
	content, filename, ok := readUploadedSQL(c)
	if !ok {
		return
	}

	if !isAllowedSQL(content, allowedSeedPrefixes) {
		slog.Warn("import-sql: rejected — contains disallowed SQL", "file", filename)
		respondError(c, http.StatusForbidden, "disallowed_sql",
			"Only INSERT, TRUNCATE, BEGIN, COMMIT, and SELECT setval() statements are allowed.")
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		slog.Error("import-sql: begin tx", "err", err)
		respondError(c, http.StatusInternalServerError, "server_error", "Could not start transaction.")
		return
	}

	if _, err := tx.Exec(content); err != nil {
		_ = tx.Rollback()
		slog.Error("import-sql: exec failed", "file", filename, "err", err)
		respondError(c, http.StatusUnprocessableEntity, "import_failed", err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		slog.Error("import-sql: commit", "err", err)
		respondError(c, http.StatusInternalServerError, "server_error", "Could not commit import.")
		return
	}

	slog.Info("import-sql: complete", "file", filename)
	respond(c, http.StatusOK, gin.H{"message": "Import complete.", "file": filename})
}

// ResetAndImport handles POST /v1/agent/import.
// Accepts a multipart INSERT-only .sql file.
// Server automatically clears all existing data before importing —
// the file should NOT contain TRUNCATE statements.
func (h *AdminHandler) ResetAndImport(c *gin.Context) {
	content, filename, ok := readUploadedSQL(c)
	if !ok {
		return
	}

	if !isAllowedSQL(content, allowedImportPrefixes) {
		slog.Warn("agent-import: rejected — contains disallowed SQL", "file", filename)
		respondError(c, http.StatusForbidden, "disallowed_sql",
			"Import file may only contain INSERT and SELECT setval() statements. Do not include TRUNCATE — the server clears existing data automatically.")
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		slog.Error("agent-import: begin tx", "err", err)
		respondError(c, http.StatusInternalServerError, "server_error", "Could not start transaction.")
		return
	}

	// Clear all existing data first.
	if _, err := tx.Exec(resetTablesSQL); err != nil {
		_ = tx.Rollback()
		slog.Error("agent-import: reset failed", "err", err)
		respondError(c, http.StatusInternalServerError, "server_error", "Could not reset existing data.")
		return
	}

	if _, err := tx.Exec(content); err != nil {
		_ = tx.Rollback()
		slog.Error("agent-import: exec failed", "file", filename, "err", err)
		respondError(c, http.StatusUnprocessableEntity, "import_failed", err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		slog.Error("agent-import: commit", "err", err)
		respondError(c, http.StatusInternalServerError, "server_error", "Could not commit import.")
		return
	}

	slog.Info("agent-import: complete", "file", filename)
	respond(c, http.StatusOK, gin.H{"message": "Data reset and import complete.", "file": filename})
}
