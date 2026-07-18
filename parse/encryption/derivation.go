package encryption

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

// generate an Initialisation Vector with specified size in bytes, for AES encryption
func generateNonce(size_bytes uint32) ([]byte, error) {
	nonce := make([]byte, size_bytes)

	// read cryptographically secure random number into nonce
	if _, err := rand.Read(nonce); err != nil {
		return nil, err // pass error up
	}

	return nonce, nil
}

func deriveKey(password PasswordHash, salt []byte, params DerivationCapabilities, size_bytes uint32) ([]byte, error) {
	key := argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, size_bytes)
	return key, nil
}
