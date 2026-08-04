package password

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/firefish111/way2fa/parse"
	"github.com/firefish111/way2fa/parse/cryptor"
)

const (
	// how much extra leeway on top of its predicted time we give the decryption to complete before timing out
	decryptionTimeoutLeeway = 5 * time.Second
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
// Returns error and a context. If context is non-nil, then program is undergoing decryption, and should wait on that
func (m *passwordModel) submit() (context.Context, error) {
	var hashed cryptor.PasswordHash

	// in this scope only. done to make it obvious that THE RAW PASSWORD IS HERE.
	// we want the ACTUAL PASSWORD to be gc'd asap, so any potential memory vulnerabilities are harder
	{
		raw := m.field.Value()
		if pLen := len(raw); pLen > PasswordMaxLen {
			return nil, PromptError{PromptErrorType: TooLong, passlen: pLen}
		} else if pLen == 0 { // is password empty?
			return nil, nil // do absolutely nothing
		}
		hashed = cryptor.HashPassword(raw)
	}

	if m.prev == nil { // if this is our first go
		// can't take pointer without a binding
		m.prev = &hashed // store current hashed password into a "previous" field
		m.prevRender()   // render previous text

		// nothing eventful happened
		return nil, nil
	} else if hashed.Matches(*m.prev) {
		// if passwords don't match.
		// we clear prev as well, as we want to reset both initial and confirmation. (the first time could've contained the mistake)
		m.prev = nil
		return nil, PromptError{PromptErrorType: NotMatch}
	} else { // they match
		ctx := context.Background()
		var cancel context.CancelFunc

		if wayacc, ok := m.acclist.(parse.WayAccountList); ok { // we add a deadline
			// we set the timeout to the cryption time estimate provided by the cryptor, plus an extra 5 seconds leeway
			ctx, cancel = context.WithTimeout(ctx, wayacc.CryptionTimeEstimate()+decryptionTimeoutLeeway)
		} else { // decryption of a pure account list CAN'T POSSIBLY fail, so we just disable its cancellation ability
			ctx = context.WithoutCancel(ctx)
		}

		// do decrypt. cancel after it's done
		go func() {
			m.acclist.Decrypt(ctx, *m.prev)

			// if can be cancelled, do so
			if cancel != nil {
				cancel()
			}
		}()

		return ctx, nil
	}
}
