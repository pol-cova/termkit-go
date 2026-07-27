package chart

import (
	"math"
	"strings"
	"testing"
	"time"
)

func TestHistoryPushAndSeries(t *testing.T) {
	h, err := NewHistory(4, []string{"a", "b"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	h.Push(1, 10)
	h.Push(2, 20)
	h.Push(3, 30)
	h.Push(4, 40)
	h.Push(5, 50)
	if h.Len() != 4 {
		t.Fatalf("len = %d", h.Len())
	}
	series := h.Series()
	if len(series) != 2 || len(series[0].Values) != 4 {
		t.Fatalf("unexpected series: %+v", series)
	}
	if series[0].Values[0] != 2 || series[0].Values[3] != 5 {
		t.Fatalf("expected rotated window, got %v", series[0].Values)
	}
	if series[1].Values[3] != 50 {
		t.Fatal(series[1].Values)
	}
}

func TestHistoryRequiresNames(t *testing.T) {
	if _, err := NewHistory(4, nil, nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestGoldenBrailleMemory(t *testing.T) {
	ts := TimeSeries{
		Title:  "Memory",
		Window: 60 * time.Second,
		Series: []Series{
			{Name: "RAM", Values: []float64{55, 62, 58, 69, 71, 68, 72, 69}, Color: "monitor-primary"},
			{Name: "SWP", Values: []float64{10, 12, 15, 18, 20, 22, 25, 25}, Color: "monitor-secondary"},
		},
		Y: AxisSpec{Min: 0, Max: 100, Formatter: PercentFormatter},
		Legend: LegendSpec{
			Entries: []LegendEntry{
				{Name: "RAM", Current: "69%", Extra: "25.0GiB/36.0GiB"},
				{Name: "SWP", Current: "25%", Extra: "0.5GiB/2.0GiB"},
			},
		},
	}
	got, err := RenderBraille(ts, BrailleOptions{
		Width: 40, Height: 10, Color: false, ShowAxes: true, ShowLegend: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "braille-memory", got)
}

func TestGoldenBrailleNetwork(t *testing.T) {
	ts := TimeSeries{
		Title:  "Network",
		Window: 60 * time.Second,
		Series: []Series{
			{Name: "RX", Values: sineWave(24, 0, 500), Color: "monitor-primary"},
			{Name: "TX", Values: sineWave(24, 3, 800), Color: "monitor-secondary"},
		},
		Y: AxisSpec{Formatter: SIFormatter},
	}
	got, err := RenderBraille(ts, BrailleOptions{
		Width: 44, Height: 8, Color: false, ShowAxes: true, ShowLegend: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "braille-network", got)
}

func TestBrailleCompactFallback(t *testing.T) {
	ts := TimeSeries{
		Series: []Series{{Name: "CPU", Values: []float64{10, 20, 30, 40}}},
		Y:      AxisSpec{Formatter: PercentFormatter},
	}
	got, err := RenderBraille(ts, BrailleOptions{Width: 20, Height: 4, Color: false})
	if err != nil {
		t.Fatal(err)
	}
	if got == "" {
		t.Fatal("expected compact output")
	}
}

func TestBrailleNoColorEnv(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	ts := TimeSeries{
		Series: []Series{{Name: "CPU", Values: []float64{10, 50, 30}}},
		Y:      AxisSpec{Formatter: PercentFormatter},
	}
	got, err := RenderBraille(ts, BrailleOptions{Width: 40, Height: 10, Color: true, ShowAxes: true})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "\x1b[") {
		t.Fatal("expected no ANSI when NO_COLOR is set")
	}
}

func TestBrailleScrubReadout(t *testing.T) {
	ts := TimeSeries{
		Window: 60 * time.Second,
		Series: []Series{{Name: "CPU", Values: []float64{10, 20, 30, 40, 50}}},
		Y:      AxisSpec{Min: 0, Max: 100, Formatter: PercentFormatter},
	}
	got, err := RenderBraille(ts, BrailleOptions{
		Width: 40, Height: 10, Color: false, ShowAxes: true, Selected: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !contains(got, "CPU: 30%") {
		t.Fatalf("missing scrub readout: %s", got)
	}
}

func sineWave(n int, phase float64, scale float64) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = (math.Sin(float64(i)/4+phase) + 1) * scale / 2
	}
	return out
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
