package creation

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// borrowed from ../help.go
var none = lipgloss.NewStyle().Foreground(lipgloss.Color("251"))

func (m formModel) modifiedHelp() string {
	keys := append(m.form.KeyBinds(), key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "return to main menu"),
	))

	modelx := m.form.Help()
	// consistency's sake with main menu
	modelx.Styles.ShortKey = none
	modelx.Styles.ShortDesc = faint
	modelx.Styles.FullKey = none
	modelx.Styles.FullDesc = faint

	return modelx.ShortHelpView(keys)
}
