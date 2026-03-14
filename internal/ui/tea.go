package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/donderom/bubblon"

	"slices"
	"strings"

	"github.com/firefish111/way2fa/internal/ui/creation"
	"github.com/firefish111/way2fa/internal/ui/msgs"
	"github.com/firefish111/way2fa/internal/ui/password"
)

// Initialise main UI.
// Checks to see whether
func (m model) Init() tea.Cmd {
	passwordPrompt := password.CreatePasswordPrompt(m.reader)

	// this is the message that we want to send.
	// depends on passwordPrompt's value.
	var toSend tea.Cmd

	if passwordPrompt != nil { // i.e. there is a password prompt we need to use
		toSend = bubblon.Open(passwordPrompt)
	} else {
		// we have proven above that if there is no password prompt needed, then it must already be decrypted, so we send the message
		toSend = msgs.SendEncryptor(msgs.DecryptedMsg)
	}

	return tea.Batch(toSend, msgs.Tick()) // first tick, and whatever we need to send.
}

// this returns the model itself, and anything we want tea to do
func (m model) Update(event tea.Msg) (tea.Model, tea.Cmd) {
	// event is whatever tea wants us to respond to, we need to see what it is
	switch event := event.(type) {
	case tea.KeyMsg: // handle keypress
		if m.dirty == nil {
			switch event.String() {
			case "q":
				return m, tea.Quit // bye bye (does broadcast message, see above)
			case "c":
				return m, bubblon.Open(creation.DefaultForm())
			case "p":
				m.peek = !m.peek
			}
		} else {
			switch event.String() { // dirty is already proven to be non-nil by this point
			case "j":
				if *m.dirty < len(m.accs)-1 { // can we move down?
					m.accs[*m.dirty], m.accs[*m.dirty+1] = m.accs[*m.dirty+1], m.accs[*m.dirty]
					*m.dirty++
				}
			case "k":
				if *m.dirty > 0 { // can we move up?
					m.accs[*m.dirty], m.accs[*m.dirty-1] = m.accs[*m.dirty-1], m.accs[*m.dirty]
					*m.dirty--
				}
			case "n":
				m.accs = slices.Replace(m.accs, *m.dirty, *m.dirty+1)
				fallthrough
			case "y":
				go m.reader.WriteAccs(m.accs)
				m.dirty = nil
			}
		}
	case msgs.TickMsg: // our own custom tick message struct (just a typedef)
		return m, msgs.Tick() // tick again. this will be executed, and after it times out, update will be called again
	case msgs.NewAccMsg: // we get this from the form
		// to group them together, we tack it on the end of an existing group (if possible)
		// this is just the default positioning, the prompt will allow user to move it
		insert_at := len(m.accs)
		for acci, accv := range slices.Backward(m.accs) { // iterators are lazily evaluated
			if accv.Name == event.Acct.Name {
				insert_at = acci + 1 // one after
				break
			}
		}

		m.dirty = &insert_at // make us show the save changes? prompt
		m.accs = slices.Insert(m.accs, insert_at, event.Acct)
	case msgs.EncryptorMsg: // If Encrypted
		if event == msgs.DecryptedMsg { // decrypted! therefore retrieve data
			var err error
			m.accs, err = m.reader.GetAccs()
			if err != nil {
				return m, bubblon.Fail(err) // pass up error
			}
		}
	}

	return m, nil
}

var source = lipgloss.NewStyle().
	Bold(true).
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("159")).
	Background(lipgloss.Color("88")).
	Padding(0, 1).
	Margin(0, 1)

// have been copied into ./creation/tea.go. if these ever change, change them there too
var app_name = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("220")).
	PaddingLeft(1)

// used for the confirmation prompt
var srvc_name = lipgloss.NewStyle().
	Italic(true).
	Foreground(lipgloss.Color("117"))

var faint = lipgloss.NewStyle().
	Faint(true).
	Foreground(lipgloss.Color("242"))

var marg = lipgloss.NewStyle().
	MarginLeft(4)

var wip = lipgloss.NewStyle().
	Margin(1, 2).
	Padding(1, 2).
	Background(lipgloss.Color("239")).
	Foreground(lipgloss.Color("15")).
	Bold(true)

// Spit it out
func (m model) View() string {
	var s strings.Builder

	s.WriteRune('\n')

	m.writeOTPs(&s)

	s.WriteRune('\n')
	s.WriteRune(' ')

	if m.dirty != nil {
		s.WriteString("Please confirm with ")
		s.WriteString(srvc_name.Render(m.accs[*m.dirty].Name))
		s.WriteString(" that the highlighted OTP is correct before proceeding.\nIf you wish, you can also place in the desired order.\n ")
	}

	helpview := m.helpModel.View(m) // using self as a help model to access internal state
	s.WriteString(helpview)

	s.WriteRune('\n')

	return s.String()
}
