# termkit-go

[![Go Reference](https://pkg.go.dev/badge/github.com/pol-cova/termkit-go.svg)](https://pkg.go.dev/github.com/pol-cova/termkit-go)
[![CI](https://github.com/pol-cova/termkit-go/actions/workflows/ci.yml/badge.svg)](https://github.com/pol-cova/termkit-go/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/pol-cova/termkit-go)](LICENSE)

![termkit-go demo](docs/demo.gif)

Terminal charts, motion, and UI components for Go CLIs and TUIs. Dither Kit is the visual inspiration for the chart package.

## Install

```bash
go get github.com/pol-cova/termkit-go
```

## Packages

- `chart` — area, line, bar, pie, and radar charts with dot, hatch, gradient, and solid fills.
- `animate` — easing, tween, spring, pulse, and repeat helpers.
- `component` — progress, gauges, spinners, badges, cards, dividers, and status bars.

## Quick start

```go
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

## Demo

```bash
go run ./cmd/demo
```

Use `1`–`5` to change chart type, `←`/`→` to scrub, and `tab` to switch series.

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
