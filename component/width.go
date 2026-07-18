package component

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

func displayWidth(value string) int {
	return runewidth.StringWidth(strip(value))
}

func padDisplay(value string, width int) string {
	visible := displayWidth(value)
	if visible >= width {
		return value
	}
	return value + strings.Repeat(" ", width-visible)
}

func truncateDisplay(value string, width int) string {
	if width <= 0 {
		return ""
	}
	plain := strip(value)
	if runewidth.StringWidth(plain) <= width {
		return plain
	}
	if width == 1 {
		return "…"
	}
	runes := []rune(plain)
	out := strings.Builder{}
	used := 0
	for _, r := range runes {
		w := runewidth.RuneWidth(r)
		if used+w > width-1 {
			break
		}
		out.WriteRune(r)
		used += w
	}
	out.WriteRune('…')
	return out.String()
}
