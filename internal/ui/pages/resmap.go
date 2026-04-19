package pages

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/whchoi98/ecs9s/internal/aws"
	"github.com/whchoi98/ecs9s/internal/theme"
	"github.com/whchoi98/ecs9s/internal/ui"
	"github.com/whchoi98/ecs9s/internal/ui/components"
)

type resmapLoadedMsg struct {
	deployments []aws.Deployment
	tgARNs      []string
	tgs         []aws.TargetGroup
	tasks       []aws.Task
	err         error
}

type ResMapPage struct {
	deployments []aws.Deployment
	tgARNs      []string
	tgs         []aws.TargetGroup
	tasks       []aws.Task
	ecs         *aws.ECSClient
	elb         *aws.ELBClient
	nav         ui.NavContext
	loading     bool
	err         error
	thm         theme.Theme
	width       int
	height      int
}

func NewResMapPage(ecs *aws.ECSClient, elb *aws.ELBClient, thm theme.Theme) ResMapPage {
	return ResMapPage{ecs: ecs, elb: elb, thm: thm}
}

func (p *ResMapPage) SetContext(nav ui.NavContext) tea.Cmd {
	p.nav = nav
	return p.fetchData()
}

func (p *ResMapPage) SetSize(w, h int) { p.width = w; p.height = h }
func (p ResMapPage) Init() tea.Cmd     { return nil }

func (p ResMapPage) Update(msg tea.Msg) (ResMapPage, tea.Cmd) {
	switch msg := msg.(type) {
	case resmapLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.deployments = msg.deployments
		p.tgARNs = msg.tgARNs
		p.tgs = msg.tgs
		p.tasks = msg.tasks
		return p, nil
	}
	return p, nil
}

func (p ResMapPage) View() string {
	if p.nav.ServiceName == "" {
		return "  Drill down from Cluster > Service first, then switch to ResMap tab"
	}
	if p.loading {
		return "  Loading resource map..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}

	title := lipgloss.NewStyle().Bold(true).Foreground(p.thm.Accent)
	label := lipgloss.NewStyle().Bold(true).Foreground(p.thm.FgPrimary)
	value := lipgloss.NewStyle().Foreground(p.thm.FgSecond)
	muted := lipgloss.NewStyle().Foreground(p.thm.FgMuted)
	ok := lipgloss.NewStyle().Foreground(p.thm.OkColor)
	warn := lipgloss.NewStyle().Foreground(p.thm.WarnColor)
	errStyle := lipgloss.NewStyle().Foreground(p.thm.ErrColor)

	var sb strings.Builder

	sb.WriteString(title.Render(fmt.Sprintf("\n  Resource Map: %s / %s", p.nav.ClusterName, p.nav.ServiceName)))
	sb.WriteString("\n\n")

	// Deployments
	sb.WriteString(label.Render("  Deployments"))
	sb.WriteString("\n")
	for _, d := range p.deployments {
		statusStyle := ok
		if d.Status == "ACTIVE" {
			statusStyle = warn
		}
		sb.WriteString(fmt.Sprintf("    %s %s  %s  desired:%s running:%s  %s\n",
			muted.Render("├─"),
			statusStyle.Render(d.Status),
			value.Render(d.TaskDefinition),
			value.Render(fmt.Sprintf("%d", d.DesiredCount)),
			value.Render(fmt.Sprintf("%d", d.RunningCount)),
			muted.Render(d.RolloutState),
		))
	}
	sb.WriteString("\n")

	// Target Groups
	if len(p.tgs) > 0 {
		sb.WriteString(label.Render("  Target Groups"))
		sb.WriteString("\n")
		for _, tg := range p.tgs {
			hStyle := ok
			if tg.UnhealthyCount > 0 {
				hStyle = errStyle
			}
			sb.WriteString(fmt.Sprintf("    %s %s  %s:%s  healthy:%s unhealthy:%s\n",
				muted.Render("├─"),
				value.Render(tg.Name),
				muted.Render(tg.Protocol),
				muted.Render(fmt.Sprintf("%d", tg.Port)),
				ok.Render(fmt.Sprintf("%d", tg.HealthyCount)),
				hStyle.Render(fmt.Sprintf("%d", tg.UnhealthyCount)),
			))
		}
		sb.WriteString("\n")
	}

	// Tasks
	if len(p.tasks) > 0 {
		sb.WriteString(label.Render(fmt.Sprintf("  Running Tasks (%d)", len(p.tasks))))
		sb.WriteString("\n")
		for _, t := range p.tasks {
			tStyle := ok
			if t.Status != "RUNNING" {
				tStyle = warn
			}
			sb.WriteString(fmt.Sprintf("    %s %s  %s  %s  ip:%s\n",
				muted.Render("├─"),
				value.Render(t.TaskID),
				tStyle.Render(t.Status),
				muted.Render(t.TaskDefinition),
				muted.Render(t.PrivateIP),
			))
		}
	}

	return sb.String()
}

func (p ResMapPage) HelpBindings() []components.HelpBinding { return nil }

func (p *ResMapPage) fetchData() tea.Cmd {
	p.loading = true
	ecsC := p.ecs
	elbC := p.elb
	nav := p.nav
	return func() tea.Msg {
		deps, tgARNs, err := ecsC.GetServiceDeployments(context.Background(), nav.ClusterARN, nav.ServiceName)
		if err != nil {
			return resmapLoadedMsg{err: err}
		}
		tasks, _ := ecsC.ListTasks(context.Background(), nav.ClusterARN, nav.ServiceName)
		var tgs []aws.TargetGroup
		if len(tgARNs) > 0 {
			allTGs, _ := elbC.ListTargetGroups(context.Background())
			tgSet := make(map[string]bool)
			for _, a := range tgARNs {
				tgSet[a] = true
			}
			for _, tg := range allTGs {
				if tgSet[tg.ARN] {
					tgs = append(tgs, tg)
				}
			}
		}
		return resmapLoadedMsg{deployments: deps, tgARNs: tgARNs, tgs: tgs, tasks: tasks}
	}
}

func (p *ResMapPage) Refresh() tea.Cmd { return p.fetchData() }
