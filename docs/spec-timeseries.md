# SPEC: Time-Series Graphs & Monitor Composition

## Architecture

```
Host app / cmd/demo
    → component.Grid (layout)
        → component.Widget (interface)
            → chart.RenderBraille (TimeSeries)
            → component.BorderTitle
            → component.Table
        → chart.History (ring buffer)
        → style monitor tokens
```

## Package responsibilities

| Package | Change |
| --- | --- |
| `chart/` | `TimeSeries`, `History`, `RenderBraille`, axis/legend options |
| `style/` | Monitor semantic tokens mapped to ANSI-256 |
| `component/` | `Widget`, `Grid`, `BorderTitle` |
| `cmd/demo/` | `--monitor` scene with synthetic series |
| `docs/` | Cookbooks and fidelity checklist |

## Public API

### chart

```go
type History struct { /* ring buffer per series */ }
func NewHistory(capacity int, names []string, colors []string) *History
func (h *History) Push(values ...float64)
func (h *History) Series() []Series

type AxisSpec struct {
    Min, Max float64 // 0,0 = auto from data
    Formatter ValueFormatter
}

type LegendEntry struct {
    Name    string
    Current string
    Extra   string // e.g. "25.0GiB/36.0GiB" or "All: 71.2GB"
}

type LegendSpec struct {
    Entries []LegendEntry // empty = auto from series + formatters
}

type TimeSeries struct {
    Title  string
    Window time.Duration
    Series []Series
    Y      AxisSpec
    Legend LegendSpec
}

type BrailleOptions struct {
    Width, Height int
    Color         bool
    ShowAxes      bool
    ShowLegend    bool
    Selected      int // scrub index; -1 = none
    Animate       bool
}

func RenderBraille(ts TimeSeries, opt BrailleOptions) (string, error)
```

### component

```go
type Widget interface {
    SizeHint() (w, h int)
    View(width, height int) string
}

type Slot struct {
    Widget Widget
    Row, Col, RowSpan, ColSpan int
    MinW, MinH, WeightW, WeightH int
}

func Grid(width, height, rows, cols int, slots ...Slot) string
func BorderTitle(title, body string, width, height int, tone Tone) string
```

### Extension recipe

Adding a "Disk I/O" graph:

1. Host maintains `chart.NewHistory(120, []string{"read", "write"}, nil)`.
2. Each tick: `history.Push(readRate, writeRate)`.
3. Implement `Widget` that calls `RenderBraille` with `history.Series()`.
4. Register a `Slot` in `Grid`—no changes to `chart` internals.

## Visual spec

| Element | Rule |
| --- | --- |
| Braille cell | Unicode `U+2800`–`U+28FF`; 2×4 logical dots per character |
| Multi-series | Independent polylines; z-order = series order |
| Axes | Left Y labels ≥4 cols; bottom X ≥3 tick labels; grid dimmer than series |
| Legend | Top-left inside plot; padded; reserved inset |
| Border | `┌─┐│└─┘` with title on top edge; ellipsis truncation |
| Min size | Below 24×6: compact sparkline fallback |
| Determinism | Same inputs ⇒ identical output (golden tests) |

## Testing

- Unit: braille packing, axis ticks, legend truncation, history wrap
- Golden: `chart/testdata/braille*.golden` (color off)
- Composition: `Grid` + `BorderTitle` width stability in `component`
- Manual: `go run ./cmd/demo --monitor`
- CI: `go test ./...`, `go vet ./...`

## Boundaries

- **Always:** pure functions; tests + docs with API
- **Ask first:** new dependencies (e.g. gopsutil), breaking `chart.Options` changes
- **Never:** goroutines in `chart`/`component`; Bubble Tea outside `bubbletea`/`cmd`

## Milestones & acceptance criteria

### M0 — Spec freeze

- [x] PRD and SPEC committed
- [x] Extension recipe documented

### M1 — Braille core

- [x] Multi-series braille with Y% and relative time X labels
- [x] In-plot legend with series colors
- [x] `History.Push` fixed window without host slicing
- [x] Golden tests; no regressions on existing charts

### M2 — Composition + demo

- [x] `Widget`, `Grid`, `BorderTitle`
- [x] Demo `--monitor` with ≥2 braille panes + ≥1 table
- [x] New pane addable with demo-local code only
- [x] Cookbook in `docs/monitor.md`

### M3 — UX polish

- [x] Compact fallback below size threshold
- [x] NO_COLOR monochrome path
- [x] Scrub readout without layout jump
- [x] README updated

## Delivery order

1. Docs (M0)
2. Braille renderer + goldens (M1)
3. History + legend/axes (M1)
4. Widget/Grid (M2)
5. Monitor demo (M2)
6. UX polish (M3)
