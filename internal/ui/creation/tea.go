package creation

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/donderom/bubblon"
	"github.com/firefish111/way2fa/account"
	"github.com/firefish111/way2fa/internal/ui/msgs"
	"strconv"
	"strings"
)

func (m formModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m formModel) Update(event tea.Msg) (tea.Model, tea.Cmd) {
	switch event := event.(type) {
	case tea.KeyMsg: // handle keypress
		switch event.String() {
		case "esc", "ctrl+c": // ctrl+c is defined, to stop it reaching the form itself, which leads to the form being inescapable
			return m, bubblon.Close
		}
	case msgs.TickMsg: // our own custom tick message struct (just a typedef)
		return m, msgs.Tick() // tick again. this will be executed, and after it times out, update will be called again
	}

	// relay. the update function returns a new copy of self, so we replace
	f, cmd := m.form.Update(event)

	// double check that it is actually a form (no doubt, but still)
	f, ok := f.(*huh.Form)
	if !ok { // if the form has ceased to be a form. should never happen, but being careful.
		// i'm just waiting for the inevitable github issue "um why does it say this? help pls"
		return m, bubblon.Fail(fmt.Errorf("The form has metamorphosed (bad)"))
	}

	if m.form.State == huh.StateCompleted { // if we have completed form
		if !m.form.GetBool("confirmation") { // if told "no, go back"
			return m, bubblon.Close // bye bye
		}

		// try and create an account from the key
		to, err := account.NewFromTextKey(keyifyKey(m.form.GetString("2fakey")))
		if err == nil { // if there's no inexplicable problem with account creation
			to.Name = m.form.GetString("name")
			to.AcctId = handlifyAcctId(m.form.GetString("acctid"))
			if intrv := m.form.GetString("interv"); len(intrv) != 0 {
				// roundabout way of converting it to a *uint
				num, err := strconv.Atoi(intrv)
				if err == nil { // if it fails, we don't really care cause it'll just be nil anyway
					unum := uint(num)
					to.Interval = &unum
				}
			}

			return m, tea.Sequence(bubblon.Close, msgs.SendAcct(*to))
		} else {
			return m, bubblon.Fail(err) // tantamount to panic
		}
	}

	return m, cmd
}

// styles. these are copied from ../tea.go, because there seriously is no point making an extra package just for them
var app_name = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("220")).
	PaddingLeft(1)

var faint = lipgloss.NewStyle().
	Faint(true).
	Foreground(lipgloss.Color("242"))

var space = lipgloss.NewStyle().
	Margin(1)

func (m formModel) View() string {
	var s strings.Builder

	s.WriteRune('\n')

	// copied from ../tea.go for consistency's sake
	s.WriteString(
		app_name.Render("way2fa") +
			faint.Render(" - New TOTP"))

	s.WriteRune('\n')

	s.WriteString(space.Render(m.form.View()))

	s.WriteString(space.Render(m.modifiedHelp()))

	return s.String()
}
