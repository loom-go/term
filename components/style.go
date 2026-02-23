package components

import (
	"fmt"
	"math"
)

type Style struct {
	Width     any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	MinWidth  any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	MaxWidth  any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	Height    any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	MinHeight any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	MaxHeight any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"

	Top    any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	Bottom any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	Left   any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"
	Right  any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch"

	ZIndex int // relative to sibling elements

	Position string // "static" | "relative" | "absolute"

	PaddingAll        any // 10 | "100pt" | "50%"
	PaddingVertical   any // 10 | "100pt" | "50%"
	PaddingHorizontal any // 10 | "100pt" | "50%"
	PaddingTop        any // 10 | "100pt" | "50%"
	PaddingBottom     any // 10 | "100pt" | "50%"
	PaddingLeft       any // 10 | "100pt" | "50%"
	PaddingRight      any // 10 | "100pt" | "50%"

	MarginAll        any // 10 | "100pt" | "50%"
	MarginVertical   any // 10 | "100pt" | "50%"
	MarginHorizontal any // 10 | "100pt" | "50%"
	MarginTop        any // 10 | "100pt" | "50%"
	MarginBottom     any // 10 | "100pt" | "50%"
	MarginLeft       any // 10 | "100pt" | "50%"
	MarginRight      any // 10 | "100pt" | "50%"

	BorderAll        string // "single" | "double" | "rounded"  | "heavy"
	BorderVertical   string // "single" | "double" | "rounded"  | "heavy"
	BorderHorizontal string // "single" | "double" | "rounded"  | "heavy"
	BorderTop        string // "single" | "double" | "rounded"  | "heavy"
	BorderBottom     string // "single" | "double" | "rounded"  | "heavy"
	BorderLeft       string // "single" | "double" | "rounded"  | "heavy"
	BorderRight      string // "single" | "double" | "rounded"  | "heavy"

	Display string // "none" | "flex" | "contents"

	AlignSelf      string // "start" | "end" | "center" | "stretch" | "baseline"
	AlignItems     string // "start" | "end" | "center" | "stretch" | "baseline"
	AlignContent   string // "start" | "end" | "center" | "stretch" | "baseline"
	JustifyContent string // "start" | "end" | "center" | "space-between" | "space-around" | "space-evenly"
	FlexDirection  string // "row" | "row-reverse" | "column" | "column-reverse"
	FlexWrap       string // "nowrap" | "wrap" | "wrap-reverse"
	FlexGrow       string // "none" | "0" | "1" | ...
	FlexShrink     string // "none" | "0" | "1" | ...

	GapAll    any // 10 | "100pt" | "50%"
	GapRow    any // 10 | "100pt" | "50%"
	GapColumn any // 10 | "100pt" | "50%"

	Overflow string // "visible" | "hidden"

	BorderColor string // "transparent" | "#RGB" | "#RRGGBBAA"

	BackgroundColor string // "transparent" | "#RGB" | "#RRGGBBAA"

	Color          string // "transparent" | "#RGB" | "#RRGGBBAA"
	DropColor      string // "transparent" | "#RGB" | "#RRGGBBAA"
	FontWeight     string // "normal" | "bold"
	FontStyle      string // "normal" | "italic"
	TextDecoration string // "none" | "underline" | "line-through"
	TextWrap       string // "none" | "word" | "char"

}

// helper for "<value>pt"
func Point(value int) string {
	return fmt.Sprintf("%dpt", value)
}

// helper for "<value>%"
func Percent(value int) string {
	return fmt.Sprintf("%d%%", value)
}

// helper for "#<rr><gg><bb>"
func RGB(r, g, b uint8) string {
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// helper for "#<rr><gg><bb><aa>"
func RGBA(r, g, b uint8, alpha float64) string {
	return fmt.Sprintf("#%02x%02x%02x%02x", r, g, b, uint8(math.Round(float64(alpha)*255)))
}

// helper for HSL to "#<rr><gg><bb>"
func HSL(h, s, l float64) string {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := l - c/2

	var r, g, b float64

	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	R := uint8(math.Round((r + m) * 255))
	G := uint8(math.Round((g + m) * 255))
	B := uint8(math.Round((b + m) * 255))

	return fmt.Sprintf("#%02x%02x%02x", R, G, B)
}

// helper for HSLA to "#<rr><gg><bb><aa>"
func HSLA(h, s, v, alpha float64) string {
	hsl := HSL(h, s, v)

	return fmt.Sprintf("%s%02x", hsl, uint8(math.Round(float64(alpha)*255)))
}
