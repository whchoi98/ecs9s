package pages

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/whchoi98/ecs9s/internal/aws"
	"github.com/whchoi98/ecs9s/internal/theme"
	"github.com/whchoi98/ecs9s/internal/ui"
	"github.com/whchoi98/ecs9s/internal/ui/components"
)

type taskLoadedMsg struct {
	tasks []aws.Task
	err   error
}

type TaskPage struct {
	table   components.Table
	tasks   []aws.Task
	ecs     *aws.ECSClient
	nav     ui.NavContext
	loading bool
	err     error
}

func NewTaskPage(ecs *aws.ECSClient, thm theme.Theme) TaskPage {
	cols := []table.Column{
		{Title: "TASK ID", Width: 20},
		{Title: "STATUS", Width: 12},
		{Title: "TASKDEF", Width: 25},
		{Title: "LAUNCH", Width: 8},
		{Title: "CONTAINERS", Width: 11},
		{Title: "IP", Width: 16},
		{Title: "AGE", Width: 10},
	}
	return TaskPage{
		table: components.NewTable(cols, thm),
		ecs:   ecs,
	}
}

func (p *TaskPage) SetContext(nav ui.NavContext) tea.Cmd {
	p.nav = nav
	return p.fetchData()
}

func (p *TaskPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p TaskPage) Init() tea.Cmd     { return nil }

func (p TaskPage) Update(msg tea.Msg) (TaskPage, tea.Cmd) {
	switch msg := msg.(type) {
	case taskLoadedMsg:
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
		case "enter":
			if row := p.table.SelectedRow(); row != nil {
				if t := p.findTaskByID(row[0]); t != nil {
					nav := p.nav
					nav.TaskARN = t.TaskARN
					nav.TaskID = t.TaskID
					return p, func() tea.Msg {
						return ui.DrillDownMsg{Page: ui.PageContainer, Context: nav}
					}
				}
			}
		case "ctrl+d":
			if row := p.table.SelectedRow(); row != nil {
				if t := p.findTaskByID(row[0]); t != nil {
					ecs := p.ecs
					clusterARN := p.nav.ClusterARN
					taskID := t.TaskID
					tARN := t.TaskARN
					return p, func() tea.Msg {
						err := ecs.StopTask(context.Background(), clusterARN, tARN)
						if err != nil {
							return ui.ErrorMsg{Err: err}
						}
						return ui.InfoMsg{Text: fmt.Sprintf("Stopped task: %s", taskID)}
					}
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

func (p TaskPage) View() string {
	if p.nav.ClusterARN == "" {
		return "  Drill down from Cluster > Service first (Enter), then tasks will be shown."
	}
	if p.loading {
		return "  Loading tasks..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	if len(p.tasks) == 0 {
		svc := p.nav.ServiceName
		if svc == "" {
			svc = "(all)"
		}
		return fmt.Sprintf("  No running tasks found for service: %s", svc)
	}
	return p.table.View()
}

func (p TaskPage) HelpBindings() []components.HelpBinding {
	return []components.HelpBinding{
		{Key: "Enter", Desc: "Drill down to containers"},
		{Key: "Ctrl+d", Desc: "Stop task"},
	}
}

func (p *TaskPage) findTaskByID(taskID string) *aws.Task {
	for i := range p.tasks {
		if p.tasks[i].TaskID == taskID {
			return &p.tasks[i]
		}
	}
	return nil
}

func (p *TaskPage) updateRows() {
	var rows []table.Row
	for _, t := range p.tasks {
		age := time.Since(t.StartedAt).Truncate(time.Second).String()
		rows = append(rows, table.Row{
			t.TaskID, t.Status, t.TaskDefinition, t.LaunchType,
			fmt.Sprintf("%d", t.ContainerCount), t.PrivateIP, age,
		})
	}
	p.table.SetRows(rows)
}

func (p *TaskPage) fetchData() tea.Cmd {
	p.loading = true
	ecs := p.ecs
	nav := p.nav
	return func() tea.Msg {
		tasks, err := ecs.ListTasks(context.Background(), nav.ClusterARN, nav.ServiceName)
		return taskLoadedMsg{tasks: tasks, err: err}
	}
}

func (p *TaskPage) Refresh() tea.Cmd { return p.fetchData() }
