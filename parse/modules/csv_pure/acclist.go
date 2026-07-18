package csv_pure

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/firefish111/way2fa/account"
	"github.com/firefish111/way2fa/parse"
	"github.com/firefish111/way2fa/parse/cryptor"
	"github.com/gocarina/gocsv"
)

// Implementation of AccountList interface

func (c *CsvPure) GetAccs() ([]account.Account, error) {
	if !c.IsDecrypted() {
		return nil, cryptor.NotDecrypted("get Accounts")
	}
	defer c.Recrypt() // recrypt at end. defer is filo stack, so this is last thing

	f, err := os.Open(c.path)
	if err != nil {
		return nil, fmt.Errorf("cannot access keyfile %s: %w", c.path, err)
	}

	defer f.Close() // wait until end of function to close

	var out []account.Account

	if err := gocsv.UnmarshalFile(f, &out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *CsvPure) WriteAccs(to_write []account.Account) error {
	if !c.IsDecrypted() {
		return cryptor.NotDecrypted("write Accounts")
	}
	defer c.Recrypt() // recrypt at end. defer is filo stack, so this is last thing

	// os.Open only opens readonly
	f, err := os.OpenFile(c.path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("cannot access keyfile %s: %w", c.path, err)
	}

	defer f.Close() // wait until end of function to close

	if err := gocsv.MarshalFile(&to_write, f); err != nil {
		return err
	}

	return nil
}

func (c *CsvPure) GetSource() (parse.DataSource, string) {
	if c.isDefaultStore {
		return parse.NamedSource, "<default unencrypted CSV>"
	} else {
		return parse.FileSource, filepath.Base(c.path)
	}
}

func (c *CsvPure) GetSourceFilePath() string {
	return filepath.Clean(c.path)
}
