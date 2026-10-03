package handler

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// getUserID extracts the authenticated user's ID from gin.Context with safe type assertion.
// Returns 0 and sends 401 response if not authenticated or type assertion fails.
func getUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return 0, false
	}
	uid, ok := val.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
		return 0, false
	}
	return uid, true
}

// getUserIDOptional extracts userID for optional auth contexts (e.g., feed).
// Returns 0 if not authenticated (does not send error response).
func getUserIDOptional(c *gin.Context) uint {
	val, exists := c.Get("userID")
	if !exists {
		return 0
	}
	uid, ok := val.(uint)
	if !ok {
		return 0
	}
	return uid
}

// currencyCodeRe matches valid 3-letter uppercase currency codes.
var currencyCodeRe = regexp.MustCompile(`^[A-Z]{3}$`)

// validateCurrencyCode checks if a currency code is a valid 3-letter uppercase string.
func validateCurrencyCode(code string) bool {
	return currencyCodeRe.MatchString(strings.ToUpper(code))
}

// sanitizeCurrencyCode normalizes and validates a currency code.
// Returns the uppercase code and true if valid, or empty string and false if invalid.
func sanitizeCurrencyCode(code string) (string, bool) {
	code = strings.TrimSpace(strings.ToUpper(code))
	if !currencyCodeRe.MatchString(code) {
		return "", false
	}
	return code, true
}

// genericError returns a generic error message to the client.
// The actual error should be logged server-side before calling this.
func genericError(c *gin.Context, status int, clientMsg string) {
	c.JSON(status, gin.H{"error": clientMsg})
}
