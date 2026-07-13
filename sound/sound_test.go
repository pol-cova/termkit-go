package sound

import (
	"bytes"
	"testing"
)

func TestBellWritesTerminalCue(t *testing.T) {
	var out bytes.Buffer
	Bell{Out: &out}.Play(Success)
	if out.String() != "\a" {
		t.Fatalf("got %q", out.String())
	}
}
