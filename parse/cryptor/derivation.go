package cryptor

import (
	"crypto/rand"

	"github.com/firefish111/way2fa/internal/config"
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

func deriveKey(password PasswordHash, salt []byte, params config.DerivationCapabilities, size_bytes uint32) []byte {
	key := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.MemKiB,
		params.Threads,
		size_bytes,
	)

	return key
}
