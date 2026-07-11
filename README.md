# termkit-go

[![Go Reference](https://pkg.go.dev/badge/github.com/pol-cova/termkit-go.svg)](https://pkg.go.dev/github.com/pol-cova/termkit-go)
[![CI](https://github.com/pol-cova/termkit-go/actions/workflows/ci.yml/badge.svg)](https://github.com/pol-cova/termkit-go/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/pol-cova/termkit-go)](LICENSE)

![termkit-go demo](docs/demo.gif)

Terminal charts, motion, and UI components for Go CLIs and TUIs. The chart palette and texture are inspired by [Dither Kit](https://www.tripwire.sh/dither-kit).

## Install

```bash
go get github.com/pol-cova/termkit-go
```

## Packages

| Package | Includes | API |
| --- | --- | --- |
| `chart` | Area, line, bar, pie, radar; stacked and percent modes; interactive selection. | [docs](docs/charts.md) · [Go reference](https://pkg.go.dev/github.com/pol-cova/termkit-go/chart) |
| `animate` | Tween, spring, pulse, repeat, and easing helpers. | [docs](docs/components.md) · [Go reference](https://pkg.go.dev/github.com/pol-cova/termkit-go/animate) |
| `component` | Progress, gauges, spinners, badges, cards, dividers, and status bars. | [docs](docs/components.md) · [Go reference](https://pkg.go.dev/github.com/pol-cova/termkit-go/component) |

Each package returns values or strings. It does not own your event loop, terminal, or rendering framework.

## Quick start

```go
import (
    "fmt"

    "github.com/pol-cova/termkit-go/chart"
)

plot := chart.Chart{
    Kind: chart.Area,
    Labels: []string{"Jan", "Feb", "Mar", "Apr"},
    StackType: chart.Stacked,
    Series: []chart.Series{
        {Name: "desktop", Values: []float64{120, 190, 230, 210}, Variant: chart.Dotted},
        {Name: "mobile", Values: []float64{80, 120, 140, 110}, Variant: chart.Hatched},
    },
}

view, err := chart.Render(plot, chart.Options{Width: 56, Height: 12, Selected: 2, Color: true})
if err != nil { panic(err) }
fmt.Println(view)
```

### Bubble Tea

Keep selection in your model and render the view from `View`:

```go
view, err := chart.Render(plot, chart.Options{
    Width:        m.width,
    Height:       10,
    Selected:     m.selected,
    ActiveSeries: m.activeSeries,
    Color:        true,
})
if err != nil { return err.Error() }
return view
```

Update `m.selected` from arrow keys, `j`/`k`, or mouse input. See the [interactive demo](cmd/demo/main.go) for the complete model.

### Motion and status components

```go
import (
    "fmt"

    "github.com/pol-cova/termkit-go/animate"
    "github.com/pol-cova/termkit-go/component"
)

progress := animate.Tween(0.65, animate.EaseInOut)
fmt.Println(component.Progress("Deploy", progress, 20, component.Accent))
fmt.Println(component.SpinnerFrame(frame, "Connecting", component.Warning))
```

## Demo

```bash
go run ./cmd/demo
```

Use `1`–`5` to change chart type, `←`/`→` to scrub, and `tab` to switch series.

## Requirements

Go 1.24 or newer. The library itself has no required UI framework; the demo uses Bubble Tea.

## Docs

- [Charts and visual variants](docs/charts.md)
- [Motion and components](docs/components.md)
- [Contributing](CONTRIBUTING.md)

## Development

```bash
go test ./...
go vet ./...
vhs demo.tape
```

MIT licensed. See [LICENSE](LICENSE).
