// Package auth provides HS256 token signing and verification.
// It is intentionally free of business dependencies so higher layers can
// wrap it with a token store and user/role resolution.
package auth

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT payload for Halo sessions.
type Claims struct {
	jwt.RegisteredClaims
	UserUUID string `json:"uuid"`
}

// JWT issues and validates HS256 tokens.
type JWT struct {
	issuer   string
	secret   string
	duration time.Duration
}

// New builds a JWT handler. An empty secret preserves the historical
// behaviour (tokens were signed with an empty key) until wiring is revisited.
func New(issuer, secret string, duration time.Duration) *JWT {
	return &JWT{issuer: issuer, secret: secret, duration: duration}
}

// NewDefault returns a JWT handler with empty issuer/secret and a 24h duration,
// matching the previous tools.New() singleton behaviour.
func NewDefault() *JWT {
	return New("", "", time.Hour*24)
}

// CreateToken signs a token for the given user UUID.
func (j *JWT) CreateToken(userUUID string) (string, error) {
	claims := &Claims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.duration)),
			Issuer:    j.issuer,
		},
		userUUID,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.secret))
}

// ParseToken validates the signature and expiry, returning the user UUID.
func (j *JWT) ParseToken(input string) (string, error) {
	token, err := jwt.ParseWithClaims(input, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	return claims.UserUUID, nil
}
