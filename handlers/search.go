package handlers

import (
	"agent-balam/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SearchHandler handles the global search endpoint.
type SearchHandler struct {
	svc *domain.SearchService
}

// NewSearchHandler creates a new SearchHandler.
func NewSearchHandler(svc *domain.SearchService) *SearchHandler {
	return &SearchHandler{svc: svc}
}

// Search handles GET /search?q=.
func (h *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if len(q) < 2 {
		respondError(c, http.StatusBadRequest, "query_too_short", "q must be at least 2 characters")
		return
	}

	results, err := h.svc.Search(q)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Search failed")
		return
	}
	respond(c, http.StatusOK, results)
}
