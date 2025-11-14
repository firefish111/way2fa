// Package containing code for the "create new" form page.
// This is a fairly thin wrapper around a huh? form, so that
// some extra details can be added (namely exit code, and a
// "this is your code, please double check" prompt)
//
// This is a submodel, all the specifics of which are handled by bubblon
package creation

import (
	"github.com/charmbracelet/huh"
)

// The form itself
type formModel struct {
	form *huh.Form
	done bool
}

func DefaultForm() formModel {
	return formModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Service name"),

				huh.NewInput().
					Title("Account name"),

				huh.NewInput().
					Title("2FA Key"),

				huh.NewInput().
					Title("Interval").
					Placeholder("Default: 30"),

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
		),
	}
}

func blankForm() {
}

/*
// write Create page.
func (m model) writeCreate(s *strings.Builder) {
	s.WriteString(
		app_name.Render("way2fa") +
			faint.Render(" - New TOTP"))

	s.WriteRune('\n')

	s.WriteString(m.createForm.View())

	s.WriteString(wip.Render("WIP, press esc to go back"))
}
*/
