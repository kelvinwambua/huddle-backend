package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"huddle-backend/internal/database/sqlc"

	"github.com/markbates/goth"
)

type Service struct {
	queries *sqlc.Queries
}

func NewService(queries *sqlc.Queries) *Service {
	return &Service{queries: queries}
}

func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *Service) FindOrCreateOAuthUser(ctx context.Context, gothUser goth.User) (sqlc.User, error) {
	log.Printf("=== FindOrCreateOAuthUser ===")
	log.Printf("Provider: %s, UserID: %s, Email: %s", gothUser.Provider, gothUser.UserID, gothUser.Email)

	user, err := s.queries.GetUserByProviderID(ctx, sqlc.GetUserByProviderIDParams{
		Provider:       sql.NullString{String: gothUser.Provider, Valid: true},
		ProviderUserID: sql.NullString{String: gothUser.UserID, Valid: true},
	})

	if err == nil {
		log.Printf("User found, updating tokens")
		return s.queries.UpdateUserOAuthTokens(ctx, sqlc.UpdateUserOAuthTokensParams{
			ID:           user.ID,
			AccessToken:  sql.NullString{String: gothUser.AccessToken, Valid: true},
			RefreshToken: sql.NullString{String: gothUser.RefreshToken, Valid: true},
			ExpiresAt:    sql.NullTime{Time: gothUser.ExpiresAt, Valid: !gothUser.ExpiresAt.IsZero()},
		})
	}

	if err != sql.ErrNoRows {
		log.Printf("Error checking existing user: %v", err)
		return sqlc.User{}, fmt.Errorf("error checking existing user: %w", err)
	}

	log.Printf("User not found, creating new user")

	username := gothUser.NickName
	if username == "" {
		username = gothUser.Email
	}

	params := sqlc.CreateOAuthUserParams{
		Username:       username,
		Email:          gothUser.Email,
		AvatarUrl:      sql.NullString{String: gothUser.AvatarURL, Valid: gothUser.AvatarURL != ""},
		Provider:       sql.NullString{String: gothUser.Provider, Valid: true},
		ProviderUserID: sql.NullString{String: gothUser.UserID, Valid: true},
		AccessToken:    sql.NullString{String: gothUser.AccessToken, Valid: true},
		RefreshToken:   sql.NullString{String: gothUser.RefreshToken, Valid: true},
		ExpiresAt:      sql.NullTime{Time: gothUser.ExpiresAt, Valid: !gothUser.ExpiresAt.IsZero()},
		Name:           sql.NullString{String: gothUser.Name, Valid: gothUser.Name != ""},
		FirstName:      sql.NullString{String: gothUser.FirstName, Valid: gothUser.FirstName != ""},
		LastName:       sql.NullString{String: gothUser.LastName, Valid: gothUser.LastName != ""},
		NickName:       sql.NullString{String: gothUser.NickName, Valid: gothUser.NickName != ""},
		Description:    sql.NullString{String: gothUser.Description, Valid: gothUser.Description != ""},
		Location:       sql.NullString{String: gothUser.Location, Valid: gothUser.Location != ""},
	}

	newUser, err := s.queries.CreateOAuthUser(ctx, params)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return sqlc.User{}, fmt.Errorf("error creating user: %w", err)
	}

	log.Printf("User created successfully: %+v", newUser)
	return newUser, nil
}

func (s *Service) CreateSession(ctx context.Context, userID int32, provider, ipAddress, userAgent string) (sqlc.Session, error) {
	sessionID, err := GenerateSessionID()
	if err != nil {
		return sqlc.Session{}, fmt.Errorf("failed to generate session ID: %w", err)
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	return s.queries.CreateSession(ctx, sqlc.CreateSessionParams{
		ID:        sessionID,
		UserID:    userID,
		Provider:  sql.NullString{String: provider, Valid: provider != ""},
		IpAddress: sql.NullString{String: ipAddress, Valid: ipAddress != ""},
		UserAgent: sql.NullString{String: userAgent, Valid: userAgent != ""},
		ExpiresAt: expiresAt,
	})
}

func (s *Service) GetSessionByID(ctx context.Context, sessionID string) (sqlc.GetSessionByIDRow, error) {
	return s.queries.GetSessionByID(ctx, sessionID)
}

func (s *Service) DeleteSession(ctx context.Context, sessionID string) error {
	return s.queries.DeleteSession(ctx, sessionID)
}

func (s *Service) DeleteUserSessions(ctx context.Context, userID int32) error {
	return s.queries.DeleteUserSessions(ctx, userID)
}

func (s *Service) GetUserByID(ctx context.Context, id int32) (sqlc.User, error) {
	return s.queries.GetUserByID(ctx, id)
}
