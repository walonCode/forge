package auth

import (
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func testAuthHandler(r Repository) *Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return newHandler(newService(r), logger, testSecret)
}

func TestLoginHandler(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)

	t.Run("200 with tokens on valid credentials", func(t *testing.T) {
		h := testAuthHandler(&fakeAuthRepo{user: &DbUser{ID: "u1", Username: "ada", Password: string(hash)}})
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"ada","password":"password123"}`))

		h.LoginHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "access_token") {
			t.Fatalf("response missing access_token: %s", rr.Body.String())
		}
	})

	t.Run("401 on wrong password", func(t *testing.T) {
		h := testAuthHandler(&fakeAuthRepo{user: &DbUser{ID: "u1", Username: "ada", Password: string(hash)}})
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"ada","password":"wrong-password"}`))

		h.LoginHandler(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", rr.Code)
		}
	})

	t.Run("401 when user not found", func(t *testing.T) {
		h := testAuthHandler(&fakeAuthRepo{findErr: sql.ErrNoRows})
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"ghost","password":"password123"}`))

		h.LoginHandler(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", rr.Code)
		}
	})
}

func TestSignupHandler(t *testing.T) {
	t.Run("201 on new user", func(t *testing.T) {
		h := testAuthHandler(&fakeAuthRepo{findErr: sql.ErrNoRows, createID: "new-id"})
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(`{"name":"Ada","username":"ada","password":"password123"}`))

		h.SignupHandler(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "access_token") {
			t.Fatalf("response missing access_token: %s", rr.Body.String())
		}
	})

	t.Run("409 when username already exists", func(t *testing.T) {
		// the repo surfaces the DB unique-constraint violation
		h := testAuthHandler(&fakeAuthRepo{createErr: &pgconn.PgError{Code: "23505"}})
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(`{"name":"Ada","username":"ada","password":"password123"}`))

		h.SignupHandler(rr, req)

		if rr.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", rr.Code)
		}
	})

	t.Run("400 on invalid signup payload", func(t *testing.T) {
		h := testAuthHandler(&fakeAuthRepo{})
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(`{"name":"A","username":"a","password":"short"}`))

		h.SignupHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})
}

func TestLogoutHandler(t *testing.T) {
	h := testAuthHandler(&fakeAuthRepo{})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)

	h.LogoutHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
}
