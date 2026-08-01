package chart

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/pol-cova/termkit-go/style"
)

const (
	brailleDotW = 2
	brailleDotH = 4
)

type brailleCanvas struct {
	// dots stores OR-ed braille bitmasks per character cell; color per cell (last wins for dots)
	chars  [][]rune
	colors [][]int
	w, h   int
}

func newBrailleCanvas(charW, charH int) *brailleCanvas {
	chars := make([][]rune, charH)
	colors := make([][]int, charH)
	for y := 0; y < charH; y++ {
		chars[y] = make([]rune, charW)
		colors[y] = make([]int, charW)
		for x := 0; x < charW; x++ {
			chars[y][x] = '\u2800'
		}
	}
	return &brailleCanvas{chars: chars, colors: colors, w: charW, h: charH}
}

func brailleBit(dx, dy int) rune {
	switch {
	case dy == 0 && dx == 0:
		return 0x01
	case dy == 1 && dx == 0:
		return 0x02
	case dy == 2 && dx == 0:
		return 0x04
	case dy == 3 && dx == 0:
		return 0x40
	case dy == 0 && dx == 1:
		return 0x08
	case dy == 1 && dx == 1:
		return 0x10
	case dy == 2 && dx == 1:
		return 0x20
	case dy == 3 && dx == 1:
		return 0x80
	}
	return 0
}

func (c *brailleCanvas) setDot(px, py int, color int) {
	if px < 0 || py < 0 {
		return
	}
	cx, cy := px/brailleDotW, py/brailleDotH
	if cx >= c.w || cy >= c.h {
		return
	}
	dx, dy := px%brailleDotW, py%brailleDotH
	c.chars[cy][cx] |= brailleBit(dx, dy)
	if color != 0 {
		c.colors[cy][cx] = color
	}
}

func (c *brailleCanvas) line(x0, y0, x1, y1 int, color int) {
	dx := abs(x1 - x0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := -abs(y1 - y0)
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		c.setDot(x0, y0, color)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func (c *brailleCanvas) render(color bool) []string {
	out := make([]string, c.h)
	for y := 0; y < c.h; y++ {
		var b strings.Builder
		last := -1
		for x := 0; x < c.w; x++ {
			ch := c.chars[y][x]
			col := c.colors[y][x]
			if color && col != last {
				if last != -1 {
					b.WriteString("\x1b[0m")
				}
				if col != 0 {
					b.WriteString(fmt.Sprintf("\x1b[38;5;%dm", col))
				}
				last = col
			}
			b.WriteRune(ch)
		}
		if color && last != -1 {
			b.WriteString("\x1b[0m")
		}
		out[y] = b.String()
	}
	return out
}

func styleResolve(name string, index int) int {
	return style.Resolve(name, index).ANSI256(false)
}

func monitorSeriesColor(index int) int {
	return style.MonitorSeriesColor(index, false)
}

func seriesMonitorColor(s Series, index int, color bool, monoStyle int) int {
	if !color {
		_ = monoStyle
		return []int{252, 245, 250}[index%3]
	}
	if s.Color != "" {
		return styleResolve(s.Color, index)
	}
	return monitorSeriesColor(index)
}

func renderBrailleFull(ts TimeSeries, opt BrailleOptions) string {
	yAxisW := 5
	if !opt.ShowAxes {
		yAxisW = 0
	}
	xAxisH := 1
	if !opt.ShowAxes {
		xAxisH = 0
	}
	plotCharW := opt.Width - yAxisW
	plotCharH := opt.Height - xAxisH
	if plotCharW < 4 {
		plotCharW = 4
	}
	if plotCharH < 3 {
		plotCharH = 3
	}

	dotW := plotCharW * brailleDotW
	dotH := plotCharH * brailleDotH
	canvas := newBrailleCanvas(plotCharW, plotCharH)

	yMin, yMax := ts.yRange()
	points := ts.pointCount()

	// faint horizontal grid
	gridColor := styleResolve(string(style.MonitorGrid), 0)
	if opt.Color {
		for row := 0; row <= 4; row++ {
			py := row * (dotH - 1) / 4
			for px := 0; px < dotW; px++ {
				canvas.setDot(px, dotH-1-py, gridColor)
			}
		}
	}

	for si, series := range ts.Series {
		col := seriesMonitorColor(series, si, opt.Color, si)
		prevX, prevY := -1, -1
		for i, v := range series.Values {
			if math.IsNaN(v) {
				prevX, prevY = -1, -1
				continue
			}
			px := 0
			if points > 1 {
				px = i * (dotW - 1) / (points - 1)
			}
			ratio := (v - yMin) / (yMax - yMin)
			py := int(ratio*float64(dotH-1) + 0.5)
			if py < 0 {
				py = 0
			}
			if py >= dotH {
				py = dotH - 1
			}
			py = dotH - 1 - py
			if prevX >= 0 {
				canvas.line(prevX, prevY, px, py, col)
			} else {
				canvas.setDot(px, py, col)
			}
			prevX, prevY = px, py
		}
	}

	// scrub line
	if opt.Selected >= 0 && points > 0 {
		sx := opt.Selected * (dotW - 1) / max(1, points-1)
		scrubColor := 73
		if opt.Color {
			scrubColor = styleResolve(string(style.MonitorSecondary), 0)
		}
		for py := 0; py < dotH; py++ {
			canvas.setDot(sx, py, scrubColor)
		}
	}

	plotLines := canvas.render(opt.Color)

	// Y axis labels
	if opt.ShowAxes {
		formatter := opt.Formatter
		if ts.Y.Formatter != nil {
			formatter = ts.Y.Formatter
		}
		for i, line := range plotLines {
			ratio := float64(len(plotLines)-1-i) / float64(max(1, len(plotLines)-1))
			val := yMin + (yMax-yMin)*ratio
			label := formatter(val)
			if i == len(plotLines)-1 {
				label = formatter(yMin)
			}
			if len(label) > yAxisW {
				label = label[:yAxisW]
			}
			plotLines[i] = paint(fmt.Sprintf("%*s ", yAxisW, label), mutedColor, opt.Color) + line
		}
	}

	var b strings.Builder
	if ts.Title != "" {
		b.WriteString(bold(ts.Title, opt.Color) + "\n")
	}

	// legend overlay (reserve first lines inside plot)
	if opt.ShowLegend {
		legendLines := buildLegendLines(ts, opt)
		for i, ll := range legendLines {
			if i >= len(plotLines) {
				break
			}
			// overlay legend on plot line with padding
			padded := " " + ll
			if len(padded) < len(stripANSI(plotLines[i]))+1 {
				plotLines[i] = padded + plotLines[i][len(padded):]
			}
		}
	}

	b.WriteString(strings.Join(plotLines, "\n"))

	if opt.ShowAxes {
		window := ts.Window
		if window <= 0 {
			window = 60 * time.Second
		}
		xLabels := fmt.Sprintf("%*s%s", yAxisW, "", xAxisLabels(window, plotCharW))
		b.WriteString("\n" + paint(xLabels, mutedColor, opt.Color))
	}

	if opt.Selected >= 0 {
		b.WriteString("\n" + brailleReadout(ts, opt))
	}

	return strings.TrimRight(b.String(), "\n")
}

func buildLegendLines(ts TimeSeries, opt BrailleOptions) []string {
	formatter := opt.Formatter
	if ts.Y.Formatter != nil {
		formatter = ts.Y.Formatter
	}
	entries := ts.Legend.Entries
	if len(entries) == 0 {
		entries = make([]LegendEntry, len(ts.Series))
		for i, s := range ts.Series {
			cur := math.NaN()
			if len(s.Values) > 0 {
				cur = s.Values[len(s.Values)-1]
			}
			entries[i] = LegendEntry{
				Name:    s.Name,
				Current: formatter(cur),
			}
		}
	}
	lines := make([]string, 0, len(entries))
	for i, e := range entries {
		col := seriesMonitorColor(ts.Series[min(i, len(ts.Series)-1)], i, opt.Color, i)
		part := e.Name + ": " + e.Current
		if e.Extra != "" {
			part += " " + e.Extra
		}
		lines = append(lines, paint(part, col, opt.Color))
	}
	return lines
}

func xAxisLabels(window time.Duration, width int) string {
	if width < 8 {
		return strings.Repeat(" ", width)
	}
	left := fmt.Sprintf("%ds", int(window.Seconds()))
	right := "0s"
	mid := fmt.Sprintf("%ds", int(window.Seconds())/2)
	inner := width - displayLen(left) - displayLen(right)
	if inner < 1 {
		return left + right
	}
	gap := inner / 2
	return left + strings.Repeat(" ", gap) + mid + strings.Repeat(" ", inner-gap-displayLen(mid)) + right
}

func brailleReadout(ts TimeSeries, opt BrailleOptions) string {
	formatter := opt.Formatter
	if ts.Y.Formatter != nil {
		formatter = ts.Y.Formatter
	}
	idx := opt.Selected
	parts := make([]string, len(ts.Series))
	for i, s := range ts.Series {
		v := math.NaN()
		if idx >= 0 && idx < len(s.Values) {
			v = s.Values[idx]
		}
		parts[i] = s.Name + ": " + formatter(v)
	}
	label := fmt.Sprintf("t=%d", idx)
	if ts.Window > 0 && ts.pointCount() > 1 {
		secs := int(ts.Window.Seconds()) * (ts.pointCount() - 1 - idx) / (ts.pointCount() - 1)
		label = fmt.Sprintf("%ds ago", secs)
	}
	return paint("◆", selectionColor, opt.Color) + " " + label + "  " + strings.Join(parts, "  ")
}

func renderBrailleCompact(ts TimeSeries, opt BrailleOptions) string {
	var b strings.Builder
	if ts.Title != "" {
		b.WriteString(bold(ts.Title, opt.Color) + "\n")
	}
	width := opt.Width
	if width < 8 {
		width = 8
	}
	levels := []rune("▁▂▃▄▅▆▇█")
	for si, series := range ts.Series {
		col := seriesMonitorColor(series, si, opt.Color, si)
		line := compactSparkline(series.Values, width, levels)
		name := series.Name
		if len(name) > 6 {
			name = name[:6]
		}
		cur := math.NaN()
		if len(series.Values) > 0 {
			cur = series.Values[len(series.Values)-1]
		}
		fmt.Fprintf(&b, "%s %s %s\n", paint(name+":", col, opt.Color), paint(line, col, opt.Color), paint(opt.Formatter(cur), col, opt.Color))
	}
	return strings.TrimRight(b.String(), "\n")
}

func compactSparkline(values []float64, width int, levels []rune) string {
	if len(values) == 0 || width < 1 {
		return ""
	}
	minimum, maximum := math.Inf(1), math.Inf(-1)
	for _, v := range values {
		if math.IsNaN(v) {
			continue
		}
		if v < minimum {
			minimum = v
		}
		if v > maximum {
			maximum = v
		}
	}
	if math.IsInf(minimum, 1) {
		return strings.Repeat(string(levels[0]), width)
	}
	out := make([]rune, width)
	for i := range out {
		idx := i * (len(values) - 1) / max(1, width-1)
		v := values[idx]
		level := 0
		if maximum > minimum && !math.IsNaN(v) {
			level = int((v - minimum) / (maximum - minimum) * float64(len(levels)-1))
		}
		if level < 0 {
			level = 0
		}
		if level >= len(levels) {
			level = len(levels) - 1
		}
		out[i] = levels[level]
	}
	return string(out)
}

func stripANSI(s string) string {
	var b strings.Builder
	esc := false
	for _, r := range s {
		if esc {
			if r == 'm' {
				esc = false
			}
			continue
		}
		if r == '\x1b' {
			esc = true
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func displayLen(s string) int {
	return len(stripANSI(s))
}
