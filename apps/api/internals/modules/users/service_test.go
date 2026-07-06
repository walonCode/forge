package users

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// fakeRepo is an in-memory Repository for exercising the service layer.
type fakeRepo struct {
	profile       *UserProfile
	profileErr    error
	passwordHash  string
	passwordErr   error
	usernameTaken bool
	existsErr     error
	updateErr     error
	updatePwErr   error
	deleteErr     error

	updatePwCalled bool
	deleteCalled   bool
}

func (f *fakeRepo) FindProfileByID(ctx context.Context, id string) (*UserProfile, error) {
	return f.profile, f.profileErr
}
func (f *fakeRepo) GetPasswordHash(ctx context.Context, id string) (string, error) {
	return f.passwordHash, f.passwordErr
}
func (f *fakeRepo) ExistsByUsernameExcept(ctx context.Context, username, excludeID string) (bool, error) {
	return f.usernameTaken, f.existsErr
}
func (f *fakeRepo) UpdateProfile(ctx context.Context, id, name, username string) (*UserProfile, error) {
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	return &UserProfile{ID: id, Name: name, Username: username}, nil
}
func (f *fakeRepo) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	f.updatePwCalled = true
	return f.updatePwErr
}
func (f *fakeRepo) DeleteUser(ctx context.Context, id string) error {
	f.deleteCalled = true
	return f.deleteErr
}

func TestGetProfile(t *testing.T) {
	t.Run("returns profile", func(t *testing.T) {
		svc := newService(&fakeRepo{profile: &UserProfile{ID: "u1", Name: "Ada", Username: "ada"}})
		got, err := svc.GetProfile(context.Background(), "u1")
		if err != nil {
			t.Fatalf("GetProfile: %v", err)
		}
		if got.Username != "ada" {
			t.Fatalf("username = %q, want ada", got.Username)
		}
	})

	t.Run("maps no rows to ErrUserNotFound", func(t *testing.T) {
		svc := newService(&fakeRepo{profileErr: sql.ErrNoRows})
		if _, err := svc.GetProfile(context.Background(), "u1"); !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("err = %v, want ErrUserNotFound", err)
		}
	})
}

func TestUpdateProfile(t *testing.T) {
	t.Run("no changes returns ErrNoFieldsToUpdate", func(t *testing.T) {
		svc := newService(&fakeRepo{profile: &UserProfile{ID: "u1", Name: "Ada", Username: "ada"}})
		_, err := svc.UpdateProfile(context.Background(), "u1", UpdateProfileRequest{})
		if !errors.Is(err, ErrNoFieldsToUpdate) {
			t.Fatalf("err = %v, want ErrNoFieldsToUpdate", err)
		}
	})

	t.Run("taken username returns ErrUsernameTaken", func(t *testing.T) {
		svc := newService(&fakeRepo{
			profile:       &UserProfile{ID: "u1", Name: "Ada", Username: "ada"},
			usernameTaken: true,
		})
		_, err := svc.UpdateProfile(context.Background(), "u1", UpdateProfileRequest{Username: "taken"})
		if !errors.Is(err, ErrUsernameTaken) {
			t.Fatalf("err = %v, want ErrUsernameTaken", err)
		}
	})

	t.Run("updates name", func(t *testing.T) {
		svc := newService(&fakeRepo{profile: &UserProfile{ID: "u1", Name: "Ada", Username: "ada"}})
		got, err := svc.UpdateProfile(context.Background(), "u1", UpdateProfileRequest{Name: "Ada L."})
		if err != nil {
			t.Fatalf("UpdateProfile: %v", err)
		}
		if got.Name != "Ada L." {
			t.Fatalf("name = %q, want %q", got.Name, "Ada L.")
		}
	})
}

func TestUpdatePassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-horse"), bcrypt.MinCost)

	t.Run("wrong current password returns ErrInvalidPassword", func(t *testing.T) {
		repo := &fakeRepo{passwordHash: string(hash)}
		svc := newService(repo)
		err := svc.UpdatePassword(context.Background(), "u1", UpdatePasswordRequest{
			CurrentPassword: "wrong",
			NewPassword:     "new-password",
		})
		if !errors.Is(err, ErrInvalidPassword) {
			t.Fatalf("err = %v, want ErrInvalidPassword", err)
		}
		if repo.updatePwCalled {
			t.Fatal("UpdatePassword repo call should not happen on wrong current password")
		}
	})

	t.Run("correct current password updates", func(t *testing.T) {
		repo := &fakeRepo{passwordHash: string(hash)}
		svc := newService(repo)
		err := svc.UpdatePassword(context.Background(), "u1", UpdatePasswordRequest{
			CurrentPassword: "correct-horse",
			NewPassword:     "new-password",
		})
		if err != nil {
			t.Fatalf("UpdatePassword: %v", err)
		}
		if !repo.updatePwCalled {
			t.Fatal("expected repo UpdatePassword to be called")
		}
	})
}

func TestDeleteAccount(t *testing.T) {
	repo := &fakeRepo{}
	svc := newService(repo)
	if err := svc.DeleteAccount(context.Background(), "u1"); err != nil {
		t.Fatalf("DeleteAccount: %v", err)
	}
	if !repo.deleteCalled {
		t.Fatal("expected repo DeleteUser to be called")
	}
}
