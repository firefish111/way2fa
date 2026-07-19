package cryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/firefish111/way2fa/internal/config"
)

// Wrapper for key and iv
type AesCryptor struct {
	Key []byte
	Iv  []byte
}

func padToBlocksize(plaintext []byte) []byte {
	// the number of bytes that we need to pad
	bytesLeft := aes.BlockSize - (len(plaintext) % aes.BlockSize)

	// our padding as an array of bytes
	to_pad := bytes.Repeat([]byte{0}, bytesLeft)

	// return the padded array
	return append(plaintext, to_pad...)
}

func (a AesCryptor) EncryptAes(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	padded := padToBlocksize(plaintext)

	// Seal is a bit strange, in that it takes a destination pointer, but we don't need that
	// TODO: this also takes an "authenticated additional data" pointer, which we could possibly use to authenticate the header?
	ciphertext := gcm.Seal(nil, a.Iv, padded, nil)
	return ciphertext, nil
}

func (a AesCryptor) DecryptAes(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// TODO: see earlier about authenticating header too
	plaintext, err := gcm.Open(nil, a.Iv, ciphertext, nil)
	if err != nil {
		// NOTE: though it's likely an error to do with the password, this doesn't eliminate
		// the possibility that the header has been tampered with, as the GCM doesn't distinguish.
		// either way, file is unopenable
		return nil, fmt.Errorf("%w; %w", FailedAuthentication(), err)
	}

	// undo pad with nulls
	unpadded := bytes.TrimRight(plaintext, "\000")

	return unpadded, nil
}

func MakeAes(passhash PasswordHash, withCapabilities *config.DerivationCapabilities) (AesCryptor, error) {
	// generate derivation salt and iv, using predetermined sizes defined in ./main.go, q.v.
	salt, err := generateNonce(Argon2SaltSize)
	if err != nil {
		// structs can't be nil, just empty
		return AesCryptor{}, fmt.Errorf("%w; %w", CouldNotGenerateNonce("Argon2 salt"), err)
	}

	iv, err := generateNonce(AesIvSize)
	if err != nil {
		return AesCryptor{}, fmt.Errorf("%w; %w", CouldNotGenerateNonce("AES IV"), err)
	}

	// if we haven't set capabilities, we get what our computer has
	if withCapabilities == nil {
		capabilities := config.GetCurrentCapabilities() // pointer foolery, as reference can't exist if variable is not on stack
		withCapabilities = &capabilities
	}

	derived := deriveKey(passhash, salt, *withCapabilities, AesKeySize)

	return AesCryptor{
		Key: derived,
		Iv:  iv,
	}, nil
}
