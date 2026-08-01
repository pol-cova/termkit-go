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

// fitDisplayWidth truncates or pads a string to an exact visible width, preserving leading ANSI codes.
func fitDisplayWidth(value string, width int) string {
	if width <= 0 {
		return ""
	}
	prefix := ansiPrefix(value)
	plain := strip(value)
	if runewidth.StringWidth(plain) <= width {
		return padDisplay(prefix+plain, width)
	}
	return prefix + truncateDisplay(plain, width)
}

func ansiPrefix(value string) string {
	var b strings.Builder
	i := 0
	for i < len(value) {
		if value[i] != '\x1b' {
			break
		}
		for j := i; j < len(value); j++ {
			b.WriteByte(value[j])
			if value[j] == 'm' {
				i = j + 1
				break
			}
		}
	}
	return b.String()
}
