package password

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/donderom/bubblon"

	"github.com/firefish111/way2fa/internal/ui/msgs"
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

var box = lipgloss.NewStyle().
	BorderForeground(lipgloss.Color("6")).
	//	Align(lipgloss.Center).
	//	Border(lipgloss.DoubleBorder()).
	Padding(1)

var title = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("157"))

var supplement = lipgloss.NewStyle().
	Italic(true).
	Padding(2, 4).
	Foreground(lipgloss.Color("9"))

// copied from ..
var faint = lipgloss.NewStyle().
	Faint(true).
	Foreground(lipgloss.Color("242"))

var source = lipgloss.NewStyle().
	Bold(true).
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("159")).
	Background(lipgloss.Color("88")).
	Padding(0, 1).
	Margin(0, 1)

func (m passwordModel) putSrc(s *strings.Builder) {
	srct, srcs := m.acclist.GetSource()

	style := source
	if srct == parse.FileSource {
		style = style.Background(lipgloss.Color("22"))
	}

	s.WriteString(style.Render(srcs))
}

func (m passwordModel) View() string {
	var s strings.Builder

	// whether this is first or second entering
	if m.prev == nil { // prevRendered is ignored, cause it's useless without prev
		s.WriteString(title.Render("Enter password: "))
		m.putSrc(&s)
	} else {
		s.WriteString(title.Render("Confirm password: "))
		m.putSrc(&s)
		s.WriteRune('\n')
		s.WriteString(faint.Render(m.prevRendered))
	}
	s.WriteRune('\n')
	s.WriteString(m.field.View())

	// if we have something to say?
	if m.supplMsg != "" {
		s.WriteString(supplement.Render(m.supplMsg))
	}

	return box.Render(s.String())
}
