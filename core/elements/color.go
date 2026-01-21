package elements

import (
	"fmt"
	"strings"

	"github.com/AnatoleLucet/go-opentui"
)

type Color struct {
	cached string
	color  opentui.RGBA
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
	color = strings.TrimSpace(strings.ToLower(color))

	if color == "" && len(fallback) > 0 {
		color = strings.TrimSpace(strings.ToLower(fallback[0]))
	}

	if color == "" || color == "transparent" || color == "none" {
		return opentui.Transparent, nil
	}

	var r, g, b, a uint8

	// todo: refactor and add support for `rgba(...)` and `rgb(...)`
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

	return rgba, fmt.Errorf("invalid color format '%s'", color)
}
