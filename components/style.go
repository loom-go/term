package components

import (
	"fmt"
	"math"

	"github.com/AnatoleLucet/loom-term/core"
)

// Style applier for elements.
// Values can be of various types and units, each field has a comment for each possible value type.
//
// For reactive style, properties can be defined as a function that returns the value.
type Style struct {
	Width     any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch" | func() any
	MinWidth  any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch" | func() any
	MaxWidth  any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch" | func() any
	Height    any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch" | func() any
	MinHeight any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch" | func() any
	MaxHeight any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch" | func() any

	Top    any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch" | func() any
	Bottom any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch" | func() any
	Left   any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch" | func() any
	Right  any // 10 | "100pt" | "50%" | "auto" | "max-content" | "fit-content" | "stretch" | func() any

	ZIndex any // relative to sibling elements

	Position any // "static" | "relative" | "absolute" | func() any

	PaddingAll        any // 10 | "100pt" | "50%" | func() any
	PaddingVertical   any // 10 | "100pt" | "50%" | func() any
	PaddingHorizontal any // 10 | "100pt" | "50%" | func() any
	PaddingTop        any // 10 | "100pt" | "50%" | func() any
	PaddingBottom     any // 10 | "100pt" | "50%" | func() any
	PaddingLeft       any // 10 | "100pt" | "50%" | func() any
	PaddingRight      any // 10 | "100pt" | "50%" | func() any

	MarginAll        any // 10 | "100pt" | "50%" | func() any
	MarginVertical   any // 10 | "100pt" | "50%" | func() any
	MarginHorizontal any // 10 | "100pt" | "50%" | func() any
	MarginTop        any // 10 | "100pt" | "50%" | func() any
	MarginBottom     any // 10 | "100pt" | "50%" | func() any
	MarginLeft       any // 10 | "100pt" | "50%" | func() any
	MarginRight      any // 10 | "100pt" | "50%" | func() any

	BorderAll        any // "single" | "double" | "rounded"  | "heavy" | func() any
	BorderVertical   any // "single" | "double" | "rounded"  | "heavy" | func() any
	BorderHorizontal any // "single" | "double" | "rounded"  | "heavy" | func() any
	BorderTop        any // "single" | "double" | "rounded"  | "heavy" | func() any
	BorderBottom     any // "single" | "double" | "rounded"  | "heavy" | func() any
	BorderLeft       any // "single" | "double" | "rounded"  | "heavy" | func() any
	BorderRight      any // "single" | "double" | "rounded"  | "heavy" | func() any

	Display any // "none" | "flex" | "contents" | func() any

	AlignSelf      any // "start" | "end" | "center" | "stretch" | "baseline" | func() any
	AlignItems     any // "start" | "end" | "center" | "stretch" | "baseline" | func() any
	AlignContent   any // "start" | "end" | "center" | "stretch" | "baseline" | func() any
	JustifyContent any // "start" | "end" | "center" | "space-between" | "space-around" | "space-evenly" | func() any
	FlexDirection  any // "row" | "row-reverse" | "column" | "column-reverse" | func() any
	FlexWrap       any // "nowrap" | "wrap" | "wrap-reverse" | func() any
	FlexGrow       any // "none" | "0" | "1" | ... | func() any
	FlexShrink     any // "none" | "0" | "1" | ... | func() any

	GapAll    any // 10 | "100pt" | "50%" | func() any
	GapRow    any // 10 | "100pt" | "50%" | func() any
	GapColumn any // 10 | "100pt" | "50%" | func() any

	Overflow any // "visible" | "hidden" | func() any

	BorderColor any // "transparent" | "#RGB" | "#RRGGBBAA" | func() any

	BackgroundColor   any // "transparent" | "#RGB" | "#RRGGBBAA" | func() any
	BackgroundOpacity any // 0.0 - 1.0 | func() any

	Color          any // "transparent" | "#RGB" | "#RRGGBBAA" | func() any
	DropColor      any // "transparent" | "#RGB" | "#RRGGBBAA" | func() any
	FontWeight     any // "normal" | "bold" | func() any
	FontStyle      any // "normal" | "italic" | func() any
	TextDecoration any // "none" | "underline" | "line-through" | func() any
	TextWrap       any // "none" | "word" | "char" | func() any

	PlaceholderColor      any // "transparent" | "#RGB" | "#RRGGBBAA" | func() any
	PlaceholderDropColor  any // "transparent" | "#RGB" | "#RRGGBBAA" | func() any
	PlaceholderFontWeight any // "normal" | "bold" | func() any
	PlaceholderFontStyle  any // "normal" | "italic" | func() any
	PlaceholderDecoration any // "none" | "underline" | "line-through" | func() any
}

func (s Style) Apply(parent any) (func() error, error) {
	elem, ok := parent.(core.Element)
	if !ok {
		return nil, fmt.Errorf("Style: parent node is not an Element")
	}

	removers := applyStyle(elem, s)
	remove := func() error {
		for _, r := range removers {
			r()
		}
		return nil
	}

	return remove, nil
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

	s = s / 100
	l = l / 100

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
