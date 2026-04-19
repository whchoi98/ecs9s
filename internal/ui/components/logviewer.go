package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/whchoi98/ecs9s/internal/theme"
)

type LogLine struct {
	Timestamp string
	Message   string
}

type LogViewer struct {
	viewport viewport.Model
	lines    []LogLine
	filter   string
	follow   bool
	thm      theme.Theme
}

func NewLogViewer(thm theme.Theme) LogViewer {
	vp := viewport.New(80, 20)
	return LogViewer{
		viewport: vp,
		follow:   true,
		thm:      thm,
	}
}

func (l *LogViewer) SetSize(w, h int) {
	l.viewport.Width = w
	l.viewport.Height = h
	l.renderContent()
}

func (l *LogViewer) AppendLines(lines []LogLine) {
	l.lines = append(l.lines, lines...)
	l.renderContent()
	if l.follow {
		l.viewport.GotoBottom()
	}
}

func (l *LogViewer) SetFilter(f string) {
	l.filter = f
	l.renderContent()
}

func (l *LogViewer) ToggleFollow() {
	l.follow = !l.follow
	if l.follow {
		l.viewport.GotoBottom()
	}
}

func (l *LogViewer) Clear() {
	l.lines = nil
	l.renderContent()
}

func (l *LogViewer) renderContent() {
	tsStyle := lipgloss.NewStyle().
		Foreground(l.thm.FgMuted)

	msgStyle := lipgloss.NewStyle().
		Foreground(l.thm.FgPrimary)

	var sb strings.Builder
	lf := strings.ToLower(l.filter)

	for _, line := range l.lines {
		if lf != "" {
			if !strings.Contains(strings.ToLower(line.Message), lf) {
				continue
			}
		}
		ts := tsStyle.Render(line.Timestamp)
		msg := msgStyle.Render(line.Message)
		fmt.Fprintf(&sb, "%s %s\n", ts, msg)
	}

	l.viewport.SetContent(sb.String())
}

func (l LogViewer) Update(msg tea.Msg) (LogViewer, tea.Cmd) {
	var cmd tea.Cmd
	l.viewport, cmd = l.viewport.Update(msg)
	return l, cmd
}

func (l LogViewer) View() string {
	header := lipgloss.NewStyle().
		Foreground(l.thm.FgMuted).
		Render(fmt.Sprintf(" Lines: %d | Follow: %v | Filter: %q",
			len(l.lines), l.follow, l.filter))

	return header + "\n" + l.viewport.View()
}
