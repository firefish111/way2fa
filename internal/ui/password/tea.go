package password

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/donderom/bubblon"

	"github.com/firefish111/way2fa/internal/ui/common/msgs"
	"github.com/firefish111/way2fa/internal/ui/common/styles"
	"github.com/firefish111/way2fa/parse"

	"strings"
)

func (m passwordModel) Init() tea.Cmd {
	// we want a blinking cursor
	return textinput.Blink
}

func (m passwordModel) Update(event tea.Msg) (tea.Model, tea.Cmd) {
	switch event := event.(type) {
	case tea.KeyMsg: // handle keypress
		switch event.String() {
		case "esc", "ctrl+c": // able to leave
			return m, bubblon.Close // just close
		case "enter":
			// submit password, and get error.
			success, err := m.submit()

			if success { // we got it decrypted!!
				return m, tea.Sequence(bubblon.Close, msgs.SendEncryptor(msgs.DecryptedMsg)) // broadcast this to entire world
			}

			// take a closer look at the error
			switch outerr := err.(type) {
			// if we got a decryption error that is because the password was wrong (that is they matched, but decryption failed),
			// then we sneakily replace the error, as that is not a fail condition: unless that happens 3 times.
			// in which case, we replace the error with our own
			case parse.DecryptError:
				if outerr.IsFaultOfPassword { // if it is the fault of the password that we got this error
					// increase try count; if we surpassed max count
					if m.tries++; m.tries >= PasswordTriesCount {
						err = PromptError{PromptErrorType: OutOfTries}
					} else {
						err = nil // we don't want to pass up password fails
					}

					// tries have actually increased, so:
					m.supplMsg = fmt.Sprintf("Password incorrect. Please try again.\n\tAttempt %d of %d", m.tries+1, PasswordTriesCount)
				}
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

	// if we have something to say?
	if m.supplMsg != "" {
		s.WriteString(styles.Supplement.Render(m.supplMsg))
	}

	return styles.Box.Render(s.String())
}
