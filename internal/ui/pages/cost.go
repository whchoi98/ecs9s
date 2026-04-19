package pages

import (
	"context"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/whchoi98/ecs9s/internal/aws"
	"github.com/whchoi98/ecs9s/internal/theme"
	"github.com/whchoi98/ecs9s/internal/ui"
	"github.com/whchoi98/ecs9s/internal/ui/components"
)

// Fargate pricing (us-east-1, on-demand) — per hour
const (
	fargateCPUPerHour = 0.04048  // per vCPU-hour
	fargateMemPerHour = 0.004445 // per GB-hour
)

type serviceCostInfo struct {
	Name         string
	RunningCount int32
	CPUUnits     int     // e.g. 256, 512, 1024
	MemoryMB     int     // e.g. 512, 1024, 2048
	Known        bool    // true if DescribeTaskDefinition succeeded
}

type costLoadedMsg struct {
	infos []serviceCostInfo
	err   error
}

type CostPage struct {
	table   components.Table
	infos   []serviceCostInfo
	ecs     *aws.ECSClient
	nav     ui.NavContext
	loading bool
	err     error
	thm     theme.Theme
	width   int
}

func NewCostPage(ecs *aws.ECSClient, thm theme.Theme) CostPage {
	cols := []table.Column{
		{Title: "SERVICE", Width: 25},
		{Title: "TASKS", Width: 8},
		{Title: "CPU (vCPU)", Width: 12},
		{Title: "MEMORY (GB)", Width: 12},
		{Title: "$/HOUR", Width: 10},
		{Title: "$/DAY", Width: 10},
		{Title: "$/MONTH", Width: 12},
	}
	return CostPage{table: components.NewTable(cols, thm), ecs: ecs, thm: thm}
}

func (p *CostPage) SetContext(nav ui.NavContext) tea.Cmd {
	p.nav = nav
	return p.fetchData()
}

func (p *CostPage) SetSize(w, h int) { p.table.SetSize(w, h); p.width = w }
func (p CostPage) Init() tea.Cmd     { return nil }

func (p CostPage) Update(msg tea.Msg) (CostPage, tea.Cmd) {
	switch msg := msg.(type) {
	case costLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.infos = msg.infos
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

func (p CostPage) View() string {
	if p.nav.ClusterARN == "" {
		return "  Select a cluster first (drill down from Cluster tab)"
	}
	if p.loading {
		return "  Fetching task definitions for cost estimation..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}

	var totalHourly, totalDaily, totalMonthly float64
	var hasUnknown bool
	for _, info := range p.infos {
		if !info.Known {
			hasUnknown = true
			continue
		}
		h := computeHourly(info)
		totalHourly += h
		totalDaily += h * 24
		totalMonthly += h * 24 * 30
	}

	summary := lipgloss.NewStyle().
		Foreground(p.thm.FgPrimary).
		Bold(true).
		Render(fmt.Sprintf("\n  Cluster: %s | Fargate Cost Estimate (on-demand, us-east-1 pricing)\n  Total: $%.2f/hour  $%.2f/day  $%.2f/month\n",
			p.nav.ClusterName, totalHourly, totalDaily, totalMonthly))

	note := lipgloss.NewStyle().Foreground(p.thm.FgMuted)
	footer := note.Render("  Pricing: vCPU $0.04048/hr, Memory $0.004445/GB-hr\n")
	if hasUnknown {
		footer += note.Render("  Services marked '-' have task definitions without CPU/Memory (EC2 launch type).\n")
	}

	return summary + "\n" + p.table.View() + "\n" + footer
}

func (p CostPage) HelpBindings() []components.HelpBinding { return nil }

func (p *CostPage) updateRows() {
	var rows []table.Row
	for _, info := range p.infos {
		if !info.Known {
			rows = append(rows, table.Row{
				info.Name,
				fmt.Sprintf("%d", info.RunningCount),
				"-", "-", "-", "-", "-",
			})
			continue
		}
		tasks := float64(info.RunningCount)
		cpuVCPU := float64(info.CPUUnits) / 1024.0
		memGB := float64(info.MemoryMB) / 1024.0
		h := computeHourly(info)

		rows = append(rows, table.Row{
			info.Name,
			fmt.Sprintf("%d", info.RunningCount),
			fmt.Sprintf("%.2f", cpuVCPU*tasks),
			fmt.Sprintf("%.2f", memGB*tasks),
			fmt.Sprintf("$%.4f", h),
			fmt.Sprintf("$%.2f", h*24),
			fmt.Sprintf("$%.2f", h*24*30),
		})
	}
	p.table.SetRows(rows)
}

func computeHourly(info serviceCostInfo) float64 {
	tasks := float64(info.RunningCount)
	cpuVCPU := float64(info.CPUUnits) / 1024.0
	memGB := float64(info.MemoryMB) / 1024.0
	return tasks * (cpuVCPU*fargateCPUPerHour + memGB*fargateMemPerHour)
}

func (p *CostPage) fetchData() tea.Cmd {
	p.loading = true
	ecs := p.ecs
	clusterARN := p.nav.ClusterARN
	return func() tea.Msg {
		services, err := ecs.ListServices(context.Background(), clusterARN)
		if err != nil {
			return costLoadedMsg{err: err}
		}

		var infos []serviceCostInfo
		for _, svc := range services {
			info := serviceCostInfo{
				Name:         svc.Name,
				RunningCount: svc.RunningCount,
			}

			// Fetch actual CPU/Memory from DescribeTaskDefinition
			res, err := ecs.GetTaskDefResources(context.Background(), svc.TaskDefinition)
			if err == nil && res.CPU != "" && res.Memory != "" {
				cpu, cpuErr := strconv.Atoi(res.CPU)
				mem, memErr := strconv.Atoi(res.Memory)
				if cpuErr == nil && memErr == nil {
					info.CPUUnits = cpu
					info.MemoryMB = mem
					info.Known = true
				}
			}

			infos = append(infos, info)
		}
		return costLoadedMsg{infos: infos}
	}
}

func (p *CostPage) Refresh() tea.Cmd { return p.fetchData() }
