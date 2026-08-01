package chart

import (
	"math"

	"github.com/pol-cova/termkit-go/pixel"
	"github.com/pol-cova/termkit-go/style"
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
	layout := NewLayout(c)
	if c.Kind == Pie {
		renderPixelPie(&out, c, layout)
	} else if c.Kind == Radar {
		renderPixelRadar(&out, c, layout)
	} else {
		renderPixelCartesian(&out, c, layout)
	}
	return out, nil
}

func seriesPixelColor(s Series, index int) pixel.RGBA {
	return style.Resolve(s.Color, index).PixelRGBA()
}

func renderPixelCartesian(out *pixel.Canvas, c Chart, layout Layout) {
	for si, s := range c.Series {
		col := seriesPixelColor(s, si)
		for i, value := range s.Values {
			if math.IsNaN(value) {
				continue
			}
			x0 := layout.PointX(i, out.Width)
			x1 := x0
			if i+1 < len(s.Values) {
				x1 = layout.PointX(i+1, out.Width)
			}
			y := layout.PlotY(value, out.Height)
			nextY := y
			if i+1 < len(s.Values) && !math.IsNaN(s.Values[i+1]) {
				nextY = layout.PlotY(s.Values[i+1], out.Height)
			}
			switch c.Kind {
			case Line:
				drawPixelLine(out, x0, y, x1, nextY, col)
			case Bar:
				barW := max(1, (out.Width/max(1, layout.Points))/2)
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
				keep = style.Bayer4[py%4][px%4] < (py-y+1)*16/max(1, h)
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

func renderPixelPie(out *pixel.Canvas, c Chart, layout Layout) {
	cx, cy := out.Width/2, out.Height/2
	radius := float64(min(out.Width, out.Height))/2 - 1
	total := 0.0
	for _, s := range c.Series {
		for _, v := range s.Values {
			if !math.IsNaN(v) {
				total += math.Max(0, v)
			}
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
					if math.IsNaN(v) {
						continue
					}
					portion := math.Max(0, v) / total * 2 * math.Pi
					if angle >= cursor && angle < cursor+portion {
						if style.Bayer4[y%4][x%4] < 12 {
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

func renderPixelRadar(out *pixel.Canvas, c Chart, layout Layout) {
	cx, cy := out.Width/2, out.Height/2
	radius := min(out.Width, out.Height)/2 - 2
	for i := 0; i < 8; i++ {
		a := float64(i) * math.Pi / 4
		drawPixelLine(out, cx, cy, cx+int(math.Cos(a)*float64(radius)), cy+int(math.Sin(a)*float64(radius)), style.Grey.PixelRGBA())
	}
	maxValue := layout.MaxValue
	for si, s := range c.Series {
		col := seriesPixelColor(s, si)
		for i, v := range s.Values {
			if math.IsNaN(v) {
				continue
			}
			a := layout.Angle(i)
			r := v / maxValue * float64(radius)
			x := cx + int(math.Cos(a)*r)
			y := cy + int(math.Sin(a)*r)
			out.Set(x, y, col)
		}
	}
}
