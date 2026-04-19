# ui

UI layer with shared message types, reusable components, and page models.

- `messages.go`: PageType enum, NavContext, shared message types (ErrorMsg, InfoMsg, DrillDownMsg, etc.)
- `components/`: Reusable widgets (table, tabs, commandbar, statusbar, help, confirm, logviewer, sparkline)
- `pages/`: One file per resource view. Each page is a bubbletea Model with Init/Update/View + fetchData/Refresh.
- `styles/`: Theme-aware lipgloss style definitions
