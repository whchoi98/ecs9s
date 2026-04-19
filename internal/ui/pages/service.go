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

type serviceLoadedMsg struct {
	services []aws.Service
	err      error
}

type ServicePage struct {
	table    components.Table
	services []aws.Service
	ecs      *aws.ECSClient
	nav      ui.NavContext
	loading  bool
	err      error
}

func NewServicePage(ecs *aws.ECSClient, thm theme.Theme) ServicePage {
	cols := []table.Column{
		{Title: "NAME", Width: 25},
		{Title: "STATUS", Width: 10},
		{Title: "DESIRED", Width: 8},
		{Title: "RUNNING", Width: 8},
		{Title: "TASKDEF", Width: 25},
		{Title: "LAUNCH", Width: 8},
		{Title: "LB", Width: 15},
	}
	return ServicePage{
		table: components.NewTable(cols, thm),
		ecs:   ecs,
	}
}

func (p *ServicePage) SetContext(nav ui.NavContext) tea.Cmd {
	p.nav = nav
	return p.fetchData()
}

func (p *ServicePage) SetSize(w, h int) { p.table.SetSize(w, h) }

func (p ServicePage) Init() tea.Cmd { return nil }

func (p ServicePage) Update(msg tea.Msg) (ServicePage, tea.Cmd) {
	switch msg := msg.(type) {
	case serviceLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.services = msg.services
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
			if idx := p.table.Cursor(); idx < len(p.services) {
				svc := p.services[idx]
				nav := p.nav
				nav.ServiceName = svc.Name
				nav.ServiceARN = svc.ARN
				return p, func() tea.Msg {
					return ui.DrillDownMsg{Page: ui.PageTask, Context: nav}
				}
			}
		case "f":
			if idx := p.table.Cursor(); idx < len(p.services) {
				svc := p.services[idx]
				ecs := p.ecs
				clusterARN := p.nav.ClusterARN
				return p, func() tea.Msg {
					err := ecs.ForceNewDeployment(context.Background(), clusterARN, svc.Name)
					if err != nil {
						return ui.ErrorMsg{Err: err}
					}
					return ui.InfoMsg{Text: fmt.Sprintf("Force deploy: %s", svc.Name)}
				}
			}
		case "e":
			// Enable ECS Exec on the selected service
			if idx := p.table.Cursor(); idx < len(p.services) {
				svc := p.services[idx]
				ecs := p.ecs
				clusterARN := p.nav.ClusterARN
				return p, func() tea.Msg {
					err := ecs.EnableExecOnService(context.Background(), clusterARN, svc.Name)
					if err != nil {
						return ui.ErrorMsg{Err: err}
					}
					return ui.InfoMsg{Text: fmt.Sprintf("ECS Exec enabled + force deploy: %s (new tasks will have ExecuteCommandAgent)", svc.Name)}
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

func (p ServicePage) View() string {
	if p.nav.ClusterARN == "" {
		return "  Drill down from Cluster first (Enter) to see services."
	}
	if p.loading {
		return "  Loading services..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	if len(p.services) == 0 {
		return fmt.Sprintf("  No services found in cluster: %s", p.nav.ClusterName)
	}
	return p.table.View()
}

func (p ServicePage) HelpBindings() []components.HelpBinding {
	return []components.HelpBinding{
		{Key: "Enter", Desc: "Drill down to tasks"},
		{Key: "f", Desc: "Force new deployment"},
		{Key: "e", Desc: "Enable ECS Exec"},
		{Key: "S", Desc: "Scale service"},
		{Key: "b", Desc: "Rollback"},
	}
}

func (p *ServicePage) updateRows() {
	var rows []table.Row
	for _, s := range p.services {
		rows = append(rows, table.Row{
			s.Name, s.Status,
			fmt.Sprintf("%d", s.DesiredCount),
			fmt.Sprintf("%d", s.RunningCount),
			s.TaskDefinition, s.LaunchType, s.LoadBalancers,
		})
	}
	p.table.SetRows(rows)
}

func (p *ServicePage) fetchData() tea.Cmd {
	p.loading = true
	ecs := p.ecs
	clusterARN := p.nav.ClusterARN
	return func() tea.Msg {
		svcs, err := ecs.ListServices(context.Background(), clusterARN)
		return serviceLoadedMsg{services: svcs, err: err}
	}
}

func (p *ServicePage) Refresh() tea.Cmd { return p.fetchData() }
