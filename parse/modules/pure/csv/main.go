// This package decodes files in "pure" csv form, that is a csv file
// with no encryption, encoding or headers.
//
// As is probably guessable, this is horrendously insecure, as all you
// need is an infostealer to empty your AppData, ~/.config, or equivalent
// and all your accounts are compromised.
//
// TODO: make whatever it is warn you when you use this that it is insecure
package csv_pure

import (
	"path/filepath"
)

const (
	CsvPureExt      = ".csv"
	CsvPureFilename = "keys" + CsvPureExt
)

// implements AccountList
type CsvPure struct {
	path           string
	isDefaultStore bool
	buffer         string
}

func (c *CsvPure) PrepopulateFromFile(path string, isDefault bool) error {
	c.isDefaultStore = isDefault
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	c.path = p

	return nil
}

func (c *CsvPure) GetDefaultFilename() string {
	return CsvPureFilename
}
