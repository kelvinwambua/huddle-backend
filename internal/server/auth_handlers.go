package server

import (
	"log"
	"net/http"
	"os"

	"huddle-backend/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

func (s *Server) beginAuthHandler(c *gin.Context) {
	provider := c.Param("provider")
	log.Printf("BeginAuth - Provider: %s", provider)

	c.Request = setProviderInContext(c.Request, provider)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (s *Server) callbackAuthHandler(c *gin.Context) {
	log.Println("=== Callback Handler Started ===")

	provider := c.Param("provider")
	c.Request = setProviderInContext(c.Request, provider)

	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Printf("CompleteUserAuth error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed", "details": err.Error()})
		return
	}

	user, err := s.authService.FindOrCreateOAuthUser(c.Request.Context(), gothUser)
	if err != nil {
		log.Printf("FindOrCreateOAuthUser error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user", "details": err.Error()})
		return
	}

	session, err := s.authService.CreateSession(
		c.Request.Context(),
		user.ID,
		provider,
		c.ClientIP(),
		c.Request.UserAgent(),
	)
	if err != nil {
		log.Printf("CreateSession error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	cookieSession, err := auth.Store.Get(c.Request, auth.SessionName)
	if err != nil {
		log.Printf("Cookie session error: %v", err)
	}

	cookieSession.Values[auth.SessionIDKey] = session.ID
	if err := cookieSession.Save(c.Request, c.Writer); err != nil {
		log.Printf("Failed to save cookie: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	log.Println("=== Authentication Successful ===")
	redirectURL := os.Getenv("FRONTEND_URL") + "/"
	c.Redirect(http.StatusFound, redirectURL)
	c.JSON(http.StatusOK, gin.H{
		"message": "authentication successful",
		"user":    user,
	})
}

func (s *Server) logoutHandler(c *gin.Context) {
	cookieSession, err := auth.Store.Get(c.Request, auth.SessionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session error"})
		return
	}

	sessionID, ok := cookieSession.Values[auth.SessionIDKey].(string)
	if ok && sessionID != "" {

		s.authService.DeleteSession(c.Request.Context(), sessionID)
	}

	cookieSession.Options.MaxAge = -1
	cookieSession.Save(c.Request, c.Writer)

	gothic.Logout(c.Writer, c.Request)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (s *Server) getCurrentUserHandler(c *gin.Context) {
	cookieSession, err := auth.Store.Get(c.Request, auth.SessionName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	sessionID, ok := cookieSession.Values[auth.SessionIDKey].(string)
	if !ok || sessionID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	sessionData, err := s.authService.GetSessionByID(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired or invalid"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         sessionData.ID,
		"username":   sessionData.Username,
		"email":      sessionData.Email,
		"avatar_url": sessionData.AvatarUrl,
		"provider":   sessionData.Provider,
	})
}

func setProviderInContext(r *http.Request, provider string) *http.Request {
	return mux.SetURLVars(r, map[string]string{"provider": provider})
}
