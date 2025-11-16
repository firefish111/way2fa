package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/firefish111/way2fa/parse"

	"fmt"
	"strings"
	"time"
)

// write OTPs page.
func (m model) writeOTPs(s *strings.Builder) {
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
}

var ditto = lipgloss.NewStyle().Align(lipgloss.Center).Foreground(lipgloss.Color("33"))

// This is the generator for the actual account table.
func (m model) getTable() *table.Table {
	curr_time := uint64(time.Now().Unix())

	// get list
	otps := make([][]string, len(m.accs))
	timeout_cols := make([]string, len(m.accs))
	for i, acc := range m.accs {
		var skey string
		k, err := acc.GenKey(curr_time / uint64(acc.GetInterval()))
		if err == nil {
			skey = fmt.Sprintf("%03d %03d", k/1000, k%1000)
		} else {
			skey = " ERROR " // 7 chars
		}

		// how many seconds until next
		leftover_secs := acc.GetInterval() - (uint(curr_time) % acc.GetInterval())

		otps[i] = []string{
			acc.Name,
			acc.AcctId,
			skey,
			fmt.Sprintf("%02ds", leftover_secs)}

		// setting the colour of the time. this is done separately, but we save the colours
		switch {
		case leftover_secs < 5:
			timeout_cols[i] = "9"
		case leftover_secs < 10:
			timeout_cols[i] = "11"
		default:
			timeout_cols[i] = "10"
		}

		// peek: basically same code as above but again
		if m.peek {
			k, err = acc.GenKey(curr_time/uint64(acc.GetInterval()) + 1)
			if err == nil {
				skey = fmt.Sprintf("[%03d %03d]", k/1000, k%1000)
			} else {
				skey = " ERROR " // 7 chars
			}

			otps[i] = append(otps[i], skey)
		}
	}

	heads := []string{"service", "account", "otp", "time"}
	if m.peek {
		heads = append(heads, "next otp")
	}

	tab := table.New().
		Border(lipgloss.HiddenBorder()).
		Rows(otps...).
		Headers(heads...).
		StyleFunc(func(row, col int) lipgloss.Style { // styles cells, takes coords as args
			style := lipgloss.NewStyle().Bold(true)

			if row == -1 {
				style = style.
					Foreground(lipgloss.Color("208"))
				return style
			}

			style = style.Margin(0, 2)

			// if recently added
			if m.dirty != nil && row == *m.dirty {
				style = style.Underline(true)
			}

			switch col {
			case 0: // service name
				if otps[row][col] == "" {
					style = style.Foreground(lipgloss.Color("96")).Bold(false).Transform(func (_ string) string {
						return "<no name>"
					})
				} else if row > 0 && otps[row][col] == otps[row-1][col] {
					style = style.Foreground(lipgloss.Color("6")).Bold(false)
				} else {
					style = style.Foreground(lipgloss.Color("14"))
				}
			case 1: // account id
				style = style.Foreground(lipgloss.Color("69"))
				if len(otps[row][col]) > 0 && !strings.Contains(otps[row][col], "@") {
					style = style.BorderLeft(true).
						BorderLeftForeground(lipgloss.Color("25")).
						BorderStyle(lipgloss.Border{ Left: "@" }) // prefix at is a border (jank)
				}
			case 2: // code
				style = style.Foreground(lipgloss.Color("15"))
			case 3: // time
				style = style.Foreground(lipgloss.Color(timeout_cols[row])).
					Padding(0, 1)
			case 4: // peek
				style = style.Foreground(lipgloss.Color("218")).
					Faint(true).
					Bold(false)
			}

			return style
		})

	return tab
}
