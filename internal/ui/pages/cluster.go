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

type clusterLoadedMsg struct {
	clusters []aws.Cluster
	err      error
}

type ClusterPage struct {
	table    components.Table
	clusters []aws.Cluster
	ecs      *aws.ECSClient
	loading  bool
	err      error
}

func NewClusterPage(ecs *aws.ECSClient, thm theme.Theme) ClusterPage {
	cols := []table.Column{
		{Title: "NAME", Width: 30},
		{Title: "STATUS", Width: 10},
		{Title: "SERVICES", Width: 10},
		{Title: "TASKS", Width: 8},
		{Title: "INSTANCES", Width: 10},
	}
	return ClusterPage{
		table: components.NewTable(cols, thm),
		ecs:   ecs,
	}
}

func (p ClusterPage) Init() tea.Cmd {
	return p.fetchData()
}

func (p *ClusterPage) SetSize(w, h int) {
	p.table.SetSize(w, h)
}

func (p ClusterPage) Update(msg tea.Msg) (ClusterPage, tea.Cmd) {
	switch msg := msg.(type) {
	case clusterLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.clusters = msg.clusters
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
				idx := p.table.Cursor()
				if idx < len(p.clusters) {
					cl := p.clusters[idx]
					return p, func() tea.Msg {
						return ui.DrillDownMsg{
							Page: ui.PageService,
							Context: ui.NavContext{
								ClusterARN:  cl.ARN,
								ClusterName: cl.Name,
							},
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

func (p ClusterPage) View() string {
	if p.loading {
		return "  Loading clusters..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	return p.table.View()
}

func (p ClusterPage) HelpBindings() []components.HelpBinding {
	return []components.HelpBinding{
		{Key: "Enter", Desc: "Drill down to services"},
	}
}

func (p *ClusterPage) updateRows() {
	var rows []table.Row
	for _, c := range p.clusters {
		rows = append(rows, table.Row{
			c.Name,
			c.Status,
			fmt.Sprintf("%d", c.ActiveServicesCount),
			fmt.Sprintf("%d", c.RunningTasksCount),
			fmt.Sprintf("%d", c.RegisteredInstances),
		})
	}
	p.table.SetRows(rows)
}

func (p *ClusterPage) fetchData() tea.Cmd {
	p.loading = true
	ecs := p.ecs
	return func() tea.Msg {
		clusters, err := ecs.ListClusters(context.Background())
		return clusterLoadedMsg{clusters: clusters, err: err}
	}
}

func (p *ClusterPage) Refresh() tea.Cmd {
	return p.fetchData()
}
