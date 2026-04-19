package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/whchoi98/ecs9s/internal/theme"
)

type StatusBar struct {
	Cluster string
	Profile string
	Region  string
	Error   string
	Info    string
	thm     theme.Theme
	width   int
}

func NewStatusBar(thm theme.Theme) StatusBar {
	return StatusBar{thm: thm}
}

func (s *StatusBar) SetWidth(w int) { s.width = w }

func (s *StatusBar) SetError(e string) {
	s.Error = e
	s.Info = ""
}

func (s *StatusBar) SetInfo(i string) {
	s.Info = i
	s.Error = ""
}

func (s *StatusBar) ClearMessage() {
	s.Error = ""
	s.Info = ""
}

func (s StatusBar) View() string {
	bg := lipgloss.NewStyle().
		Background(s.thm.StatusBg).
		Width(s.width)

	left := lipgloss.NewStyle().
		Foreground(s.thm.StatusFg).
		Background(s.thm.StatusBg).
		Bold(true).
		Padding(0, 1)

	right := lipgloss.NewStyle().
		Foreground(s.thm.FgMuted).
		Background(s.thm.StatusBg).
		Padding(0, 1)

	var leftText string
	if s.Cluster != "" {
		leftText = fmt.Sprintf("Cluster: %s", s.Cluster)
	} else {
		leftText = "ecs9s"
	}

	rightText := fmt.Sprintf("%s | %s", s.Profile, s.Region)

	var msg string
	if s.Error != "" {
		msg = lipgloss.NewStyle().
			Foreground(s.thm.ErrColor).
			Background(s.thm.StatusBg).
			Render(" " + s.Error)
	} else if s.Info != "" {
		msg = lipgloss.NewStyle().
			Foreground(s.thm.OkColor).
			Background(s.thm.StatusBg).
			Render(" " + s.Info)
	}

	leftRendered := left.Render(leftText) + msg
	rightRendered := right.Render(rightText)

	gap := s.width - lipgloss.Width(leftRendered) - lipgloss.Width(rightRendered)
	if gap < 0 {
		gap = 0
	}
	filler := lipgloss.NewStyle().
		Background(s.thm.StatusBg).
		Render(fmt.Sprintf("%*s", gap, ""))

	return bg.Render(leftRendered + filler + rightRendered)
}
