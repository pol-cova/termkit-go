package chart

import (
	"math"

	"github.com/pol-cova/termkit-go/pixel"
)

// PixelOptions controls the raster renderer. Width and Height are logical
// pixels, not terminal columns and rows. For a standard ANSI terminal use a
// height of roughly rows*2; Kitty/iTerm2 can display the raster directly.
type PixelOptions struct {
	Width, Height int
	Background    pixel.RGBA
}

// RenderPixel renders the same chart model as Render, but keeps every dither
// decision at pixel resolution. The returned canvas can be sent through ANSI
// half-blocks, Kitty, or iTerm2.
func RenderPixel(c Chart, o PixelOptions) (pixel.Canvas, error) {
	if err := c.validate(); err != nil {
		return pixel.Canvas{}, err
	}
	if o.Width < 8 {
		o.Width = 96
	}
	if o.Height < 8 {
		o.Height = 32
	}
	if o.Background.A == 0 {
		o.Background = pixel.RGBA{R: 39, G: 40, B: 53, A: 255}
	}
	out := pixel.New(o.Width, o.Height, o.Background)
	if c.Kind == Pie {
		renderPixelPie(&out, c)
	} else if c.Kind == Radar {
		renderPixelRadar(&out, c)
	} else {
		renderPixelCartesian(&out, c)
	}
	return out, nil
}

var pixelPalette = map[string]pixel.RGBA{
	"green":  {R: 40, G: 210, B: 110, A: 255},
	"blue":   {R: 53, G: 143, B: 243, A: 255},
	"purple": {R: 150, G: 110, B: 255, A: 255},
	"pink":   {R: 240, G: 90, B: 190, A: 255},
	"orange": {R: 255, G: 150, B: 50, A: 255},
	"red":    {R: 240, G: 70, B: 70, A: 255},
	"grey":   {R: 92, G: 92, B: 100, A: 255},
}

var pixelBayer = [4][4]int{{0, 8, 2, 10}, {12, 4, 14, 6}, {3, 11, 1, 9}, {15, 7, 13, 5}}

func seriesPixelColor(s Series, index int) pixel.RGBA {
	if c, ok := pixelPalette[s.Color]; ok {
		return c
	}
	colors := []string{"green", "blue", "purple", "pink", "orange", "red", "grey"}
	return pixelPalette[colors[index%len(colors)]]
}

func renderPixelCartesian(out *pixel.Canvas, c Chart) {
	maxValue := chartMaximum(c)
	if maxValue <= 0 {
		maxValue = 1
	}
	for si, s := range c.Series {
		col := seriesPixelColor(s, si)
		for i, value := range s.Values {
			x0 := i * (out.Width - 1) / max(1, len(s.Values)-1)
			x1 := x0
			if i+1 < len(s.Values) {
				x1 = (i + 1) * (out.Width - 1) / max(1, len(s.Values)-1)
			}
			y := out.Height - 1 - int(value/maxValue*float64(out.Height-2))
			nextY := y
			if i+1 < len(s.Values) {
				nextY = out.Height - 1 - int(s.Values[i+1]/maxValue*float64(out.Height-2))
			}
			switch c.Kind {
			case Line:
				drawPixelLine(out, x0, y, x1, nextY, col)
			case Bar:
				barW := max(1, (out.Width/max(1, len(s.Values)))/2)
				ditherPixelRect(out, x0-barW/2, y, barW, out.Height-y-1, col, s.Variant)
			default:
				ditherPixelRect(out, x0, y, max(1, x1-x0+1), out.Height-y-1, col, s.Variant)
			}
		}
	}
}

func ditherPixelRect(out *pixel.Canvas, x, y, w, h int, col pixel.RGBA, variant Variant) {
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			keep := true
			switch variant {
			case Dotted:
				keep = px%4 == 0 || py%4 == 0
			case Hatched:
				keep = (px+py)%6 < 2
			case Gradient:
				keep = pixelBayer[py%4][px%4] < (py-y+1)*16/max(1, h)
			}
			if keep {
				out.Set(px, py, col)
			}
		}
	}
}

func drawPixelLine(out *pixel.Canvas, x0, y0, x1, y1 int, col pixel.RGBA) {
	dx, sx := abs(x1-x0), 1
	if x0 > x1 {
		sx = -1
	}
	dy, sy := -abs(y1-y0), 1
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy
	for {
		out.Set(x0, y0, col)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func renderPixelPie(out *pixel.Canvas, c Chart) {
	cx, cy := out.Width/2, out.Height/2
	radius := float64(min(out.Width, out.Height))/2 - 1
	total := 0.0
	for _, s := range c.Series {
		for _, v := range s.Values {
			total += math.Max(0, v)
		}
	}
	if total == 0 {
		return
	}
	for y := 0; y < out.Height; y++ {
		for x := 0; x < out.Width; x++ {
			dx, dy := float64(x-cx), float64(y-cy)
			if math.Hypot(dx, dy) > radius {
				continue
			}
			angle := math.Atan2(dy, dx)
			if angle < 0 {
				angle += 2 * math.Pi
			}
			cursor := 0.0
			for i, s := range c.Series {
				for _, v := range s.Values {
					portion := math.Max(0, v) / total * 2 * math.Pi
					if angle >= cursor && angle < cursor+portion {
						if pixelBayer[y%4][x%4] < 12 {
							out.Set(x, y, seriesPixelColor(s, i))
						}
						break
					}
					cursor += portion
				}
			}
		}
	}
}

func renderPixelRadar(out *pixel.Canvas, c Chart) {
	cx, cy := out.Width/2, out.Height/2
	radius := min(out.Width, out.Height)/2 - 2
	for i := 0; i < 8; i++ {
		a := float64(i) * math.Pi / 4
		drawPixelLine(out, cx, cy, cx+int(math.Cos(a)*float64(radius)), cy+int(math.Sin(a)*float64(radius)), pixelPalette["grey"])
	}
	for si, s := range c.Series {
		col := seriesPixelColor(s, si)
		for i, v := range s.Values {
			a := float64(i) * 2 * math.Pi / float64(len(s.Values))
			r := v / chartMaximum(c) * float64(radius)
			x := cx + int(math.Cos(a)*r)
			y := cy + int(math.Sin(a)*r)
			out.Set(x, y, col)
			if i > 0 { /* sparse joins preserve the raster dither feel */
			}
		}
	}
}
