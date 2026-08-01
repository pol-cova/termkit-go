package pixel

import "testing"

func TestDetectProtocol(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "Apple_Terminal")
	if DetectProtocol() != ProtocolANSI {
		t.Fatal("expected ANSI fallback")
	}
	t.Setenv("TERM_PROGRAM", "kitty")
	if DetectProtocol() != ProtocolKitty {
		t.Fatal("expected Kitty")
	}
	t.Setenv("TERM_PROGRAM", "iTerm.app")
	if DetectProtocol() != ProtocolITerm2 {
		t.Fatal("expected iTerm2")
	}
}

func TestAutoFallsBackToANSI(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "")
	c := New(4, 4, RGBA{R: 255, G: 0, B: 0, A: 255})
	out, err := c.Auto()
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Fatal("expected ANSI output")
	}
}
