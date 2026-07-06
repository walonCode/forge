package users

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUsernameTaken    = errors.New("username already taken")
	ErrInvalidPassword  = errors.New("current password is incorrect")
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)

type Service struct {
	repo Repository
}

func newService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetProfile(ctx context.Context, userId string) (*UserProfile, error) {
	profile, err := s.repo.FindProfileByID(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return profile, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userId string, req UpdateProfileRequest) (*UserProfile, error) {
	current, err := s.repo.FindProfileByID(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// overlay only the fields the client actually provided
	name := current.Name
	if v := strings.TrimSpace(req.Name); v != "" {
		name = v
	}

	username := current.Username
	if v := strings.TrimSpace(req.Username); v != "" {
		username = v
	}

	if name == current.Name && username == current.Username {
		return nil, ErrNoFieldsToUpdate
	}

	if username != current.Username {
		taken, err := s.repo.ExistsByUsernameExcept(ctx, username, userId)
		if err != nil {
			return nil, err
		}
		if taken {
			return nil, ErrUsernameTaken
		}
	}

	updated, err := s.repo.UpdateProfile(ctx, userId, name, username)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) UpdatePassword(ctx context.Context, userId string, req UpdatePasswordRequest) error {
	hash, err := s.repo.GetPasswordHash(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.CurrentPassword)); err != nil {
		return ErrInvalidPassword
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userId, string(newHash))
}

func (s *Service) DeleteAccount(ctx context.Context, userId string) error {
	if err := s.repo.DeleteUser(ctx, userId); err != nil {
		return err
	}

	return nil
}
