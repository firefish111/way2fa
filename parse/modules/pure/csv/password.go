package csv_pure

import (
	"context"

	"github.com/firefish111/way2fa/parse/cryptor"
)

// This filetype is NEVER password protected, so always false.
// All other methods are stubs.
func (c *CsvPure) IsPasswordProtected() bool {
	return false
}

// no error, as is always decrypted. error means something went wrong
func (c *CsvPure) Decrypt(_ctx context.Context, _password cryptor.PasswordHash) {
	/* nothing */
}

func (c *CsvPure) Recrypt() {
	/* can't be recrypted */
}

// always is
func (c *CsvPure) IsDecrypted() bool {
	return true
}
