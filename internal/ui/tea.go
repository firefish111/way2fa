package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/donderom/bubblon"

	"slices"
	"strings"

	"github.com/firefish111/way2fa/internal/ui/common/msgs"
	"github.com/firefish111/way2fa/internal/ui/common/styles"
	"github.com/firefish111/way2fa/internal/ui/creation"
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
				go m.reader.SetAccs(m.accs)
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
	case msgs.CryptorMsg: // If Encrypted
		if event == msgs.DecryptedMsg { // decrypted! therefore retrieve data
			var err error
			m.accs, err = m.reader.GetAccs()
			if err != nil {
				return m, bubblon.Fail(err) // pass up error
			}

			// Repair all. Just in case our reader messed something up.
			for ix := range m.accs {
				// can't use reassignment, because this takes a mutable pointer
				m.accs[ix].RepairValues()
			}
		}
	}

	return m, nil
}

// Spit it out
func (m model) View() string {
	var s strings.Builder

	s.WriteRune('\n')

	m.writeOTPs(&s)

	s.WriteRune('\n')

	if m.dirty != nil {
		s.WriteString(
			styles.Spaced.Render(
				"Please confirm with " +
					styles.SrvcName.Render(m.accs[*m.dirty].Name) +
					" that the highlighted OTP is correct before proceeding.\n" +
					"If you wish, you can also place it in the desired order."))
		s.WriteRune('\n')
	}

	helpview := m.helpModel.View(m) // using self as a help model to access internal state
	s.WriteString(styles.SidePad.Render(helpview))

	s.WriteRune('\n')

	return s.String()
}
