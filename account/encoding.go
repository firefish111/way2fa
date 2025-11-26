package account

import (
	"encoding/base32"
)

func DecodeTextKey(text_key string) (RawKey, error) {
	return base32.StdEncoding.
		WithPadding(base32.NoPadding).
		DecodeString(text_key)
}

func EncodeTextKey(key RawKey) string {
	return base32.StdEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(key)
}

func NewFromTextKey(text_key string) (*Account, error) {
	key, err := DecodeTextKey(text_key)

	if err != nil {
		return nil, err
	}

	// new is used for custom structs, make is used for things like slices that need preallocation for length
	a := new(Account)
	a.Key = key

	return a, nil
}
