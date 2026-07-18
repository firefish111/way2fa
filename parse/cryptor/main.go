package cryptor

import (
	"slices"

	"crypto/sha256"
)

type (
	PasswordHash []byte
)

// key size of 32 bytes means AES-256.
// IV size is typically 12 bytes for AES-GCM
const (
	AesKeySize     = 32
	AesIvSize      = 12
	Argon2SaltSize = 16
)

// abstracted out equality function, as we can't assume it's a []byte
func (h PasswordHash) Matches(rhs PasswordHash) bool {
	return slices.Equal(h, rhs)
}

// Hash the password. No salting just yet, that is done as part of key derivation
func HashPassword(key string) PasswordHash {
	hash := sha256.Sum256([]byte(key)) // does not depend on AesKeySize either

	return PasswordHash(hash[:])
}
