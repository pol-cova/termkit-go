package chart

import (
	"fmt"
	"math"
)

// History stores rolling samples for one or more series without host-side slicing.
type History struct {
	capacity int
	names    []string
	colors   []string
	buffers  [][]float64
	start    int
	count    int
}

// NewHistory creates a ring buffer for named series. colors may be nil or shorter than names.
func NewHistory(capacity int, names []string, colors []string) (*History, error) {
	if capacity < 2 {
		return nil, fmt.Errorf("termkit/chart: history capacity must be >= 2")
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("termkit/chart: history requires at least one series name")
	}
	buffers := make([][]float64, len(names))
	for i := range buffers {
		buffers[i] = make([]float64, capacity)
	}
	return &History{
		capacity: capacity,
		names:    append([]string(nil), names...),
		colors:   append([]string(nil), colors...),
		buffers:  buffers,
	}, nil
}

// Push appends one sample per series. Extra values are ignored; missing values become NaN.
func (h *History) Push(values ...float64) {
	if h == nil {
		return
	}
	idx := (h.start + h.count) % h.capacity
	for i := range h.buffers {
		var v float64
		if i < len(values) {
			v = values[i]
		} else {
			v = math.NaN()
		}
		h.buffers[i][idx] = v
	}
	if h.count < h.capacity {
		h.count++
	} else {
		h.start = (h.start + 1) % h.capacity
	}
}

// Len returns the number of samples currently stored.
func (h *History) Len() int {
	if h == nil {
		return 0
	}
	return h.count
}

// Series exports the current window as chart.Series slices (oldest to newest).
func (h *History) Series() []Series {
	if h == nil {
		return nil
	}
	out := make([]Series, len(h.buffers))
	for i := range h.buffers {
		values := make([]float64, h.count)
		for j := 0; j < h.count; j++ {
			values[j] = h.buffers[i][(h.start+j)%h.capacity]
		}
		color := ""
		if i < len(h.colors) {
			color = h.colors[i]
		}
		out[i] = Series{Name: h.names[i], Values: values, Color: color}
	}
	return out
}
