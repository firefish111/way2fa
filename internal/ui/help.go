package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/firefish111/way2fa/internal/ui/common/styles"

	"strconv"
)

// implements help.KeyMap interface.
// requires two methods: ShortHelp, and FullHelp.

// usually, this would be done on a separate type, but it needs access to the state of the model for toggles
// (which i can't believe they don't have support for but whatever, this is my program)

// simple, one line help. used to show basic functions
func (m model) ShortHelp() []key.Binding {
	if m.saveState == tryingExit {
		return []key.Binding{
			m.helpDB["forceexit"],
			m.helpDB["goback"],
		}
	} else if m.selected == nil {
		var helps []key.Binding
		if m.saveState == saved {
			helps = []key.Binding{
				m.helpDB["newfalse"],
				m.helpDB["peek"+strconv.FormatBool(m.peek)], // jank, because no ternary statement
				m.helpDB["quit"+strconv.FormatBool(m.saveState == saved)],
			}
		} else {
			helps = []key.Binding{
				m.helpDB["newfalse"],
				m.helpDB["peek"+strconv.FormatBool(m.peek)], // jank, because no ternary statement
				m.helpDB["save"+strconv.FormatBool(m.reader.IsPasswordProtected())],
				m.helpDB["quit"+strconv.FormatBool(m.saveState == saved)],
			}
		}

		return helps
	} else {
		return []key.Binding{
			m.helpDB["accept"],
			m.helpDB["reject"],
			m.helpDB["down"],
			m.helpDB["up"],
		}
	}
}

// column-based help: for more complicated methods
func (m model) FullHelp() [][]key.Binding {
	// temporary, cheating
	return [][]key.Binding{m.ShortHelp()}
}

func defaultHelp() map[string]key.Binding {
	return map[string]key.Binding{
		"down": key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp(styles.None.Render("j"), "move down"),
		),
		"up": key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp(styles.None.Render("k"), "move up"),
		),
		"accept": key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp(styles.None.Render("y"), "accept + save"),
		),
		"reject": key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp(styles.None.Render("n"), "reject + remove"),
		),
		"newfalse": key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp(styles.None.Render("c"), "create TOTP"),
		),
		"newtrue": key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp(styles.NewOn.Render("esc"), "go back"),
		),
		"peektrue": key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp(styles.PeekOn.Render("p"), "unpeek"),
		),
		"peekfalse": key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp(styles.Off.Render("p"), "peek"),
		),
		"savetrue": key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp(styles.None.Render("s"), "save (asks for password)"),
		),
		"savefalse": key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp(styles.None.Render("s"), "save"),
		),
		"quittrue": key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp(styles.None.Render("q"), "quit"),
		),
		"quitfalse": key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp(styles.None.Render("q"), styles.Discard.Render("quit + discard")),
		),
		"forceexit": key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp(styles.None.Render("y"), styles.Discard.Render("force quit without saving")),
		),
		"goback": key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp(styles.None.Render("n"), "go back"),
		),
	}
}
