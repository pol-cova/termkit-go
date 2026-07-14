# Dither Kit fidelity map

This library targets the public reference at
[tripwire.sh/dither-kit](https://www.tripwire.sh/dither-kit). The reference
bundle was inspected directly so the terminal renderer uses the same visual
inputs rather than a generic chart theme.

| Reference | Go terminal equivalent | Fidelity details |
| --- | --- | --- |
| `AreaChart` | `chart.Area` | Stacked bands, Bayer gradient fill, scrub column, legend/readout, optional bloom sparkles |
| `BarChart` | `chart.Bar` | Grouped/stacked bars, same four texture variants and palette |
| `LineChart` | `chart.Line` | Connected solid/dashed terminal strokes, point markers, scrub/readout |
| `PieChart` | `chart.Pie` | Elliptical/donut slices, selected slice marker, percentage legend |
| `RadarChart` | `chart.Radar` | Independent terminal radii, rings, axes, dithered polygons, labels |
| `DitherAvatar` | `component.DitherAvatar` / `DitherAvatarHue` | FNV-seeded xorshift pattern, 8×8 mirrored logical bitmap, numeric hue |
| `DitherButton` | `component.DitherButton` | Four variants, seven colors, disabled/hovered/pressed density states |
| `DitherGradient` | `component.DitherGradient` | Four directions, transparent or two-color wash, Bayer threshold |
| decorative `Sparkline` | `component.DitherSparkline` | Axis-free dither fill, connected stroke, deterministic sparkle accents |

## Reference constants ported

- Ordered-dither matrix: 4×4 Bayer order `0,8,2,10 / 12,4,14,6 /
  3,11,1,9 / 15,7,13,5`.
- Named colors: blue `(53,143,243)`, purple `(150,110,255)`, green
  `(40,210,110)`, pink `(240,90,190)`, orange `(255,150,50)`, red
  `(240,70,70)`, grey `(92,92,100)`.
- Bloom presets: off, low, high, aura. Browser blur is represented by
  brighter/denser terminal cells because terminal output has no pixel blur
  primitive.
- Avatar seed: FNV-1a followed by xorshift32, with the same 32-bit logical
  pattern size and mirror choices.

Use `go run ./cmd/demo --dither-kit` to inspect all standalone terminal
components together. Use `go run ./cmd/demo --static` to inspect the chart
surface used by the automated render tests.
