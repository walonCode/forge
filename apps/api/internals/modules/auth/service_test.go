package auth

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// fakeAuthRepo is an in-memory Repository for the auth service/handler tests.
type fakeAuthRepo struct {
	user      *DbUser
	findErr   error
	createID  string
	createErr error

	created *CreateDbUser
}

func (f *fakeAuthRepo) CreateUser(ctx context.Context, params CreateDbUser) (string, error) {
	f.created = &params
	if f.createErr != nil {
		return "", f.createErr
	}
	return f.createID, nil
}

func (f *fakeAuthRepo) FindUserByUsername(ctx context.Context, username string) (*DbUser, error) {
	return f.user, f.findErr
}

func TestServiceCreateUser_HashesPassword(t *testing.T) {
	repo := &fakeAuthRepo{createID: "new-id"}
	svc := newService(repo)

	id, err := svc.CreateUser(context.Background(), SignupRequest{
		Name:     "Ada",
		Username: "ada",
		Password: "plaintext-password",
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if id != "new-id" {
		t.Fatalf("id = %q, want new-id", id)
	}
	if repo.created == nil {
		t.Fatal("repo CreateUser was not called")
	}
	if repo.created.Password == "plaintext-password" {
		t.Fatal("password was stored in plaintext")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repo.created.Password), []byte("plaintext-password")); err != nil {
		t.Fatalf("stored password is not a valid bcrypt hash of the input: %v", err)
	}
	if repo.created.ID == "" {
		t.Fatal("expected a generated id")
	}
}

func TestServiceFindUserByUsername(t *testing.T) {
	t.Run("returns user", func(t *testing.T) {
		svc := newService(&fakeAuthRepo{user: &DbUser{ID: "u1", Username: "ada"}})
		got, err := svc.FindUserByUsername(context.Background(), "ada")
		if err != nil {
			t.Fatalf("FindUserByUsername: %v", err)
		}
		if got.ID != "u1" {
			t.Fatalf("id = %q, want u1", got.ID)
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		svc := newService(&fakeAuthRepo{findErr: errors.New("boom")})
		if _, err := svc.FindUserByUsername(context.Background(), "ada"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
