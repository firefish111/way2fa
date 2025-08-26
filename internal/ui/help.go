package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"

	"strconv"
)

// implements help.KeyMap interface.
// requires two methods: ShortHelp, and FullHelp.

// usually, this would be done on a separate type, but it needs access to the state of the model for toggles
// (which i can't believe they don't have support for but whatever, this is my program)

// simple, one line help. used to show basic functions
func (m model) ShortHelp() []key.Binding {
	if !m.create {
		return []key.Binding{
			m.helpDB["newfalse"],
			m.helpDB["peek"+strconv.FormatBool(m.peek)], // jank, because no ternary statement
			m.helpDB["quit"],
		}
	} else {
		return []key.Binding{
			m.helpDB["newtrue"],
		}
	}
}

// column-based help: for more complicated methods
func (m model) FullHelp() [][]key.Binding {
	// temporary, cheating
	return [][]key.Binding{m.ShortHelp()}
}

var none = lipgloss.NewStyle().Foreground(lipgloss.Color("251"))
var peekon = lipgloss.NewStyle().Foreground(lipgloss.Color("218")).Bold(true)
var newon = lipgloss.NewStyle().Foreground(lipgloss.Color("112")).Bold(true)
var off = lipgloss.NewStyle().Foreground(lipgloss.Color("251"))

func defaultHelp() map[string]key.Binding {
	return map[string]key.Binding{
		"newfalse": key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp(none.Render("n"), "new TOTP"),
		),
		"newtrue": key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp(newon.Render("esc"), "go back"),
		),
		"peektrue": key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp(peekon.Render("p"), "unpeek"),
		),
		"peekfalse": key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp(off.Render("p"), "peek"),
		),
		"quit": key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp(none.Render("q"), "quit"),
		)}
}
