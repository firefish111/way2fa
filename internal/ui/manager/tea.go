package manager

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/donderom/bubblon/v2"
	"github.com/firefish111/way2fa/internal/ui/common/msgs"
	"github.com/firefish111/way2fa/internal/ui/common/styles"
	"github.com/firefish111/way2fa/internal/ui/password"
	"github.com/firefish111/way2fa/parse"
)

func (m managerModel) Init() tea.Cmd {
	return nil
}

func (m managerModel) Update(event tea.Msg) (tea.Model, tea.Cmd) {
	switch event := event.(type) {
	case tea.KeyPressMsg:
		switch event.String() {
		case "q":
			if m.phase != waitingPassword {
				return m, tea.Quit
			}
		case "j":
			if m.phase == selectStore {
				if m.selected < len(m.possibilities)-1 { // can we move down?
					m.selected++
				}
			}
		case "k":
			if m.phase == selectStore {
				if m.selected > 0 { // can we move up?
					m.selected--
				}
			}
		case "enter":
			// get store index
			if m.phase == selectStore {
				fname, isDefault, err := m.matchFilename()
				if err != nil {
					return m, bubblon.Fail(err)
				}

				m.filename = &fname

				m.possibilities[m.selected].PrepopulateFromFile(fname, isDefault)

				if _, ok := m.possibilities[m.selected].(parse.WayAccountList); ok {
					m.phase = setPassword
				} else {
					m.phase = done

					// no password to be set
					m.possibilities[m.selected].PopulateNew()

					return m, msgs.SendEncryptor(msgs.DecryptedMsg)
				}
			}
		case "y":
			if m.phase == setPassword {
				// guaranteed to succeed
				m.possibilities[m.selected].(parse.WayAccountList).SetPasswordProtected(true)

				m.phase = waitingPassword

				// populate new
				m.possibilities[m.selected].PopulateNew()

				model := password.CreatePasswordPrompt(m.possibilities[m.selected], "Set a password", true, false)
				return m, bubblon.Open(model)
			}
		case "n":
			if m.phase == setPassword {
				// guaranteed to succeed
				m.possibilities[m.selected].(parse.WayAccountList).SetPasswordProtected(false)

				m.phase = done

				m.possibilities[m.selected].PopulateNew()

				return m, msgs.SendEncryptor(msgs.DecryptedMsg)
			}
		}
	case msgs.CryptorMsg: // If Encrypted
		if event == msgs.DecryptedMsg { // decrypted! therefore retrieve data
			m.phase = done

			m.possibilities[m.selected].Save()
			*m.destination = m.possibilities[m.selected]
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m managerModel) View() tea.View {
	var s strings.Builder

	switch m.phase {
	case selectStore:
		// TODO: replace by a proper bubbletea ui
		s.WriteString(styles.SidePad.Render("Select destination store type: \n"))
		s.WriteString(MakeListTable(&m.selected))
	case setPassword:
		s.WriteString(styles.SidePad.Render("Set a password? (y/n)\n"))
	default:
		s.WriteString(styles.SidePad.Render(fmt.Sprintf("Created store at %s\n", *m.filename)))
	}

	s.WriteRune('\n')

	helpview := m.helpModel.View(m) // using self as a help model to access internal state
	s.WriteString(styles.SidePad.Render(helpview))

	s.WriteRune('\n')

	return tea.NewView(s.String())
}
