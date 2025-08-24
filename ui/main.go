package ui

import (
	"github.com/firefish111/way2fa/account"
	"github.com/firefish111/way2fa/parse"
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
	reader parse.AccountList
	accs   []account.Account
	peek   bool
}

func Create(list parse.AccountList) (model, error) {
	acclist, err := list.GetAccs()
	return model{
		reader: list,
		accs:   acclist,
		peek:   false,
	}, err
}
