package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/whchoi98/ecs9s/internal/theme"
)

type Styles struct {
	Theme theme.Theme

	TabActive   lipgloss.Style
	TabInactive lipgloss.Style
	TabBar      lipgloss.Style

	StatusBar lipgloss.Style

	TableHeader   lipgloss.Style
	TableRow      lipgloss.Style
	TableSelected lipgloss.Style

	CommandBar lipgloss.Style

	HelpKey  lipgloss.Style
	HelpDesc lipgloss.Style

	Title   lipgloss.Style
	Error   lipgloss.Style
	Warning lipgloss.Style
	Success lipgloss.Style
	Muted   lipgloss.Style

	DialogBox    lipgloss.Style
	DialogButton lipgloss.Style
}

func New(t theme.Theme) Styles {
	return Styles{
		Theme: t,
		TabActive: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.TabActive).
			Padding(0, 2).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(t.TabActive),
		TabInactive: lipgloss.NewStyle().
			Foreground(t.TabInact).
			Padding(0, 2),
		TabBar: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(t.Border),
		StatusBar: lipgloss.NewStyle().
			Background(t.StatusBg).
			Foreground(t.StatusFg).
			Padding(0, 1),
		TableHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent).
			Padding(0, 1),
		TableRow: lipgloss.NewStyle().
			Foreground(t.FgPrimary).
			Padding(0, 1),
		TableSelected: lipgloss.NewStyle().
			Bold(true).
			Background(t.SelectBg).
			Foreground(t.SelectFg).
			Padding(0, 1),
		CommandBar: lipgloss.NewStyle().
			Foreground(t.FgMuted).
			Padding(0, 1),
		HelpKey: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent),
		HelpDesc: lipgloss.NewStyle().
			Foreground(t.FgMuted),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent),
		Error: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.ErrColor),
		Warning: lipgloss.NewStyle().
			Foreground(t.WarnColor),
		Success: lipgloss.NewStyle().
			Foreground(t.OkColor),
		Muted: lipgloss.NewStyle().
			Foreground(t.FgMuted),
		DialogBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Accent).
			Padding(1, 2).
			Width(50),
		DialogButton: lipgloss.NewStyle().
			Padding(0, 2).
			Background(t.Accent).
			Foreground(t.BgPrimary),
	}
}
