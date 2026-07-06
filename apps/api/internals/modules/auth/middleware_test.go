package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"api/pkg/utils"
)

const testSecret = "test-secret-value"

func TestAuthMiddleware(t *testing.T) {
	validToken, err := utils.CreateToken("user-1", testSecret, utils.AccessToken, utils.AccessTokenTTL)
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	refreshToken, err := utils.CreateToken("user-1", testSecret, utils.RefreshToken, utils.RefreshTokenTTL)
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}

	tests := []struct {
		name       string
		header     string
		wantStatus int
		wantNext   bool
	}{
		{"valid token", "Bearer " + validToken, http.StatusOK, true},
		{"missing header", "", http.StatusUnauthorized, false},
		{"malformed header", validToken, http.StatusUnauthorized, false},
		{"bad token", "Bearer not-a-real-token", http.StatusUnauthorized, false},
		{"refresh token rejected", "Bearer " + refreshToken, http.StatusUnauthorized, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nextCalled := false
			var gotUserID any
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				gotUserID = r.Context().Value(UserIdKey)
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			rr := httptest.NewRecorder()

			AuthMiddleware(testSecret)(next).ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rr.Code, tc.wantStatus)
			}
			if nextCalled != tc.wantNext {
				t.Fatalf("next called = %v, want %v", nextCalled, tc.wantNext)
			}
			if tc.wantNext && gotUserID != "user-1" {
				t.Fatalf("context user id = %v, want user-1", gotUserID)
			}
		})
	}
}
