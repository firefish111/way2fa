package creation

import (
	"github.com/charmbracelet/huh"
)

func (m *formModel) resetForm() {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Service name").
				Validate(validateService).
				Placeholder("e.g. google, microsoft"),

			huh.NewInput().
				Key("acctid").
				Title("Account name").
				Placeholder("e.g. @myusername, myemail@provider.com"),

			huh.NewInput().
				Key("2fakey").
				Title("2FA Key").
				Validate(validateKey).
				Placeholder("Must contain letters and numbers 2-7 only"),

			huh.NewInput().
				Key("interv").
				Title("Interval (in seconds)").
				Validate(validateInterv).
				Placeholder("30 (default)"),

			huh.NewConfirm().
				Key("confirmation").
				Title("Continue?").
				Affirmative("Yes, add to list").
				Negative("No, go back"),
		),
	)
}
