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

type vpcLoadedMsg struct {
	vpcs []aws.VPC
	sgs  []aws.SecurityGroup
	err  error
}

type VPCPage struct {
	table   components.Table
	vpcs    []aws.VPC
	ec2     *aws.EC2Client
	loading bool
	err     error
}

func NewVPCPage(ec2 *aws.EC2Client, thm theme.Theme) VPCPage {
	cols := []table.Column{
		{Title: "VPC ID", Width: 25},
		{Title: "NAME", Width: 20},
		{Title: "CIDR", Width: 18},
		{Title: "STATE", Width: 12},
		{Title: "DEFAULT", Width: 8},
	}
	return VPCPage{table: components.NewTable(cols, thm), ec2: ec2}
}

func (p *VPCPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p VPCPage) Init() tea.Cmd     { return p.fetchData() }

func (p VPCPage) Update(msg tea.Msg) (VPCPage, tea.Cmd) {
	switch msg := msg.(type) {
	case vpcLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.vpcs = msg.vpcs
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

func (p VPCPage) View() string {
	if p.loading {
		return "  Loading VPCs..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	return p.table.View()
}

func (p VPCPage) HelpBindings() []components.HelpBinding { return nil }

func (p *VPCPage) updateRows() {
	var rows []table.Row
	for _, v := range p.vpcs {
		def := "No"
		if v.IsDefault {
			def = "Yes"
		}
		rows = append(rows, table.Row{v.ID, v.Name, v.CIDR, v.State, def})
	}
	p.table.SetRows(rows)
}

func (p VPCPage) fetchData() tea.Cmd {
	ec2 := p.ec2
	return func() tea.Msg {
		vpcs, err := ec2.ListVPCs(context.Background())
		if err != nil {
			return vpcLoadedMsg{err: err}
		}
		sgs, _ := ec2.ListSecurityGroups(context.Background(), "")
		return vpcLoadedMsg{vpcs: vpcs, sgs: sgs}
	}
}

func (p *VPCPage) Refresh() tea.Cmd { return p.fetchData() }
