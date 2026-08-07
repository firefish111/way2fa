package csv_way

import (
	"errors"
	"time"

	"github.com/firefish111/way2fa/cryptor"
)

// This filetype is NEVER password protected, so always false.
// All other methods are stubs.
func (c *CsvWay) IsPasswordProtected() bool {
	return c.hasPassword
}

// no error, as is always decrypted. error means something went wrong
func (c *CsvWay) Decrypt(password cryptor.PasswordHash) error {
	// idempotency: already decrypted is not an error, so does nothing
	if c.IsDecrypted() {
		return nil
	}

	// NOTE: the above already eliminates the possibility of not being password-protected, so we can assume that we are

	// we know the capabilities for sure: we use the ones that already come with the file's header
	c.crypt.DeriveKey(password, &c.capabilities)

	if !c.IsDecrypted() {
		return errors.New("CsvWay.Decrypt(): Key derivation failed")
	}

	// nothing to decrypt
	if len(c.payload) == 0 {
		return nil
	}

	// verify that password is correct
	_, err := c.crypt.Decrypt(c.payload)
	return err
}

func (c *CsvWay) Recrypt() {
	// idempotency
	if c.crypt.IsValid() {
		c.crypt.Invalidate()
	}
}

// if it isn't password-protected that can't not be decrypted.
// otherwise, check if we have a valid cryptor
func (c *CsvWay) IsDecrypted() bool {
	return !c.IsPasswordProtected() || c.crypt.IsValid()
}

// estimate how long it should take to decrypt the file
func (c *CsvWay) CryptionTimeEstimate() time.Duration {
	if !c.IsPasswordProtected() {
		return 0
	}

	// two seconds per arbitrary unit of time. could possibly be complete overkill
	estimate := time.Duration(c.capabilities.Time) * 2 * time.Second
	return estimate
}
