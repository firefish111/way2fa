package password

import (
	"fmt"
	"strings"

	"github.com/firefish111/way2fa/parse/encryption"
)

// Keeps a backup rendered password prompt, in order to show the end user to make it obvious
// that they are confirming their password.
//
// The result is placed in m.prevRendered. This property can be filled with whatever,
// as only m.prev's nillishness is an indicator of which prompt is shown (initial/confirmation).
func (m *passwordModel) prevRender() {
	// we don't use m.field.Render() as it renders its own cursor,
	// and we don't want to freeze a cursor in time. therefore, we do it ourselves
	// the prompt should look like "> ****" where the password is starred out, so this just copies that

	times := len(m.field.Value()) // number of characters in password
	m.prevRendered = fmt.Sprintf("> %s", strings.Repeat("*", times))
}

// Submit password.
// Returns error and success status: if true, then program can exit.
func (m *passwordModel) submit() (bool, error) {
	var hashed encryption.PasswordHash

	// in this scope only. done to make it obvious that THE RAW PASSWORD IS HERE.
	// we want the ACTUAL PASSWORD to be gc'd asap, so any potential memory vulnerabilities are harder
	{
		raw := m.field.Value()
		if pLen := len(raw); pLen > PasswordMaxLen {
			return false, PromptError{PromptErrorType: TooLong, passlen: pLen}
		} else if pLen == 0 { // is password empty?
			return false, nil // do absolutely nothing
		}
		hashed = encryption.HashPassword(raw)
	}

	if m.prev == nil { // if this is our first go
		// can't take pointer without a binding
		m.prev = &hashed // store current hashed password into a "previous" field
		m.prevRender()   // render previous text
	} else if hashed != *m.prev {
		// if passwords don't match.
		// we clear prev as well, as we want to reset both initial and confirmation. (the first time could've contained the mistake)
		m.prev = nil
		return false, PromptError{PromptErrorType: NotMatch}
	} else { // they match
		if err := m.acclist.Decrypt(*m.prev); err != nil {
			// decrypt account. if it fails, then return that error
			// we clear prev, as an error here likely means wrong password, and we want another chance to enter password twice
			m.prev = nil
			return false, err
		} else { // WE DECRYPTED!!! return true for success
			return true, nil
		}
	}

	return false, nil
}
