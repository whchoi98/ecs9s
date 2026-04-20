# ui

UI layer: shared message types, 8 reusable components, and 21 page models.

- `messages.go`: PageType enum (21 types, Cluster→Alarms), NavContext, shared messages (ErrorMsg, InfoMsg, DrillDownMsg, GoBackMsg, etc.)
- `components/`: table (sortable/filterable, uses bubbles/table), tabs, commandbar (: and / modes), statusbar, help overlay, confirm dialog, logviewer (viewport-based), sparkline (Unicode block chars)
- `pages/`: One file per resource view. Each implements Init/Update/View + fetchData/Refresh + SetContext (for context-dependent pages) + SetSize + HelpBindings.
- `styles/`: Lipgloss style constants generated from theme.Theme.

## Critical Convention

Use `table.SelectedRow()` + `findXxxByName/ID()` for row selection. Never use `table.Cursor()` as an index — it breaks with filtering/sorting.
