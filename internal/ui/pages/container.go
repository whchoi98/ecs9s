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

type containerLoadedMsg struct {
	containers []aws.Container
	err        error
}

type ContainerPage struct {
	table      components.Table
	containers []aws.Container
	ecs        *aws.ECSClient
	nav        ui.NavContext
	loading    bool
	err        error
}

func NewContainerPage(ecs *aws.ECSClient, thm theme.Theme) ContainerPage {
	cols := []table.Column{
		{Title: "NAME", Width: 20},
		{Title: "IMAGE", Width: 35},
		{Title: "STATUS", Width: 12},
		{Title: "HEALTH", Width: 10},
		{Title: "PORTS", Width: 20},
	}
	return ContainerPage{
		table: components.NewTable(cols, thm),
		ecs:   ecs,
	}
}

func (p *ContainerPage) SetContext(nav ui.NavContext) tea.Cmd {
	p.nav = nav
	return p.fetchData()
}

func (p *ContainerPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p ContainerPage) Init() tea.Cmd     { return nil }

func (p ContainerPage) Update(msg tea.Msg) (ContainerPage, tea.Cmd) {
	switch msg := msg.(type) {
	case containerLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.containers = msg.containers
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
		case "x":
			if row := p.table.SelectedRow(); row != nil {
				if ctr := p.findContainerByName(row[0]); ctr != nil {
					return p, func() tea.Msg {
						return ExecRequestMsg{
							ClusterARN: ctr.ClusterARN,
							TaskARN:    ctr.TaskARN,
							Container:  ctr.Name,
						}
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

func (p ContainerPage) View() string {
	if p.nav.TaskARN == "" {
		return "  Drill down from Cluster > Service > Task first (Enter) to see containers."
	}
	if p.loading {
		return "  Loading containers..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	if len(p.containers) == 0 {
		return fmt.Sprintf("  No containers found in task: %s", p.nav.TaskID)
	}
	return p.table.View()
}

func (p ContainerPage) HelpBindings() []components.HelpBinding {
	return []components.HelpBinding{
		{Key: "x", Desc: "ECS Exec (shell)"},
		{Key: "Ctrl+f", Desc: "Port forward"},
	}
}

// ExecRequestMsg is sent when the user presses 'x' to start ECS Exec.
// Exported so the app shell can handle it with tea.ExecProcess.
type ExecRequestMsg struct {
	ClusterARN string
	TaskARN    string
	Container  string
}

func (p *ContainerPage) findContainerByName(name string) *aws.Container {
	for i := range p.containers {
		if p.containers[i].Name == name {
			return &p.containers[i]
		}
	}
	return nil
}

func (p *ContainerPage) updateRows() {
	var rows []table.Row
	for _, c := range p.containers {
		rows = append(rows, table.Row{
			c.Name, c.Image, c.Status, c.Health, c.Ports,
		})
	}
	p.table.SetRows(rows)
}

func (p *ContainerPage) fetchData() tea.Cmd {
	p.loading = true
	ecs := p.ecs
	nav := p.nav
	return func() tea.Msg {
		containers, err := ecs.DescribeContainers(context.Background(), nav.ClusterARN, nav.TaskARN)
		return containerLoadedMsg{containers: containers, err: err}
	}
}

func (p *ContainerPage) Refresh() tea.Cmd { return p.fetchData() }
