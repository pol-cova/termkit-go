# Motion and components

The animation helpers are pure functions. Your application owns the timer and passes a normalized value or frame count.

```go
progress := animate.Tween(0.6, animate.EaseInOut)
fmt.Println(component.Progress("Deploy", progress, 20, component.Accent))
fmt.Println(component.SpinnerFrame(frame, "Connecting", component.Warning))
```

Use `animate.Spring`, `animate.Pulse`, and `animate.Repeat` for a natural-looking response without background goroutines. Components return ANSI strings and work with Bubble Tea, Cobra, SSH, or `fmt.Println`.
