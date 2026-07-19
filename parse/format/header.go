package format

import (
	"fmt"

	"github.com/firefish111/way2fa/internal/config"
)

// file constants
const (
	MagicNumber          = "WAY\xff"
	FormatVersion uint16 = 1
)

// FileFormat constants
const (
	Csv uint8 = iota

	// special value, for comparison check
	_MaxFileType
)

// Flags
const (
	HasPassword uint8 = 1 << iota
)

type Header struct {
	// magic number
	Magic [4]byte
	// version: incremental, starting with version 1.
	// is updated every time something changes, e.g. new formats, flags, reserved space is filled, etc
	Version uint16

	// see constant definitions
	FileType uint8
	Flags    uint8

	// when the payload was last modified, in Unix seconds
	LastModified int64
	// the size of the output payload, after this header
	PayloadSize uint64

	// packed capabilities struct
	Capabilities config.DerivationCapabilitiesPacked

	// reserved
	_ uint16

	// IV and salt, random numbers generated at encryption time, needed for decryption.
	// if Flags.HasPassword is unset, these are both all nulls.
	AesIv      [12]byte
	AesKeySalt [16]byte
}

func (h Header) Validate() error {
	if string(h.Magic[:]) != MagicNumber {
		return headerError("magic number invalid")
	}

	if h.Version > FormatVersion {
		return headerError(fmt.Sprintf("unknown format version %d, expecting at most %d", h.Version, FormatVersion))
	}

	if h.FileType >= _MaxFileType {
		return headerError(fmt.Sprintf("unknown format type #%d", h.FileType))
	}

	/* XXX: extend when needed */

	return nil
}

/*
// extracts capabilities and AES information from header
// returns an error if the header is invalid
func (h Header) Extract() (config.DerivationCapabilities, cryptor.AesCryptor, error) {
	err := h.Validate()
	if err != nil {
		return config.DerivationCapabilities{}, cryptor.AesCryptor{}, err
	}

  // FIX: figure out what to do with the salt
	caps := h.Capabilities.Unpack()
	aes := cryptor.AesCryptor{Key: []byte(), Iv: h.AesIv[:]}
	return caps, aes, nil
}
*/
