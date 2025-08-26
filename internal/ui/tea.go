package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"strings"
	"time"
)

// tea's Msg is actually an empty interface, so you can pass anything you want to Update.
type TickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { // callback
		return TickMsg(t)
	})
}

// XXX: this could be important in the future
func (m model) Init() tea.Cmd {
	return tick() // first tick.
}

// this returns the model itself, and anything we want tea to do
func (m model) Update(event tea.Msg) (tea.Model, tea.Cmd) {
	// event is whatever tea wants us to respond, we need to see what it is
	switch event := event.(type) {
	case tea.KeyMsg: // handle keypress
		if !m.create {
			switch event.String() {
			case "ctrl+c", "q":
				return m, tea.Quit // bye bye
			case "n":
				m.create = true
			case "p":
				m.peek = !m.peek
			}
		} else {
			switch event.String() {
			case "esc":
				m.create = false
			}
		}
	case TickMsg: // our own custom tick message struct (just a typedef)
		return m, tick() // tick again. this will be executed, and after it times out, update will be called again
	}

	return m, nil
}

var source = lipgloss.NewStyle().
	Bold(true).
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("159")).
	Background(lipgloss.Color("88")).
	Padding(0, 1).
	Margin(0, 1)

var app_name = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("220")).
	PaddingLeft(1)

var faint = lipgloss.NewStyle().
	Faint(true).
	Foreground(lipgloss.Color("242"))

var marg = lipgloss.NewStyle().
	MarginLeft(4)

var wip = lipgloss.NewStyle().
	Margin(1, 2).
	Padding(1, 2).
	Background(lipgloss.Color("239")).
	Foreground(lipgloss.Color("15")).
	Bold(true)

// Spit it out
func (m model) View() string {
	var s strings.Builder

	s.WriteRune('\n')

	switch {
	case m.create:
		m.writeCreate(&s)
	default:
		m.writeOTPs(&s)
	}

	s.WriteRune('\n')
	s.WriteRune(' ')

	helpview := m.helpModel.View(m) // using self as a help model to acces internal state
	s.WriteString(helpview)

	s.WriteRune('\n')

	return s.String()
}
