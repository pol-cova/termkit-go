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

func TestPieIncludesSliceDetails(t *testing.T) {
	c := example(Pie)
	got, err := Render(c, Options{Width: 40, Height: 10, Selected: 1})
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"Mon 6%", "Tue 23%", "◆ Tue"} {
		if !strings.Contains(got, expected) {
			t.Fatalf("pie is missing %q:\n%s", expected, got)
		}
	}
}

func TestRadarUsesAvailableWidth(t *testing.T) {
	got, err := Render(example(Radar), Options{Width: 40, Height: 10})
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(got, "\n")
	maxWidth := 0
	for _, line := range lines[1:11] {
		if width := len([]rune(strings.TrimRight(line, " "))); width > maxWidth {
			maxWidth = width
		}
	}
	if maxWidth < 25 {
		t.Fatalf("radar did not use terminal width: %d\n%s", maxWidth, got)
	}
	if !strings.Contains(got, "1 Mon") {
		t.Fatalf("radar omitted axis labels:\n%s", got)
	}
}

func TestAreaKeepsEdgeAndSelectionVisible(t *testing.T) {
	got, err := Render(example(Area), Options{Width: 40, Height: 10, Selected: 2})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "▄") || !strings.Contains(got, "┊") {
		t.Fatalf("area lost its edge or selection marker:\n%s", got)
	}
}

func TestRenderRejectsUnevenSeries(t *testing.T) {
	c := Chart{Kind: Line, Series: []Series{{Name: "a", Values: []float64{1, 2}}, {Name: "b", Values: []float64{1}}}}
	if _, err := Render(c, Options{}); err == nil {
		t.Fatal("expected uneven series validation error")
	}
}

func TestInteractiveLegendDimsInactiveSeries(t *testing.T) {
	c := example(Area)
	c.Series[0].Color = "blue"
	c.Series[1].Color = "purple"
	got, err := Render(c, Options{Width: 48, Height: 10, Interactive: true, HoveredSeries: 0, Color: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "38;5;240m") || !strings.Contains(got, "38;5;33m") {
		t.Fatalf("interactive rendering did not dim/retain series colours:\n%s", got)
	}
}

func TestLineConnectsSamples(t *testing.T) {
	got, err := Render(example(Line), Options{Width: 40, Height: 10, Color: false})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "╌") && !strings.Contains(got, "•") {
		t.Fatalf("line renderer omitted connected stroke:\n%s", got)
	}
}
