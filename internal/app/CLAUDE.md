# app

Root bubbletea model. Composes all UI components and pages. Handles global keybindings, page routing, navigation stack, and AWS client initialization.

Adding a new page requires 5 changes: struct field, New() init, resize(), initCurrentPage()/updateActivePage()/viewActivePage().
