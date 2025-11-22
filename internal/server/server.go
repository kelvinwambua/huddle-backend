package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"huddle-backend/internal/auth"
	"huddle-backend/internal/database"
	"huddle-backend/internal/database/sqlc"
	"huddle-backend/internal/profiles"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port           int
	db             database.Service
	queries        *sqlc.Queries
	authService    *auth.Service
	profileService *profile.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	db := database.New()
	queries := sqlc.New(db.DB())

	auth.InitAuth()

	NewServer := &Server{
		port:           port,
		db:             db,
		queries:        queries,
		authService:    auth.NewService(queries),
		profileService: profile.NewService(queries),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
