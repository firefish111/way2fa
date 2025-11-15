package creation

import (
	"github.com/charmbracelet/huh"
)

func blankForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Service name").
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
				Title("Create?").
				Affirmative("Yes").
				Negative("No"),
		),

		huh.NewGroup(
			huh.NewNote().
				Title("Confirmation").
				DescriptionFunc(func() string {
					return "This is your code: 123 456 [30s].\nPlease double check before confirming."
				}, nil), // this nil is whatever value the above closure captures (because go closures are a bit bad)

			huh.NewConfirm().
				Title("Is this correct?").
				Affirmative("Yes, add to list").
				Negative("No, go back"),
		),
	)
}
