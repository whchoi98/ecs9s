# components

8 reusable Bubbletea UI widgets. Each is a struct with Update/View methods.

| File | Widget | Notes |
|------|--------|-------|
| table.go | Sortable/filterable table | Wraps bubbles/table. Use `SelectedRow()` not `Cursor()` for data lookup. |
| tabs.go | Tab bar | Lipgloss-rendered. `SetActiveByCommand()` for command mode. |
| commandbar.go | Command input | `:` for commands, `/` for filter. Sends `CommandExecuteMsg`. |
| statusbar.go | Status display | Shows cluster, profile, region, error/info messages. |
| help.go | Keybinding overlay | Global + page-specific bindings via `SetExtra()`. |
| confirm.go | Confirmation dialog | For destructive actions. Default focus on "No" for safety. |
| logviewer.go | Scrollable log view | Viewport-based. Follow mode, filter, clear. |
| sparkline.go | Inline sparkline chart | Unicode block chars (▁▂▃▄▅▆▇█). Color changes at 60%/80% thresholds. |
