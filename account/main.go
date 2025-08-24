package account

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
)

type RawKey []byte

// All fields must be public because of the marshaller
type Account struct {
	Name     string `csv:"name"`
	AcctId   string `csv:"acc_id"`
	Interval *uint  `csv:"interv"`
	Key      RawKey `csv:"key"`
}

func (a Account) GenKey(nth uint64) (uint32, error) {
	hasher := hmac.New(sha1.New, a.Key)
	err := binary.Write(hasher, binary.BigEndian, nth)
	if err != nil {
		return 0, err
	}

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

	a := new(Account)
	a.Key = key

	return a, nil
}

func (a Account) GetInterval() uint {
	if a.Interval == nil || *a.Interval == 0 {
		return 30
	} else {
		return *a.Interval
	}
}
