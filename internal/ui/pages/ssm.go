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

type ssmLoadedMsg struct {
	params []aws.Parameter
	err    error
}

type SSMPage struct {
	table   components.Table
	params  []aws.Parameter
	ssm     *aws.SSMClient
	prefix  string
	loading bool
	err     error
}

func NewSSMPage(ssm *aws.SSMClient, thm theme.Theme) SSMPage {
	cols := []table.Column{
		{Title: "NAME", Width: 45},
		{Title: "TYPE", Width: 15},
		{Title: "VALUE", Width: 30},
		{Title: "VERSION", Width: 8},
		{Title: "LAST MODIFIED", Width: 20},
	}
	return SSMPage{table: components.NewTable(cols, thm), ssm: ssm, prefix: "/"}
}

func (p *SSMPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p SSMPage) Init() tea.Cmd     { return p.fetchData() }

func (p SSMPage) Update(msg tea.Msg) (SSMPage, tea.Cmd) {
	switch msg := msg.(type) {
	case ssmLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.params = msg.params
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

func (p SSMPage) View() string {
	if p.loading {
		return "  Loading SSM parameters..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	if len(p.params) == 0 {
		return "  No parameters found at path: " + p.prefix
	}
	return p.table.View()
}

func (p SSMPage) HelpBindings() []components.HelpBinding { return nil }

func (p *SSMPage) updateRows() {
	var rows []table.Row
	for _, param := range p.params {
		val := param.Value
		if param.Type == "SecureString" {
			val = "****"
		}
		if len(val) > 28 {
			val = val[:28] + ".."
		}
		rows = append(rows, table.Row{
			param.Name, param.Type, val,
			fmt.Sprintf("%d", param.Version),
			param.LastModified.Format("2006-01-02 15:04"),
		})
	}
	p.table.SetRows(rows)
}

func (p SSMPage) fetchData() tea.Cmd {
	ssm := p.ssm
	prefix := p.prefix
	return func() tea.Msg {
		params, err := ssm.ListParameters(context.Background(), prefix)
		return ssmLoadedMsg{params: params, err: err}
	}
}

func (p *SSMPage) Refresh() tea.Cmd { return p.fetchData() }
