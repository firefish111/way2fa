// Package containing code for the password prompt.
// This is needed every time some a major operation is performed to the data file, such that it needs to be briefly decrypted.
//
// This is a submodel, which is all handled by bubblon
package password

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/firefish111/way2fa/internal/ui/common/styles"
	"github.com/firefish111/way2fa/parse"
	"github.com/firefish111/way2fa/parse/cryptor"
)

// tries = how many attempts have failed. starts out at 0
type passwordModel struct {
	helpModel    help.Model             // the renderer. i can use self as keymap
	helpDB       map[string]key.Binding // to pick and choose which helps to use and when
	acclist      parse.AccountList
	field        textinput.Model
	warningOnly  bool // whether to only show a warning and no password
	tries        uint
	prev         *cryptor.PasswordHash // previous password prompt, nil if only on first try
	prevRendered string                // ditto, but a rendered string (for prompt)

	supplMsg string // supplementary message, to be shown beneath the password prompt
}

// TODO: move to format??
const (
	PasswordMaxLen     int  = 64 // len returns an int... good to know ig
	PasswordTriesCount uint = 3
)

// Returns a password prompt model.
// Takes in an AccountList, as it wishes to decrypt it.
// NOTE: AccountList should itself be a pointer, so it is already mutable
//
// Can return nil, in which case no password prompt is required: this happens only if the
// list is already decrypted.
func CreatePasswordPrompt(acclist parse.AccountList) *passwordModel {
	// create a textbox
	textbox := textinput.New()
	textbox.Focus() // we want it focussed, lest all keypressees will be dropped
	textbox.CharLimit = PasswordMaxLen
	textbox.Width = PasswordMaxLen
	textbox.EchoMode = textinput.EchoPassword

	// password model
	ret := passwordModel{
		acclist:     acclist,
		field:       textbox,
		warningOnly: false,
		tries:       0,
		helpModel:   help.New(),
		helpDB:      defaultHelp(),
	}

	ret.helpModel.Styles.ShortDesc = styles.Faint
	ret.helpModel.Styles.FullDesc = styles.Faint

	if !acclist.IsPasswordProtected() { // show a warning
		ret.warningOnly = true
	} else if acclist.IsDecrypted() { // this is in an else branch, as a warning still qualifies as a prompt
		// already decrypted, so no point in having a prompt.
		return nil
	}

	return &ret
}
