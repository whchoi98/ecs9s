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

type alarmsLoadedMsg struct {
	alarms []aws.Alarm
	err    error
}

type AlarmsPage struct {
	table   components.Table
	alarms  []aws.Alarm
	cw      *aws.CloudWatchClient
	loading bool
	err     error
}

func NewAlarmsPage(cw *aws.CloudWatchClient, thm theme.Theme) AlarmsPage {
	cols := []table.Column{
		{Title: "ALARM NAME", Width: 30},
		{Title: "STATE", Width: 10},
		{Title: "METRIC", Width: 20},
		{Title: "NAMESPACE", Width: 15},
		{Title: "THRESHOLD", Width: 12},
		{Title: "COMPARISON", Width: 20},
		{Title: "UPDATED", Width: 20},
	}
	return AlarmsPage{table: components.NewTable(cols, thm), cw: cw}
}

func (p *AlarmsPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p AlarmsPage) Init() tea.Cmd     { return p.fetchData() }

func (p AlarmsPage) Update(msg tea.Msg) (AlarmsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case alarmsLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.alarms = msg.alarms
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

func (p AlarmsPage) View() string {
	if p.loading {
		return "  Loading CloudWatch alarms..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	return p.table.View()
}

func (p AlarmsPage) HelpBindings() []components.HelpBinding { return nil }

func (p *AlarmsPage) updateRows() {
	var rows []table.Row
	for _, a := range p.alarms {
		rows = append(rows, table.Row{
			a.Name, a.State, a.MetricName, a.Namespace,
			fmt.Sprintf("%.2f", a.Threshold), a.Comparison,
			a.UpdatedAt.Format("2006-01-02 15:04"),
		})
	}
	p.table.SetRows(rows)
}

func (p AlarmsPage) fetchData() tea.Cmd {
	cw := p.cw
	return func() tea.Msg {
		alarms, err := cw.ListAlarms(context.Background())
		return alarmsLoadedMsg{alarms: alarms, err: err}
	}
}

func (p *AlarmsPage) Refresh() tea.Cmd { return p.fetchData() }
