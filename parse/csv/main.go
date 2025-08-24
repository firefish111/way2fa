package csv

import (
	"fmt"
	"github.com/firefish111/way2fa/parse"
	"github.com/kirsle/configdir"
	"path/filepath"
)

const (
	CsvFilename = "keys.csv"
)

// implements AccountList
type CsvParser struct {
	path           string
	isDefaultStore bool
	password       *string
}

func getDefaultCsv() (*CsvParser, error) {
	cpath := configdir.LocalConfig(parse.ConfigDirName)
	// force directory to exist
	if err := configdir.MakePath(cpath); err != nil {
		return nil, fmt.Errorf("cannot create config directory %s: %w", cpath, err)
	}

	return &CsvParser{
		path: filepath.Join(
			cpath,
			CsvFilename,
		),
		isDefaultStore: true,
		password:       nil,
	}, nil
}

func getFileByName(path string) (*CsvParser, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return &CsvParser{
		path:           p,
		isDefaultStore: false,
		password:       nil,
	}, nil
}

func GetFile(name *string) (*CsvParser, error) {
	if name != nil {
		return getFileByName(*name)
	} else {
		return getDefaultCsv()
	}
}
