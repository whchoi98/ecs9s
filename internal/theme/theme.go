package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name       string
	BgPrimary  lipgloss.Color
	BgSecond   lipgloss.Color
	FgPrimary  lipgloss.Color
	FgSecond   lipgloss.Color
	FgMuted    lipgloss.Color
	Accent     lipgloss.Color
	Border     lipgloss.Color
	TabActive  lipgloss.Color
	TabInact   lipgloss.Color
	StatusBg   lipgloss.Color
	StatusFg   lipgloss.Color
	ErrColor   lipgloss.Color
	WarnColor  lipgloss.Color
	OkColor    lipgloss.Color
	SelectBg   lipgloss.Color
	SelectFg   lipgloss.Color
}

var themes = map[string]Theme{
	"dark":  Dark,
	"light": Light,
	"blue":  Blue,
}

var Dark = Theme{
	Name:      "dark",
	BgPrimary: lipgloss.Color("#1a1b26"),
	BgSecond:  lipgloss.Color("#24283b"),
	FgPrimary: lipgloss.Color("#c0caf5"),
	FgSecond:  lipgloss.Color("#a9b1d6"),
	FgMuted:   lipgloss.Color("#565f89"),
	Accent:    lipgloss.Color("#7aa2f7"),
	Border:    lipgloss.Color("#3b4261"),
	TabActive: lipgloss.Color("#7aa2f7"),
	TabInact:  lipgloss.Color("#565f89"),
	StatusBg:  lipgloss.Color("#1f2335"),
	StatusFg:  lipgloss.Color("#7aa2f7"),
	ErrColor:  lipgloss.Color("#f7768e"),
	WarnColor: lipgloss.Color("#e0af68"),
	OkColor:   lipgloss.Color("#9ece6a"),
	SelectBg:  lipgloss.Color("#3b4261"),
	SelectFg:  lipgloss.Color("#c0caf5"),
}

var Light = Theme{
	Name:      "light",
	BgPrimary: lipgloss.Color("#f5f5f5"),
	BgSecond:  lipgloss.Color("#e8e8e8"),
	FgPrimary: lipgloss.Color("#1a1a2e"),
	FgSecond:  lipgloss.Color("#3a3a5e"),
	FgMuted:   lipgloss.Color("#8888aa"),
	Accent:    lipgloss.Color("#2563eb"),
	Border:    lipgloss.Color("#ccccdd"),
	TabActive: lipgloss.Color("#2563eb"),
	TabInact:  lipgloss.Color("#8888aa"),
	StatusBg:  lipgloss.Color("#dddde8"),
	StatusFg:  lipgloss.Color("#2563eb"),
	ErrColor:  lipgloss.Color("#dc2626"),
	WarnColor: lipgloss.Color("#d97706"),
	OkColor:   lipgloss.Color("#16a34a"),
	SelectBg:  lipgloss.Color("#dbeafe"),
	SelectFg:  lipgloss.Color("#1a1a2e"),
}

var Blue = Theme{
	Name:      "blue",
	BgPrimary: lipgloss.Color("#0d1b2a"),
	BgSecond:  lipgloss.Color("#1b2838"),
	FgPrimary: lipgloss.Color("#e0e1dd"),
	FgSecond:  lipgloss.Color("#a8b2c1"),
	FgMuted:   lipgloss.Color("#5c6b7f"),
	Accent:    lipgloss.Color("#00b4d8"),
	Border:    lipgloss.Color("#2a3a4e"),
	TabActive: lipgloss.Color("#00b4d8"),
	TabInact:  lipgloss.Color("#5c6b7f"),
	StatusBg:  lipgloss.Color("#0a1628"),
	StatusFg:  lipgloss.Color("#00b4d8"),
	ErrColor:  lipgloss.Color("#ef476f"),
	WarnColor: lipgloss.Color("#ffd166"),
	OkColor:   lipgloss.Color("#06d6a0"),
	SelectBg:  lipgloss.Color("#1b3a5c"),
	SelectFg:  lipgloss.Color("#e0e1dd"),
}

func Get(name string) Theme {
	if t, ok := themes[name]; ok {
		return t
	}
	return Dark
}

func Names() []string {
	return []string{"dark", "light", "blue"}
}
