// Package auth ports the legacy JWT / password / OTP / reset-token utilities.
// The JWT claim structure ({id, type, exp, iat}, HS256) is kept byte-compatible
// with the Python implementation so existing tokens keep working after cutover.
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Manager issues and verifies the platform's JWTs.
type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
	resetSecret   []byte
}

// New builds a Manager from the configured secrets.
func New(accessSecret, refreshSecret, resetSecret string) *Manager {
	return &Manager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		resetSecret:   []byte(resetSecret),
	}
}

func sign(secret []byte, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func parse(secret []byte, tokenStr string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// GenerateAccessToken mirrors generate_access_token (exp: +1 day).
func (m *Manager) GenerateAccessToken(userID, userType string) (string, error) {
	now := time.Now().UTC()
	return sign(m.accessSecret, jwt.MapClaims{
		"id":   userID,
		"type": userType,
		"exp":  now.Add(24 * time.Hour).Unix(),
		"iat":  now.Unix(),
	})
}

// GenerateRefreshToken mirrors generate_refresh_token (exp: +20 days).
func (m *Manager) GenerateRefreshToken(userID, userType string) (string, error) {
	now := time.Now().UTC()
	return sign(m.refreshSecret, jwt.MapClaims{
		"id":   userID,
		"type": userType,
		"exp":  now.Add(20 * 24 * time.Hour).Unix(),
		"iat":  now.Unix(),
	})
}

// VerifyAccessToken returns the claims or nil on any failure (matches Python).
func (m *Manager) VerifyAccessToken(token string) jwt.MapClaims {
	claims, err := parse(m.accessSecret, token)
	if err != nil {
		return nil
	}
	return claims
}

// VerifyRefreshToken returns the claims or nil on any failure.
func (m *Manager) VerifyRefreshToken(token string) jwt.MapClaims {
	claims, err := parse(m.refreshSecret, token)
	if err != nil {
		return nil
	}
	return claims
}
