package parse

import (
	"fmt"
)

// Custom error type.
//
// wasDecrypted is whether the error was decrypted or not
// private field action contains what was attempted
// public field IsFaultOfPassword signifies whether this was the password's fault, or some other issue
type DecryptError struct {
	wasDecrypted      bool
	action            string
	IsFaultOfPassword bool
}

func (e DecryptError) Error() string {
	switch e.wasDecrypted {
	case false:
		return fmt.Sprintf("DecryptError: could not %s, was still encrypted", e.action)
	case true:
		fallthrough // TODO
	default:
		return fmt.Sprintf("DecryptError: unknown cause, given %s", e.action)
	}
}

func NotDecrypted(act string) DecryptError {
	return DecryptError{
		wasDecrypted:      false,
		action:            act,
		IsFaultOfPassword: true, // this is a password problem
	}
}
