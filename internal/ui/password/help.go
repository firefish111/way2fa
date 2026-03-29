package password

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/firefish111/way2fa/internal/ui/common/styles"

	"strconv"
)

// implements help.KeyMap interface.
// requires two methods: ShortHelp, and FullHelp.

// usually, this would be done on a separate type, but it needs access to the state of the model for toggles
// (which i can't believe they don't have support for but whatever, this is my program)

// refer to ../help.go for some of the rationale behnid this.

// simple, one line help. used to show basic functions
func (m passwordModel) ShortHelp() []key.Binding {
	if m.warningOnly {
		return []key.Binding{
			m.helpDB["acknowledge"],
			m.helpDB["quit"],
			m.helpDB["back"],
		}
	} else {
		return []key.Binding{
			m.helpDB["submit"+strconv.FormatBool(m.prev != nil)], // jank, because no ternary statement
			m.helpDB["back"],
		}
	}
}

// column-based help: for more complicated methods
func (m passwordModel) FullHelp() [][]key.Binding {
	// temporary, cheating
	return [][]key.Binding{m.ShortHelp()}
}

func defaultHelp() map[string]key.Binding {
	return map[string]key.Binding{
		"back": key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp(styles.None.Render("esc"), "go back"),
		),
		"acknowledge": key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp(styles.None.Render("enter"), "acknowledge"),
		),
		"quit": key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp(styles.None.Render("q"), "quit way2fa"),
		),
		"submitfalse": key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp(styles.None.Render("enter"), "next"),
		),
		"submittrue": key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp(styles.None.Render("enter"), "submit"),
		),
	}
}
