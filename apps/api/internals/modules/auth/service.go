package auth

import (
	"context"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func newService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) FindUserByUsername(ctx context.Context, email string) (*DbUser, error) {
	user, err := s.repo.FindUserByUsername(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, param SignupRequest) (string, error) {
	id, _ := gonanoid.New(20)

	//hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(param.Password), 10)
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
