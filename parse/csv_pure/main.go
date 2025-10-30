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
	"fmt"
	"github.com/firefish111/way2fa/parse"
	"github.com/kirsle/configdir"
	"path/filepath"
)

const (
	CsvPureExt = ".csv"
	CsvPureFilename = "keys" + CsvPureExt
)

// implements AccountList
type CsvPure struct {
	path           string
	isDefaultStore bool
}

func getDefaultCsv() (*CsvPure, error) {
	cpath := configdir.LocalConfig(parse.ConfigDirName)
	// force directory to exist
	if err := configdir.MakePath(cpath); err != nil {
		return nil, fmt.Errorf("cannot create config directory %s: %w", cpath, err)
	}

	return &CsvPure{
		path: filepath.Join(
			cpath,
			CsvPureFilename,
		),
		isDefaultStore: true,
	}, nil
}

func getFileByName(path string) (*CsvPure, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return &CsvPure{
		path:           p,
		isDefaultStore: false,
	}, nil
}

func GetFile(name *string) (*CsvPure, error) {
	if name != nil {
		return getFileByName(*name)
	} else {
		return getDefaultCsv()
	}
}
