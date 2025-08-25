package ui

import (
	"github.com/firefish111/way2fa/account"
	"github.com/firefish111/way2fa/parse"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

const (
	ReportedAppName = "way2fa"
)

// The current state of the app. This implements tea.Model
//
// - reader: where to fetch the `Account`s from
// - curr: the current list of accounts, in memory
// - peek: whether peek mode is active
type model struct {
	helpModel help.Model // the renderer. i can use self as keymap
	helpDB map[string]key.Binding // to pick and choose which helps to use and when
	reader parse.AccountList
	accs   []account.Account
	peek   bool
}

func Create(list parse.AccountList) (model, error) {
	acclist, err := list.GetAccs()
	if err != nil {
		return model{}, err
	}

	ret := model{
		helpModel: help.New(),
		helpDB: map[string]key.Binding{
			"new": key.NewBinding(
				key.WithKeys("new"),
				key.WithHelp(none.Render("n"), "new TOTP"),
			),
			"peektrue": key.NewBinding(
				key.WithKeys("p"),
				key.WithHelp(on.Render("p"), "unpeek"),
			),
			"peekfalse": key.NewBinding(
				key.WithKeys("p"),
				key.WithHelp(off.Render("p"), "peek"),
			),
			"quit": key.NewBinding(
				key.WithKeys("q", "esc"),
				key.WithHelp(none.Render("q"), "quit"),
			) },
		reader: list,
		accs:   acclist,
		peek:   false,
	}

	ret.helpModel.Styles.ShortDesc = faint
	ret.helpModel.Styles.FullDesc = faint

	return ret, nil
}
