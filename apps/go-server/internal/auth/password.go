package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plaintext password with bcrypt.
// Go's bcrypt produces/verifies $2a$/$2b$ hashes and can verify the
// passlib-generated hashes already stored in the database.
func HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword reports whether plain matches the stored bcrypt hash.
func VerifyPassword(plain, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
