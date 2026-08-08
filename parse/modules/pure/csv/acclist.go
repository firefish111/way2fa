package csv_pure

import (
	"path/filepath"

	"github.com/firefish111/way2fa/account"
	"github.com/firefish111/way2fa/cryptor"
	"github.com/firefish111/way2fa/parse"
	"github.com/gocarina/gocsv"
)

// Implementation of AccountList interface

func (c *CsvPure) GetAccs() ([]account.Account, error) {
	if !c.IsDecrypted() {
		return nil, cryptor.NotDecrypted("get Accounts")
	}
	defer c.Recrypt() // recrypt at end. defer is filo stack, so this is last thing

	var out []account.Account

	// empty guard
	if len(c.buffer) == 0 {
		return out, nil
	}

	if err := gocsv.UnmarshalString(c.buffer, &out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *CsvPure) SetAccs(to_set []account.Account) error {
	if !c.IsDecrypted() {
		return cryptor.NotDecrypted("set Accounts")
	}
	defer c.Recrypt() // recrypt at end. defer is filo stack, so this is last thing

	buf, err := gocsv.MarshalString(&to_set)
	if err != nil {
		return err
	}

	c.buffer = buf
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
