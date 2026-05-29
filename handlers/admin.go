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

// allowedPrefixes are the only SQL statement types the import endpoint will execute.
// DROP, ALTER, CREATE, GRANT, COPY etc. are all rejected.
var allowedPrefixes = []string{
	"begin", "commit", "truncate", "insert", "select setval(", "--", "\n", "\r",
}

func isSeedSQL(sql string) bool {
	for _, line := range strings.Split(sql, "\n") {
		trimmed := strings.ToLower(strings.TrimSpace(line))
		if trimmed == "" {
			continue
		}
		allowed := false
		for _, pfx := range allowedPrefixes {
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

// ImportSQL handles POST /v1/admin/import-sql.
// Accepts a multipart .sql file produced by the Balam migrate script.
// Only INSERT, TRUNCATE, BEGIN, and COMMIT statements are permitted.
func (h *AdminHandler) ImportSQL(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		respondError(c, http.StatusBadRequest, "missing_file", "Attach a .sql file in the 'file' field.")
		return
	}
	defer file.Close()

	if header.Size > 50<<20 {
		respondError(c, http.StatusRequestEntityTooLarge, "file_too_large", "SQL file must be under 50 MB.")
		return
	}

	raw, err := io.ReadAll(file)
	if err != nil {
		slog.Error("import-sql: read file", "err", err)
		respondError(c, http.StatusInternalServerError, "server_error", "Could not read uploaded file.")
		return
	}

	content := string(raw)
	if !isSeedSQL(content) {
		slog.Warn("import-sql: rejected — contains non-seed SQL", "file", header.Filename)
		respondError(c, http.StatusForbidden, "disallowed_sql",
			"Only INSERT, TRUNCATE, BEGIN, and COMMIT statements are allowed in import files.")
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
		slog.Error("import-sql: exec failed", "file", header.Filename, "err", err)
		respondError(c, http.StatusUnprocessableEntity, "import_failed", "SQL execution failed. Check server logs for details.")
		return
	}

	if err := tx.Commit(); err != nil {
		slog.Error("import-sql: commit", "err", err)
		respondError(c, http.StatusInternalServerError, "server_error", "Could not commit import.")
		return
	}

	slog.Info("import-sql: complete", "file", header.Filename, "bytes", header.Size)
	respond(c, http.StatusOK, gin.H{"message": "Import complete.", "file": header.Filename, "bytes": header.Size})
}
