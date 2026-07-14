# Pixel backend

The screenshot-style Dither Kit charts are best understood as a raster, not as
Unicode bar characters. `chart.RenderPixel` keeps the chart model unchanged,
but renders its geometry into `pixel.Canvas` at logical-pixel resolution.

```go
c, _ := chart.RenderPixel(chart.Chart{Kind: chart.Area, Series: series}, chart.PixelOptions{
    Width: 192, Height: 48,
})
fmt.Print(c.ANSI())   // works in ordinary truecolor terminals
```

There are three output strategies:

- `Canvas.ANSI()` uses one `▀`/`▄` cell for two vertical pixels. This is the
  portable fallback and preserves the dither texture, palette, and geometry,
  but terminal cells still have character-grid proportions.
- `Canvas.Kitty()` emits a PNG through Kitty's graphics protocol. This is the
  closest match because the terminal displays the actual raster.
- `Canvas.ITerm2()` emits the same raster through iTerm2's inline-image
  protocol.

Run `go run ./cmd/demo --pixel` to see all five chart geometries through the
ANSI fallback. In a Kitty or iTerm2 session, send the protocol string directly
to stdout instead of `ANSI()`.

This is the practical “someone actually implemented it” path: terminal image
viewers such as `chafa`, `viu`, and Kitty's graphics protocol use the same
concept—encode a small raster and negotiate the richest display protocol. A
strict 1:1 browser match is only possible on a terminal with an inline-raster
protocol. Plain ANSI cannot reproduce CSS blur, rounded anti-aliased corners,
or arbitrary browser pixels; it can reproduce the visual language and dither
pattern at the terminal's available resolution.
