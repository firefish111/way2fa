package creation

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/donderom/bubblon"
	"github.com/firefish111/way2fa/internal/ui/msgs"
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
		// TODO: send a message to parent containing Account
		return m, bubblon.Close
	}

	return m, cmd
}

func (m formModel) View() string {
	return m.form.View()
}
