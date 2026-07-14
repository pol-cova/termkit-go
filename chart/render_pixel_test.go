package chart

import "testing"

func TestRenderPixelSupportsEveryDitherKitChartKind(t *testing.T) {
	for _, kind := range []Kind{Area, Line, Bar, Pie, Radar} {
		c, err := RenderPixel(Chart{Kind: kind, Series: []Series{{Values: []float64{1, 4, 2, 5}, Variant: Dotted}}}, PixelOptions{Width: 40, Height: 20})
		if err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
		if c.Width != 40 || c.Height != 20 || len(c.ANSI()) == 0 {
			t.Fatalf("%s: empty or incorrectly sized raster", kind)
		}
	}
}
