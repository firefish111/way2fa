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
	return []key.Binding {
		m.helpDB["new"],
		m.helpDB["peek" + strconv.FormatBool(m.peek)], // jank, because no ternary statement
		m.helpDB["quit"] }
}

// column-based help: for more complicated methods
func (m model) FullHelp() [][]key.Binding {
	// temporary, cheating
	return [][]key.Binding { m.ShortHelp() }
}

var none = lipgloss.NewStyle().Foreground(lipgloss.Color("251"))
var on = lipgloss.NewStyle().Foreground(lipgloss.Color("218")).Bold(true)
var off = lipgloss.NewStyle().Foreground(lipgloss.Color("251"))
