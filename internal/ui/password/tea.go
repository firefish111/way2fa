package password

import (
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

			// if we got a decryption error that is because the password was wrong (that is they matched, but decryption failed),
			// then we sneakily replace the error, as that is not a fail condition: unless that happens 3 times.
			// in which case, we replace the error with our own
			if decErr, ok := err.(parse.DecryptError); ok {
				if decErr.IsFaultOfPassword { // if it is the fault of the password that we got this error
					// increase try count; if we surpassed max count
					if m.tries++; m.tries >= PasswordTriesCount {
						err = PromptError{PromptErrorType: OutOfTries}
					} else {
						err = nil // we don't want to pass up password fails
					}
				}
			}

			// pass up error iff we need to
			if err != nil {
				return m, bubblon.Fail(err)
			}

			// otherwise, after sumbit, we clear input field
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
	Foreground(lipgloss.Color("162"))

func (m passwordModel) View() string {
	var s strings.Builder

	s.WriteString(title.Render("Password: "))
	s.WriteRune('\n')
	s.WriteString(m.field.View())

	return box.Render(s.String())
}
