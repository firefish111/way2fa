package password

import (
	"errors"
	"fmt"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/donderom/bubblon/v2"

	"github.com/firefish111/way2fa/cryptor"
	"github.com/firefish111/way2fa/internal/ui/common/msgs"
	"github.com/firefish111/way2fa/internal/ui/common/styles"

	"strings"
)

func (m passwordModel) Init() tea.Cmd {
	// we want a blinking cursor
	return textinput.Blink
}

func (m passwordModel) Update(event tea.Msg) (tea.Model, tea.Cmd) {
	switch event := event.(type) {
	case tea.KeyPressMsg: // handle keypress
		if m.isDecrypting {
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

			// submit, and if obtained a password, we start decryption.
			hashed, err := m.submit()

			// after submit, we clear input field
			m.field.Reset()

			if hashed != nil { // we have a possible password
				m.isDecrypting = true
				// do the decrypt; works as a bubbletea command, which is run concurrently
				return m, doDecrypt(m.acclist, *hashed)
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
		}
	case msgs.TickMsg: // our own custom tick message struct (just a typedef)
		m.ticks++             // timekeeping for loading animation
		return m, msgs.Tick() // tick again. this will be executed, and after it times out, update will be called again
	case didDecryptResultMsg:
		// no matter what happens, after this we won't be decrypting
		m.isDecrypting = false

		err := event.err
		if err == nil { // we did it!!
			if event.didTimeout {
				m.supplMsg = fmt.Sprintf("Decryption timed out (waited %dms) Try again?", event.haveWaited/time.Millisecond)
			} else {
				return m, tea.Sequence(bubblon.Close, msgs.SendEncryptor(msgs.DecryptedMsg)) // broadcast this to entire world
			}
		}

		// if we got a decryption error that is because the password was wrong (that is they matched, but decryption failed),
		// then we sneakily replace the error, as that is not a fail condition: unless that happens 3 times.
		// in which case, we replace the error with our own
		if cryptErr, ok := errors.AsType[cryptor.CryptError](err); ok {
			if cryptErr.IsFaultOfPassword { // if it is the fault of the password that we got this error
				// increase try count; if we surpassed max count
				if m.tries++; m.tries >= PasswordTriesCount {
					err = PromptError{PromptErrorType: OutOfTries}
				} else {
					err = nil // we don't want to pass up password fails
				}

				// tries have actually increased, so:
				m.prev = nil
				m.supplMsg = fmt.Sprintf("Password incorrect.\nPlease try again.\n\tAttempt %d of %d", m.tries+1, PasswordTriesCount)
			}
		}

		// decryption didn't work, undo what we tried to do
		m.acclist.Recrypt()

		// pass up error if we need to
		if err != nil {
			return m, bubblon.Fail(err)
		}

		/* otherwise do nothing */
		return m, nil
	}

	// update textfield.
	// NOTE: textinput.Model is NOT a tea.Model, as it doesn't implement Init().
	var cmd tea.Cmd // cmd so as to not lose sanity over := operator
	m.field, cmd = m.field.Update(event)

	return m, cmd
}

const (
	spinningWheel = "|/-\\"
)

func (m passwordModel) View() tea.View {
	var s strings.Builder

	if m.isDecrypting {
		s.WriteString(styles.Error.Render(
			fmt.Sprintf("Please wait, decrypting... %c", spinningWheel[m.ticks%uint8(len(spinningWheel))]),
		))
	} else if m.warningOnly {
		s.WriteString(styles.Title.Render("Retrieving from: "))
		s.WriteString(styles.RenderSource(m.acclist.GetSource()))
		s.WriteRune('\n')

		s.WriteString(styles.Error.Render(
			styles.Title.Render("WARNING: ") +
				"\n\nThis account list is not password protected.\n" +
				"Please switch to a password-protected format, using the -export flag.",
		))
	} else {
		// whether this is first or second entering
		if m.prev == nil { // prevRendered is ignored, cause it's useless without prev
			s.WriteString(styles.Title.Render(m.title + ": "))
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

	return tea.NewView(styles.Box.Render(s.String()))
}
