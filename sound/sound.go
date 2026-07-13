// Package sound provides optional, terminal-native interaction cues.
package sound

import "io"

// Cue names mirror the small interaction vocabulary popularized by Cuelume.
type Cue string

const (
	Idle    Cue = "idle"
	Chime   Cue = "chime"
	Tick    Cue = "tick"
	Press   Cue = "press"
	Release Cue = "release"
	Toggle  Cue = "toggle"
	Success Cue = "success"
	Error   Cue = "error"
)

// Player is intentionally tiny so applications can connect a real synth,
// desktop notification, or test recorder without a platform dependency.
type Player interface{ Play(Cue) }

// Bell writes the terminal alert character for every cue. Terminals decide
// whether that becomes an audible bell, visual flash, or nothing.
type Bell struct{ Out io.Writer }

func (b Bell) Play(Cue) {
	if b.Out != nil {
		_, _ = io.WriteString(b.Out, "\a")
	}
}

// Silent is useful for accessibility settings and tests.
type Silent struct{}

func (Silent) Play(Cue) {}
