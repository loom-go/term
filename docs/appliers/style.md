---
weight: 1
title: Style{}
---

```go {style=tokyonight-moon}
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
```
