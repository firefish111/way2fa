package csv_way

import (
	"encoding/base64"
	"fmt"
	"path/filepath"

	"github.com/firefish111/way2fa/account"
	"github.com/firefish111/way2fa/cryptor"
	"github.com/firefish111/way2fa/parse"
	"github.com/gocarina/gocsv"
)

// Implementation of AccountList interface

func (c *CsvWay) GetAccs() ([]account.Account, error) {
	if !c.IsDecrypted() {
		return nil, cryptor.NotDecrypted("get Accounts")
	}
	defer c.Recrypt() // recrypt at end. defer is filo stack, so this is last thing

	// if no payload or payload is nil (check included in len), just return empty array
	if len(c.payload) == 0 {
		return []account.Account{}, nil
	}

	// this is the csv text
	var plaintext []byte
	var err error

	if c.IsPasswordProtected() { // if password protected, decrypt
		plaintext, err = c.crypt.Decrypt(c.payload)
	} else { // otherwise, just base64 decode
		plaintext, err = base64.StdEncoding.DecodeString(string(c.payload))
	}

	if err != nil {
		return nil, fmt.Errorf("GetAccs: decryption failed; %w", err)
	}

	var out []account.Account

	if err := gocsv.UnmarshalBytes(plaintext, &out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *CsvWay) SetAccs(to_set []account.Account) error {
	if !c.IsDecrypted() {
		return cryptor.NotDecrypted("set Accounts")
	}
	defer c.Recrypt() // recrypt at end. defer is filo stack, so this is last thing

	var encrypted []byte
	// error has to be declared up here because := sucks
	var err error

	// marshal to CSV
	plaintext, err := gocsv.MarshalBytes(to_set)
	if err != nil {
		return err
	}

	if c.IsPasswordProtected() {
		encrypted, err = c.crypt.Encrypt([]byte(plaintext))
		if err != nil {
			return err
		}
	} else {
		encrypted = []byte(
			base64.StdEncoding.EncodeToString(
				[]byte(plaintext),
			),
		)
	}

	c.payload = encrypted
	return nil
}

func (c *CsvWay) GetSource() (parse.DataSource, string) {
	if c.isDefaultStore {
		return parse.NamedSource, "<default AES-256 CSV>"
	} else {
		return parse.FileSource, filepath.Base(c.path)
	}
}

func (c *CsvWay) GetSourceFilePath() string {
	return filepath.Clean(c.path)
}
