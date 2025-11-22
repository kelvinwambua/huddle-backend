package middleware

import (
	"context"
	"net/http"

	"huddle-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := auth.Store.Get(c.Request, auth.SessionName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		authenticated, ok := session.Values["authenticated"].(bool)
		if !ok || !authenticated {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userID, ok := session.Values["user_id"].(int32)
		if !ok || userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}

func GetUserIDFromContext(ctx context.Context) (int32, bool) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		userID, exists := ginCtx.Get(UserIDKey)
		if !exists {
			return 0, false
		}
		id, ok := userID.(int32)
		return id, ok
	}
	return 0, false
}
