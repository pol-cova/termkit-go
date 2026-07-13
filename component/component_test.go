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

func TestInputEditingAndSelection(t *testing.T) {
	input := Input{Placeholder: "type here"}
	input.HandleKey("go")
	input.HandleKey("left")
	input.HandleKey("backspace")
	if got := string(input.Value); got != "o" {
		t.Fatalf("input = %q", got)
	}
	selectable := Select{Options: []Option{{Label: "one"}, {Label: "two"}}, Height: 1}
	selectable.HandleKey("down")
	if !selectable.HandleKey("enter") || selectable.Selected != 1 {
		t.Fatal("selection was not confirmed")
	}
}

func TestTableAndPanelRespectWidth(t *testing.T) {
	if got := Table([]string{"Name", "State"}, [][]string{{"worker", "running"}}, 20); len(strings.Split(got, "\n")) != 5 {
		t.Fatal("table shape changed")
	}
	if got := Panel("Logs", "a very long line", 16, Accent); !strings.Contains(got, "Logs") {
		t.Fatal("panel omitted title")
	}
}
