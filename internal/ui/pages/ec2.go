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

type ec2LoadedMsg struct {
	instances []aws.Instance
	err       error
}

type EC2Page struct {
	table     components.Table
	instances []aws.Instance
	ec2       *aws.EC2Client
	loading   bool
	err       error
}

func NewEC2Page(ec2 *aws.EC2Client, thm theme.Theme) EC2Page {
	cols := []table.Column{
		{Title: "INSTANCE ID", Width: 22},
		{Title: "NAME", Width: 20},
		{Title: "TYPE", Width: 14},
		{Title: "STATE", Width: 10},
		{Title: "PRIVATE IP", Width: 16},
		{Title: "PUBLIC IP", Width: 16},
		{Title: "AMI", Width: 22},
	}
	return EC2Page{table: components.NewTable(cols, thm), ec2: ec2}
}

func (p *EC2Page) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p EC2Page) Init() tea.Cmd     { return p.fetchData() }

func (p EC2Page) Update(msg tea.Msg) (EC2Page, tea.Cmd) {
	switch msg := msg.(type) {
	case ec2LoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.instances = msg.instances
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

func (p EC2Page) View() string {
	if p.loading {
		return "  Loading EC2 instances..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	return p.table.View()
}

func (p EC2Page) HelpBindings() []components.HelpBinding { return nil }

func (p *EC2Page) updateRows() {
	var rows []table.Row
	for _, i := range p.instances {
		rows = append(rows, table.Row{
			i.ID, i.Name, i.Type, i.State, i.PrivateIP, i.PublicIP, i.AMI,
		})
	}
	p.table.SetRows(rows)
}

func (p EC2Page) fetchData() tea.Cmd {
	ec2 := p.ec2
	return func() tea.Msg {
		instances, err := ec2.ListInstances(context.Background())
		return ec2LoadedMsg{instances: instances, err: err}
	}
}

func (p *EC2Page) Refresh() tea.Cmd { return p.fetchData() }
