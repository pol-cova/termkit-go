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

func TestSparklineAndBreadcrumb(t *testing.T) {
	if got := Sparkline([]float64{1, 4, 2, 8}, 4, Accent); len([]rune(got)) == 0 {
		t.Fatal("sparkline is empty")
	}
	if got := Breadcrumb([]string{"home", "reports"}, 1); !strings.Contains(got, "reports") {
		t.Fatal("breadcrumb omitted active item")
	}
}

func TestDitherAvatarIsSeeded(t *testing.T) {
	a := DitherAvatar("dan", 8, DitherPurple, BloomAura)
	if a == "" || a != DitherAvatar("dan", 8, DitherPurple, BloomAura) {
		t.Fatal("avatar should be deterministic")
	}
	if a == DitherAvatar("ada", 8, DitherPurple, BloomAura) {
		t.Fatal("different names should produce different avatars")
	}
}

func TestDitherStandaloneComponents(t *testing.T) {
	if DitherButton("save", DitherBlue, DitherGradientVariant, BloomLow, false, false, false) == "" {
		t.Fatal("button should render")
	}
	if got := DitherGradient(12, 3, DitherPurple, DitherBlue, "up", BloomLow); len(strings.Split(got, "\n")) != 3 {
		t.Fatalf("gradient rows = %d, want 3", len(strings.Split(got, "\n")))
	}
	if got := DitherSparkline([]float64{3, 7, 5, 9, 8, 12}, 16, 5, DitherGreen, DitherGradientVariant, BloomAura); len(strings.Split(got, "\n")) != 5 {
		t.Fatalf("sparkline rows = %d, want 5", len(strings.Split(got, "\n")))
	}
}
