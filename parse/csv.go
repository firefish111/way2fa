package parse

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/kirsle/configdir"
	"github.com/gocarina/gocsv"
	"github.com/firefish111/way2fa/account"
)

const (
	CsvFilename = "keys.csv"
)

// implements AccountList
type CsvParser struct {
	path string
	isDefaultStore bool
	password *string
}

func GetDefaultCsv() (*CsvParser, error) {
	cpath := configdir.LocalConfig(ConfigDirName)
	// force directory to exist
	if err := configdir.MakePath(cpath); err != nil {
		return nil, fmt.Errorf("cannot create config directory %s: %w", cpath, err)
	}

	return &CsvParser {
		path: filepath.Join(
			cpath,
			CsvFilename,
		),
		isDefaultStore: false,
		password: nil,
	}, nil
}

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

func (c CsvParser) GetSource() (DataSource, string) {
        return 0, ""
}

func (c CsvParser) WriteAccs(to_write []account.Account) { }
