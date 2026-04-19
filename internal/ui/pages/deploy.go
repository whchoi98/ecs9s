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

type deployLoadedMsg struct {
	deployments []aws.Deployment
	err         error
}

type DeployPage struct {
	table       components.Table
	deployments []aws.Deployment
	ecs         *aws.ECSClient
	nav         ui.NavContext
	loading     bool
	err         error
}

func NewDeployPage(ecs *aws.ECSClient, thm theme.Theme) DeployPage {
	cols := []table.Column{
		{Title: "ID", Width: 15},
		{Title: "STATUS", Width: 10},
		{Title: "TASKDEF", Width: 25},
		{Title: "DESIRED", Width: 8},
		{Title: "RUNNING", Width: 8},
		{Title: "ROLLOUT", Width: 15},
		{Title: "CREATED", Width: 20},
	}
	return DeployPage{table: components.NewTable(cols, thm), ecs: ecs}
}

func (p *DeployPage) SetContext(nav ui.NavContext) tea.Cmd {
	p.nav = nav
	return p.fetchData()
}

func (p *DeployPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p DeployPage) Init() tea.Cmd     { return nil }

func (p DeployPage) Update(msg tea.Msg) (DeployPage, tea.Cmd) {
	switch msg := msg.(type) {
	case deployLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.deployments = msg.deployments
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

func (p DeployPage) View() string {
	if p.nav.ServiceName == "" {
		return "  Drill down from Cluster > Service first, then switch to Deploy tab"
	}
	if p.loading {
		return "  Loading deployment history..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	header := fmt.Sprintf("  Service: %s/%s\n\n", p.nav.ClusterName, p.nav.ServiceName)
	return header + p.table.View()
}

func (p DeployPage) HelpBindings() []components.HelpBinding { return nil }

func (p *DeployPage) updateRows() {
	var rows []table.Row
	for _, d := range p.deployments {
		rows = append(rows, table.Row{
			d.ID, d.Status, d.TaskDefinition,
			fmt.Sprintf("%d", d.DesiredCount),
			fmt.Sprintf("%d", d.RunningCount),
			d.RolloutState,
			d.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	p.table.SetRows(rows)
}

func (p *DeployPage) fetchData() tea.Cmd {
	p.loading = true
	ecs := p.ecs
	nav := p.nav
	return func() tea.Msg {
		deps, _, err := ecs.GetServiceDeployments(context.Background(), nav.ClusterARN, nav.ServiceName)
		return deployLoadedMsg{deployments: deps, err: err}
	}
}

func (p *DeployPage) Refresh() tea.Cmd { return p.fetchData() }
