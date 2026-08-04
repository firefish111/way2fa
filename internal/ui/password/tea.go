package password

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/donderom/bubblon"

	"github.com/firefish111/way2fa/internal/ui/common/msgs"
	"github.com/firefish111/way2fa/internal/ui/common/styles"
	"github.com/firefish111/way2fa/parse/cryptor"

	"strings"
)

func (m passwordModel) Init() tea.Cmd {
	// we want a blinking curs"Please wait
	return textinput.Blink
}

func (m passwordModel) Update(event tea.Msg) (tea.Model, tea.Cmd) {
	if m.decryptContext != nil {
		select {
		case <-m.decryptContext.Done():
			if m.decryptContext.Err() == context.DeadlineExceeded {
				m.supplMsg = "Decryption timed out."
			} else if context.Cause(m.decryptContext) == context.Canceled {
				// if cancelled with no reason. this is only done if it completed successfully, so we exit
				if !m.acclist.IsDecrypted() {
					return m, bubblon.Fail(fmt.Errorf("ui: password: account list decryption claimed it succeeded; it didn't"))
				} else {
					// we did it!!
					return m, tea.Sequence(bubblon.Close, msgs.SendEncryptor(msgs.DecryptedMsg)) // broadcast this to entire world
				}
			} else {
				err := context.Cause(m.decryptContext)
				// take a closer look at the error
				switch outerr := err.(type) {
				// if we got a decryption error that is because the password was wrong (that is they matched, but decryption failed),
				// then we sneakily replace the error, as that is not a fail condition: unless that happens 3 times.
				// in which case, we replace the error with our own
				case cryptor.CryptError:
					if outerr.IsFaultOfPassword { // if it is the fault of the password that we got this error
						// increase try count; if we surpassed max count
						if m.tries++; m.tries >= PasswordTriesCount {
							err = PromptError{PromptErrorType: OutOfTries}
						} else {
							err = nil // we don't want to pass up password fails
						}

						// tries have actually increased, so:
						m.supplMsg = fmt.Sprintf("%v\nPlease try again.\n\tAttempt %d of %d", err, m.tries+1, PasswordTriesCount)
					}
				}

				// pass up error if we need to
				if err != nil {
					return m, bubblon.Fail(err)
				}
			}

			// remove the context, as it will be created by another submit
			m.decryptContext = nil
		default:
			/* do nothing */
		}
	}

	switch event := event.(type) {
	case tea.KeyMsg: // handle keypress
		if m.decryptContext != nil {
			// please do NOT submit again, or handle any other keypresses for that matter, if already submitted
			break
		}

		switch event.String() {
		case "q":
			if m.warningOnly {
				return m, tea.Quit // goodbye
			}
		/* TODO: make it so that it can actually handle not decrypting
		case "esc", "ctrl+c": // able to leave
			return m, bubblon.Close // just close
		*/
		case "enter":
			if m.warningOnly { // not password protected.
				// as it's not password protected, we still have to send a "decrypted" message, as rest of code expects that it is
				return m, tea.Sequence(bubblon.Close, msgs.SendEncryptor(msgs.DecryptedMsg)) // broadcast this to entire world
			}

			// submit password, and get error.
			ctx, err := m.submit()

			if ctx != nil { // we got it decrypted!!
				// get context and save it
				m.decryptContext = ctx
			}

			// take a closer look at the error
			// NOTE: as of yet, submit() only returns PromptErrors
			switch outerr := err.(type) {
			case PromptError:
				// if passwords don't match: we ignore error.
				if outerr.PromptErrorType == NotMatch {
					m.supplMsg = "Passwords didn't match, please try again"
					// TODO: may want to do something with this later
					err = nil
				}
			}

			// pass up error if we need to
			if err != nil {
				return m, bubblon.Fail(err)
			}

			// otherwise, after submit, we clear input field
			m.field.Reset()
		}
	case msgs.TickMsg: // our own custom tick message struct (just a typedef)
		m.ticks++
		return m, msgs.Tick() // tick again. this will be executed, and after it times out, update will be called again
	}

	// update textfield.
	// NOTE: textinput.Model is NOT a tea.Model, as it doesn't implement Init().
	var cmd tea.Cmd // cmd so as to not lose sanity over := operatir
	m.field, cmd = m.field.Update(event)

	return m, cmd
}

func (m passwordModel) View() string {
	var s strings.Builder

	if m.decryptContext != nil {
		s.WriteString(styles.Error.Render(
			fmt.Sprintf("Please wait, decrypting... %c", "-\\|/"[m.ticks&3]),
		))
	} else if m.warningOnly {
		s.WriteString(styles.Title.Render("Retrieving from: "))
		s.WriteString(styles.RenderSource(m.acclist.GetSource()))
		s.WriteRune('\n')

		s.WriteString(styles.Error.Render(
			styles.Title.Render("WARNING: ") +
				"\n\nThis account list is not password protected.\n" +
				// TODO: once passwording is done
				//"Please switch to a password-protected format!",
				"Password-protected formats are coming soon.",
		))
	} else {
		// whether this is first or second entering
		if m.prev == nil { // prevRendered is ignored, cause it's useless without prev
			s.WriteString(styles.Title.Render("Enter password: "))
			s.WriteString(styles.RenderSource(m.acclist.GetSource()))
		} else {
			s.WriteString(styles.Title.Render("Confirm password: "))
			s.WriteString(styles.RenderSource(m.acclist.GetSource()))
			s.WriteRune('\n')
			s.WriteString(styles.Faint.Render(m.prevRendered))
		}
		s.WriteRune('\n')
		s.WriteString(m.field.View())
	}

	s.WriteRune('\n')

	// if we have something to say?
	if m.supplMsg != "" {
		s.WriteString(styles.Supplement.Render(m.supplMsg))
	}

	s.WriteRune('\n')

	helpview := m.helpModel.View(m) // using self as a help model to access internal state
	s.WriteString(styles.SidePad.Render(helpview))

	return styles.Box.Render(s.String())
}
