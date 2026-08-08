package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/firefish111/way2fa/detector"
	"github.com/firefish111/way2fa/internal/config"
	"github.com/firefish111/way2fa/parse"
)

// Fail if the file exists, for creation routines
func FailOnExist(name string) error {
	_, err := os.Stat(name)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return fmt.Errorf("Cannot check existence of file; %w", err)
	} else {
		return fmt.Errorf("File %s already exists", name)
	}
}

// returns whether the name is a default or not and an error
func (m *managerModel) matchFilename() (fname string, isDefault bool, err error) {
	// deal with filename

	if m.filename != nil && *m.filename != "" {
		isDefault = false
		fname, err = filepath.Abs(*m.filename)
		if err != nil {
			return fname, false, err
		}
	} else if pure, ok := m.possibilities[m.selected].(parse.PureAccountList); ok {
		isDefault = true
		fname = filepath.Join(
			config.ConfPath,
			pure.GetDefaultFilename(),
		)
	} else if _, ok := m.possibilities[m.selected].(parse.WayAccountList); ok {
		isDefault = true
		fname = filepath.Join(
			config.ConfPath,
			detector.DefaultWayFilename,
		)

		m.filename = &fname
	} else {
		return fname, isDefault, fmt.Errorf("Can't match account list to file")
	}

	if err := FailOnExist(fname); err != nil {
		return fname, isDefault, err
	}

	return
}
