package creation

import (
	"charm.land/huh/v2"
	"github.com/firefish111/way2fa/internal/ui/common/styles"
)

func (m *formModel) resetForm() {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title(styles.Header.Render("service name")).
				Validate(validateService).
				Placeholder("e.g. google, microsoft"),

			huh.NewInput().
				Key("acctid").
				Title(styles.Header.Render("account name")).
				Placeholder("e.g. @myusername, myemail@provider.com"),

			huh.NewInput().
				Key("2fakey").
				Title(styles.Header.Render("2FA key")).
				Validate(validateKey).
				Placeholder("must contain letters and numbers 2-7 only"),

			huh.NewInput().
				Key("interv").
				Title(styles.Header.Render("interval (in seconds)")).
				Validate(validateInterv).
				Placeholder("30 (default)"),

			huh.NewConfirm().
				Key("confirmation").
				Title("continue?").
				Affirmative("yes, add to list").
				Negative("no, go back"),
		),
	).WithShowHelp(false)
}
