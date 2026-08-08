// Package containing code for the "create new" form page.
// This is a fairly thin wrapper around a huh? form, so that
// some extra details can be added (namely exit code, and a
// "this is your code, please double check" prompt)
//
// This is a submodel, all the specifics of which are handled by bubblon
package creation

import (
	"charm.land/huh/v2"
)

// The form itself
type formModel struct {
	form *huh.Form
	done bool
}

func DefaultForm() *formModel {
	m := new(formModel)
	m.resetForm() // form was never initialised

	return m
}
