package chart

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoldenAreaChart(t *testing.T) {
	got, err := Render(example(Area), Options{Width: 40, Height: 10, Selected: 2, Color: false})
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "area", got)
}

func TestGoldenLineChart(t *testing.T) {
	got, err := Render(example(Line), Options{Width: 40, Height: 10, Color: false})
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "line", got)
}

func TestSelectedPointStructured(t *testing.T) {
	c := example(Line)
	sp := c.SelectedPoint(Options{Selected: 2})
	if sp.Label != "Wed" {
		t.Fatalf("label = %q", sp.Label)
	}
	if len(sp.Series) != 2 || sp.Series[0].Value != 28 {
		t.Fatalf("unexpected readings: %+v", sp.Series)
	}
}

func TestMissingValuesRender(t *testing.T) {
	c := Chart{
		Kind: Line, Labels: []string{"a", "b", "c"},
		Series: []Series{{Name: "x", Values: []float64{1, math.NaN(), 3}}},
	}
	got, err := Render(c, Options{Width: 20, Height: 6, Color: false})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "x:") {
		t.Fatalf("expected series in readout: %s", got)
	}
	sp := c.SelectedPoint(Options{Selected: 1})
	if !sp.Series[0].Missing {
		t.Fatal("expected missing value at index 1")
	}
	if DefaultFormatter(sp.Series[0].Value) != "—" {
		t.Fatal("expected em dash for missing value")
	}
}

func TestSIFormatter(t *testing.T) {
	if SIFormatter(1500) != "1.5K" {
		t.Fatal(SIFormatter(1500))
	}
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", name+".golden")
	if os.Getenv("UPDATE_GOLDEN") != "" {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("missing golden file %s (run with UPDATE_GOLDEN=1): %v", path, err)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch for %s\n--- got ---\n%s\n--- want ---\n%s", name, got, want)
	}
}
