package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/whchoi98/ecs9s/internal/theme"
)

type CommandExecuteMsg struct {
	Command string
}

type CommandBar struct {
	input   textinput.Model
	active  bool
	mode    string // "command" or "filter"
	thm     theme.Theme
	width   int
}

func NewCommandBar(thm theme.Theme) CommandBar {
	ti := textinput.New()
	ti.Prompt = ":"
	ti.CharLimit = 100
	return CommandBar{input: ti, thm: thm}
}

func (c *CommandBar) SetWidth(w int) { c.width = w }

func (c *CommandBar) Activate(mode string) {
	c.active = true
	c.mode = mode
	c.input.Reset()
	if mode == "filter" {
		c.input.Prompt = "/"
	} else {
		c.input.Prompt = ":"
	}
	c.input.Focus()
}

func (c *CommandBar) Deactivate() {
	c.active = false
	c.input.Blur()
	c.input.Reset()
}

func (c CommandBar) Active() bool { return c.active }
func (c CommandBar) Mode() string { return c.mode }
func (c CommandBar) Value() string { return c.input.Value() }

func (c CommandBar) Update(msg tea.Msg) (CommandBar, tea.Cmd) {
	if !c.active {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			val := strings.TrimSpace(c.input.Value())
			c.Deactivate()
			if val != "" {
				return c, func() tea.Msg {
					return CommandExecuteMsg{Command: val}
				}
			}
			return c, nil
		case "esc":
			c.Deactivate()
			return c, nil
		}
	}

	var cmd tea.Cmd
	c.input, cmd = c.input.Update(msg)
	return c, cmd
}

func (c CommandBar) View() string {
	if c.active {
		return c.input.View()
	}

	hints := lipgloss.NewStyle().
		Foreground(c.thm.FgMuted).
		Width(c.width).
		Padding(0, 1)

	return hints.Render(":command  /filter  ?help  q:quit")
}
