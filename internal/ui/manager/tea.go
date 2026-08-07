package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/donderom/bubblon"
	"github.com/firefish111/way2fa/detector"
	"github.com/firefish111/way2fa/internal/config"
	"github.com/firefish111/way2fa/internal/ui/common/msgs"
	"github.com/firefish111/way2fa/internal/ui/password"
	"github.com/firefish111/way2fa/parse"
)

// returns whether the name is a default or not and an error
func (m *managerModel) matchFilename() (isDefault bool, err error) {
	// deal with filename
	var fname string

	if m.filename != nil && *m.filename != "" {
		isDefault = false
		fname, err = filepath.Abs(*m.filename)
		if err != nil {
			return false, err
		}
	} else if pure, ok := m.selected.(parse.PureAccountList); ok {
		isDefault = true
		fname = filepath.Join(
			config.ConfPath,
			pure.GetDefaultFilename(),
		)
	} else if _, ok := m.selected.(parse.WayAccountList); ok {
		isDefault = true
		fname = filepath.Join(
			config.ConfPath,
			detector.DefaultWayFilename,
		)

		m.filename = &fname
	} else {
		return isDefault, fmt.Errorf("Can't match account list to file")
	}

	f, err := os.OpenFile(fname, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		if os.IsExist(err) {
			return isDefault, fmt.Errorf("File %s already exists: %w", fname, err)
		}

		return isDefault, err
	}
	f.Close() // leave it to the interface itself to do that

	return
}

func (m managerModel) Init() tea.Cmd {
	return nil
}

func (m managerModel) Update(event tea.Msg) (tea.Model, tea.Cmd) {
	switch event := event.(type) {
	case tea.KeyMsg:
		switch event.String() {
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// get store index
			if m.phase == selectStore {
				m.selected = m.possibilities[int(event.String()[0]-'1')]

				isDefault, err := m.matchFilename()
				if err != nil {
					return m, bubblon.Fail(err)
				}

				m.selected.PrepopulateFromFile(*m.filename, isDefault)

				if _, ok := m.selected.(parse.WayAccountList); ok {
					m.phase = setPassword
				} else {
					m.phase = done

					// no password to be set
					m.selected.PopulateNew()

					return m, msgs.SendEncryptor(msgs.DecryptedMsg)
				}
			}
		case "y":
			if m.phase == setPassword {
				// guaranteed to succeed
				m.selected.(parse.WayAccountList).SetPasswordProtected(true)

				m.phase = waitingPassword

				// populate new
				m.selected.PopulateNew()

				model := password.CreatePasswordPrompt(m.selected, "Set a password", true)
				return m, bubblon.Open(model)
			}
		case "n":
			if m.phase == setPassword {
				// guaranteed to succeed
				m.selected.(parse.WayAccountList).SetPasswordProtected(false)

				m.phase = done

				m.selected.PopulateNew()

				return m, msgs.SendEncryptor(msgs.DecryptedMsg)
			}
		}
	case msgs.CryptorMsg: // If Encrypted
		if event == msgs.DecryptedMsg { // decrypted! therefore retrieve data
			m.phase = done

			m.selected.Save()
			*m.destination = m.selected
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m managerModel) View() string {
	var s strings.Builder

	switch m.phase {
	case selectStore:
		// TODO: replace by a proper bubbletea ui
		s.WriteString("Select destination store type: \n")
		for i, typ := range m.possibilities { // print all of them
			fmt.Fprintf(&s, "\t[%d]: %s\n", i+1, reflect.TypeOf(typ).Elem().Name())
		}
	case setPassword:
		s.WriteString("Set a password? (y/n)\n")
	default:
		fmt.Fprintf(&s, "Created store at %s", *m.filename)
	}

	return s.String()
}
