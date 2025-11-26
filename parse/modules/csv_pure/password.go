package csv_pure

// This filetype is NEVER password protected, so always false.
// All other methods are stubs.
func (c *CsvPure) IsPasswordProtected() bool {
	return false
}

// no error, as is always decrypted. error means something went wrong
func (c *CsvPure) Decrypt(_password string) error {
	return nil
}

func (c *CsvPure) Recrypt() {
	/* can't be recrypted */
}

// always is
func (c *CsvPure) IsDecrypted() bool {
	return true
}
