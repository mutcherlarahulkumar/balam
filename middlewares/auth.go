// Package middlewares provides Gin middleware for the Balam API.
package middlewares

import (
	"agent-balam/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const agentIDKey = "agentID"
const agentNameKey = "agentName"

// Auth validates the Bearer JWT and injects agentID into the context.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "missing_token",
				"message": "Authorization header is required",
			})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		agentID, err := domain.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token",
				"message": "Token is invalid or expired",
			})
			return
		}

		c.Set(agentIDKey, agentID)
		c.Next()
	}
}

// AgentID retrieves the authenticated agent ID from context.
func AgentID(c *gin.Context) int {
	id, _ := c.Get(agentIDKey)
	v, _ := id.(int)
	return v
}
