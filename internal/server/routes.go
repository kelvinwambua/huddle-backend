package server

import (
    "net/http"

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

    api := r.Group("/api")
    api.Use(middleware.RequireAuth())
    {
        api.GET("/me", s.getCurrentUserHandler)
    }

    return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}

func (s *Server) healthHandler(c *gin.Context) {
    c.JSON(http.StatusOK, s.db.Health())
}
