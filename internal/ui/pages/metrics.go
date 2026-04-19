package pages

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/whchoi98/ecs9s/internal/aws"
	"github.com/whchoi98/ecs9s/internal/theme"
	"github.com/whchoi98/ecs9s/internal/ui"
	"github.com/whchoi98/ecs9s/internal/ui/components"
)

type metricsLoadedMsg struct {
	cpuPoints []aws.MetricDatapoint
	memPoints []aws.MetricDatapoint
	err       error
}

type MetricsPage struct {
	cpuSpark components.Sparkline
	memSpark components.Sparkline
	cw       *aws.CloudWatchClient
	nav      ui.NavContext
	duration time.Duration
	loading  bool
	err      error
	thm      theme.Theme
	width    int
	height   int
}

func NewMetricsPage(cw *aws.CloudWatchClient, thm theme.Theme) MetricsPage {
	return MetricsPage{
		cpuSpark: components.NewSparkline("CPU", thm),
		memSpark: components.NewSparkline("Memory", thm),
		cw:       cw,
		duration: 1 * time.Hour,
		thm:      thm,
	}
}

func (p *MetricsPage) SetContext(nav ui.NavContext) tea.Cmd {
	p.nav = nav
	return p.fetchData()
}

func (p *MetricsPage) SetSize(w, h int) {
	p.width = w
	p.height = h
	p.cpuSpark.Width = w - 20
	p.memSpark.Width = w - 20
}

func (p MetricsPage) Init() tea.Cmd { return nil }

func (p MetricsPage) Update(msg tea.Msg) (MetricsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case metricsLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		sort.Slice(msg.cpuPoints, func(i, j int) bool {
			return msg.cpuPoints[i].Timestamp.Before(msg.cpuPoints[j].Timestamp)
		})
		sort.Slice(msg.memPoints, func(i, j int) bool {
			return msg.memPoints[i].Timestamp.Before(msg.memPoints[j].Timestamp)
		})
		var cpuVals, memVals []float64
		for _, dp := range msg.cpuPoints {
			cpuVals = append(cpuVals, dp.Value)
		}
		for _, dp := range msg.memPoints {
			memVals = append(memVals, dp.Value)
		}
		p.cpuSpark.SetValues(cpuVals)
		p.memSpark.SetValues(memVals)
		return p, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			p.duration = 1 * time.Hour
			return p, p.fetchData()
		case "3":
			p.duration = 3 * time.Hour
			return p, p.fetchData()
		case "6":
			p.duration = 6 * time.Hour
			return p, p.fetchData()
		}
	}
	return p, nil
}

func (p MetricsPage) View() string {
	if p.loading {
		return "  Loading metrics..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	if p.nav.ClusterName == "" {
		return "  Select a cluster and service first (drill down from Cluster > Service)"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n  Service: %s/%s | Range: %v\n\n",
		p.nav.ClusterName, p.nav.ServiceName, p.duration))
	sb.WriteString("  " + p.cpuSpark.View() + "\n\n")
	sb.WriteString("  " + p.memSpark.View() + "\n\n")
	sb.WriteString("  [1] 1h  [3] 3h  [6] 6h")
	return sb.String()
}

func (p MetricsPage) HelpBindings() []components.HelpBinding {
	return []components.HelpBinding{
		{Key: "1/3/6", Desc: "Set time range (hours)"},
	}
}

func (p *MetricsPage) fetchData() tea.Cmd {
	p.loading = true
	cw := p.cw
	cluster := p.nav.ClusterName
	service := p.nav.ServiceName
	dur := p.duration
	return func() tea.Msg {
		cpuPts, err := cw.GetECSMetrics(context.Background(), cluster, service, "CPUUtilization", dur)
		if err != nil {
			return metricsLoadedMsg{err: err}
		}
		memPts, _ := cw.GetECSMetrics(context.Background(), cluster, service, "MemoryUtilization", dur)
		return metricsLoadedMsg{cpuPoints: cpuPts, memPoints: memPts}
	}
}

func (p *MetricsPage) Refresh() tea.Cmd { return p.fetchData() }
