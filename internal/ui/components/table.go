package components

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/whchoi98/ecs9s/internal/theme"
)

type Table struct {
	inner      table.Model
	allRows    []table.Row
	filter     string
	filtering  bool
	sortCol    int
	sortAsc    bool
	columns    []table.Column
	thm        theme.Theme
	width      int
	height     int
}

func NewTable(cols []table.Column, thm theme.Theme) Table {
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(thm.Accent).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(thm.Border)
	s.Selected = lipgloss.NewStyle().
		Bold(true).
		Background(thm.SelectBg).
		Foreground(thm.SelectFg).
		Padding(0, 1)
	s.Cell = lipgloss.NewStyle().
		Foreground(thm.FgPrimary).
		Padding(0, 1)
	t.SetStyles(s)

	return Table{
		inner:   t,
		columns: cols,
		thm:     thm,
		sortAsc: true,
	}
}

func (t *Table) SetRows(rows []table.Row) {
	t.allRows = rows
	t.applyFilter()
}

func (t *Table) SetSize(w, h int) {
	t.width = w
	t.height = h
	t.inner.SetWidth(w)
	t.inner.SetHeight(h)
}

func (t *Table) SelectedRow() table.Row {
	return t.inner.SelectedRow()
}

func (t *Table) Cursor() int {
	return t.inner.Cursor()
}

func (t Table) Filtering() bool {
	return t.filtering
}

func (t Table) Filter() string {
	return t.filter
}

func (t *Table) StartFilter() {
	t.filtering = true
	t.filter = ""
}

func (t *Table) StopFilter() {
	t.filtering = false
}

func (t *Table) ClearFilter() {
	t.filter = ""
	t.filtering = false
	t.applyFilter()
}

func (t *Table) SetFilter(f string) {
	t.filter = f
	t.applyFilter()
}

func (t *Table) CycleSort() {
	if t.sortAsc {
		t.sortAsc = false
	} else {
		t.sortCol = (t.sortCol + 1) % len(t.columns)
		t.sortAsc = true
	}
	t.applyFilter()
}

func (t *Table) applyFilter() {
	rows := t.allRows
	if t.filter != "" {
		lf := strings.ToLower(t.filter)
		var filtered []table.Row
		for _, r := range rows {
			for _, cell := range r {
				if strings.Contains(strings.ToLower(cell), lf) {
					filtered = append(filtered, r)
					break
				}
			}
		}
		rows = filtered
	}

	if len(rows) > 0 && t.sortCol < len(rows[0]) {
		col := t.sortCol
		asc := t.sortAsc
		sort.SliceStable(rows, func(i, j int) bool {
			if asc {
				return rows[i][col] < rows[j][col]
			}
			return rows[i][col] > rows[j][col]
		})
	}

	t.inner.SetRows(rows)
}

func (t Table) Update(msg tea.Msg) (Table, tea.Cmd) {
	var cmd tea.Cmd
	t.inner, cmd = t.inner.Update(msg)
	return t, cmd
}

func (t Table) View() string {
	return t.inner.View()
}
