package component

import (
	"strings"
)

// Widget is a composable terminal surface that renders into an exact box.
type Widget interface {
	SizeHint() (w, h int)
	View(width, height int) string
}

// Slot places a widget inside a Grid cell with optional spanning and weights.
type Slot struct {
	Widget Widget
	Row, Col, RowSpan, ColSpan int
	MinW, MinH, WeightW, WeightH int
}

type gridRegion struct {
	x, y, w, h int
	view       []string
}

// Grid lays out widgets in a rows×cols matrix with weighted sizing.
func Grid(width, height, rows, cols int, slots ...Slot) string {
	if width < 1 || height < 1 || rows < 1 || cols < 1 {
		return ""
	}
	rowHeights := distribute(height, rows, slots, true)
	colWidths := distribute(width, cols, slots, false)

	regions := make([]gridRegion, 0, len(slots))
	for _, slot := range slots {
		if slot.Widget == nil {
			continue
		}
		r0, c0 := slot.Row, slot.Col
		rSpan, cSpan := maxInt(1, slot.RowSpan), maxInt(1, slot.ColSpan)
		cellW, cellH := 0, 0
		for c := c0; c < c0+cSpan && c < cols; c++ {
			cellW += colWidths[c]
		}
		for r := r0; r < r0+rSpan && r < rows; r++ {
			cellH += rowHeights[r]
		}
		if cellW < 1 {
			cellW = 1
		}
		if cellH < 1 {
			cellH = 1
		}
		view := slot.Widget.View(cellW, cellH)
		regions = append(regions, gridRegion{
			x: colOffset(colWidths, c0), y: rowOffset(rowHeights, r0),
			w: cellW, h: cellH, view: splitLines(view, cellH, cellW),
		})
	}

	out := make([]string, height)
	for y := 0; y < height; y++ {
		out[y] = strings.Repeat(" ", width)
	}
	for _, region := range regions {
		for dy := 0; dy < region.h && region.y+dy < height; dy++ {
			line := region.view[dy]
			out[region.y+dy] = mergeLine(out[region.y+dy], line, region.x, region.w, width)
		}
	}
	return strings.Join(out, "\n")
}

func splitLines(view string, height, width int) []string {
	raw := strings.Split(view, "\n")
	out := make([]string, height)
	for i := 0; i < height; i++ {
		line := ""
		if i < len(raw) {
			line = raw[i]
		}
		out[i] = fitDisplayWidth(line, width)
	}
	return out
}

func mergeLine(base, overlay string, x, w, total int) string {
	overlay = fitDisplayWidth(overlay, w)
	if x >= total {
		return base
	}
	end := x + w
	if end > total {
		end = total
	}
	baseRunes := []rune(base)
	return string(baseRunes[:x]) + overlay + string(baseRunes[end:])
}

func distribute(total, count int, slots []Slot, vertical bool) []int {
	sizes := make([]int, count)
	for i := range sizes {
		sizes[i] = 1
	}
	if total <= count {
		return sizes
	}
	remaining := total - count
	weights := make([]int, count)
	for i := range weights {
		weights[i] = 1
	}
	for _, slot := range slots {
		if vertical {
			for r := slot.Row; r < slot.Row+maxInt(1, slot.RowSpan) && r < count; r++ {
				if slot.WeightH > 0 {
					weights[r] = slot.WeightH
				}
			}
		} else {
			for c := slot.Col; c < slot.Col+maxInt(1, slot.ColSpan) && c < count; c++ {
				if slot.WeightW > 0 {
					weights[c] = slot.WeightW
				}
			}
		}
	}
	sumWeight := 0
	for _, wt := range weights {
		sumWeight += wt
	}
	for i := 0; i < count; i++ {
		share := remaining * weights[i] / maxInt(1, sumWeight)
		sizes[i] += share
	}
	sum := 0
	for _, s := range sizes {
		sum += s
	}
	for sum > total {
		for i := count - 1; i >= 0 && sum > total; i-- {
			if sizes[i] > 1 {
				sizes[i]--
				sum--
			}
		}
	}
	for sum < total {
		for i := 0; sum < total; i = (i + 1) % count {
			sizes[i]++
			sum++
		}
	}
	return sizes
}

func rowOffset(heights []int, row int) int {
	o := 0
	for i := 0; i < row && i < len(heights); i++ {
		o += heights[i]
	}
	return o
}

func colOffset(widths []int, col int) int {
	o := 0
	for i := 0; i < col && i < len(widths); i++ {
		o += widths[i]
	}
	return o
}

// BorderTitle wraps body content in a box with the title integrated into the top border.
func BorderTitle(title, body string, width, height int, tone Tone) string {
	if width < 4 {
		width = 4
	}
	if height < 3 {
		height = 3
	}
	innerW := width - 2
	innerH := height - 2
	title = truncateDisplay(title, innerW-2)
	top := "┌─ " + paint(title, tone) + " " + strings.Repeat("─", maxInt(0, innerW-3-displayWidth(title))) + "┐"

	bodyLines := strings.Split(body, "\n")
	content := make([]string, innerH)
	for i := 0; i < innerH; i++ {
		line := ""
		if i < len(bodyLines) {
			line = bodyLines[i]
		}
		content[i] = "│" + padDisplay(truncateDisplay(line, innerW), innerW) + "│"
	}
	bottom := "└" + strings.Repeat("─", innerW) + "┘"
	out := append([]string{top}, content...)
	out = append(out, bottom)
	return strings.Join(out, "\n")
}

// StaticWidget renders fixed content and implements Widget.
type StaticWidget struct {
	Content string
	HintW   int
	HintH   int
}

func (s StaticWidget) SizeHint() (w, h int) {
	return s.HintW, s.HintH
}

func (s StaticWidget) View(width, height int) string {
	lines := strings.Split(s.Content, "\n")
	out := make([]string, height)
	for i := 0; i < height; i++ {
		line := ""
		if i < len(lines) {
			line = lines[i]
		}
		out[i] = padDisplay(truncateDisplay(line, width), width)
	}
	return strings.Join(out, "\n")
}

// FuncWidget adapts a render function to Widget.
type FuncWidget struct {
	Render func(width, height int) string
	MinW   int
	MinH   int
}

func (f FuncWidget) SizeHint() (w, h int) {
	return f.MinW, f.MinH
}

func (f FuncWidget) View(width, height int) string {
	if f.Render == nil {
		return strings.Repeat(" ", width)
	}
	return f.Render(width, height)
}
