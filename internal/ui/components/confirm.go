package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/whchoi98/ecs9s/internal/theme"
)

type ConfirmResultMsg struct {
	Confirmed bool
}

type Confirm struct {
	Visible bool
	Title   string
	Message string
	focused int // 0=yes, 1=no
	thm     theme.Theme
	width   int
	height  int
}

func NewConfirm(thm theme.Theme) Confirm {
	return Confirm{thm: thm, focused: 1}
}

func (c *Confirm) Show(title, message string) {
	c.Visible = true
	c.Title = title
	c.Message = message
	c.focused = 1 // default to "No" for safety
}

func (c *Confirm) Hide() {
	c.Visible = false
}

func (c *Confirm) SetSize(w, h int) {
	c.width = w
	c.height = h
}

func (c Confirm) Update(msg tea.Msg) (Confirm, tea.Cmd) {
	if !c.Visible {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			c.focused = 0
		case "right", "l":
			c.focused = 1
		case "enter":
			c.Visible = false
			confirmed := c.focused == 0
			return c, func() tea.Msg {
				return ConfirmResultMsg{Confirmed: confirmed}
			}
		case "esc", "n":
			c.Visible = false
			return c, func() tea.Msg {
				return ConfirmResultMsg{Confirmed: false}
			}
		case "y":
			c.Visible = false
			return c, func() tea.Msg {
				return ConfirmResultMsg{Confirmed: true}
			}
		}
	}
	return c, nil
}

func (c Confirm) View() string {
	if !c.Visible {
		return ""
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(c.thm.WarnColor).
		MarginBottom(1).
		Render(c.Title)

	msg := lipgloss.NewStyle().
		Foreground(c.thm.FgPrimary).
		MarginBottom(1).
		Render(c.Message)

	btnActive := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 3).
		Background(c.thm.Accent).
		Foreground(c.thm.BgPrimary)

	btnInactive := lipgloss.NewStyle().
		Padding(0, 3).
		Foreground(c.thm.FgMuted)

	var yes, no string
	if c.focused == 0 {
		yes = btnActive.Render("Yes")
		no = btnInactive.Render("No")
	} else {
		yes = btnInactive.Render("Yes")
		no = btnActive.Render("No")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, yes, "  ", no)

	content := lipgloss.JoinVertical(lipgloss.Left, title, msg, "", buttons)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.thm.WarnColor).
		Background(c.thm.BgSecond).
		Padding(1, 3).
		Width(50)

	return lipgloss.Place(c.width, c.height,
		lipgloss.Center, lipgloss.Center,
		box.Render(content))
}
