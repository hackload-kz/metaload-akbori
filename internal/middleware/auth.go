package middleware

import (
	"biletter-service/internal/models"
	"biletter-service/internal/services"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	UserContextKey = "current_user"
)

func BasicAuth(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Basic ") {
			c.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		encoded := strings.TrimPrefix(authHeader, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			c.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid base64 encoding"})
			c.Abort()
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			c.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials format"})
			c.Abort()
			return
		}

		email := credentials[0]
		password := credentials[1]

		user, err := userService.ValidateCredentials(email, password)
		if err != nil {
			c.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			c.Abort()
			return
		}

		// Сохраняем пользователя в контексте
		c.Set(UserContextKey, user)
		c.Next()
	}
}

// GetCurrentUser извлекает текущего пользователя из контекста
func GetCurrentUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil, false
	}

	currentUser, ok := user.(*models.User)
	return currentUser, ok
}
