package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"huddle-backend/internal/database/sqlc"

	"github.com/markbates/goth"
)

type Service struct {
	queries *sqlc.Queries
}

func NewService(queries *sqlc.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) FindOrCreateOAuthUser(ctx context.Context, gothUser goth.User) (sqlc.User, error) {
	log.Printf("=== FindOrCreateOAuthUser ===")
	log.Printf("Provider: %s, UserID: %s, Email: %s", gothUser.Provider, gothUser.UserID, gothUser.Email)

	// Try to find existing user by provider ID
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

	log.Printf("Creating user with params: %+v", params)

	newUser, err := s.queries.CreateOAuthUser(ctx, params)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return sqlc.User{}, fmt.Errorf("error creating user: %w", err)
	}

	log.Printf("User created successfully: %+v", newUser)
	return newUser, nil
}

func (s *Service) GetUserByID(ctx context.Context, id int32) (sqlc.User, error) {
	return s.queries.GetUserByID(ctx, id)
}
