package chart

import (
	"fmt"
	"math"
	"strings"
)

// Options controls the terminal canvas. Width and Height exclude the title,
// legend, and selection readout.
type Options struct {
	Width, Height int
	Selected      int
	ActiveSeries  int
	Color         bool
}

// Render draws a chart using Unicode cells and ANSI colour when requested.
func Render(chart Chart, options Options) (string, error) {
	if err := chart.validate(); err != nil {
		return "", err
	}
	options = normalized(options, chart)
	var plot []string
	switch chart.Kind {
	case Area:
		plot = renderCartesian(chart, options, true, false)
	case Line:
		plot = renderCartesian(chart, options, false, false)
	case Bar:
		plot = renderCartesian(chart, options, true, true)
	case Pie:
		plot = renderPie(chart, options)
	case Radar:
		plot = renderRadar(chart, options)
	}
	var b strings.Builder
	if chart.Title != "" {
		b.WriteString(bold(chart.Title, options.Color) + "\n")
	}
	b.WriteString(strings.Join(plot, "\n"))
	if chart.Kind == Area || chart.Kind == Line || chart.Kind == Bar {
		b.WriteString("\n" + xAxis(chart, options))
	}
	if chart.Kind == Radar && len(chart.Labels) > 0 {
		b.WriteString("\n" + radarLabels(chart))
	}
	b.WriteString("\n" + legend(chart, options))
	b.WriteString("\n" + readout(chart, options))
	return b.String(), nil
}

func normalized(o Options, c Chart) Options {
	if o.Width < 12 {
		o.Width = 44
	}
	if o.Height < 4 {
		o.Height = 12
	}
	if o.Selected < 0 {
		o.Selected = 0
	}
	if o.Selected >= c.pointCount() {
		o.Selected = c.pointCount() - 1
	}
	if o.ActiveSeries < 0 || o.ActiveSeries >= len(c.Series) {
		o.ActiveSeries = 0
	}
	return o
}

func renderCartesian(c Chart, o Options, area, bars bool) []string {
	grid := newGrid(o.Width, o.Height)
	maxValue := chartMaximum(c)
	for y := 0; y < o.Height; y++ {
		for x := 0; x < o.Width; x++ {
			if y == o.Height-1 {
				grid.put(x, y, '─', mutedColor)
			} else if y == o.Height/2 {
				grid.put(x, y, '┈', mutedColor)
			}
		}
	}
	if bars {
		renderBars(grid, c, o, maxValue)
	} else {
		renderSeries(grid, c, o, maxValue, area)
	}
	sx := pointX(o.Selected, c.pointCount(), o.Width)
	for y := 0; y < o.Height; y++ {
		if y < o.Height-1 {
			grid.put(sx, y, '┊', selectionColor)
		}
	}
	return grid.lines(o.Color)
}

func renderSeries(g grid, c Chart, o Options, maximum float64, fill bool) {
	stacked := c.StackType == Stacked || c.StackType == Percent
	for seriesIndex, series := range c.Series {
		for x := 0; x < o.Width; x++ {
			value := sampledValue(c, seriesIndex, x, o.Width)
			base := 0.0
			if stacked {
				for before := 0; before < seriesIndex; before++ {
					base += sampledValue(c, before, x, o.Width)
				}
			}
			if c.StackType == Percent {
				total := sampledTotal(c, x, o.Width)
				if total > 0 {
					value = value / total * 100
					base = sampledBase(c, seriesIndex, x, o.Width) / total * 100
				}
			}
			topY := plotY(base+value, maximum, o.Height)
			baseY := plotY(base, maximum, o.Height)
			if fill {
				for y := topY; y < baseY; y++ {
					g.put(x, y, fillRune(series.Variant, x, y, seriesIndex), seriesColor(seriesIndex))
				}
			}
			edge := edgeRune(false)
			if fill {
				edge = '▄'
			}
			g.put(x, topY, edge, seriesColor(seriesIndex))
		}
	}
}

func renderBars(g grid, c Chart, o Options, maximum float64) {
	stacked := c.StackType == Stacked || c.StackType == Percent
	count := c.pointCount()
	barWidth := max(1, o.Width/(count*2))
	if !stacked {
		barWidth = max(1, o.Width/(count*(len(c.Series)+1)))
	}
	for point := 0; point < count; point++ {
		x := pointX(point, count, o.Width)
		rawBase := 0.0
		total := pointTotal(c, point)
		for si, s := range c.Series {
			value := s.Values[point]
			base := rawBase
			if c.StackType == Percent && total > 0 {
				value = value / total * 100
				base = rawBase / total * 100
			}
			topY := plotY(base+value, maximum, o.Height)
			baseY := plotY(base, maximum, o.Height)
			start := x - barWidth/2
			if !stacked {
				start = x - (len(c.Series)*barWidth)/2 + si*barWidth
			}
			for px := start; px < start+barWidth; px++ {
				for y := topY; y < baseY; y++ {
					g.put(px, y, fillRune(s.Variant, px, y, si), seriesColor(si))
				}
				g.put(px, topY, '▀', seriesColor(si))
			}
			if stacked {
				rawBase += s.Values[point]
			}
		}
	}
}

func chartMaximum(c Chart) float64 {
	if c.StackType == Percent {
		return 100
	}
	maximum := maxSeries(c.Series)
	if c.StackType == Stacked {
		maximum = 0
		for point := 0; point < c.pointCount(); point++ {
			if total := pointTotal(c, point); total > maximum {
				maximum = total
			}
		}
	}
	if maximum == 0 {
		return 1
	}
	return maximum
}
func pointTotal(c Chart, point int) float64 {
	total := 0.0
	for _, series := range c.Series {
		total += maxFloat(series.Values[point], 0)
	}
	return total
}
func sampledTotal(c Chart, x, width int) float64 {
	total := 0.0
	for i := range c.Series {
		total += sampledValue(c, i, x, width)
	}
	return total
}
func sampledBase(c Chart, until, x, width int) float64 {
	base := 0.0
	for i := 0; i < until; i++ {
		base += sampledValue(c, i, x, width)
	}
	return base
}
func sampledValue(c Chart, series, x, width int) float64 {
	values := c.Series[series].Values
	if len(values) == 1 {
		return values[0]
	}
	position := float64(x) * float64(len(values)-1) / float64(max(1, width-1))
	left := int(math.Floor(position))
	right := min(left+1, len(values)-1)
	return values[left] + (values[right]-values[left])*(position-float64(left))
}
func plotY(value, maximum float64, height int) int {
	y := height - 1 - int(math.Round(value/maximum*float64(height-1)))
	return min(height-1, max(0, y))
}
func fillRune(variant Variant, x, y, series int) rune {
	switch variant {
	case Solid:
		return '█'
	case Hatched:
		return '▒'
	case Dotted:
		return '░'
	default:
		return []rune{'░', '░', '▒', '▓'}[(x+y+series)%4]
	}
}

func edgeRune(filled bool) rune {
	if filled {
		return '▀'
	}
	return '•'
}

func xAxis(c Chart, o Options) string {
	labels := make([]rune, o.Width)
	for i := range labels {
		labels[i] = ' '
	}
	for i, label := range c.Labels {
		x := pointX(i, len(c.Labels), o.Width)
		start := max(0, x-len([]rune(label))/2)
		for offset, r := range []rune(label) {
			if start+offset < len(labels) {
				labels[start+offset] = r
			}
		}
	}
	return string(labels)
}

func renderPie(c Chart, o Options) []string {
	g := newGrid(o.Width, o.Height)
	s := c.Series[o.ActiveSeries]
	total := 0.0
	for _, v := range s.Values {
		total += maxFloat(v, 0)
	}
	if total == 0 {
		total = 1
	}
	cx, cy := float64(o.Width-1)/2, float64(o.Height-1)/2
	rx, ry := float64(o.Width-2)/2.15, float64(o.Height-1)/2
	for y := 0; y < o.Height; y++ {
		for x := 0; x < o.Width; x++ {
			dx, dy := (float64(x)-cx)/maxFloat(rx, 1), (float64(y)-cy)/maxFloat(ry, 1)
			distance := dx*dx + dy*dy
			if distance > 1 {
				continue
			}
			angle := math.Atan2(dy, dx) + math.Pi
			running := 0.0
			index := 0
			for i, v := range s.Values {
				running += maxFloat(v, 0) / total * 2 * math.Pi
				if angle <= running {
					index = i
					break
				}
			}
			glyph := fillRune(Dotted, x, y, index)
			if index == o.Selected {
				glyph = '▓'
			}
			g.put(x, y, glyph, seriesColor(index))
		}
	}
	return g.lines(o.Color)
}

func renderRadar(c Chart, o Options) []string {
	g := newGrid(o.Width, o.Height)
	n := c.pointCount()
	cx, cy := o.Width/2, o.Height/2
	rx := float64(max(2, o.Width-2)) / 2
	ry := float64(max(2, o.Height-2)) / 2
	maxValue := maxSeries(c.Series)
	if maxValue == 0 {
		maxValue = 1
	}
	for ring := 1; ring <= 3; ring++ {
		scale := float64(ring) / 3
		points := make([][2]int, n)
		for i := range points {
			a := angle(i, n)
			points[i] = radarPoint(cx, cy, rx*scale, ry*scale, a)
		}
		for i, point := range points {
			next := points[(i+1)%n]
			g.line(point[0], point[1], next[0], next[1], '·', mutedColor)
		}
	}
	for i := 0; i < n; i++ {
		a := angle(i, n)
		x, y := radarPoint(cx, cy, rx, ry, a)[0], radarPoint(cx, cy, rx, ry, a)[1]
		g.line(cx, cy, x, y, '·', mutedColor)
	}
	for si, s := range c.Series {
		points := make([][2]int, n)
		for i, v := range s.Values {
			a := angle(i, n)
			r := maxFloat(v, 0) / maxValue
			points[i] = radarPoint(cx, cy, rx*r, ry*r, a)
		}
		for y := 0; y < o.Height; y++ {
			for x := 0; x < o.Width; x++ {
				if insidePolygon(x, y, points) {
					g.put(x, y, fillRune(s.Variant, x, y, si), seriesColor(si))
				}
			}
		}
		for i := range points {
			next := points[(i+1)%n]
			g.line(points[i][0], points[i][1], next[0], next[1], '•', seriesColor(si))
			g.put(points[i][0], points[i][1], '●', seriesColor(si))
		}
	}
	return g.lines(o.Color)
}

func insidePolygon(x, y int, points [][2]int) bool {
	inside := false
	for i, j := 0, len(points)-1; i < len(points); j, i = i, i+1 {
		xi, yi := points[i][0], points[i][1]
		xj, yj := points[j][0], points[j][1]
		crosses := (yi > y) != (yj > y) && x < (xj-xi)*(y-yi)/(yj-yi)+xi
		if crosses {
			inside = !inside
		}
	}
	return inside
}

func pointX(i, n, w int) int {
	if n < 2 {
		return w / 2
	}
	return int(math.Round(float64(i) * float64(w-1) / float64(n-1)))
}
func angle(i, n int) float64 { return -math.Pi/2 + 2*math.Pi*float64(i)/float64(n) }
func maxSeries(ss []Series) float64 {
	m := 0.0
	for _, s := range ss {
		for _, v := range s.Values {
			if v > m {
				m = v
			}
		}
	}
	return m
}
func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func dither(x, y, series int) rune {
	patterns := []rune{'░', '▒', '▓', '█'}
	return patterns[(x*3+y*5+series*2)%len(patterns)]
}

func legend(c Chart, o Options) string {
	if c.Kind == Pie {
		return pieLegend(c, o)
	}
	p := make([]string, len(c.Series))
	for i, s := range c.Series {
		p[i] = paint("●", seriesColor(i), o.Color) + " " + s.Name
	}
	return strings.Join(p, "  ")
}

func pieLegend(c Chart, o Options) string {
	s := c.Series[o.ActiveSeries]
	total := 0.0
	for _, value := range s.Values {
		total += maxFloat(value, 0)
	}
	if total == 0 {
		total = 1
	}
	parts := make([]string, len(s.Values))
	for i, value := range s.Values {
		label := fmt.Sprintf("%s %.0f%%", labelAt(c, i), maxFloat(value, 0)/total*100)
		if i == o.Selected {
			label = "◆ " + label
		} else {
			label = "● " + label
		}
		parts[i] = paint(label, seriesColor(i), o.Color)
	}
	return strings.Join(parts, "  ")
}

func radarLabels(c Chart) string {
	labels := make([]string, 0, len(c.Labels))
	for i, label := range c.Labels {
		labels = append(labels, fmt.Sprintf("%d %s", i+1, label))
	}
	return strings.Join(labels, "  ")
}

func labelAt(c Chart, index int) string {
	if index >= 0 && index < len(c.Labels) && c.Labels[index] != "" {
		return c.Labels[index]
	}
	return fmt.Sprintf("point %d", index+1)
}

func radarPoint(cx, cy int, rx, ry, a float64) [2]int {
	return [2]int{cx + int(math.Round(math.Cos(a)*rx)), cy + int(math.Round(math.Sin(a)*ry))}
}
func readout(c Chart, o Options) string {
	label := fmt.Sprintf("point %d", o.Selected+1)
	if len(c.Labels) > o.Selected {
		label = c.Labels[o.Selected]
	}
	s := c.Series[o.ActiveSeries]
	return paint("◆", selectionColor, o.Color) + " " + label + "  " + s.Name + ": " + fmt.Sprintf("%.2f", s.Values[o.Selected])
}
func bold(s string, color bool) string {
	if !color {
		return s
	}
	return "\x1b[1m" + s + "\x1b[0m"
}

type cell struct {
	r     rune
	color int
}
type grid struct{ cells [][]cell }

func newGrid(w, h int) grid {
	c := make([][]cell, h)
	for y := range c {
		c[y] = make([]cell, w)
		for x := range c[y] {
			c[y][x] = cell{' ', 0}
		}
	}
	return grid{c}
}
func (g grid) put(x, y int, r rune, color int) {
	if y >= 0 && y < len(g.cells) && x >= 0 && x < len(g.cells[0]) {
		g.cells[y][x] = cell{r, color}
	}
}
func (g grid) line(x0, y0, x1, y1 int, r rune, color int) {
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
	e := dx + dy
	for {
		g.put(x0, y0, r, color)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := 2 * e
		if e2 >= dy {
			e += dy
			x0 += sx
		}
		if e2 <= dx {
			e += dx
			y0 += sy
		}
	}
}
func (g grid) lines(color bool) []string {
	out := make([]string, len(g.cells))
	for y, row := range g.cells {
		var b strings.Builder
		last := -1
		for _, c := range row {
			if color && c.color != last {
				if last != -1 {
					b.WriteString("\x1b[0m")
				}
				if c.color != 0 {
					b.WriteString(fmt.Sprintf("\x1b[38;5;%dm", c.color))
				}
				last = c.color
			}
			b.WriteRune(c.r)
		}
		if color && last != -1 {
			b.WriteString("\x1b[0m")
		}
		out[y] = strings.TrimRight(b.String(), " ")
	}
	return out
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

const (
	mutedColor     = 244
	selectionColor = 226
)

func seriesColor(i int) int { return []int{39, 99, 46, 214, 141}[i%5] }
func paint(s string, c int, on bool) string {
	if !on {
		return s
	}
	return fmt.Sprintf("\x1b[38;5;%dm%s\x1b[0m", c, s)
}
