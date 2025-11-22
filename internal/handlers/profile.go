package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"huddle-backend/internal/database/sqlc"
	"huddle-backend/internal/middleware"
	"huddle-backend/internal/profiles"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService *profile.Service
}

func NewProfileHandler(profileService *profile.Service) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		Username    string  `json:"username" binding:"required"`
		DisplayName *string `json:"display_name"`
		Bio         *string `json:"bio"`
		Website     *string `json:"website"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	if err := h.profileService.ValidateUsername(c.Request.Context(), req.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := sqlc.CreateProfileParams{
		UserID:   userID.(int32),
		Username: req.Username,
		DisplayName: sql.NullString{
			String: getStringValue(req.DisplayName),
			Valid:  req.DisplayName != nil,
		},
		Bio: sql.NullString{
			String: getStringValue(req.Bio),
			Valid:  req.Bio != nil,
		},
		Website: sql.NullString{
			String: getStringValue(req.Website),
			Valid:  req.Website != nil,
		},
	}

	profile, err := h.profileService.CreateProfile(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	profile, err := h.profileService.GetProfileByUserID(c.Request.Context(), userID.(int32))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) GetProfileByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	profile, err := h.profileService.GetProfileByUsername(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) CheckUsernameAvailability(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username query parameter is required"})
		return
	}

	available, err := h.profileService.CheckUsernameAvailability(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check username availability"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  username,
		"available": available,
	})
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		Username    string  `json:"username" binding:"required"`
		DisplayName *string `json:"display_name"`
		Bio         *string `json:"bio"`
		Website     *string `json:"website"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	params := sqlc.UpdateProfileParams{
		UserID:   userID.(int32),
		Username: req.Username,
		DisplayName: sql.NullString{
			String: getStringValue(req.DisplayName),
			Valid:  req.DisplayName != nil,
		},
		Bio: sql.NullString{
			String: getStringValue(req.Bio),
			Valid:  req.Bio != nil,
		},
		Website: sql.NullString{
			String: getStringValue(req.Website),
			Valid:  req.Website != nil,
		},
	}

	profile, err := h.profileService.UpdateProfile(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateUsername(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	profile, err := h.profileService.UpdateUsername(c.Request.Context(), userID.(int32), req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update username", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	if err := h.profileService.DeleteProfile(c.Request.Context(), userID.(int32)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile deleted successfully"})
}

func (h *ProfileHandler) ListProfiles(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)

	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 20
	}

	profiles, err := h.profileService.ListProfiles(c.Request.Context(), int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list profiles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profiles": profiles,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *ProfileHandler) SearchProfiles(c *gin.Context) {
	searchTerm := c.Query("q")
	if searchTerm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query parameter 'q' is required"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)

	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 20
	}

	profiles, err := h.profileService.SearchProfilesByUsername(c.Request.Context(), searchTerm, int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search profiles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profiles": profiles,
		"query":    searchTerm,
		"limit":    limit,
		"offset":   offset,
	})
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
