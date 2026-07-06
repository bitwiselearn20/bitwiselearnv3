package auth

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateResetToken mirrors generate_reset_token (exp: +10 minutes).
func (m *Manager) GenerateResetToken(email, userType, userID string) (string, error) {
	now := time.Now().UTC()
	return sign(m.resetSecret, jwt.MapClaims{
		"email": email,
		"type":  userType,
		"id":    userID,
		"exp":   now.Add(10 * time.Minute).Unix(),
		"iat":   now.Unix(),
	})
}

// VerifyResetToken returns claims, or nil if invalid/expired/already used.
func (m *Manager) VerifyResetToken(token string, blocklist TokenBlocklist) jwt.MapClaims {
	if blocklist != nil && blocklist.IsUsed(token) {
		return nil
	}
	claims, err := parse(m.resetSecret, token)
	if err != nil {
		return nil
	}
	return claims
}

// TokenBlocklist tracks consumed reset tokens. Back with Redis (with TTL equal
// to the token lifetime) when running multiple replicas.
type TokenBlocklist interface {
	IsUsed(token string) bool
	Invalidate(token string)
}

// MemoryBlocklist is a thread-safe in-process blocklist.
type MemoryBlocklist struct {
	mu   sync.Mutex
	used map[string]struct{}
}

// NewMemoryBlocklist creates an empty in-memory token blocklist.
func NewMemoryBlocklist() *MemoryBlocklist {
	return &MemoryBlocklist{used: make(map[string]struct{})}
}

// IsUsed reports whether the token was already consumed.
func (b *MemoryBlocklist) IsUsed(token string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, ok := b.used[token]
	return ok
}

// Invalidate marks the token as consumed.
func (b *MemoryBlocklist) Invalidate(token string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.used[token] = struct{}{}
}
