package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/firefish111/way2fa/account"
	"github.com/firefish111/way2fa/parse"
)

const (
	ReportedAppName = "way2fa"
)

// The current state of the app. This implements tea.Model
// - helpModel: the help renderer
// - helpDB:	the list of help bindings
// - reader:	where to fetch the `Account`s from
// - accs: 		the current list of accounts, in memory
// - peek: 		whether peek mode is active
// - create: 	whether create mode is active
type model struct {
	helpModel help.Model             // the renderer. i can use self as keymap
	helpDB    map[string]key.Binding // to pick and choose which helps to use and when
	reader    parse.AccountList
	accs      []account.Account
	peek      bool // is in peek mode
	create    bool // is in new account mode
}

func Create(list parse.AccountList) (model, error) {
	acclist, err := list.GetAccs()
	if err != nil {
		return model{}, err
	}

	ret := model{
		helpModel: help.New(),
		helpDB:    defaultHelp(),
		reader:    list,
		accs:      acclist,
		peek:      false,
		create:    false,
	}

	ret.helpModel.Styles.ShortDesc = faint
	ret.helpModel.Styles.FullDesc = faint

	return ret, nil
}
