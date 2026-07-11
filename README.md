# ditherkit-go

![Interactive dithered charts in a terminal](docs/demo.gif)

`ditherkit-go` is a small, renderer-independent chart engine for Go CLI and TUI applications. It brings the composable, dithered-chart feel of [Dither Kit](https://www.tripwire.sh/dither-kit) to terminals: animated or interactive apps can render area, line, bar, pie, and radar plots without requiring a browser, canvas, or a web dependency.

This is a clean-room Go port of the *product idea and chart coverage*, written from scratch for terminal output. The upstream project is currently published without a license, so this repository does not copy or translate its source code.

## Why a standalone module?

Observe needs better live plots, but the useful abstraction is broader than one monitoring dashboard. Keeping the engine separate means Boba, Observe, and any other Bubble Tea, Cobra, SSH, or plain CLI project can share one stable chart API. The core has no UI-framework dependency; Bubble Tea is used only by the interactive demo.

## Install

```bash
go get github.com/pol-cova/ditherkit-go
```

## Use

```go
chart := ditherkit.Chart{
    Kind:   ditherkit.Area,
    Title:  "Requests per minute",
    Labels: []string{"10:00", "10:01", "10:02", "10:03"},
    Series: []ditherkit.Series{
        {Name: "api", Values: []float64{22, 47, 31, 64}},
        {Name: "worker", Values: []float64{14, 19, 53, 39}},
    },
}

view, err := ditherkit.Render(chart, ditherkit.Options{
    Width:  56,
    Height: 12,
    Selected: 2,
    Color: true,
})
if err != nil { panic(err) }
fmt.Println(view)
```

`Chart.Kind` supports `Area`, `Line`, `Bar`, `Pie`, and `Radar`. `Selected` and `ActiveSeries` are deliberately passed as renderer state, so a host TUI owns keyboard and mouse interactions while this module remains portable.

## Demo

```bash
go run ./cmd/demo
```

Press `1`–`5` to switch chart types, `←`/`→` (or `h`/`l`) to scrub points, and `tab` to switch the active series. Re-record the README demo with `vhs demo.tape`.

## Development

```bash
go test ./...
go run ./cmd/demo --static
```

MIT licensed. See [LICENSE](LICENSE).
