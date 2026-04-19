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

type iamLoadedMsg struct {
	roles []aws.Role
	err   error
}

type IAMPage struct {
	table   components.Table
	roles   []aws.Role
	iam     *aws.IAMClient
	loading bool
	err     error
}

func NewIAMPage(iam *aws.IAMClient, thm theme.Theme) IAMPage {
	cols := []table.Column{
		{Title: "ROLE NAME", Width: 35},
		{Title: "PATH", Width: 20},
		{Title: "CREATED", Width: 12},
		{Title: "ARN", Width: 50},
	}
	return IAMPage{table: components.NewTable(cols, thm), iam: iam}
}

func (p *IAMPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p IAMPage) Init() tea.Cmd     { return p.fetchData() }

func (p IAMPage) Update(msg tea.Msg) (IAMPage, tea.Cmd) {
	switch msg := msg.(type) {
	case iamLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.roles = msg.roles
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

func (p IAMPage) View() string {
	if p.loading {
		return "  Loading IAM roles..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	return p.table.View()
}

func (p IAMPage) HelpBindings() []components.HelpBinding { return nil }

func (p *IAMPage) updateRows() {
	var rows []table.Row
	for _, r := range p.roles {
		rows = append(rows, table.Row{r.Name, r.Path, r.CreateDate, r.ARN})
	}
	p.table.SetRows(rows)
}

func (p IAMPage) fetchData() tea.Cmd {
	iam := p.iam
	return func() tea.Msg {
		roles, err := iam.ListRoles(context.Background(), "")
		return iamLoadedMsg{roles: roles, err: err}
	}
}

func (p *IAMPage) Refresh() tea.Cmd { return p.fetchData() }
