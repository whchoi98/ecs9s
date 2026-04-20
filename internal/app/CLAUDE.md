# app

Root bubbletea model. Composes all UI components and 21 page models. Handles global keybindings, page routing, navigation stack, and AWS client initialization (10 clients).

## Key Patterns

- **Page registration**: 5 changes needed — struct field, New() init, resize(), initCurrentPage()/updateActivePage()/viewActivePage(). Plus PageType in messages.go + tab entry.
- **Value receiver**: `Init()` and `Update()` use value receivers. Set state in `New()`, not `Init()`. Mutations in `Update()` are returned via `return a, cmd`.
- **Tab switch**: Preserves drill-down NavContext. Global pages clear status bar cluster; context pages display it.
- **DrillDown**: Only `Enter` key sets NavContext (via `DrillDownMsg`). Tab never changes nav.
- **ECS Exec**: `ExecRequestMsg` → prerequisite check → `execReadyMsg` → `tea.ExecProcess` suspends TUI.
- **resetContext()**: Clears nav + navStack + statusBar.Cluster. Currently unused (removed to preserve context on tab switch).
