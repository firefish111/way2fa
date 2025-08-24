package csv

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/gocarina/gocsv"
	"github.com/firefish111/way2fa/account"
	"github.com/firefish111/way2fa/parse"
)

// Implementation of AccountList interface

func (c CsvParser) GetAccs() ([]account.Account, error) {
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

func (c CsvParser) GetSource() (parse.DataSource, string) {
	if c.isDefaultStore {
		return parse.NamedSource, "<default store>"
	} else {
		return parse.FileSource, filepath.Base(c.path)
	}
}

func (c CsvParser) WriteAccs(to_write []account.Account) { }
