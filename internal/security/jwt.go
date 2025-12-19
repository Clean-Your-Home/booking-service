package security

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret: []byte(secret)}
}

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (m *JWTManager) Parse(tokenStr string) (*Claims, error) {
	var claims Claims
	t, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	return &claims, nil
}
