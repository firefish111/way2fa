// Package containing custom messages implementing tea.Msg
package msgs

import (
	tea "github.com/charmbracelet/bubbletea"

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
