// Package pixel provides a tiny terminal pixel canvas.
//
// It renders at two vertical pixels per terminal cell using truecolor ANSI,
// with Kitty and iTerm2 inline-image encodings available when the terminal can
// display an actual raster. This makes dense dithered UI possible without
// pretending a block character is a browser canvas pixel.
package pixel

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strings"
)

// RGBA is a compact pixel color. Alpha zero means transparent.
type RGBA struct{ R, G, B, A uint8 }

var Transparent = RGBA{}

// Canvas is a row-major raster.
type Canvas struct {
	Width, Height int
	pixels        []RGBA
}

func New(width, height int, background RGBA) Canvas {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	c := Canvas{Width: width, Height: height, pixels: make([]RGBA, width*height)}
	for i := range c.pixels {
		c.pixels[i] = background
	}
	return c
}

func (c Canvas) At(x, y int) RGBA {
	if x < 0 || y < 0 || x >= c.Width || y >= c.Height {
		return Transparent
	}
	return c.pixels[y*c.Width+x]
}

func (c *Canvas) Set(x, y int, value RGBA) {
	if x >= 0 && y >= 0 && x < c.Width && y < c.Height {
		c.pixels[y*c.Width+x] = value
	}
}

func (c *Canvas) FillRect(x, y, width, height int, value RGBA) {
	for py := max(0, y); py < min(c.Height, y+height); py++ {
		for px := max(0, x); px < min(c.Width, x+width); px++ {
			c.Set(px, py, value)
		}
	}
}

// PNG encodes the canvas for Kitty/iTerm2 or for a normal image viewer.
func (c Canvas) PNG() ([]byte, error) {
	img := image.NewNRGBA(image.Rect(0, 0, c.Width, c.Height))
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			p := c.At(x, y)
			img.SetNRGBA(x, y, color.NRGBA{R: p.R, G: p.G, B: p.B, A: p.A})
		}
	}
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// ANSI renders the raster as truecolor ANSI half-block cells. It works in
// terminals that support 24-bit color and degrades to terminal background for
// transparent pixels.
func (c Canvas) ANSI() string {
	var out strings.Builder
	for y := 0; y < c.Height; y += 2 {
		for x := 0; x < c.Width; x++ {
			top, bottom := c.At(x, y), c.At(x, y+1)
			if top.A == 0 && bottom.A == 0 {
				out.WriteByte(' ')
				continue
			}
			switch {
			case bottom.A == 0:
				fmt.Fprintf(&out, "\x1b[38;2;%d;%d;%dm▀", top.R, top.G, top.B)
			case top.A == 0:
				fmt.Fprintf(&out, "\x1b[38;2;%d;%d;%dm▄", bottom.R, bottom.G, bottom.B)
			default:
				fmt.Fprintf(&out, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", top.R, top.G, top.B, bottom.R, bottom.G, bottom.B)
			}
		}
		out.WriteString("\x1b[0m")
		if y+2 < c.Height {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

// Kitty encodes the canvas using Kitty's graphics protocol. The returned
// string is deliberately not emitted automatically; callers should only send
// it to a terminal known to support Kitty graphics.
func (c Canvas) Kitty() (string, error) {
	pngBytes, err := c.PNG()
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(pngBytes)
	var out strings.Builder
	const chunkSize = 4096
	for start := 0; start < len(encoded); start += chunkSize {
		end := min(len(encoded), start+chunkSize)
		more := 0
		if end < len(encoded) {
			more = 1
		}
		if start == 0 {
			fmt.Fprintf(&out, "\x1b_Ga=T,f=100,s=%d,v=%d,m=%d;", c.Width, c.Height, more)
		} else {
			fmt.Fprintf(&out, "\x1b_Gm=%d;", more)
		}
		out.WriteString(encoded[start:end])
		out.WriteString("\x1b\\")
	}
	return out.String(), nil
}

// ITerm2 encodes the canvas using iTerm2's inline image protocol.
func (c Canvas) ITerm2() (string, error) {
	pngBytes, err := c.PNG()
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(pngBytes)
	return fmt.Sprintf("\x1b]1337;File=inline=1;width=%dpx;height=%dpx;preserveAspectRatio=0:%s\a", c.Width, c.Height, encoded), nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
