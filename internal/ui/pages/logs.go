package pages

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/whchoi98/ecs9s/internal/aws"
	"github.com/whchoi98/ecs9s/internal/theme"
	"github.com/whchoi98/ecs9s/internal/ui"
	"github.com/whchoi98/ecs9s/internal/ui/components"
)

type logsLoadedMsg struct {
	events []aws.LogEvent
	err    error
}

type LogsPage struct {
	viewer   components.LogViewer
	cw       *aws.CloudWatchClient
	logGroup string
	stream   string
	loading  bool
}

func NewLogsPage(cw *aws.CloudWatchClient, thm theme.Theme) LogsPage {
	return LogsPage{
		viewer: components.NewLogViewer(thm),
		cw:     cw,
	}
}

func (p *LogsPage) SetLogGroup(group, stream string) tea.Cmd {
	p.logGroup = group
	p.stream = stream
	return p.fetchData()
}

func (p *LogsPage) SetSize(w, h int) { p.viewer.SetSize(w, h) }
func (p LogsPage) Init() tea.Cmd     { return nil }

func (p LogsPage) Update(msg tea.Msg) (LogsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case logsLoadedMsg:
		p.loading = false
		if msg.err != nil {
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		var lines []components.LogLine
		for _, e := range msg.events {
			lines = append(lines, components.LogLine{
				Timestamp: e.Timestamp.Format("15:04:05"),
				Message:   e.Message,
			})
		}
		p.viewer.AppendLines(lines)
		return p, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "F":
			p.viewer.ToggleFollow()
			return p, nil
		case "c":
			p.viewer.Clear()
			return p, nil
		}
	}

	var cmd tea.Cmd
	p.viewer, cmd = p.viewer.Update(msg)
	return p, cmd
}

func (p LogsPage) View() string {
	if p.loading {
		return "  Loading logs..."
	}
	return p.viewer.View()
}

func (p LogsPage) HelpBindings() []components.HelpBinding {
	return []components.HelpBinding{
		{Key: "F", Desc: "Toggle follow"},
		{Key: "c", Desc: "Clear logs"},
	}
}

func (p *LogsPage) fetchData() tea.Cmd {
	p.loading = true
	cw := p.cw
	group := p.logGroup
	stream := p.stream
	return func() tea.Msg {
		events, err := cw.GetLogEvents(context.Background(), group, stream, time.Now().Add(-1*time.Hour), 100)
		return logsLoadedMsg{events: events, err: err}
	}
}

func (p *LogsPage) Refresh() tea.Cmd { return p.fetchData() }
