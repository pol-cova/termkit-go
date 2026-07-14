package component

import (
	"fmt"
	"strings"
)

// DitherColor is the seven-colour palette used by dither-kit.
type DitherColor string

const (
	DitherBlue   DitherColor = "blue"
	DitherPurple DitherColor = "purple"
	DitherGreen  DitherColor = "green"
	DitherPink   DitherColor = "pink"
	DitherOrange DitherColor = "orange"
	DitherRed    DitherColor = "red"
	DitherGrey   DitherColor = "grey"
)

// Short aliases mirror dither-kit's public vocabulary.
const (
	VariantGradient = DitherGradientVariant
	VariantDotted   = DitherDottedVariant
	VariantHatched  = DitherHatchedVariant
	VariantSolid    = DitherSolidVariant
)

// DitherVariant controls the ordered-dither texture of a standalone component.
type DitherVariant string

const (
	DitherGradientVariant DitherVariant = "gradient"
	DitherDottedVariant   DitherVariant = "dotted"
	DitherHatchedVariant  DitherVariant = "hatched"
	DitherSolidVariant    DitherVariant = "solid"
)

// Bloom controls how much visual emphasis a component receives. Terminals
// cannot blur pixels, so bloom is represented by denser/brighter glyphs.
type Bloom string

const (
	BloomOff  Bloom = "off"
	BloomLow  Bloom = "low"
	BloomHigh Bloom = "high"
	BloomAura Bloom = "aura"
)

// DitherButton renders the compact pixel-button treatment from dither-kit.
// hovered and pressed let event-loop owners map their own input state.
func DitherButton(label string, color DitherColor, variant DitherVariant, bloom Bloom, hovered, pressed, disabled bool) string {
	if variant == "" {
		variant = DitherGradientVariant
	}
	if bloom == "" {
		bloom = BloomOff
	}
	left, right := "[", "]"
	if disabled {
		color = DitherGrey
	}
	fill := ditherFill(variant, len([]rune(label))+2, int(bloomDensity(bloom, hovered, pressed)))
	text := left + fill + " " + label + " " + fill + right
	return paint256(text, ditherANSI(color, hovered || pressed))
}

// DitherGradient renders a directional dither wash as terminal rows.
func DitherGradient(width, height int, from, to DitherColor, direction string, bloom Bloom) string {
	if width < 1 || height < 1 {
		return ""
	}
	if from == "" {
		from = DitherPurple
	}
	if direction == "" {
		direction = "up"
	}
	rows := make([]string, height)
	for y := 0; y < height; y++ {
		var row strings.Builder
		for x := 0; x < width; x++ {
			var progress float64
			switch direction {
			case "down":
				progress = 1 - (float64(y)+.5)/float64(height)
			case "left":
				progress = (float64(x) + .5) / float64(width)
			case "right":
				progress = 1 - (float64(x)+.5)/float64(width)
			default:
				progress = (float64(y) + .5) / float64(height)
			}
			threshold := bayer4[y&3][x&3] / 16
			if progress <= threshold {
				row.WriteByte(' ')
				continue
			}
			choice := from
			if to != "" && to != from && progress < .5 {
				choice = to
			}
			if to == "" && progress < .22 {
				row.WriteByte(' ')
				continue
			}
			glyph := "░"
			if progress > .72 {
				glyph = "▒"
			}
			if progress > .9 {
				glyph = "▓"
			}
			row.WriteString(paintRGB(glyph, ditherRGB(choice), bloom != BloomOff))
		}
		rows[y] = row.String()
	}
	return strings.Join(rows, "\n")
}

// DitherSparkline is the decorative, axis-free sparkline exposed by dither-kit.
// It uses one terminal cell per column and keeps a small dithered area under
// the interpolated stroke so it remains legible in plain terminals.
func DitherSparkline(values []float64, width, height int, color DitherColor, variant DitherVariant, bloom Bloom) string {
	if len(values) == 0 || width < 1 || height < 2 {
		return ""
	}
	if width < 4 {
		width = 4
	}
	if height < 3 {
		height = 3
	}
	if color == "" {
		color = DitherGreen
	}
	if variant == "" {
		variant = DitherGradientVariant
	}
	points := resample(values, width)
	lo, hi := points[0], points[0]
	for _, v := range points[1:] {
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	plot := make([][]rune, height)
	for y := range plot {
		plot[y] = []rune(strings.Repeat(" ", width))
	}
	pointY := make([]int, width)
	for x, v := range points {
		ratio := 0.5
		if hi > lo {
			ratio = (v - lo) / (hi - lo)
		}
		pointY[x] = height - 1 - int(ratio*float64(height-2))
	}
	for x, top := range pointY {
		for y := top + 1; y < height; y++ {
			plot[y][x] = []rune(ditherGlyph(variant, x, y))[0]
		}
		plot[top][x] = '•'
		if x > 0 {
			connectSpark(plot, pointY[x-1], top, x)
		}
	}
	if bloom == BloomAura || bloom == BloomHigh {
		for x, y := range pointY {
			if (x*7+y*3)%11 == 0 {
				plot[y][x] = '✦'
			}
		}
	}
	var b strings.Builder
	for y, row := range plot {
		if y > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(paintRGB(string(row), ditherRGB(color), bloom != BloomOff))
	}
	return b.String()
}

func resample(values []float64, width int) []float64 {
	out := make([]float64, width)
	for x := range out {
		p := float64(x) * float64(len(values)-1) / float64(maxInt(width-1, 1))
		left := int(p)
		right := minInt(left+1, len(values)-1)
		f := p - float64(left)
		out[x] = values[left] + (values[right]-values[left])*f
	}
	return out
}

func connectSpark(plot [][]rune, from, to, x int) {
	if from == to {
		plot[to][x] = '•'
		return
	}
	if to < from {
		plot[to][x] = '╱'
	} else {
		plot[to][x] = '╲'
	}
}

func ditherRGB(c DitherColor) [3]int {
	switch c {
	case DitherGreen:
		return [3]int{40, 210, 110}
	case DitherBlue:
		return [3]int{53, 143, 243}
	case DitherPurple:
		return [3]int{150, 110, 255}
	case DitherPink:
		return [3]int{240, 90, 190}
	case DitherOrange:
		return [3]int{255, 150, 50}
	case DitherRed:
		return [3]int{240, 70, 70}
	}
	return [3]int{92, 92, 100}
}

// DitherAvatar returns a deterministic mirrored pixel avatar. Equal names
// always produce equal output, matching dither-kit's seeded avatar contract.
func DitherAvatar(name string, size int, hue DitherColor, bloom Bloom) string {
	return DitherAvatarHue(name, size, ditherHue(hue), "auto", bloom)
}

// DitherAvatarHue is the numeric-hue form used by the web component. mirror
// accepts "auto", "horizontal", or "vertical".
func DitherAvatarHue(name string, size, hue int, mirror string, bloom Bloom) string {
	if size < 4 {
		size = 4
	}
	if size > 32 {
		size = 32
	}
	seed := fnvSeed(name)
	rng := func() float64 {
		seed ^= seed << 13
		seed ^= seed >> 17
		seed ^= seed << 5
		return float64(seed) / float64(uint64(1)<<32)
	}
	on := make([]bool, 32)
	for i := range on {
		on[i] = rng() < .5
	}
	autoMirror := rng() < .5
	rotation := int(rng()*180) * 2
	density := make([]float64, 32)
	for i := range density {
		density[i] = .55 + rng()*.45
	}
	vertical := mirror == "vertical" || (mirror == "auto" && autoMirror)
	rows := make([]string, size)
	for y := 0; y < size; y++ {
		var b strings.Builder
		for x := 0; x < size; x++ {
			// Map terminal cells onto the reference's 8×8 logical bitmap.
			py, px := y*8/size, x*8/size
			if vertical {
				py = minInt(py, 7-py)
			} else {
				px = minInt(px, 7-px)
			}
			idx := avatarIndex(py, px, vertical)
			if idx >= len(on) || !on[idx] {
				b.WriteRune(' ')
				continue
			}
			threshold := bayer4[py&3][px&3] / 16.0
			bright := density[idx]
			if bright < threshold*.7 {
				b.WriteRune('░')
			} else if rotation%4 == 0 {
				b.WriteRune('█')
			} else {
				b.WriteRune('▓')
			}
		}
		rows[y] = paintRGB(b.String(), hueRGB(hue), bloom != BloomOff)
	}
	return strings.Join(rows, "\n")
}

var bayer4 = [4][4]float64{{.5, 8.5, 2.5, 10.5}, {12.5, 4.5, 14.5, 6.5}, {3.5, 11.5, 1.5, 9.5}, {15.5, 7.5, 13.5, 5.5}}

func fnvSeed(name string) uint32 {
	// JavaScript's reference hashes UTF-16 charCodeAt units, not UTF-8 bytes.
	h := uint32(2166136261)
	for _, r := range name {
		if r <= 0xffff {
			h ^= uint32(r)
			h *= 16777619
			continue
		}
		r -= 0x10000
		for _, unit := range []rune{0xd800 + (r >> 10), 0xdc00 + (r & 0x3ff)} {
			h ^= uint32(unit)
			h *= 16777619
		}
	}
	return h
}

func avatarIndex(y, x int, vertical bool) int {
	if vertical {
		return minInt(y, 7-y)*8 + x
	}
	return y*4 + minInt(x, 7-x)
}

func ditherHue(c DitherColor) int {
	switch c {
	case DitherBlue:
		return 215
	case DitherGreen:
		return 145
	case DitherPink:
		return 325
	case DitherOrange:
		return 30
	case DitherRed:
		return 0
	case DitherGrey:
		return 0
	}
	return 270
}

func hueRGB(h int) [3]int {
	h = ((h % 360) + 360) % 360
	switch {
	case h < 20 || h >= 345:
		return [3]int{240, 70, 70}
	case h < 65:
		return [3]int{255, 150, 50}
	case h < 180:
		return [3]int{40, 210, 110}
	case h < 240:
		return [3]int{53, 143, 243}
	case h < 310:
		return [3]int{150, 110, 255}
	default:
		return [3]int{240, 90, 190}
	}
}

func ditherFill(v DitherVariant, width, density int) string {
	g := ditherGlyph(v, density, width)
	return strings.Repeat(g, maxInt(1, width/3))
}

func ditherGlyph(v DitherVariant, x, y int) string {
	switch v {
	case DitherSolidVariant:
		return "█"
	case DitherHatchedVariant:
		return "▒"
	case DitherDottedVariant:
		return "·"
	default:
		return []string{"░", "▒", "▓", "█"}[(x*3+y*5)&3]
	}
}

func bloomDensity(b Bloom, hovered, pressed bool) int {
	d := 0
	if b == BloomLow {
		d = 1
	}
	if b == BloomHigh {
		d = 2
	}
	if b == BloomAura {
		d = 3
	}
	if hovered {
		d++
	}
	if pressed {
		d++
	}
	return d
}

func ditherANSI(c DitherColor, bright bool) int {
	colors := map[DitherColor]int{DitherGreen: 46, DitherBlue: 33, DitherPurple: 99, DitherPink: 205, DitherOrange: 208, DitherRed: 196, DitherGrey: 244}
	v, ok := colors[c]
	if !ok {
		v = colors[DitherPurple]
	}
	if bright && v < 240 {
		v++
	}
	return v
}

// DitherColorCode is useful when composing a component with custom ANSI.
func DitherColorCode(c DitherColor) string { return fmt.Sprintf("\x1b[38;5;%dm", ditherANSI(c, false)) }

func paint256(value string, color int) string {
	return fmt.Sprintf("\x1b[38;5;%dm%s\x1b[0m", color, value)
}

func paintRGB(value string, rgb [3]int, bright bool) string {
	if bright {
		rgb[0] = minInt(255, rgb[0]+20)
		rgb[1] = minInt(255, rgb[1]+20)
		rgb[2] = minInt(255, rgb[2]+20)
	}
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", rgb[0], rgb[1], rgb[2], value)
}
