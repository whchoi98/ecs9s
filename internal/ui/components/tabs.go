package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/whchoi98/ecs9s/internal/theme"
)

type TabItem struct {
	Name    string
	Command string // for command mode (:cluster, :service, etc.)
}

type Tabs struct {
	Items  []TabItem
	Active int
	thm    theme.Theme
	width  int
}

func NewTabs(items []TabItem, thm theme.Theme) Tabs {
	return Tabs{Items: items, thm: thm}
}

func (t *Tabs) SetWidth(w int) {
	t.width = w
}

func (t *Tabs) Next() {
	t.Active = (t.Active + 1) % len(t.Items)
}

func (t *Tabs) Prev() {
	t.Active = (t.Active - 1 + len(t.Items)) % len(t.Items)
}

func (t *Tabs) SetActive(idx int) {
	if idx >= 0 && idx < len(t.Items) {
		t.Active = idx
	}
}

func (t *Tabs) SetActiveByCommand(cmd string) bool {
	for i, item := range t.Items {
		if item.Command == cmd {
			t.Active = i
			return true
		}
	}
	return false
}

func (t Tabs) ActiveItem() TabItem {
	return t.Items[t.Active]
}

func (t Tabs) View() string {
	activeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.thm.TabActive).
		Background(t.thm.BgSecond).
		Padding(0, 2)

	inactiveStyle := lipgloss.NewStyle().
		Foreground(t.thm.TabInact).
		Padding(0, 1)

	var tabs []string
	for i, item := range t.Items {
		if i == t.Active {
			tabs = append(tabs, activeStyle.Render(item.Name))
		} else {
			tabs = append(tabs, inactiveStyle.Render(item.Name))
		}
	}

	row := strings.Join(tabs, " ")
	bar := lipgloss.NewStyle().
		Width(t.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(t.thm.Border)

	return bar.Render(row)
}
