package component

import (
	"strings"
	"testing"
)

func TestBorderTitle(t *testing.T) {
	got := BorderTitle("Memory", "plot line", 20, 5, Accent)
	if !strings.Contains(got, "Memory") {
		t.Fatalf("missing title: %q", got)
	}
	if !strings.Contains(got, "plot line") {
		t.Fatalf("missing body: %q", got)
	}
	lines := strings.Split(got, "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(lines))
	}
}

func TestGridPlacesWidgets(t *testing.T) {
	left := FuncWidget{
		Render: func(w, h int) string {
			return BorderTitle("Left", "content", w, h, Muted)
		},
	}
	right := StaticWidget{Content: "right pane"}
	got := Grid(30, 8, 1, 2,
		Slot{Widget: left, Row: 0, Col: 0, WeightW: 2},
		Slot{Widget: right, Row: 0, Col: 1, WeightW: 1},
	)
	if !strings.Contains(got, "Left") {
		t.Fatalf("missing left widget: %q", got)
	}
	if !strings.Contains(got, "right pane") {
		t.Fatalf("missing right widget: %q", got)
	}
}

func TestGridStableWidth(t *testing.T) {
	widget := StaticWidget{Content: "hello"}
	got := Grid(24, 3, 1, 1, Slot{Widget: widget, Row: 0, Col: 0})
	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines")
	}
	for _, line := range lines {
		if len([]rune(line)) != 24 {
			t.Fatalf("line width = %d", len([]rune(line)))
		}
	}
}
