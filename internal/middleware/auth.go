package middleware

import (
	"net/http"

	"huddle-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey    = "user_id"
	SessionIDKey = "session_id"
)

func RequireAuth(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieSession, err := auth.Store.Get(c.Request, auth.SessionName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		sessionID, ok := cookieSession.Values[auth.SessionIDKey].(string)
		if !ok || sessionID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		sessionData, err := authService.GetSessionByID(c.Request.Context(), sessionID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired or invalid"})
			c.Abort()
			return
		}

		c.Set(UserIDKey, sessionData.UserID)
		c.Set(SessionIDKey, sessionID)
		c.Next()
	}
}
