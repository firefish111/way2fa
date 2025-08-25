package account

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
)

// go treats a simple typedef as a different type:
// e.g. `type MyInt int` creates my own int, but it is not interchangeable with the real int and must be explicitly casted.
// additionally, i can actually implement methods for, and embed my phony int into other stucts
type RawKey []byte

// All fields must be public because of the marshaller
type Account struct {
	Name     string `csv:"name"`
	AcctId   string `csv:"acc_id"`
	Interval *uint  `csv:"interv"`
	Key      RawKey `csv:"key"`
}

// Generates an RFC 6238 compliant TOTP
func (a Account) GenKey(nth uint64) (uint32, error) {
	hasher := hmac.New(sha1.New, a.Key)
	err := binary.Write(hasher, binary.BigEndian, nth)
	if err != nil {
		return 1_000_000, err // 1 million is smallest invalid value
	}

	// sum is some sort of hash concatenation function (idk), but is also the only way to extract the digest
	total := hasher.Sum(nil)
	ix := total[19] & 0xf

	return (binary.BigEndian.Uint32(total[ix:ix+4]) & 0x7f_ff_ff_ff) % 1_000_000, nil
}

func decodeTextKey(text_key string) (RawKey, error) {
	return base32.StdEncoding.DecodeString(text_key)
}

func NewFromTextKey(text_key string) (*Account, error) {
	key, err := decodeTextKey(text_key)

	if err != nil {
		return nil, err
	}

	// new is used for custom structs, make is used for things like slices that need preallocation for length
	a := new(Account)
	a.Key = key

	return a, nil
}

// Interval default checking thunk. Account.Interval must still be public though, as the marshaller requires all fields be public
func (a Account) GetInterval() uint {
	if a.Interval == nil || *a.Interval == 0 {
		return 30
	} else {
		return *a.Interval
	}
}
