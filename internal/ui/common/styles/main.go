// Contains style variables and routines common to all ui models.
package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/firefish111/way2fa/parse"
)

// used for source location
var Source = lipgloss.NewStyle().
	Bold(true).
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("159")).
	Background(lipgloss.Color("88")).
	Padding(0, 1).
	Margin(0, 1)

func RenderSource(which parse.DataSource, name string) string {
	style := Source
	if which == parse.FileSource {
		style = style.Background(lipgloss.Color("22"))
	}

	return style.Render(name)
}

// the name of this app, i.e. way2fa
var AppName = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("220")).
	PaddingLeft(1)

// used for the confirmation prompt
var SrvcName = lipgloss.NewStyle().
	Italic(true).
	Foreground(lipgloss.Color("14"))

var None = lipgloss.NewStyle().Foreground(lipgloss.Color("251"))
var PeekOn = lipgloss.NewStyle().Foreground(lipgloss.Color("218")).Bold(true)
var NewOn = lipgloss.NewStyle().Foreground(lipgloss.Color("112")).Bold(true)
var Off = lipgloss.NewStyle().Foreground(lipgloss.Color("251"))

var Faint = lipgloss.NewStyle().
	Faint(true).
	Foreground(lipgloss.Color("242"))

var BigIndent = lipgloss.NewStyle().
	MarginLeft(4)

var SidePad = lipgloss.NewStyle().
	Margin(0, 1)

var Wip = lipgloss.NewStyle().
	Margin(1, 2).
	Padding(1, 2).
	Background(lipgloss.Color("239")).
	Foreground(lipgloss.Color("15")).
	Bold(true)

var Error = Wip.
	Background(lipgloss.Color("53"))

var Spaced = lipgloss.NewStyle().
	Margin(1)

var Box = lipgloss.NewStyle().
	BorderForeground(lipgloss.Color("6")).
	//	Align(lipgloss.Center).
	//	Border(lipgloss.DoubleBorder()).
	Padding(1)

var Title = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("157"))

var Supplement = lipgloss.NewStyle().
	Italic(true).
	Padding(2, 4).
	Foreground(lipgloss.Color("9"))

var Header = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("208"))
