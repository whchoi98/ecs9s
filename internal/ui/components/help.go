package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/whchoi98/ecs9s/internal/theme"
)

type HelpBinding struct {
	Key  string
	Desc string
}

type Help struct {
	Visible  bool
	Bindings []HelpBinding
	Extra    []HelpBinding // page-specific bindings
	thm      theme.Theme
	width    int
	height   int
}

func NewHelp(thm theme.Theme) Help {
	return Help{
		thm: thm,
		Bindings: []HelpBinding{
			{"Tab/[/]", "Switch tabs"},
			{"Enter", "Drill down"},
			{"Esc/Backspace", "Go back"},
			{":", "Command mode"},
			{"/", "Filter"},
			{"s", "Sort"},
			{"R", "Refresh"},
			{"p", "Switch profile"},
			{"r", "Switch region"},
			{"j/k ↑/↓", "Navigate rows"},
			{"?", "Toggle help"},
			{"q", "Quit"},
		},
	}
}

func (h *Help) SetSize(w, hh int) {
	h.width = w
	h.height = hh
}

func (h *Help) Toggle() {
	h.Visible = !h.Visible
}

func (h *Help) SetExtra(bindings []HelpBinding) {
	h.Extra = bindings
}

func (h Help) View() string {
	if !h.Visible {
		return ""
	}

	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(h.thm.Accent).
		Width(18).
		Align(lipgloss.Right).
		PaddingRight(2)

	descStyle := lipgloss.NewStyle().
		Foreground(h.thm.FgPrimary)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(h.thm.Accent).
		Underline(true).
		MarginBottom(1)

	var lines []string
	lines = append(lines, titleStyle.Render("Keybindings"))
	lines = append(lines, "")

	all := append(h.Bindings, h.Extra...)
	for _, b := range all {
		lines = append(lines, keyStyle.Render(b.Key)+descStyle.Render(b.Desc))
	}

	content := strings.Join(lines, "\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(h.thm.Accent).
		Background(h.thm.BgSecond).
		Padding(1, 3).
		Width(50)

	return lipgloss.Place(h.width, h.height,
		lipgloss.Center, lipgloss.Center,
		box.Render(content))
}
