// Package animate provides deterministic easing and spring helpers for TUI state.
package animate

import "math"

// Easing maps a normalised progress value (0 through 1) to another normalised value.
type Easing func(float64) float64

// Linear leaves progress unchanged.
func Linear(t float64) float64 { return clamp(t) }

// EaseInOut accelerates gently, then decelerates into its destination.
func EaseInOut(t float64) float64 {
	t = clamp(t)
	return t * t * (3 - 2*t)
}

// EaseOut is a quick-starting, gentle-arriving curve.
func EaseOut(t float64) float64 { t = clamp(t); return 1 - (1-t)*(1-t) }

// Tween applies an easing curve to normalised progress.
func Tween(progress float64, easing Easing) float64 {
	if easing == nil {
		easing = Linear
	}
	return easing(progress)
}

// Pulse returns a smooth 0-to-1-to-0 waveform for repeating emphasis.
func Pulse(progress float64) float64 {
	progress = clamp(progress)
	return .5 - .5*math.Cos(2*math.Pi*progress)
}

// Repeat converts an ever-increasing frame counter into normalised loop progress.
// A non-positive period falls back to one frame.
func Repeat(frame, period int) float64 {
	if period <= 0 {
		period = 1
	}
	frame %= period
	if frame < 0 {
		frame += period
	}
	return float64(frame) / float64(period)
}

// Spring returns a stable, normalised spring interpolation. response controls
// how quickly it settles; damping values near 1 avoid overshoot.
func Spring(progress, response, damping float64) float64 {
	progress = clamp(progress)
	if response <= 0 {
		response = 12
	}
	damping = clamp(damping)
	if damping == 0 {
		damping = 0.8
	}
	omega := response
	if damping < 1 {
		beta := omega * math.Sqrt(1-damping*damping)
		return clamp(1 - math.Exp(-damping*omega*progress)*(math.Cos(beta*progress)+damping*omega*math.Sin(beta*progress)/beta))
	}
	return clamp(1 - math.Exp(-omega*progress)*(1+omega*progress))
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
