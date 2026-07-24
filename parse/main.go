package parse

import (
	"github.com/firefish111/way2fa/account"
	"github.com/firefish111/way2fa/parse/cryptor"
	"github.com/firefish111/way2fa/parse/format"
)

// Where the data was obtained from
type DataSource uint

const (
	NamedSource DataSource = iota
	FileSource
)

// A list of Accounts as an interface. This does not store the accounts themselves,
// merely how to retrieve them from and place them back into some store.
// This allows us to abstract all the file reading code behind this interface.
//
// This also does, as a by-product, some error-checking. Any problems flagged by Validate constitute *detection issues*,
// whereas anything by GetAccs constitutes *content errors*.
//
// detection issues are NOT errors, and ought not be treated as such.
type AccountList interface {
	// Retrieve accounts from storage
	GetAccs() ([]account.Account, error)

	// Returns a string detailing the source of the data, to go on the titlebar
	GetSource() (DataSource, string)

	// Returns a string containing the path of the file this is attached to.
	GetSourceFilePath() string

	// Write accounts to storage
	WriteAccs(to_write []account.Account) error

	// Prepopulate all fields with those of the given file.
	// Completely erases what was already there.
	PrepopulateFromFile(path string, isDefault bool) error

	// Whether is password protected: is used to trigger the password prompt
	//
	// Even the CsvEnc module can *not* have a password, so this isn't all that useless
	IsPasswordProtected() bool

	// Decrypt. This needs to be done (ideally immediately) before any of the meaningful operations (i.e. GetAccs, WriteAccs, etc).
	// If IsPasswordProtected() returns false, this should not error, as it is perpetually decrypted.
	//
	// Is idempotent: in that if already in decrypted state, should do nothing and not error.
	//
	// This should set a password flag in the struct containing the password in one form or another.
	// After any operation, this password flag is CLEARED by the Recrypt() method, such that multiple operations can't be chained and so that the password doesn't reside in memory for too long.
	Decrypt(password cryptor.PasswordHash) error

	// Clears whatever password flag is set.
	// Called by all meaningful operations anyway, so user need not worry, unless explicit Recryption is required (i.e. if Decrypted but no operation was executed).
	//
	// Therefore is idempotent, if is already encrypted, does nothing
	//
	// See Decrypt(string) for use.
	Recrypt()

	// Checks that state of the password flag.
	// If IsPasswordProtected() returns false, this should always return true, as it is prepetually decrypted.
	// See Decrypt(string) for use.
	IsDecrypted() bool
}

// the other, more specific interfaces

type PureAccountList interface {
	AccountList

	// Detect whether the account list is of the correct format.
	//
	// NOTE: this may not necessarily mean that there are no errors in the format, only that basic header checks and filetype checks pass.
	// Therefore, does not return error: as an invalid filetype ought not be an error, only a signal to move to the next filetype.
	Detect() bool

	// obtain the default filename it's expecting. see detector for more info
	GetDefaultFilename() string
}

type WayAccountList interface {
	AccountList

	// Gets the ID assigned to it in the .way file.
	GetWayTypeId() format.FileTypeId
}
