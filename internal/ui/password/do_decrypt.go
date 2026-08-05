package password

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/firefish111/way2fa/parse"
	"github.com/firefish111/way2fa/parse/cryptor"
)

const (
	// how much extra leeway on top of its predicted time we give the decryption to complete before timing out
	decryptionTimeoutLeeway = 5 * time.Second
)

// a message to fit with the decryptor bubbletea command.
// this is different from internal/ui/common/msgs.CryptorMsg because that just tells the
type didDecryptResultMsg struct {
	didTimeout bool
	err        error

	// wait time: put here to pass up call chain
	haveWaited time.Duration
}

func doDecrypt(acclist parse.AccountList, hash cryptor.PasswordHash) tea.Cmd {
	return func() tea.Msg {
		// make asynchronous error read.
		// in a goroutine because timeout happens at the same time
		errChan := make(chan error, 1)

		go func() {
			// go do the decrypt.
			errChan <- acclist.Decrypt(hash)
		}()

		// in the meantime, set the timeout.
		var timeout time.Duration
		if wayacc, ok := acclist.(parse.WayAccountList); ok {
			timeout = wayacc.CryptionTimeEstimate() + decryptionTimeoutLeeway
		} else {
			// we are a pure format, and decryption is instant.
			// forgo the timeout, and return the error.
			// the read from the channel blocks until decryption is done.
			return didDecryptResultMsg{didTimeout: false, err: <-errChan}
		}

		// use a timer to time it
		timer := time.NewTimer(timeout)
		defer timer.Stop() // stop the timer just in case

		select {
		// if it yields first
		case err := <-errChan:
			return didDecryptResultMsg{didTimeout: false, err: err}
		// it runs out of time first
		case <-timer.C:
			return didDecryptResultMsg{didTimeout: true, err: nil, haveWaited: timeout}
			/* there is no default, so this blocks until something happens */
		}
	}
}
