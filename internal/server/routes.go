package server

import (
    "net/http"

    "huddle-backend/internal/handlers"
    "huddle-backend/internal/middleware"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
        AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
        AllowCredentials: true,
    }))

    r.GET("/", s.HelloWorldHandler)
    r.GET("/health", s.healthHandler)

    auth := r.Group("/auth")
    {
        auth.GET("/:provider", s.beginAuthHandler)
        auth.GET("/:provider/callback", s.callbackAuthHandler)
        auth.POST("/logout", s.logoutHandler)
    }

    profileHandler := handlers.NewProfileHandler(s.profileService)

    api := r.Group("/api")
    api.Use(middleware.RequireAuth(s.authService))
    {
        api.GET("/me", s.getCurrentUserHandler)

        profiles := api.Group("/profiles")
        {
            profiles.POST("", profileHandler.CreateProfile)
            profiles.GET("/me", profileHandler.GetMyProfile)
            profiles.GET("/check-username", profileHandler.CheckUsernameAvailability)
            profiles.GET("/search", profileHandler.SearchProfiles)
            profiles.GET("", profileHandler.ListProfiles)
            profiles.GET("/:username", profileHandler.GetProfileByUsername)
            profiles.PUT("", profileHandler.UpdateProfile)
            profiles.PATCH("/username", profileHandler.UpdateUsername)
            profiles.DELETE("", profileHandler.DeleteProfile)
        }
    }

    return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}

func (s *Server) healthHandler(c *gin.Context) {
    c.JSON(http.StatusOK, s.db.Health())
}
