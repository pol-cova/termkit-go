package ditherkit

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
	maxValue := maxSeries(c.Series)
	if maxValue == 0 {
		maxValue = 1
	}
	for y := 0; y < o.Height; y++ {
		for x := 0; x < o.Width; x++ {
			if y == o.Height-1 {
				grid.put(x, y, '─', mutedColor)
			} else if y%3 == 0 {
				grid.put(x, y, '·', mutedColor)
			}
		}
	}
	for si, s := range c.Series {
		for i, v := range s.Values {
			x := pointX(i, len(s.Values), o.Width)
			y := o.Height - 1 - int(math.Round((v/maxValue)*float64(o.Height-1)))
			if bars {
				barWidth := max(1, o.Width/(len(s.Values)*len(c.Series)+1))
				start := x - (len(c.Series)*barWidth)/2 + si*barWidth
				for bx := start; bx < start+barWidth; bx++ {
					for by := y; by < o.Height-1; by++ {
						grid.put(bx, by, dither(bx, by, si), seriesColor(si))
					}
				}
				continue
			}
			if area {
				for by := y; by < o.Height-1; by++ {
					grid.put(x, by, dither(x, by, si), seriesColor(si))
				}
			}
			grid.put(x, y, '●', seriesColor(si))
			if i > 0 {
				px := pointX(i-1, len(s.Values), o.Width)
				py := o.Height - 1 - int(math.Round((s.Values[i-1]/maxValue)*float64(o.Height-1)))
				grid.line(px, py, x, y, '•', seriesColor(si))
			}
		}
	}
	sx := pointX(o.Selected, c.pointCount(), o.Width)
	for y := 0; y < o.Height; y++ {
		if grid.cells[y][sx].r == ' ' || grid.cells[y][sx].r == '·' {
			grid.put(sx, y, '┊', selectionColor)
		}
	}
	return grid.lines(o.Color)
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
	rx, ry := float64(o.Width)/2.5, float64(o.Height)/2.1
	for y := 0; y < o.Height; y++ {
		for x := 0; x < o.Width; x++ {
			dx, dy := (float64(x)-cx)/rx, (float64(y)-cy)/ry
			if dx*dx+dy*dy > 1 {
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
			glyph := dither(x, y, index)
			if index == o.Selected {
				glyph = '█'
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
	radius := float64(min(o.Width, o.Height)) / 2.3
	maxValue := maxSeries(c.Series)
	if maxValue == 0 {
		maxValue = 1
	}
	for ring := 1; ring <= 3; ring++ {
		r := radius * float64(ring) / 3
		for i := 0; i < n; i++ {
			a := angle(i, n)
			x, y := cx+int(math.Round(math.Cos(a)*r)), cy+int(math.Round(math.Sin(a)*r*.55))
			g.put(x, y, '·', mutedColor)
		}
	}
	for si, s := range c.Series {
		points := make([][2]int, n)
		for i, v := range s.Values {
			a := angle(i, n)
			r := radius * v / maxValue
			points[i] = [2]int{cx + int(math.Round(math.Cos(a)*r)), cy + int(math.Round(math.Sin(a)*r*.55))}
		}
		for i := range points {
			next := points[(i+1)%n]
			g.line(points[i][0], points[i][1], next[0], next[1], dither(i, si, si), seriesColor(si))
			g.put(points[i][0], points[i][1], '●', seriesColor(si))
		}
	}
	return g.lines(o.Color)
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
	p := make([]string, len(c.Series))
	for i, s := range c.Series {
		p[i] = paint("●", seriesColor(i), o.Color) + " " + s.Name
	}
	return strings.Join(p, "  ")
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

func seriesColor(i int) int { return []int{81, 213, 42, 208, 141}[i%5] }
func paint(s string, c int, on bool) string {
	if !on {
		return s
	}
	return fmt.Sprintf("\x1b[38;5;%dm%s\x1b[0m", c, s)
}
