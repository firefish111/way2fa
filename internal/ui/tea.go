package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/donderom/bubblon"

	"github.com/firefish111/way2fa/internal/ui/creation"
	"github.com/firefish111/way2fa/internal/ui/msgs"
	"slices"
	"strings"
)

// XXX: this could be important in the future
func (m model) Init() tea.Cmd {
	return msgs.Tick() // first tick.
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
		s.WriteString("Please check that the highlighted OTP is correct, and move it to the desired location, before proceeding.\n ")
	}

	helpview := m.helpModel.View(m) // using self as a help model to access internal state
	s.WriteString(helpview)

	s.WriteRune('\n')

	return s.String()
}
