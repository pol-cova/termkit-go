package component

import (
	"fmt"
	"strings"

	"github.com/pol-cova/termkit-go/style"
)

// DitherColor is the seven-colour palette used by dither-kit.
type DitherColor = style.Color

const (
	DitherBlue   = style.Blue
	DitherPurple = style.Purple
	DitherGreen  = style.Green
	DitherPink   = style.Pink
	DitherOrange = style.Orange
	DitherRed    = style.Red
	DitherGrey   = style.Grey
)

// Short aliases mirror dither-kit's public vocabulary.
const (
	VariantGradient = DitherGradientVariant
	VariantDotted   = DitherDottedVariant
	VariantHatched  = DitherHatchedVariant
	VariantSolid    = DitherSolidVariant
)

// DitherVariant controls the ordered-dither texture of a standalone component.
type DitherVariant = style.Variant

const (
	DitherGradientVariant = style.Gradient
	DitherDottedVariant   = style.Dotted
	DitherHatchedVariant  = style.Hatched
	DitherSolidVariant    = style.Solid
)

// Bloom controls how much visual emphasis a component receives. Terminals
// cannot blur pixels, so bloom is represented by denser/brighter glyphs.
type Bloom = style.Bloom

const (
	BloomOff  = style.BloomOff
	BloomLow  = style.BloomLow
	BloomHigh = style.BloomHigh
	BloomAura = style.BloomAura
)

type DitherButtonOptions struct {
	Label            string
	Color            DitherColor
	Variant          DitherVariant
	Bloom            Bloom
	Hovered, Pressed bool
	Disabled         bool
}

type DitherGradientOptions struct {
	Width, Height int
	From, To      DitherColor
	Direction     string
	Bloom         Bloom
}

type DitherSparklineOptions struct {
	Values        []float64
	Width, Height int
	Color         DitherColor
	Variant       DitherVariant
	Bloom         Bloom
}

type DitherAvatarOptions struct {
	Name  string
	Size  int
	Hue   DitherColor
	Bloom Bloom
}

type DitherAvatarHueOptions struct {
	Name, Mirror string
	Size, Hue    int
	Bloom        Bloom
}

// DitherButton renders the compact pixel-button treatment from dither-kit.
// hovered and pressed let event-loop owners map their own input state.
func DitherButton(label string, color DitherColor, variant DitherVariant, bloom Bloom, hovered, pressed, disabled bool) string {
	return DitherButtonWith(DitherButtonOptions{
		Label: label, Color: color, Variant: variant, Bloom: bloom,
		Hovered: hovered, Pressed: pressed, Disabled: disabled,
	})
}

func DitherButtonWith(o DitherButtonOptions) string {
	if o.Variant == "" {
		o.Variant = DitherGradientVariant
	}
	if o.Bloom == "" {
		o.Bloom = BloomOff
	}
	left, right := "[", "]"
	color := o.Color
	if o.Disabled {
		color = DitherGrey
	}
	fill := ditherFill(o.Variant, runeWidth(o.Label)+2, style.BloomDensity(o.Bloom, o.Hovered, o.Pressed))
	text := left + fill + " " + o.Label + " " + fill + right
	return paint256(text, color.ANSI256(o.Hovered || o.Pressed))
}

// DitherGradient renders a directional dither wash as terminal rows.
func DitherGradient(width, height int, from, to DitherColor, direction string, bloom Bloom) string {
	return DitherGradientWith(DitherGradientOptions{
		Width: width, Height: height, From: from, To: to, Direction: direction, Bloom: bloom,
	})
}

func DitherGradientWith(o DitherGradientOptions) string {
	if o.Width < 1 || o.Height < 1 {
		return ""
	}
	if o.From == "" {
		o.From = DitherPurple
	}
	if o.Direction == "" {
		o.Direction = "up"
	}
	rows := make([]string, o.Height)
	for y := 0; y < o.Height; y++ {
		var row strings.Builder
		for x := 0; x < o.Width; x++ {
			var progress float64
			switch o.Direction {
			case "down":
				progress = 1 - (float64(y)+.5)/float64(o.Height)
			case "left":
				progress = (float64(x) + .5) / float64(o.Width)
			case "right":
				progress = 1 - (float64(x)+.5)/float64(o.Width)
			default:
				progress = (float64(y) + .5) / float64(o.Height)
			}
			threshold := float64(style.Bayer4[y&3][x&3]) / 16
			if progress <= threshold {
				row.WriteByte(' ')
				continue
			}
			choice := o.From
			if o.To != "" && o.To != o.From && progress < .5 {
				choice = o.To
			}
			if o.To == "" && progress < .22 {
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
			row.WriteString(paintRGB(glyph, choice.RGB(), o.Bloom != BloomOff))
		}
		rows[y] = row.String()
	}
	return strings.Join(rows, "\n")
}

// DitherSparkline renders a decorative, axis-free sparkline.
func DitherSparkline(values []float64, width, height int, color DitherColor, variant DitherVariant, bloom Bloom) string {
	return DitherSparklineWith(DitherSparklineOptions{
		Values: values, Width: width, Height: height, Color: color, Variant: variant, Bloom: bloom,
	})
}

func DitherSparklineWith(o DitherSparklineOptions) string {
	if len(o.Values) == 0 || o.Width < 1 || o.Height < 2 {
		return ""
	}
	if o.Width < 4 {
		o.Width = 4
	}
	if o.Height < 3 {
		o.Height = 3
	}
	if o.Color == "" {
		o.Color = DitherGreen
	}
	if o.Variant == "" {
		o.Variant = DitherGradientVariant
	}
	points := resample(o.Values, o.Width)
	lo, hi := points[0], points[0]
	for _, v := range points[1:] {
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	plot := make([][]rune, o.Height)
	for y := range plot {
		plot[y] = []rune(strings.Repeat(" ", o.Width))
	}
	pointY := make([]int, o.Width)
	for x, v := range points {
		ratio := 0.5
		if hi > lo {
			ratio = (v - lo) / (hi - lo)
		}
		pointY[x] = o.Height - 1 - int(ratio*float64(o.Height-2))
	}
	for x, top := range pointY {
		for y := top + 1; y < o.Height; y++ {
			plot[y][x] = []rune(ditherGlyph(o.Variant, x, y))[0]
		}
		plot[top][x] = '•'
		if x > 0 {
			connectSpark(plot, pointY[x-1], top, x)
		}
	}
	if o.Bloom == BloomAura || o.Bloom == BloomHigh {
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
		b.WriteString(paintRGB(string(row), o.Color.RGB(), o.Bloom != BloomOff))
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

// DitherAvatar returns a deterministic mirrored pixel avatar.
func DitherAvatar(name string, size int, hue DitherColor, bloom Bloom) string {
	return DitherAvatarWith(DitherAvatarOptions{Name: name, Size: size, Hue: hue, Bloom: bloom})
}

func DitherAvatarWith(o DitherAvatarOptions) string {
	return DitherAvatarHueWith(DitherAvatarHueOptions{
		Name: o.Name, Size: o.Size, Hue: o.Hue.Hue(), Mirror: "auto", Bloom: o.Bloom,
	})
}

// DitherAvatarHue is the numeric-hue form used by the web component.
func DitherAvatarHue(name string, size, hue int, mirror string, bloom Bloom) string {
	return DitherAvatarHueWith(DitherAvatarHueOptions{Name: name, Size: size, Hue: hue, Mirror: mirror, Bloom: bloom})
}

func DitherAvatarHueWith(o DitherAvatarHueOptions) string {
	size := o.Size
	if size < 4 {
		size = 4
	}
	if size > 32 {
		size = 32
	}
	seed := fnvSeed(o.Name)
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
	vertical := o.Mirror == "vertical" || (o.Mirror == "auto" && autoMirror)
	rows := make([]string, size)
	for y := 0; y < size; y++ {
		var b strings.Builder
		for x := 0; x < size; x++ {
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
			threshold := float64(style.Bayer4[py&3][px&3]) / 16.0
			bright := density[idx]
			if bright < threshold*.7 {
				b.WriteRune('░')
			} else if rotation%4 == 0 {
				b.WriteRune('█')
			} else {
				b.WriteRune('▓')
			}
		}
		rows[y] = paintRGB(b.String(), style.HueRGB(o.Hue), o.Bloom != BloomOff)
	}
	return strings.Join(rows, "\n")
}

func fnvSeed(name string) uint32 {
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

// DitherColorCode is useful when composing a component with custom ANSI.
func DitherColorCode(c DitherColor) string {
	return fmt.Sprintf("\x1b[38;5;%dm", c.ANSI256(false))
}

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

func runeWidth(s string) int {
	return displayWidth(s)
}
