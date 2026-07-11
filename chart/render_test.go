package chart

import (
	"strings"
	"testing"
)

func example(kind Kind) Chart {
	return Chart{Kind: kind, Title: "Traffic", Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri"}, Series: []Series{{Name: "desktop", Values: []float64{12, 45, 28, 62, 51}}, {Name: "mobile", Values: []float64{18, 25, 55, 39, 70}}}}
}
func TestRenderEveryChartKind(t *testing.T) {
	for _, kind := range []Kind{Area, Line, Bar, Pie, Radar} {
		t.Run(string(kind), func(t *testing.T) {
			got, err := Render(example(kind), Options{Width: 40, Height: 10, Selected: 2})
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(got, "Traffic") || !strings.Contains(got, "Wed") {
				t.Fatalf("incomplete chart: %q", got)
			}
		})
	}
}
func TestRenderRejectsInvalidInput(t *testing.T) {
	_, err := Render(Chart{Kind: Line}, Options{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
func TestRenderClampsSelection(t *testing.T) {
	got, err := Render(example(Line), Options{Selected: 99})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "Fri") {
		t.Fatalf("selection was not clamped: %s", got)
	}
}

func TestRenderSupportsDitherKitStyleStacks(t *testing.T) {
	c := example(Area)
	c.StackType = Stacked
	c.Series[0].Variant = Dotted
	c.Series[1].Variant = Hatched
	got, err := Render(c, Options{Width: 48, Height: 10, Color: false})
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"Mon", "desktop", "mobile", "▒"} {
		if !strings.Contains(got, expected) {
			t.Fatalf("stacked chart is missing %q:\n%s", expected, got)
		}
	}
}
