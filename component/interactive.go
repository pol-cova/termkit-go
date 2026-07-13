package component

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Input is a small, framework-neutral single-line editor. Key names use the
// spelling commonly returned by Bubble Tea and termbox-style readers.
type Input struct {
	Value       []rune
	Cursor      int
	Placeholder string
	MaxLength   int
}

// HandleKey applies a key event and reports whether it changed the input or
// was submitted. Printable strings are inserted as-is; named keys are e.g.
// "backspace", "left", "ctrl+a", or "enter".
func (i *Input) HandleKey(key string) (changed, submitted bool) {
	if i.Cursor < 0 {
		i.Cursor = 0
	}
	if i.Cursor > len(i.Value) {
		i.Cursor = len(i.Value)
	}
	switch key {
	case "backspace":
		if i.Cursor > 0 {
			i.Value = append(i.Value[:i.Cursor-1], i.Value[i.Cursor:]...)
			i.Cursor--
			return true, false
		}
	case "delete":
		if i.Cursor < len(i.Value) {
			i.Value = append(i.Value[:i.Cursor], i.Value[i.Cursor+1:]...)
			return true, false
		}
	case "left", "ctrl+b":
		if i.Cursor > 0 {
			i.Cursor--
		}
	case "right", "ctrl+f":
		if i.Cursor < len(i.Value) {
			i.Cursor++
		}
	case "home", "ctrl+a":
		i.Cursor = 0
	case "end", "ctrl+e":
		i.Cursor = len(i.Value)
	case "enter":
		return false, true
	default:
		if key == "" || strings.ContainsAny(key, "\x1b\r\n") {
			return false, false
		}
		incoming := []rune(key)
		if i.MaxLength > 0 && len(i.Value)+len(incoming) > i.MaxLength {
			return false, false
		}
		oldLength := len(i.Value)
		i.Value = append(i.Value, make([]rune, len(incoming))...)
		copy(i.Value[i.Cursor+len(incoming):], i.Value[i.Cursor:oldLength])
		copy(i.Value[i.Cursor:], incoming)
		i.Cursor += len(incoming)
		return true, false
	}
	return false, false
}

// Text returns the current value, or the placeholder when empty.
func (i Input) Text() string {
	if len(i.Value) == 0 {
		return i.Placeholder
	}
	return string(i.Value)
}

// View renders the input with a visible cursor when focused.
func (i Input) View(focused bool, tone Tone) string {
	text := i.Text()
	if focused {
		cursor := i.Cursor
		if cursor > utf8.RuneCountInString(text) {
			cursor = utf8.RuneCountInString(text)
		}
		r := []rune(text)
		if len(i.Value) == 0 {
			r = []rune(" ")
		}
		if cursor >= len(r) {
			r = append(r, ' ')
		}
		r[cursor] = '▌'
		text = string(r)
	}
	return paint(text, tone)
}

// Option is an item in a Select list.
type Option struct{ Label, Description string }

// Select is a keyboard-friendly list selector with OpenTUI-like j/k support.
type Select struct {
	Options                  []Option
	Selected, Offset, Height int
}

// HandleKey moves the highlight and reports whether enter confirmed it.
func (s *Select) HandleKey(key string) (selected bool) {
	if len(s.Options) == 0 {
		return false
	}
	if s.Selected < 0 {
		s.Selected = 0
	}
	if s.Selected >= len(s.Options) {
		s.Selected = len(s.Options) - 1
	}
	switch key {
	case "up", "k":
		if s.Selected > 0 {
			s.Selected--
		}
	case "down", "j":
		if s.Selected < len(s.Options)-1 {
			s.Selected++
		}
	case "home":
		s.Selected = 0
	case "end":
		s.Selected = len(s.Options) - 1
	case "enter":
		return true
	}
	s.Height = maxInt(s.Height, 1)
	if s.Selected < s.Offset {
		s.Offset = s.Selected
	}
	if s.Selected >= s.Offset+s.Height {
		s.Offset = s.Selected - s.Height + 1
	}
	return false
}

// View renders the visible portion of the list.
func (s Select) View(width int, tone Tone) string {
	if width < 8 {
		width = 8
	}
	height := s.Height
	if height < 1 {
		height = len(s.Options)
	}
	if s.Offset < 0 {
		s.Offset = 0
	}
	lines := make([]string, 0, height)
	for n := 0; n < height; n++ {
		index := s.Offset + n
		if index >= len(s.Options) {
			break
		}
		marker := "  "
		itemTone := Muted
		if index == s.Selected {
			marker = "▸ "
			itemTone = tone
		}
		label := s.Options[index].Label
		if s.Options[index].Description != "" {
			label += " · " + s.Options[index].Description
		}
		lines = append(lines, paint(marker+truncate(label, width-2), itemTone))
	}
	return strings.Join(lines, "\n")
}

// Tabs renders a compact tab strip and clamps the active index.
func Tabs(labels []string, active, width int) string {
	if active < 0 {
		active = 0
	}
	if active >= len(labels) && len(labels) > 0 {
		active = len(labels) - 1
	}
	parts := make([]string, len(labels))
	for n, label := range labels {
		tone := Muted
		if n == active {
			tone = Accent
			label = "[" + label + "]"
		}
		parts[n] = paint(label, tone)
	}
	result := strings.Join(parts, "  ")
	if width > 0 {
		result = truncate(result, width)
	}
	return result
}

// Table renders aligned rows. Width is the total visible width budget.
func Table(headers []string, rows [][]string, width int) string {
	if len(headers) == 0 {
		return ""
	}
	cols := len(headers)
	widths := make([]int, cols)
	for n, h := range headers {
		widths[n] = len([]rune(h))
	}
	for _, row := range rows {
		for n := 0; n < len(row) && n < cols; n++ {
			if w := len([]rune(row[n])); w > widths[n] {
				widths[n] = w
			}
		}
	}
	total := 3*cols + 1
	for _, w := range widths {
		total += w
	}
	if width > 0 && total > width {
		for total > width {
			largest := 0
			for n := range widths {
				if widths[n] > widths[largest] {
					largest = n
				}
			}
			if widths[largest] <= 4 {
				break
			}
			widths[largest]--
			total--
		}
	}
	line := func(row []string) string {
		cells := make([]string, cols)
		for n := range cells {
			value := ""
			if n < len(row) {
				value = row[n]
			}
			cells[n] = " " + pad(truncate(value, widths[n]), widths[n]) + " "
		}
		return "│" + strings.Join(cells, "│") + "│"
	}
	border := func(left, mid, right, fill string) string {
		cells := make([]string, cols)
		for n, w := range widths {
			cells[n] = strings.Repeat(fill, w+2)
		}
		return left + strings.Join(cells, mid) + right
	}
	lines := []string{border("┌", "┬", "┐", "─"), paint(line(headers), Accent), border("├", "┼", "┤", "─")}
	for _, row := range rows {
		lines = append(lines, line(row))
	}
	lines = append(lines, border("└", "┴", "┘", "─"))
	return strings.Join(lines, "\n")
}

// Panel wraps multiple lines with a title and a consistent border.
func Panel(title, content string, width int, tone Tone) string {
	if width < 12 {
		width = 12
	}
	inner := width - 4
	lines := strings.Split(content, "\n")
	title = truncate(title, width-6)
	out := []string{"╭─ " + paint(title, tone) + " " + strings.Repeat("─", maxInt(0, width-len([]rune(title))-6)) + "╮"}
	for _, line := range lines {
		out = append(out, "│ "+pad(truncate(line, inner), inner)+" │")
	}
	out = append(out, "╰"+strings.Repeat("─", width-2)+"╯")
	return strings.Join(out, "\n")
}

func truncate(value string, width int) string {
	if width <= 0 {
		return ""
	}
	r := []rune(strip(value))
	if len(r) <= width {
		return string(r)
	}
	if width == 1 {
		return "…"
	}
	return string(r[:width-1]) + "…"
}
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// KeyHint is a compact shortcut label suitable for footers and command bars.
func KeyHint(key, label string) string {
	return fmt.Sprintf("%s %s", paint("["+key+"]", Accent), paint(label, Muted))
}

// Sparkline renders a compact trend line using one terminal cell per sample.
// Values are sampled evenly when width is smaller than the input length.
func Sparkline(values []float64, width int, tone Tone) string {
	if width < 1 || len(values) == 0 {
		return ""
	}
	if width > len(values) {
		width = len(values)
	}
	minimum, maximum := values[0], values[0]
	for _, value := range values[1:] {
		if value < minimum {
			minimum = value
		}
		if value > maximum {
			maximum = value
		}
	}
	levels := []rune("▁▂▃▄▅▆▇█")
	result := make([]rune, width)
	for i := range result {
		index := i * (len(values) - 1) / maxInt(width-1, 1)
		value := values[index]
		level := 0
		if maximum > minimum {
			level = int((value - minimum) / (maximum - minimum) * float64(len(levels)-1))
		}
		result[i] = levels[maxInt(0, minInt(level, len(levels)-1))]
	}
	return paint(string(result), tone)
}

// Breadcrumb renders a compact navigation path, highlighting the active item.
func Breadcrumb(items []string, active int) string {
	parts := make([]string, len(items))
	for i, item := range items {
		tone := Muted
		if i == active {
			tone = Accent
		}
		parts[i] = paint(item, tone)
	}
	return strings.Join(parts, paint(" › ", Muted))
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
