package components

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/whchoi98/ecs9s/internal/theme"
)

var sparks = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

type Sparkline struct {
	Label  string
	Values []float64
	Width  int
	thm    theme.Theme
}

func NewSparkline(label string, thm theme.Theme) Sparkline {
	return Sparkline{Label: label, Width: 30, thm: thm}
}

func (s *Sparkline) SetValues(v []float64) {
	s.Values = v
}

func (s Sparkline) View() string {
	if len(s.Values) == 0 {
		return ""
	}

	minVal, maxVal := s.Values[0], s.Values[0]
	var sum float64
	for _, v := range s.Values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
		sum += v
	}
	avg := sum / float64(len(s.Values))
	last := s.Values[len(s.Values)-1]

	valRange := maxVal - minVal
	if valRange == 0 {
		valRange = 1
	}

	// Take last N values that fit the width
	vals := s.Values
	if len(vals) > s.Width {
		vals = vals[len(vals)-s.Width:]
	}

	var sb strings.Builder
	for _, v := range vals {
		idx := int(math.Round((v - minVal) / valRange * float64(len(sparks)-1)))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(sparks) {
			idx = len(sparks) - 1
		}
		sb.WriteRune(sparks[idx])
	}

	color := s.thm.OkColor
	if last > 80 {
		color = s.thm.ErrColor
	} else if last > 60 {
		color = s.thm.WarnColor
	}

	label := lipgloss.NewStyle().
		Foreground(s.thm.FgPrimary).
		Bold(true).
		Render(s.Label)

	chart := lipgloss.NewStyle().
		Foreground(color).
		Render(sb.String())

	stats := lipgloss.NewStyle().
		Foreground(s.thm.FgMuted).
		Render(fmt.Sprintf(" %.1f%% (avg: %.1f%%)", last, avg))

	return label + " " + chart + stats
}
