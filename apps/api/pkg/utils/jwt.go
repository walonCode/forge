package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

var secretKey = GetEnvString("JWT_SECRET", "e8d5cbe9c3a1c605d61d21908054ebb54d095d3ee28c5a4c7b4ac1620dc6fd119836e85e54f9a5a2e9f70cd1e42fb30bebefb4d1096041a475555723b191fd5a")

func CreateToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": id,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

type JwtClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func VerifyToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(secretKey), nil
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
