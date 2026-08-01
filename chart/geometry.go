package chart

import "math"

type Layout struct {
	Chart    Chart
	MaxValue float64
	Points   int
}

func NewLayout(c Chart) Layout {
	return Layout{
		Chart:    c,
		MaxValue: chartMaximum(c),
		Points:   c.pointCount(),
	}
}

func (l Layout) PointX(index, width int) int { return pointX(index, l.Points, width) }

func (l Layout) PlotY(value float64, height int) int { return plotY(value, l.MaxValue, height) }

func (l Layout) SampledValue(seriesIndex, x, width int) float64 {
	return sampledValue(l.Chart, seriesIndex, x, width)
}

func (l Layout) SampledTotal(x, width int) float64 { return sampledTotal(l.Chart, x, width) }

func (l Layout) SampledBase(seriesIndex, x, width int) float64 {
	return sampledBase(l.Chart, seriesIndex, x, width)
}

func (l Layout) PointTotal(point int) float64 { return pointTotal(l.Chart, point) }

func (l Layout) Angle(i int) float64 { return angle(i, l.Points) }

func (l Layout) RadarPoint(cx, cy int, rx, ry, radius, axisIndex float64) [2]int {
	a := angle(int(axisIndex), l.Points)
	return radarPoint(cx, cy, rx*radius, ry*radius, a)
}

func pointX(i, n, w int) int {
	if n < 2 {
		return w / 2
	}
	return int(math.Round(float64(i) * float64(w-1) / float64(n-1)))
}

func angle(i, n int) float64 { return -math.Pi/2 + 2*math.Pi*float64(i)/float64(n) }

func radarPoint(cx, cy int, rx, ry, a float64) [2]int {
	return [2]int{cx + int(math.Round(math.Cos(a)*rx)), cy + int(math.Round(math.Sin(a)*ry))}
}

func plotY(value, maximum float64, height int) int {
	if math.IsNaN(value) {
		value = 0
	}
	y := height - 1 - int(math.Round(value/maximum*float64(height-1)))
	return min(height-1, max(0, y))
}

func sampledValue(c Chart, series, x, width int) float64 {
	values := c.Series[series].Values
	if len(values) == 1 {
		return finiteOrZero(values[0])
	}
	position := float64(x) * float64(len(values)-1) / float64(max(1, width-1))
	left := int(math.Floor(position))
	right := min(left+1, len(values)-1)
	leftVal := finiteOrZero(values[left])
	rightVal := finiteOrZero(values[right])
	if math.IsNaN(values[left]) && math.IsNaN(values[right]) {
		return math.NaN()
	}
	if math.IsNaN(values[left]) {
		return rightVal
	}
	if math.IsNaN(values[right]) {
		return leftVal
	}
	return leftVal + (rightVal-leftVal)*(position-float64(left))
}

func finiteOrZero(v float64) float64 {
	if math.IsNaN(v) {
		return 0
	}
	return v
}

func pointTotal(c Chart, point int) float64 {
	total := 0.0
	for _, series := range c.Series {
		if point < len(series.Values) && !math.IsNaN(series.Values[point]) {
			total += maxFloat(series.Values[point], 0)
		}
	}
	return total
}

func sampledTotal(c Chart, x, width int) float64 {
	total := 0.0
	for i := range c.Series {
		v := sampledValue(c, i, x, width)
		if !math.IsNaN(v) {
			total += v
		}
	}
	return total
}

func sampledBase(c Chart, until, x, width int) float64 {
	base := 0.0
	for i := 0; i < until; i++ {
		v := sampledValue(c, i, x, width)
		if !math.IsNaN(v) {
			base += v
		}
	}
	return base
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

func maxSeries(ss []Series) float64 {
	m := 0.0
	for _, s := range ss {
		for _, v := range s.Values {
			if !math.IsNaN(v) && v > m {
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
