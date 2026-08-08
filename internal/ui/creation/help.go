package creation

import (
	"charm.land/bubbles/v2/key"
	"github.com/firefish111/way2fa/internal/ui/common/styles"
)

func (m formModel) modifiedHelp() string {
	keys := append(m.form.KeyBinds(), key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "return to main menu"),
	))

	modelx := m.form.Help()
	// consistency's sake with main menu
	modelx.Styles.ShortKey = styles.None
	modelx.Styles.ShortDesc = styles.Faint
	modelx.Styles.FullKey = styles.None
	modelx.Styles.FullDesc = styles.Faint

	return modelx.ShortHelpView(keys)
}
