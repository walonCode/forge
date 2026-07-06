package users

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api/internals/modules/auth"

	"golang.org/x/crypto/bcrypt"
)

func testHandler(r Repository) *Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return newHandler(newService(r), logger)
}

// withUser attaches an authenticated user id the way AuthMiddleware would.
func withUser(req *http.Request, userID string) *http.Request {
	ctx := context.WithValue(req.Context(), auth.UserIdKey, userID)
	return req.WithContext(ctx)
}

func TestGetProfileHandler(t *testing.T) {
	t.Run("200 and no password leaked", func(t *testing.T) {
		h := testHandler(&fakeRepo{profile: &UserProfile{ID: "u1", Name: "Ada", Username: "ada"}})
		rr := httptest.NewRecorder()
		req := withUser(httptest.NewRequest(http.MethodGet, "/user/profile", nil), "u1")

		h.GetProfile(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
		if strings.Contains(strings.ToLower(rr.Body.String()), "password") {
			t.Fatalf("response leaked password field: %s", rr.Body.String())
		}
	})

	t.Run("404 when user missing", func(t *testing.T) {
		h := testHandler(&fakeRepo{profileErr: sql.ErrNoRows})
		rr := httptest.NewRecorder()
		req := withUser(httptest.NewRequest(http.MethodGet, "/user/profile", nil), "u1")

		h.GetProfile(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rr.Code)
		}
	})
}

func TestUpdateProfileHandler(t *testing.T) {
	t.Run("400 on empty update", func(t *testing.T) {
		h := testHandler(&fakeRepo{profile: &UserProfile{ID: "u1", Name: "Ada", Username: "ada"}})
		rr := httptest.NewRecorder()
		req := withUser(httptest.NewRequest(http.MethodPatch, "/user/profile", strings.NewReader(`{}`)), "u1")

		h.UpdateProfile(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("409 when username taken", func(t *testing.T) {
		h := testHandler(&fakeRepo{
			profile:       &UserProfile{ID: "u1", Name: "Ada", Username: "ada"},
			usernameTaken: true,
		})
		rr := httptest.NewRecorder()
		req := withUser(httptest.NewRequest(http.MethodPatch, "/user/profile", strings.NewReader(`{"username":"taken"}`)), "u1")

		h.UpdateProfile(rr, req)

		if rr.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", rr.Code)
		}
	})
}

func TestUpdatePasswordHandler(t *testing.T) {
	t.Run("400 on short new password", func(t *testing.T) {
		h := testHandler(&fakeRepo{})
		rr := httptest.NewRecorder()
		req := withUser(httptest.NewRequest(http.MethodPatch, "/user/password", strings.NewReader(`{"current_password":"whatever1","new_password":"short"}`)), "u1")

		h.UpdatePassword(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("400 when current password wrong", func(t *testing.T) {
		hash, _ := bcrypt.GenerateFromPassword([]byte("correct-horse"), bcrypt.MinCost)
		h := testHandler(&fakeRepo{passwordHash: string(hash)})
		rr := httptest.NewRecorder()
		req := withUser(httptest.NewRequest(http.MethodPatch, "/user/password", strings.NewReader(`{"current_password":"wrong-password","new_password":"new-password"}`)), "u1")

		h.UpdatePassword(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})
}

func TestDeleteAccountHandler(t *testing.T) {
	repo := &fakeRepo{}
	h := testHandler(repo)
	rr := httptest.NewRecorder()
	req := withUser(httptest.NewRequest(http.MethodDelete, "/user", nil), "u1")

	h.DeleteAccount(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if !repo.deleteCalled {
		t.Fatal("expected account to be deleted")
	}
}
