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

type taskDefLoadedMsg struct {
	defs []aws.TaskDefinition
	err  error
}

type TaskDefPage struct {
	table   components.Table
	defs    []aws.TaskDefinition
	ecs     *aws.ECSClient
	loading bool
	err     error
}

func NewTaskDefPage(ecs *aws.ECSClient, thm theme.Theme) TaskDefPage {
	cols := []table.Column{
		{Title: "FAMILY", Width: 25},
		{Title: "REVISION", Width: 10},
		{Title: "STATUS", Width: 10},
		{Title: "CPU", Width: 8},
		{Title: "MEMORY", Width: 8},
		{Title: "COMPAT", Width: 15},
	}
	return TaskDefPage{table: components.NewTable(cols, thm), ecs: ecs}
}

func (p *TaskDefPage) SetSize(w, h int) { p.table.SetSize(w, h) }

func (p TaskDefPage) Init() tea.Cmd { return p.fetchData() }

func (p TaskDefPage) Update(msg tea.Msg) (TaskDefPage, tea.Cmd) {
	switch msg := msg.(type) {
	case taskDefLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.defs = msg.defs
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
		case "ctrl+d":
			if idx := p.table.Cursor(); idx < len(p.defs) {
				d := p.defs[idx]
				ecs := p.ecs
				return p, func() tea.Msg {
					err := ecs.DeregisterTaskDefinition(context.Background(), d.ARN)
					if err != nil {
						return ui.ErrorMsg{Err: err}
					}
					return ui.InfoMsg{Text: fmt.Sprintf("Deregistered: %s:%d", d.Family, d.Revision)}
				}
			}
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

func (p TaskDefPage) View() string {
	if p.loading {
		return "  Loading task definitions..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	return p.table.View()
}

func (p TaskDefPage) HelpBindings() []components.HelpBinding {
	return []components.HelpBinding{
		{Key: "Ctrl+d", Desc: "Deregister task definition"},
	}
}

func (p *TaskDefPage) updateRows() {
	var rows []table.Row
	for _, d := range p.defs {
		rows = append(rows, table.Row{
			d.Family, fmt.Sprintf("%d", d.Revision), d.Status,
			d.CPU, d.Memory, d.Compatibility,
		})
	}
	p.table.SetRows(rows)
}

func (p TaskDefPage) fetchData() tea.Cmd {
	ecs := p.ecs
	return func() tea.Msg {
		defs, err := ecs.ListTaskDefinitions(context.Background())
		return taskDefLoadedMsg{defs: defs, err: err}
	}
}

func (p *TaskDefPage) Refresh() tea.Cmd { return p.fetchData() }
