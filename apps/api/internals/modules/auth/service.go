package auth

import (
	"context"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

// bcryptCost is the work factor used when hashing passwords.
const bcryptCost = 10

type Service struct {
	repo Repository
}

func newService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) FindUserByUsername(ctx context.Context, username string) (*DbUser, error) {
	user, err := s.repo.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, param SignupRequest) (string, error) {
	id, err := gonanoid.New(20)
	if err != nil {
		return "", err
	}

	//hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcryptCost)
	if err != nil {
		return "", err
	}

	userId, err := s.repo.CreateUser(ctx, CreateDbUser{
		Username: param.Username,
		Password: string(passwordHash),
		Name:     param.Name,
		ID:       id,
	})

	if err != nil {
		return "", err
	}

	return userId, nil
}
