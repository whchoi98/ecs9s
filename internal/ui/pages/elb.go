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

type elbLoadedMsg struct {
	lbs []aws.LoadBalancer
	tgs []aws.TargetGroup
	err error
}

type ELBPage struct {
	table   components.Table
	lbs     []aws.LoadBalancer
	elb     *aws.ELBClient
	loading bool
	err     error
}

func NewELBPage(elb *aws.ELBClient, thm theme.Theme) ELBPage {
	cols := []table.Column{
		{Title: "NAME", Width: 25},
		{Title: "TYPE", Width: 12},
		{Title: "SCHEME", Width: 12},
		{Title: "STATE", Width: 10},
		{Title: "DNS", Width: 40},
		{Title: "VPC", Width: 15},
	}
	return ELBPage{table: components.NewTable(cols, thm), elb: elb}
}

func (p *ELBPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p ELBPage) Init() tea.Cmd     { return p.fetchData() }

func (p ELBPage) Update(msg tea.Msg) (ELBPage, tea.Cmd) {
	switch msg := msg.(type) {
	case elbLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.lbs = msg.lbs
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

func (p ELBPage) View() string {
	if p.loading {
		return "  Loading load balancers..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	return p.table.View()
}

func (p ELBPage) HelpBindings() []components.HelpBinding { return nil }

func (p *ELBPage) updateRows() {
	var rows []table.Row
	for _, lb := range p.lbs {
		rows = append(rows, table.Row{
			lb.Name, lb.Type, lb.Scheme, lb.State, lb.DNSName, lb.VPCID,
		})
	}
	p.table.SetRows(rows)
}

func (p ELBPage) fetchData() tea.Cmd {
	elb := p.elb
	return func() tea.Msg {
		lbs, err := elb.ListLoadBalancers(context.Background())
		return elbLoadedMsg{lbs: lbs, err: err}
	}
}

func (p *ELBPage) Refresh() tea.Cmd { return p.fetchData() }
