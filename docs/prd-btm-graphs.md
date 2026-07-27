# PRD: btm-Faithful Terminal Time-Series Graphs

## Objective

Add dense, live-monitor graph rendering to termkit-go that matches the visual grammar of btm-style system monitors (braille multi-series time plots, axes, in-plot legends, bordered titled panes) while staying framework-agnostic. Host applications push samples; the library renders strings.

## Problem

termkit-go ships dithered area/line charts and string UI primitives, but they do not match the dense monitor graph language users expect:

- No braille/dot-matrix time-series backend (current glyphs are block/dither).
- No first-class rolling history, dual-metric overlays, or in-plot legend boxes.
- No reusable layout/composition model—hosts manually concatenate strings.
- Adding a new pane today means ad-hoc formatting, not registering a composable widget.

## Goals

| Goal | Success signal |
| --- | --- |
| Visual parity | Side-by-side comparison with the reference shows the same graph grammar (braille lines, axes, legend overlay, titled borders) |
| Terminal-native UX | Readable at small sizes, graceful resize, color-safe fallbacks |
| Scalable components | New pane = implement one interface + register; no chart-package forks |
| Library purity | Zero Bubble Tea / OS coupling in `chart` and `component`; optional helpers stay in `bubbletea` / `cmd/demo` |

## Non-goals

- Cloning btm as a full system monitor product
- Replacing the existing Dither Kit chart aesthetic (coexist via separate render mode)
- OpenTUI runtime dependency
- Mouse-only interaction (keyboard remains primary)
- OS metrics collectors inside the library

## Users

- **CLI/TUI authors** embedding live metrics graphs in custom dashboards
- **Demo consumers** validating look-and-feel against the reference image
- **Contributors** adding panes (CPU sidebar, network, custom tables) without touching renderer internals

## Reference UI grammar

Graphs must support:

1. **Braille time-series** — multi-color polylines on a character grid (2×4 dots per cell).
2. **Cartesian frame** — Y ticks (0–100% or SI rates), X ticks as relative window (`60s` … `0s`).
3. **In-plot legend box** — e.g. `RAM: 69% 25.0GiB/36.0GiB` / `RX: … All: …` with series color.
4. **Border chrome** — box-drawing pane with title spliced into the top border.
5. **Semantic colors** — primary series warm (orange), secondary cool (teal); selection cyan bar for companion lists.
6. **Companion surfaces** — reuse/extend `component.Table` + `component.BorderTitle`.

## Terminal UX optimizations

| Optimization | Rationale |
| --- | --- |
| Adaptive density | Prefer braille when size allows; fall back to sparkline when below threshold |
| Stable layout slots | Panes keep reserved rows so legend/axes don't jump |
| Rolling buffer API | `History` so hosts push samples without reallocating each tick |
| NO_COLOR support | Monochrome braille + bold/dim hierarchy |
| Narrow terminals | Collapse legend; hide secondary Y ticks before clipping plot |
| Scrub + focus | Keyboard scrub with readout that doesn't cover the whole plot |
| Reduced motion | Host sets `Animate: false` to skip tick updates |

## Product principles

1. **Look like btm; feel like a library** — visual fidelity without app ownership.
2. **Compose, don't fork** — graphs are widgets; layout is orthogonal.
3. **Density with dignity** — braille for information; never sacrifice axis legibility.
4. **Host owns time** — library stores samples the host pushes; no goroutines inside render packages.

## Milestones

See [spec-timeseries.md](spec-timeseries.md) for architecture, APIs, acceptance criteria, and delivery order.
