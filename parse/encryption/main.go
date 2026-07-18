package encryption

import (
	"slices"
)

type (
	PasswordHash []byte
)

// abstracted out equality function, as we can't assume it's a []byte
func (h PasswordHash) Matches(rhs PasswordHash) bool {
	return slices.Equal(h, rhs)
}

// TODO: just a stub
func HashPassword(key string) PasswordHash {
	return PasswordHash(key)
}
