# theme

Lipgloss color palette presets. Each Theme struct defines 16 colors (BgPrimary, FgPrimary, Accent, etc.) used by the styles package.

Presets: Dark (Tokyo Night), Light, Blue (Navy). `Get(name)` falls back to Dark for unknown names.

Files: `theme.go` (Theme struct, Get, Names, preset vars), `theme_test.go` (known/unknown theme lookup, Names count).
