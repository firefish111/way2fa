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
	"github.com/firefish111/way2fa/internal/config"
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
}

func (c *CsvPure) PrepopulateDefault() error {
	c.isDefaultStore = true // is default store

	c.path = filepath.Join(
		config.ConfPath,
		CsvPureFilename,
	)

	return nil // no way to error
}

func (c *CsvPure) PrepopulateFromFile(path string) error {
	c.isDefaultStore = false // is not default store
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	c.path = p

	return nil
}
