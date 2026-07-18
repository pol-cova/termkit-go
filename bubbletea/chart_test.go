package bubbletea

import "testing"

func TestScrubSelection(t *testing.T) {
	selected := 2
	if !ScrubSelection("left", &selected, 6) || selected != 1 {
		t.Fatalf("left = %d", selected)
	}
	if !ScrubSelection("right", &selected, 6) || selected != 2 {
		t.Fatalf("right = %d", selected)
	}
	if !ScrubSelection("end", &selected, 6) || selected != 5 {
		t.Fatalf("end = %d", selected)
	}
}

func TestCycleSeries(t *testing.T) {
	active := 0
	CycleSeries(&active, 2)
	if active != 1 {
		t.Fatal(active)
	}
	CycleSeries(&active, 2)
	if active != 0 {
		t.Fatal(active)
	}
}
