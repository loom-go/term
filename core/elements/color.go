package elements

import (
	"fmt"
	"strings"

	"github.com/AnatoleLucet/go-opentui"
)

var namedColors = map[string]opentui.RGBA{
	"":            opentui.Transparent,
	"transparent": opentui.Transparent,
	"none":        opentui.Transparent,

	"black": {R: 0, G: 0, B: 0, A: 1},
	"white": {R: 1, G: 1, B: 1, A: 1},

	"lightgray": {R: 0.75, G: 0.75, B: 0.75, A: 1},
	"gray":      {R: 0.5, G: 0.5, B: 0.5, A: 1},
	"darkgray":  {R: 0.25, G: 0.25, B: 0.25, A: 1},

	"lightred": {R: 1, G: 0.5, B: 0.5, A: 1},
	"red":      {R: 1, G: 0, B: 0, A: 1},
	"darkred":  {R: 0.5, G: 0, B: 0, A: 1},

	"lightgreen": {R: 0.5, G: 1, B: 0.5, A: 1},
	"green":      {R: 0, G: 1, B: 0, A: 1},
	"darkgreen":  {R: 0, G: 0.5, B: 0, A: 1},

	"lightblue": {R: 0.5, G: 0.5, B: 1, A: 1},
	"blue":      {R: 0, G: 0, B: 1, A: 1},
	"darkblue":  {R: 0, G: 0, B: 0.5, A: 1},

	"lightyellow": {R: 1, G: 1, B: 0.5, A: 1},
	"yellow":      {R: 1, G: 1, B: 0, A: 1},
	"darkyellow":  {R: 0.5, G: 0.5, B: 0, A: 1},

	"lightcyan": {R: 0.5, G: 1, B: 1, A: 1},
	"cyan":      {R: 0, G: 1, B: 1, A: 1},
	"darkcyan":  {R: 0, G: 0.5, B: 0.5, A: 1},

	"lightmagenta": {R: 1, G: 0.5, B: 1, A: 1},
	"magenta":      {R: 1, G: 0, B: 1, A: 1},
	"darkmagenta":  {R: 0.5, G: 0, B: 0.5, A: 1},
}

type Color struct {
	cached string
	color  opentui.RGBA
}

func NewColor(color string) *Color {
	c := &Color{}
	c.Set(color)
	return c
}

func (c *Color) Set(color string) error {
	if c.cached == color {
		return nil
	}

	rgba, err := toOpenTUIColor(color)
	if err != nil {
		return err
	}

	c.cached = color
	c.color = rgba
	return nil
}

func (c Color) RGBA() opentui.RGBA {
	return c.color
}

func toOpenTUIColor(color string, fallback ...string) (rgba opentui.RGBA, err error) {
	color = strings.ReplaceAll(strings.ToLower(color), " ", "")

	if color == "" && len(fallback) > 0 {
		color = strings.ReplaceAll(strings.ToLower(fallback[0]), " ", "")
	}

	if named, ok := namedColors[color]; ok {
		return named, nil
	}

	if strings.HasPrefix(color, "#") {
		return praseHexColor(color)
	}

	if strings.HasPrefix(color, "rgb(") || strings.HasPrefix(color, "rgba(") {
		return parseRGBColor(color)
	}

	return rgba, fmt.Errorf("invalid color format '%s'", color)
}

func praseHexColor(color string) (rgba opentui.RGBA, err error) {
	var r, g, b, a uint8

	switch len(color) {
	case 9: // #RRGGBBAA
		_, err = fmt.Sscanf(color, "#%02x%02x%02x%02x", &r, &g, &b, &a)
		if err == nil {
			return opentui.RGBA{
				R: float32(r) / 255,
				G: float32(g) / 255,
				B: float32(b) / 255,
				A: float32(a) / 255,
			}, nil
		}

	case 7: // #RRGGBB
		_, err = fmt.Sscanf(color, "#%02x%02x%02x", &r, &g, &b)
		if err == nil {
			return opentui.RGBA{
				R: float32(r) / 255,
				G: float32(g) / 255,
				B: float32(b) / 255,
				A: 1,
			}, nil
		}

	case 5: // #RGBA
		var rr, gg, bb, aa uint8
		_, err = fmt.Sscanf(color, "#%1x%1x%1x%1x", &rr, &gg, &bb, &aa)
		if err == nil {
			return opentui.RGBA{
				R: float32(rr*17) / 255,
				G: float32(gg*17) / 255,
				B: float32(bb*17) / 255,
				A: float32(aa*17) / 255,
			}, nil
		}

	case 4: // #RGB
		var rr, gg, bb uint8
		_, err = fmt.Sscanf(color, "#%1x%1x%1x", &rr, &gg, &bb)
		if err == nil {
			return opentui.RGBA{
				R: float32(rr*17) / 255,
				G: float32(gg*17) / 255,
				B: float32(bb*17) / 255,
				A: 1,
			}, nil
		}
	}

	return rgba, fmt.Errorf("invalid hex color format '%s'", color)
}

// todo: that's brittle. 'rgb(0,0,0,0)' 'rgb(0.0,0,0)' etc all works but shouldn't
func parseRGBColor(color string) (rgba opentui.RGBA, err error) {
	var r, g, b, a float32
	var n int

	// rgba()
	n, _ = fmt.Sscanf(color, "rgba(%f,%f,%f,%f)", &r, &g, &b, &a)
	if n == 4 {
		if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 || a < 0 || a > 1 {
			return opentui.RGBA{}, fmt.Errorf("rgba values out of range")
		}

		return opentui.RGBA{
			R: r / 255,
			G: g / 255,
			B: b / 255,
			A: a,
		}, nil
	}

	// rgb()
	n, _ = fmt.Sscanf(color, "rgb(%f,%f,%f)", &r, &g, &b)
	if n == 3 {
		if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
			return opentui.RGBA{}, fmt.Errorf("rgb values out of range")
		}
		return opentui.RGBA{
			R: r / 255,
			G: g / 255,
			B: b / 255,
			A: 1,
		}, nil
	}

	return rgba, fmt.Errorf("invalid rgb/rgba color format '%s'", color)
}
