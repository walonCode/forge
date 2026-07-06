package utils

import (
	"testing"
	"time"
)

const testSecret = "test-secret-value"

func TestCreateAndVerifyToken(t *testing.T) {
	token, err := CreateToken("user-123", testSecret, AccessToken, AccessTokenTTL)
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}

	claims, err := VerifyToken(token, testSecret)
	if err != nil {
		t.Fatalf("VerifyToken: %v", err)
	}
	if claims.UserId != "user-123" {
		t.Fatalf("got user id %q, want %q", claims.UserId, "user-123")
	}
	if claims.Type != string(AccessToken) {
		t.Fatalf("got type %q, want %q", claims.Type, AccessToken)
	}
}

func TestTokenTypeIsPreserved(t *testing.T) {
	token, err := CreateToken("user-123", testSecret, RefreshToken, RefreshTokenTTL)
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	claims, err := VerifyToken(token, testSecret)
	if err != nil {
		t.Fatalf("VerifyToken: %v", err)
	}
	if claims.Type != string(RefreshToken) {
		t.Fatalf("got type %q, want %q", claims.Type, RefreshToken)
	}
}

func TestVerifyToken_WrongSecret(t *testing.T) {
	token, err := CreateToken("user-123", testSecret, AccessToken, AccessTokenTTL)
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	if _, err := VerifyToken(token, "different-secret"); err == nil {
		t.Fatal("expected error for wrong secret, got nil")
	}
}

func TestVerifyToken_Expired(t *testing.T) {
	token, err := CreateToken("user-123", testSecret, AccessToken, -time.Hour)
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	if _, err := VerifyToken(token, testSecret); err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestVerifyToken_Tampered(t *testing.T) {
	token, err := CreateToken("user-123", testSecret, AccessToken, AccessTokenTTL)
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	if _, err := VerifyToken(token+"tampered", testSecret); err == nil {
		t.Fatal("expected error for tampered token, got nil")
	}
}

func TestCreateToken_EmptySecret(t *testing.T) {
	if _, err := CreateToken("user-123", "", AccessToken, AccessTokenTTL); err == nil {
		t.Fatal("expected error for empty secret, got nil")
	}
}

func TestVerifyToken_EmptySecret(t *testing.T) {
	if _, err := VerifyToken("some.token.value", ""); err == nil {
		t.Fatal("expected error for empty secret, got nil")
	}
}

func TestAccessTokenShorterThanRefresh(t *testing.T) {
	if AccessTokenTTL >= RefreshTokenTTL {
		t.Fatalf("access ttl (%s) should be shorter than refresh ttl (%s)", AccessTokenTTL, RefreshTokenTTL)
	}
}
