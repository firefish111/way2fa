package parse

import (
	"github.com/firefish111/way2fa/account"
)

const (
	ConfigDirName string = "way2fa"
)

// Where the data was obtained from
type DataSource uint
const (
	NamedSource DataSource = iota
	FileSource
)

type AccountList interface {
	// Retrieve accounts from storage
	GetAccs() ([]account.Account, error)

	// Returns a string detailing the source of the data, to go on the titlebar
	GetSource() (DataSource, string)

	// Write accounts to storage
	WriteAccs(to_write []account.Account)
}
