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

type eventsLoadedMsg struct {
	events []aws.ServiceEvent
	err    error
}

type EventsPage struct {
	table   components.Table
	events  []aws.ServiceEvent
	ecs     *aws.ECSClient
	nav     ui.NavContext
	loading bool
	err     error
}

func NewEventsPage(ecs *aws.ECSClient, thm theme.Theme) EventsPage {
	cols := []table.Column{
		{Title: "TIME", Width: 20},
		{Title: "MESSAGE", Width: 100},
	}
	return EventsPage{table: components.NewTable(cols, thm), ecs: ecs}
}

func (p *EventsPage) SetContext(nav ui.NavContext) tea.Cmd {
	p.nav = nav
	return p.fetchData()
}

func (p *EventsPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p EventsPage) Init() tea.Cmd     { return nil }

func (p EventsPage) Update(msg tea.Msg) (EventsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case eventsLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.events = msg.events
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

func (p EventsPage) View() string {
	if p.nav.ServiceName == "" {
		return "  Drill down from Cluster > Service first, then switch to Events tab"
	}
	if p.loading {
		return "  Loading service events..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	header := fmt.Sprintf("  Service: %s/%s\n\n", p.nav.ClusterName, p.nav.ServiceName)
	return header + p.table.View()
}

func (p EventsPage) HelpBindings() []components.HelpBinding { return nil }

func (p *EventsPage) updateRows() {
	var rows []table.Row
	for _, e := range p.events {
		rows = append(rows, table.Row{
			e.CreatedAt.Format("2006-01-02 15:04:05"),
			e.Message,
		})
	}
	p.table.SetRows(rows)
}

func (p *EventsPage) fetchData() tea.Cmd {
	p.loading = true
	ecs := p.ecs
	nav := p.nav
	return func() tea.Msg {
		events, err := ecs.GetServiceEvents(context.Background(), nav.ClusterARN, nav.ServiceName)
		return eventsLoadedMsg{events: events, err: err}
	}
}

func (p *EventsPage) Refresh() tea.Cmd { return p.fetchData() }
