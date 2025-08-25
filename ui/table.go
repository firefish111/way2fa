package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"fmt"
	"time"
)

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
			fmt.Sprintf("%02ds", leftover_secs) }

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

			switch col {
			case 0: // username
				style = style.Foreground(lipgloss.Color("14"))
			case 1: // account id
				style = style.Foreground(lipgloss.Color("12"))
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
