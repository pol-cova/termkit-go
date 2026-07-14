package pixel

import (
	"strings"
	"testing"
)

func TestANSIUsesTwoVerticalPixelsPerCell(t *testing.T) {
	c := New(4, 4, Transparent)
	c.Set(0, 0, RGBA{R: 255, A: 255})
	c.Set(0, 1, RGBA{B: 255, A: 255})
	got := c.ANSI()
	if !strings.Contains(got, "▀") || !strings.Contains(got, "48;2;0;0;255") {
		t.Fatalf("unexpected half-block output: %q", got)
	}
	if strings.Count(got, "\n") != 1 {
		t.Fatalf("rows = %d, want 2", strings.Count(got, "\n")+1)
	}
}

func TestPNGAndProtocols(t *testing.T) {
	c := New(2, 2, RGBA{R: 20, G: 30, B: 40, A: 255})
	pngBytes, err := c.PNG()
	if err != nil || len(pngBytes) == 0 {
		t.Fatalf("PNG: %v", err)
	}
	kitty, err := c.Kitty()
	if err != nil || !strings.HasPrefix(kitty, "\x1b_Ga=T") {
		t.Fatalf("Kitty: %v", err)
	}
	iterm, err := c.ITerm2()
	if err != nil || !strings.Contains(iterm, "1337;File=inline=1") {
		t.Fatalf("iTerm2: %v", err)
	}
}
