package ui

import (
	"strings"
)

// write Create page.
func (m model) writeCreate(s *strings.Builder) {
	s.WriteString(
		app_name.Render("way2fa") +
			faint.Render(" - New TOTP"))

	s.WriteRune('\n')

	s.WriteString(wip.Render("WIP, press esc to go back"))
}
