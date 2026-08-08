package manager

import (
	"charm.land/bubbles/v2/key"
	"github.com/firefish111/way2fa/internal/ui/common/styles"
)

// implements help.KeyMap interface.
// requires two methods: ShortHelp, and FullHelp.

// usually, this would be done on a separate type, but it needs access to the state of the model for toggles
// (which i can't believe they don't have support for but whatever, this is my program)

// simple, one line help. used to show basic functions
func (m managerModel) ShortHelp() []key.Binding {
	switch m.phase {
	case selectStore:
		return []key.Binding{
			m.helpDB["down"],
			m.helpDB["up"],
			m.helpDB["select"],
			m.helpDB["quit"],
		}
	case setPassword:
		return []key.Binding{
			m.helpDB["accept"],
			m.helpDB["reject"],
			m.helpDB["quit"],
		}
	default:
		return []key.Binding{}
	}
}

// column-based help: for more complicated methods
func (m managerModel) FullHelp() [][]key.Binding {
	// temporary, cheating
	return [][]key.Binding{m.ShortHelp()}
}

func defaultHelp() map[string]key.Binding {
	return map[string]key.Binding{
		"down": key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp(styles.None.Render("j"), "down"),
		),
		"up": key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp(styles.None.Render("k"), "up"),
		),
		"select": key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp(styles.None.Render("enter"), "select"),
		),
		"accept": key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp(styles.None.Render("y"), "yes, set password"),
		),
		"reject": key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp(styles.None.Render("n"), "no thanks"),
		),
		"quit": key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp(styles.None.Render("q"), "cancel + quit"),
		),
	}
}
