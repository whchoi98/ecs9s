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

type stoppedLoadedMsg struct {
	tasks []aws.StoppedTask
	err   error
}

type StoppedPage struct {
	table   components.Table
	tasks   []aws.StoppedTask
	ecs     *aws.ECSClient
	nav     ui.NavContext
	loading bool
	err     error
}

func NewStoppedPage(ecs *aws.ECSClient, thm theme.Theme) StoppedPage {
	cols := []table.Column{
		{Title: "TASK ID", Width: 20},
		{Title: "STATUS", Width: 10},
		{Title: "TASKDEF", Width: 22},
		{Title: "GROUP", Width: 20},
		{Title: "STOPPED AT", Width: 20},
		{Title: "STOPPED REASON", Width: 50},
	}
	return StoppedPage{table: components.NewTable(cols, thm), ecs: ecs}
}

func (p *StoppedPage) SetContext(nav ui.NavContext) tea.Cmd {
	p.nav = nav
	return p.fetchData()
}

func (p *StoppedPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p StoppedPage) Init() tea.Cmd     { return nil }

func (p StoppedPage) Update(msg tea.Msg) (StoppedPage, tea.Cmd) {
	switch msg := msg.(type) {
	case stoppedLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.tasks = msg.tasks
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

func (p StoppedPage) View() string {
	if p.nav.ClusterARN == "" {
		return "  Select a cluster first (drill down from Cluster tab)"
	}
	if p.loading {
		return "  Loading stopped tasks..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	if len(p.tasks) == 0 {
		return "  No stopped tasks found"
	}
	return p.table.View()
}

func (p StoppedPage) HelpBindings() []components.HelpBinding { return nil }

func (p *StoppedPage) updateRows() {
	var rows []table.Row
	for _, t := range p.tasks {
		rows = append(rows, table.Row{
			t.TaskID, t.Status, t.TaskDef, t.Group,
			t.StoppedAt.Format("2006-01-02 15:04:05"),
			t.StoppedReason,
		})
	}
	p.table.SetRows(rows)
}

func (p *StoppedPage) fetchData() tea.Cmd {
	p.loading = true
	ecs := p.ecs
	clusterARN := p.nav.ClusterARN
	return func() tea.Msg {
		tasks, err := ecs.ListStoppedTasks(context.Background(), clusterARN)
		return stoppedLoadedMsg{tasks: tasks, err: err}
	}
}

func (p *StoppedPage) Refresh() tea.Cmd { return p.fetchData() }
