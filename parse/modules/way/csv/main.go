// Way format of csv.
// Much more secure than pure/csv, due to encryption
// Stored as follows:
//
// if password protected:
//
//	Csv (as in CsvPure) -> AES-256 encrypt with GCM
//
// otherwise:
//
//	Csv (as in CsvPure) -> Base64 encode
package csv_way

import (
	"path/filepath"

	"github.com/firefish111/way2fa/cryptor"
	"github.com/firefish111/way2fa/format"
	"github.com/firefish111/way2fa/internal/config"
)

// holds the header, the cryptor, and the raw payload held by the file.
// holds the filename exclusively for reading/writing to file.
type CsvWay struct {
	path string

	header  format.Header
	crypt   cryptor.AesCryptor
	payload []byte

	isDefaultStore bool
	hasPassword    bool
	capabilities   config.DerivationCapabilities
}

func (c *CsvWay) PrepopulateFromFile(path string, isDefault bool) error {
	c.isDefaultStore = isDefault
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	c.path = p

	return nil
}

func (c *CsvWay) GetWayTypeId() format.FileTypeId {
	return format.Csv
}

// Used for creating new files
func (c *CsvWay) PopulateNew() error {
	// set cryptor
	var err error
	c.crypt, err = cryptor.RandomisedAes()
	if err != nil {
		return err
	}

	// set capabilities
	c.capabilities = config.GetCurrentCapabilities()

	// make header sensible
	c.updateHeader()

	return nil
}

func (c *CsvWay) SetPasswordProtected(isPasswordProtected bool) {
	c.hasPassword = isPasswordProtected
}
