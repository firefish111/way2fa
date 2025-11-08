package parse

import (
	"fmt"

	"github.com/kirsle/configdir"
)

const (
	ConfigDirName string = "way2fa"
)

var ConfPath string = configdir.LocalConfig(ConfigDirName)

// Thunk that creates ConfPath for us
func InitConfPath() error {
	if err := configdir.MakePath(ConfPath); err != nil {
		return fmt.Errorf("cannot create config directory %s: %w", ConfPath, err)
	}
	return nil
}
