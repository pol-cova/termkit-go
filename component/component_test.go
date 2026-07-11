package component

import (
	"strings"
	"testing"
)

func TestProgressClampsValues(t *testing.T) {
	for _, value := range []float64{-1, 2} {
		if got := Progress("CPU", value, 8, Accent); got == "" {
			t.Fatal("empty progress")
		}
	}
}
func TestSpinnerCycles(t *testing.T) {
	if SpinnerFrame(0, "loading", Accent) == SpinnerFrame(1, "loading", Accent) {
		t.Fatal("spinner did not advance")
	}
}

func TestCardAndDivider(t *testing.T) {
	if got := Card("status", "ready", 20); !strings.Contains(got, "ready") {
		t.Fatal("card omitted content")
	}
	if got := Divider(4); got == "" {
		t.Fatal("divider is empty")
	}
}
