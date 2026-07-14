# Charts

`chart.Render` accepts a `chart.Chart` and returns a string ready for Bubble Tea or standard output.

```go
chart.Chart{
    Kind: chart.Bar,
    StackType: chart.Stacked,
    Labels: []string{"Jan", "Feb", "Mar"},
    Series: []chart.Series{
        {Name: "desktop", Values: []float64{120, 190, 230}, Variant: chart.Dotted},
        {Name: "mobile", Values: []float64{80, 120, 140}, Variant: chart.Hatched},
    },
}
```

## Kinds

`chart.Area`, `chart.Line`, `chart.Bar`, `chart.Pie`, and `chart.Radar` are available.

## Fill variants

- `chart.Gradient` is a dense dither fill.
- `chart.Dotted` uses a light terminal dither.
- `chart.Hatched` uses a denser terminal texture for a second stacked series.
- `chart.Solid` uses a dense block fill.

## Stacking

- `chart.Default` overlays cartesian series.
- `chart.Stacked` builds bands or bars cumulatively.
- `chart.Percent` normalizes each cartesian point to 100%.

`Options.Selected` controls the scrub column and readout. `Options.ActiveSeries` controls the selected series in pie charts and the readout. Set `Options.Color` to `false` for plain text snapshots.

For the dither-kit interaction model, set `Options.Interactive` and provide
`Options.HoveredSeries` while handling pointer/keyboard events in the host
TUI. The active series is retained at full colour and other series are dimmed,
matching the reference legend spotlight. `Series.Color` accepts `blue`,
`purple`, `green`, `pink`, `orange`, `red`, or `grey`.

## Terminal-friendly geometry

- **Area** charts use a stronger lower-half block edge (`▄`) and keep the scrub column visible over filled regions, so the trend remains readable in monochrome terminals.
- **Pie** charts use the full plotting ellipse and include a slice legend with the label, percentage, and selected-slice marker. Values below zero are treated as zero.
- **Radar** charts size their horizontal and vertical radii independently, using the available terminal width instead of collapsing to the shorter height. Axis labels are listed beneath the plot.

For stable layouts, start with `Width: 48` and `Height: 12`, then reduce those values after checking the smallest terminal you support. All series must contain the same number of values; this prevents a resize or selection event from producing a partial chart or panic.
