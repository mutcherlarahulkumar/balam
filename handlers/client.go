package handlers

import (
	"agent-balam/domain"
	"agent-balam/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ClientHandler handles client endpoints.
type ClientHandler struct {
	svc *domain.ClientService
}

// NewClientHandler creates a new ClientHandler.
func NewClientHandler(svc *domain.ClientService) *ClientHandler {
	return &ClientHandler{svc: svc}
}

// List handles GET /clients.
func (h *ClientHandler) List(c *gin.Context) {
	search := c.Query("search")
	familyCode := c.Query("familyCode")
	clientType := c.Query("clientType")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	items, total, err := h.svc.List(search, familyCode, clientType, page, limit)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Failed to list clients")
		return
	}
	respond(c, http.StatusOK, models.PaginatedResponse{Data: items, Total: total, Page: page, Limit: limit})
}

// Get handles GET /clients/:id.
func (h *ClientHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_id", "Client ID must be an integer")
		return
	}

	detail, err := h.svc.GetByID(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "client_not_found", "Client not found")
		return
	}
	respond(c, http.StatusOK, detail)
}

// Search handles GET /clients/search.
func (h *ClientHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		respondError(c, http.StatusBadRequest, "missing_query", "Query parameter 'q' is required")
		return
	}

	clients, err := h.svc.Search(q)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "server_error", "Search failed")
		return
	}
	respond(c, http.StatusOK, gin.H{"data": clients})
}

// Create handles POST /clients.
func (h *ClientHandler) Create(c *gin.Context) {
	var req models.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	client, err := h.svc.Create(req)
	if err != nil {
		switch err.Error() {
		case "family_not_found":
			respondError(c, http.StatusNotFound, "family_not_found", "Family not found")
		default:
			respondError(c, http.StatusBadRequest, "create_failed", err.Error())
		}
		return
	}
	respond(c, http.StatusCreated, client)
}

// Update handles PUT /clients/:id.
func (h *ClientHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_id", "Client ID must be an integer")
		return
	}

	var req models.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	client, err := h.svc.Update(id, req)
	if err != nil {
		switch err.Error() {
		case "client_not_found":
			respondError(c, http.StatusNotFound, "client_not_found", "Client not found")
		default:
			respondError(c, http.StatusBadRequest, "update_failed", err.Error())
		}
		return
	}
	respond(c, http.StatusOK, client)
}
