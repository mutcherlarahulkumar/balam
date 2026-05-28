// Package handlers contains all HTTP handler implementations for the Balam API.
package handlers

import (
	"agent-balam/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// respond sends a JSON response with the given status code and body.
func respond(c *gin.Context, status int, body interface{}) {
	c.JSON(status, body)
}

// respondError sends a standardised error response.
func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, models.ErrorResponse{
		Error:   code,
		Message: message,
	})
}

// respondFieldError sends a standardised validation error for a specific field.
func respondFieldError(c *gin.Context, field, message string) {
	c.JSON(http.StatusBadRequest, models.ErrorResponse{
		Error:   "validation_error",
		Message: message,
		Field:   field,
	})
}
