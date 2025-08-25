package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/firefish111/way2fa/parse"

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
		switch event.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit // bye bye
		case "p":
			m.peek = !m.peek
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

// Spit it out
func (m model) View() string {
	var s strings.Builder

	s.WriteRune('\n')

	s.WriteString(
		app_name.Render("way2fa") +
			faint.Render(" - My TOTPs: "))

	{
		srct, srcs := m.reader.GetSource()

		style := source
		if srct == parse.FileSource {
			style = style.Background(lipgloss.Color("22"))
		}

		s.WriteString(style.Render(srcs))
	}

	s.WriteRune('\n')

	// separated the table generation code away from here, see ./table.go
	s.WriteString(marg.Render(m.getTable().String()))

	s.WriteRune('\n')
	s.WriteRune(' ')

	helpview := m.helpModel.View(m) // using self as a help model to acces internal state
	s.WriteString(helpview)

	s.WriteRune('\n')

	return s.String()
}
