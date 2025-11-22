package profile

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"huddle-backend/internal/database/sqlc"
)

type Service struct {
	queries *sqlc.Queries
}

func NewService(queries *sqlc.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) CreateProfile(ctx context.Context, params sqlc.CreateProfileParams) (sqlc.Profile, error) {

	exists, err := s.queries.CheckUsernameExists(ctx, params.Username)
	if err != nil {
		return sqlc.Profile{}, fmt.Errorf("error checking username availability: %w", err)
	}
	if exists {
		return sqlc.Profile{}, fmt.Errorf("username '%s' is already taken", params.Username)
	}

	profile, err := s.queries.CreateProfile(ctx, params)
	if err != nil {
		return sqlc.Profile{}, fmt.Errorf("error creating profile: %w", err)
	}

	return profile, nil
}

func (s *Service) GetProfileByUserID(ctx context.Context, userID int32) (sqlc.Profile, error) {
	profile, err := s.queries.GetProfileByUserID(ctx, userID)
	if err == sql.ErrNoRows {
		return sqlc.Profile{}, fmt.Errorf("profile not found for user ID %d", userID)
	}
	if err != nil {
		return sqlc.Profile{}, fmt.Errorf("error getting profile: %w", err)
	}
	return profile, nil
}

func (s *Service) GetProfileByUsername(ctx context.Context, username string) (sqlc.Profile, error) {
	profile, err := s.queries.GetProfileByUsername(ctx, username)
	if err == sql.ErrNoRows {
		return sqlc.Profile{}, fmt.Errorf("profile not found for username '%s'", username)
	}
	if err != nil {
		return sqlc.Profile{}, fmt.Errorf("error getting profile: %w", err)
	}
	return profile, nil
}

func (s *Service) CheckUsernameAvailability(ctx context.Context, username string) (bool, error) {
	exists, err := s.queries.CheckUsernameExists(ctx, username)
	if err != nil {
		return false, fmt.Errorf("error checking username availability: %w", err)
	}
	return !exists, nil
}

func (s *Service) UpdateProfile(ctx context.Context, params sqlc.UpdateProfileParams) (sqlc.Profile, error) {

	currentProfile, err := s.queries.GetProfileByUserID(ctx, params.UserID)
	if err != nil {
		return sqlc.Profile{}, fmt.Errorf("error getting current profile: %w", err)
	}

	if currentProfile.Username != params.Username {
		exists, err := s.queries.CheckUsernameExists(ctx, params.Username)
		if err != nil {
			return sqlc.Profile{}, fmt.Errorf("error checking username availability: %w", err)
		}
		if exists {
			return sqlc.Profile{}, fmt.Errorf("username '%s' is already taken", params.Username)
		}
	}

	profile, err := s.queries.UpdateProfile(ctx, params)
	if err != nil {
		return sqlc.Profile{}, fmt.Errorf("error updating profile: %w", err)
	}

	return profile, nil
}

func (s *Service) UpdateUsername(ctx context.Context, userID int32, newUsername string) (sqlc.Profile, error) {

	exists, err := s.queries.CheckUsernameExists(ctx, newUsername)
	if err != nil {
		return sqlc.Profile{}, fmt.Errorf("error checking username availability: %w", err)
	}
	if exists {
		return sqlc.Profile{}, fmt.Errorf("username '%s' is already taken", newUsername)
	}

	profile, err := s.queries.UpdateUsername(ctx, sqlc.UpdateUsernameParams{
		UserID:   userID,
		Username: newUsername,
	})
	if err != nil {
		return sqlc.Profile{}, fmt.Errorf("error updating username: %w", err)
	}

	return profile, nil
}

func (s *Service) DeleteProfile(ctx context.Context, userID int32) error {
	err := s.queries.DeleteProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("error deleting profile: %w", err)
	}
	return nil
}

func (s *Service) ListProfiles(ctx context.Context, limit, offset int32) ([]sqlc.Profile, error) {
	profiles, err := s.queries.ListProfiles(ctx, sqlc.ListProfilesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("error listing profiles: %w", err)
	}
	return profiles, nil
}

func (s *Service) SearchProfilesByUsername(ctx context.Context, searchTerm string, limit, offset int32) ([]sqlc.Profile, error) {

	searchPattern := fmt.Sprintf("%%%s%%", strings.ToLower(searchTerm))

	profiles, err := s.queries.SearchProfilesByUsername(ctx, sqlc.SearchProfilesByUsernameParams{
		Username: searchPattern,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("error searching profiles: %w", err)
	}
	return profiles, nil
}

func (s *Service) ValidateUsername(ctx context.Context, username string) error {

	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if len(username) > 30 {
		return fmt.Errorf("username must not exceed 30 characters")
	}

	for _, char := range username {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
		}
	}

	available, err := s.CheckUsernameAvailability(ctx, username)
	if err != nil {
		return err
	}
	if !available {
		return fmt.Errorf("username '%s' is already taken", username)
	}

	return nil
}
