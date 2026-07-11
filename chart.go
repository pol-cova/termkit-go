// Package ditherkit renders compact, interactive-friendly charts for terminals.
//
// It is a clean-room Go implementation inspired by the category of dithered
// chart interfaces; it does not include source code from Dither Kit.
package ditherkit

import "fmt"

// Kind describes a chart geometry.
type Kind string

const (
	Area  Kind = "area"
	Line  Kind = "line"
	Bar   Kind = "bar"
	Pie   Kind = "pie"
	Radar Kind = "radar"
)

// Series supplies one named measure. Values are aligned with Chart.Labels.
type Series struct {
	Name   string
	Values []float64
}

// Chart is renderer-independent chart input.
type Chart struct {
	Kind   Kind
	Title  string
	Labels []string
	Series []Series
}

func (c Chart) validate() error {
	if len(c.Series) == 0 {
		return fmt.Errorf("ditherkit: chart requires at least one series")
	}
	if c.Kind != Area && c.Kind != Line && c.Kind != Bar && c.Kind != Pie && c.Kind != Radar {
		return fmt.Errorf("ditherkit: unsupported chart kind %q", c.Kind)
	}
	for _, s := range c.Series {
		if len(s.Values) == 0 {
			return fmt.Errorf("ditherkit: series %q has no values", s.Name)
		}
		if len(c.Labels) > 0 && len(s.Values) != len(c.Labels) {
			return fmt.Errorf("ditherkit: series %q has %d values for %d labels", s.Name, len(s.Values), len(c.Labels))
		}
	}
	return nil
}

func (c Chart) pointCount() int { return len(c.Series[0].Values) }
