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

## Interaction cues

The `sound` package provides a tiny `Player` interface and a terminal-native
implementation. `Bell` emits `\a`; terminals may play it, flash, or ignore it.
Applications can implement `Player` to synthesize richer cues like `tick`,
`toggle`, or `success`, or use `Silent` for quiet/accessibility modes.
