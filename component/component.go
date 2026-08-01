// Package component contains composable terminal UI primitives without a TUI framework dependency.
package component

import (
	"fmt"
	"strings"
)

// Tone selects a semantic ANSI colour.
type Tone int

const (
	Muted Tone = iota
	Accent
	Success
	Warning
	Danger
)

// Segment is a small labelled portion of a StatusBar.
type Segment struct {
	Text string
	Tone Tone
}

// StatusBarOptions describes a responsive-looking command/status row.
type StatusBarOptions struct {
	Left  []Segment
	Right string
	Width int
}

// StatusBar renders a compact segmented status row.
func StatusBar(options StatusBarOptions) string {
	parts := make([]string, 0, len(options.Left))
	for _, s := range options.Left {
		parts = append(parts, paint(s.Text, s.Tone))
	}
	left := strings.Join(parts, "  ")
	if options.Right == "" {
		return left
	}
	if options.Width > displayWidth(left)+displayWidth(options.Right)+3 {
		return left + strings.Repeat(" ", options.Width-displayWidth(left)-displayWidth(options.Right)) + paint(options.Right, Muted)
	}
	return left + "  " + paint(options.Right, Muted)
}

// Progress renders a dithered progress indicator with an optional label.
func Progress(label string, value float64, width int, tone Tone) string {
	if width < 4 {
		width = 4
	}
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	filled := int(value*float64(width) + .5)
	return fmt.Sprintf("%s %s %3.0f%%", paint(label, tone), paint(strings.Repeat("▓", filled)+strings.Repeat("░", width-filled), tone), value*100)
}

// Gauge renders a small labelled utilisation gauge.
func Gauge(label string, value float64, width int, tone Tone) string {
	return Progress(label, value, width, tone)
}

// Badge renders a concise, coloured state label.
func Badge(label string, tone Tone) string { return paint("["+label+"]", tone) }

// Divider renders a quiet horizontal separator.
func Divider(width int) string {
	if width < 1 {
		return ""
	}
	return paint(strings.Repeat("─", width), Muted)
}

// Card wraps a title and content in a compact terminal panel.
func Card(title, content string, width int) string {
	if width < 12 {
		width = 12
	}
	border := "┌" + strings.Repeat("─", width-2) + "┐"
	footer := "└" + strings.Repeat("─", width-2) + "┘"
	line := func(value string) string { return "│ " + padDisplay(value, width-4) + " │" }
	return border + "\n" + line(paint(title, Accent)) + "\n" + line(content) + "\n" + footer
}

// SpinnerFrame returns a deterministic spinner frame for a host animation tick.
func SpinnerFrame(frame int, label string, tone Tone) string {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	index := frame % len(frames)
	if index < 0 {
		index += len(frames)
	}
	return paint(frames[index]+" "+label, tone)
}

func paint(value string, tone Tone) string {
	return fmt.Sprintf("\x1b[38;5;%dm%s\x1b[0m", toneColor(tone), value)
}

func toneColor(tone Tone) int {
	switch tone {
	case Accent:
		return 81
	case Success:
		return 42
	case Warning:
		return 214
	case Danger:
		return 196
	default:
		return 244
	}
}
func pad(value string, width int) string { return padDisplay(value, width) }
func strip(value string) string {
	for {
		start := strings.Index(value, "\x1b[")
		if start < 0 {
			return value
		}
		end := strings.Index(value[start:], "m")
		if end < 0 {
			return value
		}
		value = value[:start] + value[start+end+1:]
	}
}
