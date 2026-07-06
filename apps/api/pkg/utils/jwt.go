package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType distinguishes short-lived access tokens from long-lived refresh tokens.
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"

	// AccessTokenTTL is the lifetime of a short-lived access token.
	AccessTokenTTL = 24 * time.Hour
	// RefreshTokenTTL is the lifetime of a long-lived refresh token.
	RefreshTokenTTL = 7 * 24 * time.Hour
)

type JwtClaims struct {
	UserId string `json:"userId"`
	Type   string `json:"typ"`
	jwt.RegisteredClaims
}

// CreateToken signs a JWT for the given user id and token type using secret,
// expiring after ttl.
func CreateToken(id, secret string, tokenType TokenType, ttl time.Duration) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("jwt secret is empty")
	}

	claims := JwtClaims{
		UserId: id,
		Type:   string(tokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// VerifyToken parses and validates a signed token string using secret.
func VerifyToken(tokenString, secret string) (*JwtClaims, error) {
	if secret == "" {
		return nil, fmt.Errorf("jwt secret is empty")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtClaims{},
		func(token *jwt.Token) (any, error) {
			// reject anything that is not HMAC to prevent alg-confusion attacks
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
