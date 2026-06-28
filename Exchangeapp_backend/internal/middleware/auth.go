package middleware

import (
	"exchangeapp/internal/repository"
	"exchangeapp/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwt *auth.JWTManager, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		username, err := jwt.ParseToken(token)
		if err != nil || username == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("username", username)

		// Look up user — abort if user not found (deleted, DB error, etc.)
		user, err := userRepo.FindByUsername(username)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		c.Set("userID", user.ID)
		c.Next()
	}
}
