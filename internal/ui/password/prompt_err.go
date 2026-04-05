package password

import (
	"fmt"
)

// Custom error type.
type (
	PromptErrorType uint
	PromptError     struct {
		PromptErrorType
		passlen int
	}
)

const (
	NotMatch PromptErrorType = iota
	OutOfTries
	TooLong
)

func (e PromptError) Error() string {
	switch e.PromptErrorType {
	case NotMatch:
		return "PromptError: Passwords did not match"
	case OutOfTries:
		return fmt.Sprintf("PromptError: Out of tries (used all %d)", PasswordTriesCount)
	case TooLong:
		return fmt.Sprintf(
			"PromptError: Password was %d bytes long (exceeded max of %d)",
			e.passlen,
			PasswordMaxLen,
		)
	default:
		return fmt.Sprintf("PromptError: unknown cause, given %d", e)
	}
}

func (t PromptErrorType) ToError() PromptError {
	return PromptError{PromptErrorType: t}
}
