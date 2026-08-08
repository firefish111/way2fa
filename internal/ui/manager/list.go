package manager

import (
	"reflect"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/firefish111/way2fa/parse"
)

const (
	wayEntryValue  = ".way"
	pureEntryValue = "pure"
)

// Generate table listing all available stores
func MakeListTable(selected *int) string {
	possibilities := GetPossibilities()
	rows := make([][]string, len(possibilities))
	for i, store := range possibilities {
		rows[i] = []string{
			strconv.Itoa(i + 1),
			"",
			reflect.TypeOf(store).Elem().Name(),
			"",
		}
		if way, ok := store.(parse.WayAccountList); ok {
			rows[i][1] = wayEntryValue
			rows[i][3] = strconv.Itoa(int(way.GetWayTypeId()))
		} else {
			rows[i][1] = pureEntryValue
		}
	}

	tab := table.New().
		Border(lipgloss.HiddenBorder()).
		Rows(rows...).
		Headers("prio", "type", "name", "id #").
		StyleFunc(func(row, col int) lipgloss.Style { // styles cells, takes coords as args
			style := lipgloss.NewStyle()

			if row == -1 {
				return style.Foreground(lipgloss.Color("208")).Bold(true)
			}

			if col == 0 {
				style = style.PaddingLeft(3)
			} else if col < 3 {
				style = style.PaddingRight(2)
			}

			// if selected. code copied from ../otps.go
			if selected != nil && row == *selected {
				style = style.Underline(true)

				if col == 0 { // highlighting.
					// remove paddingleft, and replace it with a smaller paddingleft and an asterisk "border".
					// see below for diagram:
					//    0
					// *1   0
					//    0
					style = style.
						PaddingLeft(2).
						BorderLeft(true).
						BorderLeftForeground(lipgloss.Color("199")).
						BorderStyle(lipgloss.Border{Left: "*"}) // prefix is a border (jank)
				}
			}

			switch col {
			case 0: // priority
				style = style.Foreground(lipgloss.Color("220")).
					Align(lipgloss.Left)
				// FIX: there's a bug where align right will cut off the last character if borderleft present, even if padding is managed
			case 1: // type
				switch rows[row][col] {
				case wayEntryValue:
					style = style.Foreground(lipgloss.Color("81")).Bold(true)
				case pureEntryValue:
					style = style.Foreground(lipgloss.Color("196")).Bold(true)
				default:
				}
			case 2: // name
				style = style.Foreground(lipgloss.Color("253")).Italic(true)
			case 3: // way id
				style = style.Align(lipgloss.Right).
					Foreground(lipgloss.Color("149"))
			}

			return style
		})

	return tab.Render()
}
