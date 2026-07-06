package auth

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
)

// otpExpiry mirrors OTP_EXPIRY_SECONDS (10 minutes).
const otpExpiry = 10 * time.Minute

// OTPStore abstracts OTP persistence. The in-memory implementation is fine for
// a single replica; under autoscaling, back this with Redis so OTPs are shared
// across replicas (the interface keeps that swap a one-line change).
type OTPStore interface {
	Generate(email string) string
	Verify(email, otp string) bool
}

type otpEntry struct {
	otp    string
	expiry time.Time
}

// MemoryOTPStore is a thread-safe in-process OTP store.
type MemoryOTPStore struct {
	mu    sync.Mutex
	store map[string]otpEntry
}

// NewMemoryOTPStore creates an empty in-memory OTP store.
func NewMemoryOTPStore() *MemoryOTPStore {
	return &MemoryOTPStore{store: make(map[string]otpEntry)}
}

func (s *MemoryOTPStore) cleanup(now time.Time) {
	for k, v := range s.store {
		if v.expiry.Before(now) {
			delete(s.store, k)
		}
	}
}

// Generate creates and stores a 6-digit OTP for the email.
func (s *MemoryOTPStore) Generate(email string) string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	otp := strconv.FormatInt(n.Int64()+100000, 10)

	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.cleanup(now)
	s.store[strings.ToLower(email)] = otpEntry{otp: otp, expiry: now.Add(otpExpiry)}
	return otp
}

// Verify checks and consumes an OTP for the email.
func (s *MemoryOTPStore) Verify(email, otp string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.cleanup(now)
	key := strings.ToLower(email)
	entry, ok := s.store[key]
	if !ok {
		return false
	}
	if entry.otp == otp && entry.expiry.After(now) {
		delete(s.store, key)
		return true
	}
	return false
}
