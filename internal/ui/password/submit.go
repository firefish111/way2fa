package password

import (
	"github.com/firefish111/way2fa/parse/encryption"
)

// Submit password.
// Returns error and success status: if true, then program can exit.
func (m passwordModel) submit() (bool, error) {
	var hashed encryption.PasswordHash

	{
		// raw password value, only in this scope
		raw := m.field.Value()
		if plen := len(raw); plen > PasswordMaxLen {
			return false, PromptError{PromptErrorType: TooLong, passlen: plen}
		}
		hashed = encryption.HashPassword(raw)
	}

	if m.prev == nil {
		// can't take pointer without a binding
		m.prev = &hashed
	} else if hashed != *m.prev {
		// if passwords don't match.
		// we clear prev as well, as that could be the wrong one
		m.prev = nil
		return false, PromptError{PromptErrorType: NotMatch}
	} else {
		if err := m.acclist.Decrypt(*m.prev); err != nil {
			// decrypt account. if it fails, then return that error
			return false, err
		} else { // WE DECRYTPED!!! return true for success
			return true, nil
		}
	}

	return false, nil
}
