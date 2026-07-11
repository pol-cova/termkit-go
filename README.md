# termkit-go

![Interactive terminal charts and components](docs/demo.gif)

`termkit-go` is a composable Go toolkit for polished CLI and TUI interfaces. It provides dithered charts, motion primitives, and terminal components that work in Bubble Tea, Cobra commands, SSH tools, or plain ANSI output—without forcing an application framework on the host.

The project began as a clean-room Go implementation of the *idea and chart coverage* behind [Dither Kit](https://www.tripwire.sh/dither-kit): rich, interactive-looking plots outside the browser. It now deliberately has a broader job: make Go terminal applications feel alive, legible, and intentional.

## Packages

| Package | What it provides |
| --- | --- |
| `chart` | Area, line, bar, pie, and radar charts with selection state and dithered fills. |
| `animate` | Deterministic tween, easing, springs, pulses, and repeat helpers for application-owned ticks. |
| `component` | Progress indicators, gauges, spinner frames, badges, cards, dividers, and segmented status bars. |

The core packages do not require Bubble Tea. The interactive demo uses it only to show how a host TUI can own keyboard input and animation ticks.

## Install

```bash
go get github.com/pol-cova/termkit-go
```

## Use a chart

```go
import "github.com/pol-cova/termkit-go/chart"

plot := chart.Chart{
    Kind:   chart.Area,
    Title:  "Requests per minute",
    Labels: []string{"10:00", "10:01", "10:02", "10:03"},
    Series: []chart.Series{
        {Name: "api", Values: []float64{22, 47, 31, 64}},
        {Name: "worker", Values: []float64{14, 19, 53, 39}},
    },
}

view, err := chart.Render(plot, chart.Options{
    Width: 56, Height: 12, Selected: 2, Color: true,
})
if err != nil { panic(err) }
fmt.Println(view)
```

## Add motion and components

```go
import (
    "github.com/pol-cova/termkit-go/animate"
    "github.com/pol-cova/termkit-go/component"
)

// Progress is typically your Bubble Tea frame number divided by its duration.
progress := animate.Tween(0.65, animate.EaseInOut)
fmt.Println(component.Progress("Deploy", progress, 20, component.Accent))
fmt.Println(component.SpinnerFrame(frame, "Connecting", component.Warning))
```

## Demo

```bash
go run ./cmd/demo
```

Press `1`–`5` to switch chart types, `←`/`→` (or `h`/`l`) to scrub points, and `tab` to switch the active series. The demo also animates progress, a gauge, and a spinner to demonstrate host-controlled motion. Re-record the README demo with `vhs demo.tape`.

## Design principles

- **Portable:** output is strings and explicit state; you choose the TUI framework.
- **Predictable:** animations are pure functions of a progress value or frame, not hidden timers.
- **Terminal-native:** Unicode density, dither patterns, ANSI colour, and low visual noise.
- **Clean-room:** no upstream Dither Kit source was copied or translated; its public repository was unlicensed when this project began.

## Development

```bash
go test ./...
go vet ./...
go run ./cmd/demo --static
```

MIT licensed. See [LICENSE](LICENSE).
