package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/firefish111/way2fa/parse"

	"fmt"
	"strings"
	"time"
)

// tea's Msg is actually an empty interface, that simply gets passed on to Update.
type TickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { // callback
		return TickMsg(t)
	})
}

// TODO: this could be important in the future
func (m model) Init() tea.Cmd {
	return tick() // first tick.
}

// this returns the model itself, and anything we want tea to do
func (m model) Update(event tea.Msg) (tea.Model, tea.Cmd) {
	// event is whatever tea wants us to respond, we need to see what it is
	switch event := event.(type) {
	case tea.KeyMsg: // handle keypress
		switch event.String() {
		case "ctrl+c", "q":
			return m, tea.Quit // bye bye
		case "p":
			m.peek = !m.peek
		}
	case TickMsg: // outr own custom tick
		return m, tick() // tick again. this will be executed, and after it times out, update will be called again
	}

	return m, nil
}

var source = lipgloss.NewStyle().
	Bold(true).
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("159")).
	Background(lipgloss.Color("1")).
	Padding(0, 1).
	Margin(0, 1)

var app_name = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("220")).
	PaddingLeft(1)

var faint = lipgloss.NewStyle().
	Faint(true).
	Foreground(lipgloss.Color("242"))

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

	curr_time := uint64(time.Now().Unix())

	otps := make([][]string, len(m.accs))
	timeout_cols := make([]string, len(m.accs))
	for i, acc := range m.accs {
		var skey string
		k, err := acc.GenKey(curr_time / uint64(acc.GetInterval()))
		if err == nil {
			skey = fmt.Sprintf("%03d %03d", k/1000, k%1000)
		} else {
			skey = " ERROR "
		}

		id := ""
		if acc.AcctId != "" {
			id = fmt.Sprintf("<@%s>", acc.AcctId)
		}

		// how many seconds until next
		leftover_secs := acc.GetInterval() - (uint(curr_time) % acc.GetInterval())

		otps[i] = []string{
			acc.Name,
			id,
			skey,
			fmt.Sprintf("%02ds", leftover_secs)}

		switch {
		case leftover_secs < 5:
			timeout_cols[i] = "9"
		case leftover_secs < 10:
			timeout_cols[i] = "11"
		default:
			timeout_cols[i] = "10"
		}

		if m.peek {
			k, err = acc.GenKey(curr_time/uint64(acc.GetInterval()) + 1)
			if err == nil {
				skey = fmt.Sprintf("%03d %03d", k/1000, k%1000)
			} else {
				skey = " ERROR "
			}

			otps[i] = append(otps[i], skey)
		}
	}

	tab := table.New().
		Border(lipgloss.HiddenBorder()).
		Rows(otps...).
		StyleFunc(func(row, col int) lipgloss.Style { // styler
			style := lipgloss.NewStyle().Bold(true)

			switch col {
			case 0: // username
				style = style.Foreground(lipgloss.Color("14"))
			case 1: // account id
				style = style.Foreground(lipgloss.Color("12"))
			case 2: // code
				style = style.Foreground(lipgloss.Color("15"))
			case 3: // time
				style = style.Foreground(lipgloss.Color(timeout_cols[row])).
					Padding(0, 2)
			case 4: // peek
				style = style.Foreground(lipgloss.Color("230")).
					UnsetBold().
					Italic(true)
			}

			return style
		})

	s.WriteString(tab.String())

	return s.String()
}
