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

type ecrLoadedMsg struct {
	repos []aws.Repository
	err   error
}

type ECRPage struct {
	table   components.Table
	repos   []aws.Repository
	ecr     *aws.ECRClient
	loading bool
	err     error
}

func NewECRPage(ecr *aws.ECRClient, thm theme.Theme) ECRPage {
	cols := []table.Column{
		{Title: "REPOSITORY", Width: 30},
		{Title: "URI", Width: 50},
		{Title: "IMAGES", Width: 8},
		{Title: "CREATED", Width: 12},
	}
	return ECRPage{table: components.NewTable(cols, thm), ecr: ecr}
}

func (p *ECRPage) SetSize(w, h int) { p.table.SetSize(w, h) }
func (p ECRPage) Init() tea.Cmd     { return p.fetchData() }

func (p ECRPage) Update(msg tea.Msg) (ECRPage, tea.Cmd) {
	switch msg := msg.(type) {
	case ecrLoadedMsg:
		p.loading = false
		if msg.err != nil {
			p.err = msg.err
			return p, func() tea.Msg { return ui.ErrorMsg{Err: msg.err} }
		}
		p.repos = msg.repos
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

func (p ECRPage) View() string {
	if p.loading {
		return "  Loading ECR repositories..."
	}
	if p.err != nil {
		return fmt.Sprintf("  Error: %v", p.err)
	}
	return p.table.View()
}

func (p ECRPage) HelpBindings() []components.HelpBinding { return nil }

func (p *ECRPage) updateRows() {
	var rows []table.Row
	for _, r := range p.repos {
		rows = append(rows, table.Row{
			r.Name, r.URI, fmt.Sprintf("%d", r.ImageCount),
			r.CreatedAt.Format("2006-01-02"),
		})
	}
	p.table.SetRows(rows)
}

func (p ECRPage) fetchData() tea.Cmd {
	ecr := p.ecr
	return func() tea.Msg {
		repos, err := ecr.ListRepositories(context.Background())
		return ecrLoadedMsg{repos: repos, err: err}
	}
}

func (p *ECRPage) Refresh() tea.Cmd { return p.fetchData() }
