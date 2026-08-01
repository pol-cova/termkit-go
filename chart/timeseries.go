package chart

import (
	"fmt"
	"math"
	"os"
	"time"
)

// AxisSpec configures axis ticks and value formatting.
type AxisSpec struct {
	Min, Max float64 // both zero = auto-scale from data
	Formatter ValueFormatter
}

// LegendEntry is one row inside the in-plot legend box.
type LegendEntry struct {
	Name    string
	Current string
	Extra   string
}

// LegendSpec configures the overlay legend. Empty Entries auto-builds from series.
type LegendSpec struct {
	Entries []LegendEntry
}

// TimeSeries is renderer-independent input for braille monitor graphs.
type TimeSeries struct {
	Title  string
	Window time.Duration
	Series []Series
	Y      AxisSpec
	Legend LegendSpec
}

// BrailleOptions controls braille time-series rendering.
type BrailleOptions struct {
	Width, Height int
	Color         bool
	ShowAxes      bool
	ShowLegend    bool
	Selected      int
	Animate       bool
	Formatter     ValueFormatter
}

const (
	brailleCompactWidth  = 24
	brailleCompactHeight = 6
)

// RenderBraille draws a btm-style braille time-series graph.
func RenderBraille(ts TimeSeries, opt BrailleOptions) (string, error) {
	if len(ts.Series) == 0 {
		return "", fmt.Errorf("termkit/chart: time series requires at least one series")
	}
	for _, s := range ts.Series {
		if len(s.Values) == 0 {
			return "", fmt.Errorf("termkit/chart: series %q has no values", s.Name)
		}
	}
	opt = normalizeBraille(opt, ts)
	if opt.Width < brailleCompactWidth || opt.Height < brailleCompactHeight {
		return renderBrailleCompact(ts, opt), nil
	}
	return renderBrailleFull(ts, opt), nil
}

func normalizeBraille(opt BrailleOptions, ts TimeSeries) BrailleOptions {
	if os.Getenv("NO_COLOR") != "" {
		opt.Color = false
	}
	if opt.Width < 8 {
		opt.Width = 44
	}
	if opt.Height < 4 {
		opt.Height = 10
	}
	if opt.Formatter == nil {
		opt.Formatter = DefaultFormatter
	}
	points := ts.pointCount()
	if opt.Selected < -1 {
		opt.Selected = -1
	}
	if opt.Selected >= points {
		opt.Selected = points - 1
	}
	return opt
}

func (ts TimeSeries) pointCount() int {
	max := 0
	for _, s := range ts.Series {
		if len(s.Values) > max {
			max = len(s.Values)
		}
	}
	return max
}

func (ts TimeSeries) yRange() (min, max float64) {
	min, max = math.Inf(1), math.Inf(-1)
	if ts.Y.Min != 0 || ts.Y.Max != 0 {
		return ts.Y.Min, ts.Y.Max
	}
	for _, s := range ts.Series {
		for _, v := range s.Values {
			if math.IsNaN(v) {
				continue
			}
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
	}
	if !math.IsInf(min, 1) && min == max {
		min -= 1
		max += 1
	}
	if math.IsInf(min, 1) {
		return 0, 100
	}
	padding := (max - min) * 0.05
	if padding == 0 {
		padding = 1
	}
	return min - padding, max + padding
}
