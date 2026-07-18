package chart

import (
	"fmt"
	"math"
	"strings"
)

type ValueFormatter func(float64) string

func DefaultFormatter(v float64) string {
	if math.IsNaN(v) {
		return "—"
	}
	return fmt.Sprintf("%.2f", v)
}

func PercentFormatter(v float64) string {
	if math.IsNaN(v) {
		return "—"
	}
	return fmt.Sprintf("%.0f%%", v)
}

func SIFormatter(v float64) string {
	if math.IsNaN(v) {
		return "—"
	}
	switch {
	case math.Abs(v) >= 1e9:
		return fmt.Sprintf("%.1fG", v/1e9)
	case math.Abs(v) >= 1e6:
		return fmt.Sprintf("%.1fM", v/1e6)
	case math.Abs(v) >= 1e3:
		return fmt.Sprintf("%.1fK", v/1e3)
	default:
		return fmt.Sprintf("%.0f", v)
	}
}

type SeriesReading struct {
	Name    string
	Value   float64
	Missing bool
}

type SelectedPoint struct {
	Index  int
	Label  string
	Series []SeriesReading
}

func (c Chart) SelectedPoint(o Options) SelectedPoint {
	o = normalized(o, c)
	label := fmt.Sprintf("point %d", o.Selected+1)
	if o.Selected >= 0 && o.Selected < len(c.Labels) && c.Labels[o.Selected] != "" {
		label = c.Labels[o.Selected]
	}
	readings := make([]SeriesReading, len(c.Series))
	for i, s := range c.Series {
		value := math.NaN()
		if o.Selected >= 0 && o.Selected < len(s.Values) {
			value = s.Values[o.Selected]
		}
		readings[i] = SeriesReading{
			Name:    s.Name,
			Value:   value,
			Missing: math.IsNaN(value),
		}
	}
	return SelectedPoint{Index: o.Selected, Label: label, Series: readings}
}

func Readout(c Chart, o Options) string {
	formatter := o.Formatter
	if formatter == nil {
		formatter = DefaultFormatter
	}
	sp := c.SelectedPoint(o)
	parts := make([]string, len(sp.Series))
	for i, s := range sp.Series {
		parts[i] = s.Name + ": " + formatter(s.Value)
	}
	return paint("◆", selectionColor, o.Color) + " " + sp.Label + "  " + strings.Join(parts, "  ")
}
