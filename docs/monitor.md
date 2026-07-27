# Monitor graphs cookbook

Add a btm-style braille pane in ~15 lines:

```go
hist, _ := chart.NewHistory(120, []string{"RX", "TX"}, []string{"monitor-primary", "monitor-secondary"})

// each tick:
hist.Push(rxRate, txRate)

widget := component.FuncWidget{
    Render: func(w, h int) string {
        body, _ := chart.RenderBraille(chart.TimeSeries{
            Window: 60 * time.Second,
            Series: hist.Series(),
            Y:      chart.AxisSpec{Formatter: chart.SIFormatter},
        }, chart.BrailleOptions{
            Width: w - 2, Height: h - 2,
            Color: true, ShowAxes: true, ShowLegend: true,
        })
        return component.BorderTitle("Network", body, w, h, component.Accent)
    },
}

view := component.Grid(termW, termH, 2, 2,
    component.Slot{Widget: widget, Row: 1, Col: 0, WeightW: 2},
)
```

## Layout

Use `component.Grid(width, height, rows, cols, slots...)` with weighted slots:

```go
component.Grid(w, h, 3, 2,
    component.Slot{Widget: cpuPane, Row: 0, Col: 0, RowSpan: 2, WeightW: 3, WeightH: 2},
    component.Slot{Widget: sideTable, Row: 0, Col: 1, WeightW: 1},
)
```

Implement `component.Widget` or use `component.FuncWidget` for one-off panes.

## Demo

```bash
go run ./cmd/demo --monitor
```

Synthetic CPU, memory, network, disk, and process tables update every 500ms. Use `←`/`→` to scrub, `q` to quit.

## Performance

- Push samples at **1–2 Hz** for system metrics; render on every host frame or at the same rate.
- Braille charts are deterministic: identical `(data, size, options)` always produce the same string.
- Set `NO_COLOR=1` or `BrailleOptions{Color: false}` for monochrome terminals.

## Compact fallback

When width &lt; 24 or height &lt; 6, `RenderBraille` switches to block sparklines automatically.

## Adding a new pane

1. Create a `chart.History` (or static `[]chart.Series`).
2. Wrap `chart.RenderBraille` in `component.BorderTitle`.
3. Register a `component.Slot` in `Grid`.

No changes to `chart.RenderBraille` are required.
