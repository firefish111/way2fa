package creation

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/donderom/bubblon"
	"github.com/firefish111/way2fa/account"
	"github.com/firefish111/way2fa/internal/ui/msgs"
	"strconv"
)

func (m formModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m formModel) Update(event tea.Msg) (tea.Model, tea.Cmd) {
	switch event := event.(type) {
	case tea.KeyMsg: // handle keypress
		switch event.String() {
		case "esc":
			return m, bubblon.Close
		}
	case msgs.TickMsg: // our own custom tick message struct (just a typedef)
		return m, msgs.Tick() // tick again. this will be executed, and after it times out, update will be called again
	}

	// relay. the update function returns a new copy of self, so we replace
	f, cmd := m.form.Update(event)

	// double check that it is actually a form (no doubt, but still)
	if f, ok := f.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		to, err := account.NewFromTextKey(keyifyKey(m.form.GetString("2fakey")))
		if err == nil {
			to.Name = m.form.GetString("name")
			to.AcctId = handlifyAcctId(m.form.GetString("acctid"))
			if intrv := m.form.GetString("interv"); len(intrv) != 0 {
				num, err := strconv.Atoi(intrv)
				if err == nil {
					unum := uint(num)
					to.Interval = &unum
				}
			}

			return m, tea.Sequence(bubblon.Close, msgs.SendAcct(*to))
		} else {
			return m, bubblon.Close
		}
	}

	return m, cmd
}

func (m formModel) View() string {
	/*
	   	s.WriteString(
	   		app_name.Render("way2fa") +
	   			faint.Render(" - New TOTP"))

	   	s.WriteRune('\n')

	   	s.WriteString(m.createForm.View())

	   	s.WriteString(wip.Render("WIP, press esc to go back"))
	   }
	*/

	return m.form.View()
}
