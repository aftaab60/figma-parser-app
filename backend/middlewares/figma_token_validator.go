package middlewares

import (
	"context"
	"net/http"
	"parser-service/internal/figma_manager"
	"strings"

	"github.com/gin-gonic/gin"
)

// middleware to apply in routes for Figma token validation
func ValidateFigmaToken(figmaManager figma_manager.IFigmaManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// token from Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			// if no token in header, check query parameter
			token = c.Query("figma_token")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Figma token is required",
				})
				c.Abort()
				return
			}
		}
		// Remove "Bearer " prefix if present
		if key, ok := strings.CutPrefix(token, "Bearer "); ok {
			token = key
		}
		// if token is empty after removing "Bearer ", return unauthorized
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Figma token is required",
			})
			c.Abort()
			return
		}

		// Validate token format
		if !strings.HasPrefix(token, "figd_") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Figma token format. Token must start with 'figd_'",
			})
			c.Abort()
			return
		}

		// Validate token with Figma API
		if err := figmaManager.ValidateFigmaToken(token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid Figma token",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// Store the validated token in context for use by handlers
		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), "figma_token", token),
		)

		c.Next()
	}
}
