package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/donderom/bubblon/v2"

	"slices"
	"strings"

	"github.com/firefish111/way2fa/internal/ui/common/msgs"
	"github.com/firefish111/way2fa/internal/ui/common/styles"
	"github.com/firefish111/way2fa/internal/ui/creation"
	"github.com/firefish111/way2fa/internal/ui/password"
)

func (m model) openPasswordPrompt(title string) (cmd tea.Cmd) {
	passwordPrompt := password.CreatePasswordPrompt(m.reader, title, false)

	// send a different command based on whether we need the prompty or not
	if passwordPrompt != nil { // i.e. there is a password prompt we need to use
		cmd = bubblon.Open(passwordPrompt)
	} else {
		// we have proven above that if there is no password prompt needed, then it must already be decrypted, so we send the message
		cmd = msgs.SendEncryptor(msgs.DecryptedMsg)
	}

	return
}

// Initialise main UI.
// Checks to see whether
func (m model) Init() tea.Cmd {
	return tea.Batch(m.openPasswordPrompt("Enter password to decrypt"), msgs.Tick()) // first tick, and whatever we need to send.
}

// this returns the model itself, and anything we want tea to do
func (m model) Update(event tea.Msg) (tea.Model, tea.Cmd) {
	// event is whatever tea wants us to respond to, we need to see what it is
	switch event := event.(type) {
	case tea.KeyPressMsg: // handle keypress
		if m.saveState == tryingExit {
			switch event.String() {
			case "y":
				return m, tea.Quit
			case "n":
				m.saveState = unsaved
			}
		} else if m.selected == nil {
			switch event.String() {
			case "s":
				if m.saveState != saved {
					m.saveState = saveOngoing

					return m, m.openPasswordPrompt("Re-enter existing password to save")
				}
			case "q":
				// if we want to quit, but have unsaved changes, show the exit prompt
				switch m.saveState {
				case saveOngoing:
					/* do nothing. don't quit, but don't set to tryingExit to avoid race condition */
				case unsaved:
					m.saveState = tryingExit
				default:
					// if saved. the "tryingExit" case is already dealt with further up
					return m, tea.Quit
				}
			case "c":
				return m, bubblon.Open(creation.DefaultForm())
			case "p":
				m.peek = !m.peek
			}
		} else {
			switch event.String() { // dirty is already proven to be non-nil by this point
			case "j":
				if *m.selected < len(m.accs)-1 { // can we move down?
					m.accs[*m.selected], m.accs[*m.selected+1] = m.accs[*m.selected+1], m.accs[*m.selected]
					*m.selected++
				}
			case "k":
				if *m.selected > 0 { // can we move up?
					m.accs[*m.selected], m.accs[*m.selected-1] = m.accs[*m.selected-1], m.accs[*m.selected]
					*m.selected--
				}
			case "n":
				m.accs = slices.Replace(m.accs, *m.selected, *m.selected+1)
				m.selected = nil
			case "y":
				m.saveState = unsaved
				m.selected = nil
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

		m.selected = &insert_at // make us show the save changes? prompt
		m.accs = slices.Insert(m.accs, insert_at, event.Acct)
	case msgs.CryptorMsg: // If Encrypted
		if event == msgs.DecryptedMsg { // decrypted! therefore retrieve data
			if m.saveState == saveOngoing {
				m.saveState = saved
				err := m.reader.SetAccs(m.accs)
				if err != nil {
					return m, bubblon.Fail(err) // pass up error
				}
			} else {
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
	}

	return m, nil
}

// Spit it out
func (m model) View() tea.View {
	var s strings.Builder

	s.WriteRune('\n')

	m.writeOTPs(&s)

	s.WriteRune('\n')

	if m.selected != nil {
		s.WriteString(
			styles.Spaced.Render(
				"Please confirm with " +
					styles.SrvcName.Render(m.accs[*m.selected].Name) +
					" that the highlighted OTP is correct before proceeding.\n" +
					"If you wish, you can also place it in the desired order."))
		s.WriteRune('\n')
	}

	if m.saveState == tryingExit {
		s.WriteString(styles.Supplement.Render("Are you sure you want to quit without saving?"))
		s.WriteRune('\n')
	}

	helpview := m.helpModel.View(m) // using self as a help model to access internal state
	s.WriteString(styles.SidePad.Render(helpview))

	s.WriteRune('\n')

	return tea.NewView(s.String())
}
