// Package style holds shared Dither Kit palette and dither constants.
package style

import "github.com/pol-cova/termkit-go/pixel"

type Color string

const (
	Blue   Color = "blue"
	Purple Color = "purple"
	Green  Color = "green"
	Pink   Color = "pink"
	Orange Color = "orange"
	Red    Color = "red"
	Grey   Color = "grey"
)

type Variant string

const (
	Gradient Variant = "gradient"
	Dotted   Variant = "dotted"
	Hatched  Variant = "hatched"
	Solid    Variant = "solid"
)

type Bloom string

const (
	BloomOff  Bloom = "off"
	BloomLow  Bloom = "low"
	BloomHigh Bloom = "high"
	BloomAura Bloom = "aura"
)

var Bayer4 = [4][4]int{{0, 8, 2, 10}, {12, 4, 14, 6}, {3, 11, 1, 9}, {15, 7, 13, 5}}

var PaletteOrder = []Color{Green, Blue, Purple, Pink, Orange, Red, Grey}

var rgb = map[Color][3]int{
	Green:  {40, 210, 110},
	Blue:   {53, 143, 243},
	Purple: {150, 110, 255},
	Pink:   {240, 90, 190},
	Orange: {255, 150, 50},
	Red:    {240, 70, 70},
	Grey:   {92, 92, 100},
}

var ansi256 = map[Color]int{
	Green: 46, Blue: 33, Purple: 99, Pink: 205, Orange: 208, Red: 196, Grey: 244,
}

func (c Color) RGB() [3]int {
	if v, ok := rgb[c]; ok {
		return v
	}
	return rgb[Purple]
}

func (c Color) PixelRGBA() pixel.RGBA {
	r := c.RGB()
	return pixel.RGBA{R: uint8(r[0]), G: uint8(r[1]), B: uint8(r[2]), A: 255}
}

func (c Color) ANSI256(bright bool) int {
	v, ok := ansi256[c]
	if !ok {
		v = ansi256[Purple]
	}
	if bright && v < 240 {
		v++
	}
	return v
}

func Resolve(name string, index int) Color {
	if _, ok := rgb[Color(name)]; ok {
		return Color(name)
	}
	return PaletteOrder[index%len(PaletteOrder)]
}

func BloomDensity(b Bloom, hovered, pressed bool) int {
	d := 0
	switch b {
	case BloomLow:
		d = 1
	case BloomHigh:
		d = 2
	case BloomAura:
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

func (c Color) Hue() int {
	switch c {
	case Blue:
		return 215
	case Green:
		return 145
	case Pink:
		return 325
	case Orange:
		return 30
	case Red:
		return 0
	case Grey:
		return 0
	default:
		return 270
	}
}

func HueRGB(h int) [3]int {
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
