# Motion and components

The animation helpers are pure functions. Your application owns the timer and passes a normalized value or frame count.

```go
progress := animate.Tween(0.6, animate.EaseInOut)
fmt.Println(component.Progress("Deploy", progress, 20, component.Accent))
fmt.Println(component.SpinnerFrame(frame, "Connecting", component.Warning))
```

Use `animate.Spring`, `animate.Pulse`, and `animate.Repeat` for a natural-looking response without background goroutines. Components return ANSI strings and work with Bubble Tea, Cobra, SSH, or `fmt.Println`.

## Interactive primitives

`component.Input` and `component.Select` keep state but do not read stdin or
start a render loop. Pass the key name from your TUI's event loop to
`HandleKey`, then render with `View`:

```go
input := component.Input{Placeholder: "Search…", MaxLength: 80}
changed, submitted := input.HandleKey("hello")
_ = changed
_ = submitted

menu := component.Select{
    Options: []component.Option{
        {Label: "Dashboard", Description: "overview"},
        {Label: "Settings", Description: "preferences"},
    },
    Height: 4,
}
menu.HandleKey("j") // up/down and j/k both work
fmt.Println(component.Panel("Command palette", input.View(true, component.Accent), 36, component.Accent))
fmt.Println(menu.View(32, component.Accent))
```

`Tabs`, `Table`, `Sparkline`, `Breadcrumb`, `Panel`, and `KeyHint` cover common layout and command-bar
patterns. They are intentionally string-returning so the same components can
be used with Bubble Tea, SSH sessions, or a plain CLI.

## Dither-kit components

The `component` package also contains terminal equivalents for dither-kit’s
standalone pieces. They preserve the same seven-colour palette, four texture
variants, seeded avatar behaviour, and bloom states. A terminal cannot blur a
cell like a browser canvas, so bloom is represented by brighter/denser glyphs.

```go
fmt.Println(component.DitherAvatar("dan", 12, component.DitherPurple, component.BloomAura))
fmt.Println(component.DitherButton("deploy →", component.DitherBlue, component.DitherGradientVariant, component.BloomLow, false, false, false))
fmt.Println(component.DitherGradient(32, 4, component.DitherPurple, component.DitherBlue, "up", component.BloomLow))
fmt.Println(component.DitherSparkline([]float64{3, 7, 5, 9, 8, 12}, 32, 5, component.DitherGreen, component.DitherGradientVariant, component.BloomAura))
```

The chart renderer maps dither-kit’s `AreaChart`, `BarChart`, `LineChart`,
`PieChart`, and `RadarChart` to `chart.Area`, `chart.Bar`, `chart.Line`,
`chart.Pie`, and `chart.Radar`. `Options.Selected` is the keyboard/mouse
scrub marker; `ActiveSeries` is the selected legend item. `Gradient`,
`Dotted`, `Hatched`, and `Solid` use deterministic cell textures, so two
renders of the same data do not shimmer.

## Interaction cues

The `sound` package provides a tiny `Player` interface and a terminal-native
implementation. `Bell` emits `\a`; terminals may play it, flash, or ignore it.
Applications can implement `Player` to synthesize richer cues like `tick`,
`toggle`, or `success`, or use `Silent` for quiet/accessibility modes.
