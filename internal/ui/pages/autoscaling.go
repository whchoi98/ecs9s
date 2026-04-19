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

type asgLoadedMsg struct {
	targets  []aws.ScalableTarget
	policies []aws.ScalingPolicy
	err      error
}

type ASGPage struct {
	table    components.Table
	targets  []aws.ScalableTarget
	asg      *aws.AutoScalingClient
	loading  bool
	err      error
}

func NewASGPage(asg *aws.AutoScalingClient, thm theme.Theme) ASGPage {
	cols := []table.Column{
		{Title: "RESOURCE", Width: 40},
		{Title: "DIMENSION", Width: 30},
		{Title: "MIN", Width: 6},
		{Title: "MAX", Width: 6},
	}
	return ASGPage{table: components.NewTable(cols, thm), asg: asg}
}

func (p *ASGPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p ASGPage) Init() tea.Cmd     { return p.fetchData() }

func (p ASGPage) Update(msg tea.Msg) (ASGPage, tea.Cmd) {
	switch msg := msg.(type) {
	case asgLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.targets = msg.targets
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

func (p ASGPage) View() string {
	if p.loading {
		return "  Loading auto scaling..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	return p.table.View()
}

func (p ASGPage) HelpBindings() []components.HelpBinding {
	return []components.HelpBinding{{Key: "S", Desc: "Adjust scale"}}
}

func (p *ASGPage) updateRows() {
	var rows []table.Row
	for _, t := range p.targets {
		rows = append(rows, table.Row{
			t.ResourceID, t.ScalableDim,
			fmt.Sprintf("%d", t.MinCapacity),
			fmt.Sprintf("%d", t.MaxCapacity),
		})
	}
	p.table.SetRows(rows)
}

func (p ASGPage) fetchData() tea.Cmd {
	asg := p.asg
	return func() tea.Msg {
		targets, err := asg.ListScalableTargets(context.Background())
		if err != nil {
			return asgLoadedMsg{err: err}
		}
		policies, _ := asg.ListScalingPolicies(context.Background())
		return asgLoadedMsg{targets: targets, policies: policies}
	}
}

func (p *ASGPage) Refresh() tea.Cmd { return p.fetchData() }
