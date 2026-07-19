package cryptor

import (
	"fmt"
)

type (
	cryptStatus uint
	// Custom error type.
	//
	// didDecrypt is whether the error was decrypted or not
	// private field action contains what was attempted
	// public field IsFaultOfPassword signifies whether this was the password's fault, or some other issue
	CryptError struct {
		IsFaultOfPassword bool
		status            cryptStatus
		verb              string
	}
)

const (
	failedAuthentication cryptStatus = iota
	failedActionWhilstEncrypted
	failedToGenerateNonce
)

func (e CryptError) Error() string {
	switch e.status {
	case failedAuthentication:
		return "CryptError: failed to authenticate, password incorrect or file has been tampered with"
	case failedActionWhilstEncrypted:
		return fmt.Sprintf("CryptError: could not %s, was still encrypted", e.verb)
	case failedToGenerateNonce:
		return fmt.Sprintf("CryptError: could not generate nonce for %s", e.verb)
	default:
		return fmt.Sprintf("CryptError: unknown cause, given %s", e.verb)
		// fallthrough // TODO
	}
}

func FailedAuthentication() CryptError {
	return CryptError{
		IsFaultOfPassword: true, // perhaps due to password, we can't know for sure
		status:            failedAuthentication,
		verb:              "", // verb is not used, so empty string will suffice
	}
}

func NotDecrypted(action string) CryptError {
	return CryptError{
		IsFaultOfPassword: true, // this is a password problem
		status:            failedActionWhilstEncrypted,
		verb:              action,
	}
}

func CouldNotGenerateNonce(object string) CryptError {
	return CryptError{
		IsFaultOfPassword: false, // this is not a password problem
		status:            failedToGenerateNonce,
		verb:              object,
	}
}
