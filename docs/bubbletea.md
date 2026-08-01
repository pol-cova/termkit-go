# Bubble Tea cookbook

termkit-go stays event-loop neutral, but most interactive dashboards use [Bubble Tea](https://github.com/charmbracelet/bubbletea). The optional `bubbletea` package and this guide show the smallest wiring patterns.

## Chart dashboard

Keep chart state in your model and render through `bubbletea.ChartView`:

```go
type model struct {
    plot     chart.Chart
    selected int
    width    int
    height   int
}

func (m model) View() string {
    return bubbletea.ChartView{
        Data: m.plot,
        Options: chart.Options{
            Width: m.width, Height: m.height,
            Selected: m.selected, Color: true,
            Formatter: chart.SIFormatter,
        },
    }.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if key, ok := msg.(tea.KeyMsg); ok {
        bubbletea.ScrubSelection(key.String(), &m.selected, len(m.plot.Labels))
    }
    return m, nil
}
```

Use `Selected()` when you need structured values instead of parsing the readout string:

```go
point := bubbletea.ChartView{Data: m.plot, Options: opts}.Selected()
for _, reading := range point.Series {
    if reading.Missing {
        continue
    }
    // reading.Name, reading.Value
}
```

## Pixel charts with auto protocol

Let the terminal pick Kitty, iTerm2, or ANSI half-blocks:

```go
canvas, err := chart.RenderPixel(plot, chart.PixelOptions{Width: 192, Height: 48})
if err != nil { return err.Error() }
out, err := canvas.Auto()
if err != nil { return err.Error() }
return out
```

Run `go run ./cmd/demo --pixel` to preview all five geometries.

## Inputs and lists

Use the framework-neutral editors directly:

```go
var field component.Input
changed, submitted := field.HandleKey(msg.String())

var menu component.Select
menu.Options = []component.Option{{Label: "Dashboard"}, {Label: "Settings"}}
if menu.HandleKey(msg.String()) {
    // enter confirmed current item
}
```

## Motion and status rows

Animate with `animate` and compose status widgets in `View`:

```go
motion := animate.Tween(float64(frame%40)/39, animate.EaseInOut)
widgets := component.Badge("LIVE", component.Success) + "  " +
    component.Progress("CPU", motion, 14, component.Accent)
footer := component.StatusBar(component.StatusBarOptions{
    Left:  []component.Segment{{Text: "←/→ scrub", Tone: component.Accent}},
    Right: "q quit",
    Width: width,
})
```

See [`cmd/demo/main.go`](../cmd/demo/main.go) for a complete Bubble Tea program.
