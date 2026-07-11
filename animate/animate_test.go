package animate

import "testing"

func TestEasingIsBounded(t *testing.T) {
	for _, fn := range []Easing{Linear, EaseInOut, EaseOut} {
		for _, value := range []float64{-1, 0, .5, 1, 2} {
			got := fn(value)
			if got < 0 || got > 1 {
				t.Fatalf("%f escaped bounds", got)
			}
		}
	}
}
func TestSpringSettles(t *testing.T) {
	if got := Spring(1, 12, 0.85); got < 0.9 {
		t.Fatalf("spring did not settle: %f", got)
	}
}

func TestRepeatWrapsBothDirections(t *testing.T) {
	if got := Repeat(12, 10); got != .2 {
		t.Fatalf("Repeat(12) = %f", got)
	}
	if got := Repeat(-1, 10); got != .9 {
		t.Fatalf("Repeat(-1) = %f", got)
	}
}
