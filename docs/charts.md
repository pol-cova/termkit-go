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
- `chart.Dotted` gives the blue dot-field treatment.
- `chart.Hatched` uses a diagonal violet hatch.
- `chart.Solid` uses a dense block fill.

## Stacking

- `chart.Default` overlays cartesian series.
- `chart.Stacked` builds bands or bars cumulatively.
- `chart.Percent` normalizes each cartesian point to 100%.

`Options.Selected` controls the scrub column and readout. `Options.ActiveSeries` controls the selected series in pie charts and the readout. Set `Options.Color` to `false` for plain text snapshots.
