package server

import (
	"log"
	"net/http"

	"huddle-backend/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func (s *Server) beginAuthHandler(c *gin.Context) {
	provider := c.Param("provider")
	log.Printf("BeginAuth - Provider from URL: %s", provider)

	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider required"})
		return
	}

	// Gothic expects the provider in the URL query params
	q := c.Request.URL.Query()
	q.Set("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	log.Printf("BeginAuth - Full URL: %s", c.Request.URL.String())

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (s *Server) callbackAuthHandler(c *gin.Context) {
	log.Println("=== Callback Handler Started ===")

	provider := c.Param("provider")
	log.Printf("Callback - Provider from URL: %s", provider)

	// Gothic needs the provider in query params
	q := c.Request.URL.Query()
	q.Set("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	log.Printf("Callback - Full URL: %s", c.Request.URL.String())

	// Complete OAuth authentication
	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Printf("CompleteUserAuth error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed", "details": err.Error()})
		return
	}
	log.Printf("Gothic user retrieved: %+v", gothUser)

	// Find or create user in database
	user, err := s.authService.FindOrCreateOAuthUser(c.Request.Context(), gothUser)
	if err != nil {
		log.Printf("FindOrCreateOAuthUser error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user", "details": err.Error()})
		return
	}
	log.Printf("User saved/found: %+v", user)

	// Create session
	session, err := auth.Store.Get(c.Request, auth.SessionName)
	if err != nil {
		log.Printf("Session Get error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session error", "details": err.Error()})
		return
	}

	session.Values["user_id"] = user.ID
	session.Values["authenticated"] = true

	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("Session Save error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session", "details": err.Error()})
		return
	}

	log.Println("=== Authentication Successful ===")
	c.JSON(http.StatusOK, gin.H{
		"message": "authentication successful",
		"user":    user,
	})
}

func (s *Server) logoutHandler(c *gin.Context) {
	session, err := auth.Store.Get(c.Request, auth.SessionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session error"})
		return
	}

	session.Values["user_id"] = nil
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1

	if err := session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear session"})
		return
	}

	gothic.Logout(c.Writer, c.Request)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (s *Server) getCurrentUserHandler(c *gin.Context) {
	session, err := auth.Store.Get(c.Request, auth.SessionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session error"})
		return
	}

	userID, ok := session.Values["user_id"].(int32)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	user, err := s.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
