// Package handlers contains all HTTP handler implementations for the Balam API.
package handlers

import (
	"agent-balam/models"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// respond sends a JSON response with the given status code and body.
func respond(c *gin.Context, status int, body interface{}) {
	c.JSON(status, body)
}

// respondError sends a standardised error envelope.
func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, models.ErrorResponse{Error: code, Message: message})
}

// respondValidation parses a ShouldBindJSON error and returns one friendly
// validation error per field. Internal/unexpected errors are logged and mapped
// to a generic server error so raw Go internals never reach the frontend.
func respondValidation(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		type fieldErr struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}
		errs := make([]fieldErr, 0, len(ve))
		for _, fe := range ve {
			errs = append(errs, fieldErr{
				Field:   jsonFieldName(fe),
				Message: friendlyMsg(fe),
			})
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "validation_error",
			"errors": errs,
		})
		return
	}
	// JSON syntax / type errors — tell client their payload is malformed.
	if strings.Contains(err.Error(), "invalid character") ||
		strings.Contains(err.Error(), "cannot unmarshal") {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_json",
			Message: "Request body is not valid JSON.",
		})
		return
	}
	// Unexpected bind error — log it, return generic message.
	slog.Error("unexpected bind error", "err", err, "path", c.FullPath())
	c.JSON(http.StatusBadRequest, models.ErrorResponse{
		Error:   "bad_request",
		Message: "Could not read request body.",
	})
}

// jsonFieldName returns the JSON key name from the struct tag so the frontend
// sees "agentCode" not "AgentCode".
func jsonFieldName(fe validator.FieldError) string {
	parts := strings.Split(fe.Namespace(), ".")
	// Drop the struct name prefix (first element).
	if len(parts) > 1 {
		parts = parts[1:]
	}
	return strings.Join(parts, ".")
}

// friendlyMsg converts a validator.FieldError into a plain-English sentence.
func friendlyMsg(fe validator.FieldError) string {
	field := jsonFieldName(fe)
	param := fe.Param()
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required.", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address.", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters.", field, param)
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters.", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters.", field, param)
	case "numeric":
		return fmt.Sprintf("%s must contain digits only.", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and digits.", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s.", field, strings.ReplaceAll(param, " ", ", "))
	case "gt":
		return fmt.Sprintf("%s must be greater than %s.", field, param)
	case "gte":
		return fmt.Sprintf("%s must be %s or more.", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s.", field, param)
	case "lte":
		return fmt.Sprintf("%s must be %s or less.", field, param)
	case "url":
		return fmt.Sprintf("%s must be a valid URL.", field)
	default:
		return fmt.Sprintf("%s is invalid.", field)
	}
}
