package pages

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/whchoi98/ecs9s/internal/aws"
	"github.com/whchoi98/ecs9s/internal/theme"
	"github.com/whchoi98/ecs9s/internal/ui"
	"github.com/whchoi98/ecs9s/internal/ui/components"
)

type secretsLoadedMsg struct {
	secrets []aws.Secret
	err     error
}

type SecretsPage struct {
	table   components.Table
	secrets []aws.Secret
	sm      *aws.SecretsClient
	loading bool
	err     error
}

func NewSecretsPage(sm *aws.SecretsClient, thm theme.Theme) SecretsPage {
	cols := []table.Column{
		{Title: "NAME", Width: 30},
		{Title: "DESCRIPTION", Width: 30},
		{Title: "ROTATION", Width: 10},
		{Title: "LAST ROTATED", Width: 20},
		{Title: "LAST CHANGED", Width: 20},
	}
	return SecretsPage{table: components.NewTable(cols, thm), sm: sm}
}

func (p *SecretsPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p SecretsPage) Init() tea.Cmd     { return p.fetchData() }

func (p SecretsPage) Update(msg tea.Msg) (SecretsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case secretsLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.secrets = msg.secrets
		p.updateRows()
		return p, nil
	case tea.KeyMsg:
		if p.table.Filtering() {
			switch msg.String() {
			case "esc":
				p.table.ClearFilter()
				return p, nil
			case "enter":
				p.table.StopFilter()
				return p, nil
			default:
				val := p.table.Filter() + msg.String()
				if msg.String() == "backspace" && len(p.table.Filter()) > 0 {
					val = p.table.Filter()[:len(p.table.Filter())-1]
				}
				p.table.SetFilter(val)
				return p, nil
			}
		}
		switch msg.String() {
		case "/":
			p.table.StartFilter()
			return p, nil
		case "s":
			p.table.CycleSort()
			return p, nil
		}
	}
	var cmd tea.Cmd
	p.table, cmd = p.table.Update(msg)
	return p, cmd
}

func (p SecretsPage) View() string {
	if p.loading {
		return "  Loading secrets..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	return p.table.View()
}

func (p SecretsPage) HelpBindings() []components.HelpBinding { return nil }

func (p *SecretsPage) updateRows() {
	var rows []table.Row
	for _, s := range p.secrets {
		rotation := "No"
		if s.RotationEnabled {
			rotation = "Yes"
		}
		lastRotated := "-"
		if !s.LastRotated.IsZero() {
			lastRotated = s.LastRotated.Format("2006-01-02 15:04")
		}
		lastChanged := "-"
		if !s.LastChanged.IsZero() {
			lastChanged = s.LastChanged.Format("2006-01-02 15:04")
		}
		desc := s.Description
		if len(desc) > 28 {
			desc = desc[:28] + ".."
		}
		rows = append(rows, table.Row{
			s.Name, desc, rotation, lastRotated, lastChanged,
		})
	}
	p.table.SetRows(rows)
}

func (p SecretsPage) fetchData() tea.Cmd {
	sm := p.sm
	return func() tea.Msg {
		secrets, err := sm.ListSecrets(context.Background())
		return secretsLoadedMsg{secrets: secrets, err: err}
	}
}

func (p *SecretsPage) Refresh() tea.Cmd { return p.fetchData() }
