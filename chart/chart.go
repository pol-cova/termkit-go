// Package chart renders dithered, interactive-friendly charts for terminals.
package chart

import (
	"fmt"
	"math"

	"github.com/pol-cova/termkit-go/style"
)

// Kind describes a chart geometry.
type Kind string

const (
	Area  Kind = "area"
	Line  Kind = "line"
	Bar   Kind = "bar"
	Pie   Kind = "pie"
	Radar Kind = "radar"
)

// Variant controls the texture used for a series fill.
type Variant = style.Variant

const (
	Gradient = style.Gradient
	Dotted   = style.Dotted
	Hatched  = style.Hatched
	Solid    = style.Solid
)

// StackType controls how cartesian series share vertical space.
type StackType string

const (
	Default StackType = "default"
	Stacked StackType = "stacked"
	Percent StackType = "percent"
)

// Series supplies one named measure. Values are aligned with Chart.Labels.
type Series struct {
	Name    string
	Values  []float64
	Variant Variant
	// Color optionally selects one of dither-kit's named palette colours.
	// When empty, the stable palette order is used.
	Color string
}

// Chart is renderer-independent chart input.
type Chart struct {
	Kind      Kind
	Title     string
	Labels    []string
	Series    []Series
	StackType StackType
	Bloom     style.Bloom
}

func (c Chart) validate() error {
	if len(c.Series) == 0 {
		return fmt.Errorf("termkit/chart: chart requires at least one series")
	}
	if c.Kind != Area && c.Kind != Line && c.Kind != Bar && c.Kind != Pie && c.Kind != Radar {
		return fmt.Errorf("termkit/chart: unsupported chart kind %q", c.Kind)
	}
	for _, s := range c.Series {
		if len(s.Values) == 0 {
			return fmt.Errorf("termkit/chart: series %q has no values", s.Name)
		}
		allMissing := true
		for _, v := range s.Values {
			if !math.IsNaN(v) {
				allMissing = false
				break
			}
		}
		if allMissing {
			return fmt.Errorf("termkit/chart: series %q has only missing values", s.Name)
		}
		if len(c.Labels) > 0 && len(s.Values) != len(c.Labels) {
			return fmt.Errorf("termkit/chart: series %q has %d values for %d labels", s.Name, len(s.Values), len(c.Labels))
		}
	}
	for i := 1; i < len(c.Series); i++ {
		if len(c.Series[i].Values) != len(c.Series[0].Values) {
			return fmt.Errorf("termkit/chart: series %q has %d values; expected %d", c.Series[i].Name, len(c.Series[i].Values), len(c.Series[0].Values))
		}
	}
	if c.StackType != "" && c.StackType != Default && c.StackType != Stacked && c.StackType != Percent {
		return fmt.Errorf("termkit/chart: unsupported stack type %q", c.StackType)
	}
	return nil
}

func (c Chart) pointCount() int { return len(c.Series[0].Values) }
