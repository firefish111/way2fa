// Package containing custom messages implementing tea.Msg
package msgs

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/firefish111/way2fa/account"
	"time"
)

// tea's Msg is actually an empty interface, so you can pass anything you want to Update.
type TickMsg time.Time

// Our Tick function. Every single submodel MUST acknowledge this
func Tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { // callback
		return TickMsg(t)
	})
}

// Signal signifying a new account.
// TODO: add other properties, such as where to put it
type NewAccMsg struct {
	Acct account.Account
}

func SendAcct(acct account.Account) tea.Cmd {
	return func() tea.Msg {
		return NewAccMsg{
			Acct: acct,
		}
	}
}
