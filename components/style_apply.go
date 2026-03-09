package components

import (
	"reflect"

	"github.com/loom-go/loom/signals"
	"github.com/loom-go/term/core"
)

func applyStyle(elem core.Element, style Style) (removers []func()) {
	removers = append(removers, applyStyleDimension(elem, style)...)
	removers = append(removers, applyStylePositionning(elem, style)...)
	removers = append(removers, applyStyleSpacing(elem, style)...)
	removers = append(removers, applyStyleBorder(elem, style)...)
	removers = append(removers, applyStyleDisplay(elem, style)...)
	removers = append(removers, applyStyleFlexbox(elem, style)...)
	removers = append(removers, applyStyleOverflow(elem, style)...)
	removers = append(removers, applyStyleBackground(elem, style)...)
	removers = append(removers, applyStyleText(elem, style)...)
	removers = append(removers, applyStylePlaceholder(elem, style)...)

	return
}

func applyStyleDimension(elem core.Element, style Style) (removers []func()) {
	if e, v, ok := matchMethod[interface {
		SetWidth(any)
		UnsetWidth()
	}, any](elem, style.Width); ok {
		e.SetWidth(v)
		removers = append(removers, e.UnsetWidth)
	}
	if e, v, ok := matchMethod[interface {
		SetMinWidth(any)
		UnsetMinWidth()
	}, any](elem, style.MinWidth); ok {
		e.SetMinWidth(v)
		removers = append(removers, e.UnsetMinWidth)
	}
	if e, v, ok := matchMethod[interface {
		SetMaxWidth(any)
		UnsetMaxWidth()
	}, any](elem, style.MaxWidth); ok {
		e.SetMaxWidth(v)
		removers = append(removers, e.UnsetMaxWidth)
	}

	if e, v, ok := matchMethod[interface {
		SetHeight(any)
		UnsetHeight()
	}, any](elem, style.Height); ok {
		e.SetHeight(v)
		removers = append(removers, e.UnsetHeight)
	}
	if e, v, ok := matchMethod[interface {
		SetMinHeight(any)
		UnsetMinHeight()
	}, any](elem, style.MinHeight); ok {
		e.SetMinHeight(v)
		removers = append(removers, e.UnsetMinHeight)
	}
	if e, v, ok := matchMethod[interface {
		SetMaxHeight(any)
		UnsetMaxHeight()
	}, any](elem, style.MaxHeight); ok {
		e.SetMaxHeight(v)
		removers = append(removers, e.UnsetMaxHeight)
	}

	return
}

func applyStylePositionning(elem core.Element, style Style) (removers []func()) {
	if e, v, ok := matchMethod[interface {
		SetTop(any)
		UnsetTop()
	}, any](elem, style.Top); ok {
		e.SetTop(v)
		removers = append(removers, e.UnsetTop)
	}
	if e, v, ok := matchMethod[interface {
		SetBottom(any)
		UnsetBottom()
	}, any](elem, style.Bottom); ok {
		e.SetBottom(v)
		removers = append(removers, e.UnsetBottom)
	}
	if e, v, ok := matchMethod[interface {
		SetLeft(any)
		UnsetLeft()
	}, any](elem, style.Left); ok {
		e.SetLeft(v)
		removers = append(removers, e.UnsetLeft)
	}
	if e, v, ok := matchMethod[interface {
		SetRight(any)
		UnsetRight()
	}, any](elem, style.Right); ok {
		e.SetRight(v)
		removers = append(removers, e.UnsetRight)
	}

	if e, v, ok := matchMethod[interface {
		SetZIndex(int)
		UnsetZIndex()
	}, int](elem, style.ZIndex); ok {
		e.SetZIndex(v)
		removers = append(removers, e.UnsetZIndex)
	}

	if e, v, ok := matchMethod[interface {
		SetPosition(string)
		UnsetPosition()
	}, string](elem, style.Position); ok && v != "" {
		e.SetPosition(v)
		removers = append(removers, e.UnsetPosition)
	}

	return
}

func applyStyleSpacing(elem core.Element, style Style) (removers []func()) {
	if e, v, ok := matchMethod[interface {
		SetPaddingAll(any)
		UnsetPaddingAll()
	}, any](elem, style.PaddingAll); ok {
		e.SetPaddingAll(v)
		removers = append(removers, e.UnsetPaddingAll)
	}
	if e, v, ok := matchMethod[interface {
		SetPaddingVertical(any)
		UnsetPaddingVertical()
	}, any](elem, style.PaddingVertical); ok {
		e.SetPaddingVertical(v)
		removers = append(removers, e.UnsetPaddingVertical)
	}
	if e, v, ok := matchMethod[interface {
		SetPaddingHorizontal(any)
		UnsetPaddingHorizontal()
	}, any](elem, style.PaddingHorizontal); ok {
		e.SetPaddingHorizontal(v)
		removers = append(removers, e.UnsetPaddingHorizontal)
	}
	if e, v, ok := matchMethod[interface {
		SetPaddingTop(any)
		UnsetPaddingTop()
	}, any](elem, style.PaddingTop); ok {
		e.SetPaddingTop(v)
		removers = append(removers, e.UnsetPaddingTop)
	}
	if e, v, ok := matchMethod[interface {
		SetPaddingBottom(any)
		UnsetPaddingBottom()
	}, any](elem, style.PaddingBottom); ok {
		e.SetPaddingBottom(v)
		removers = append(removers, e.UnsetPaddingBottom)
	}
	if e, v, ok := matchMethod[interface {
		SetPaddingLeft(any)
		UnsetPaddingLeft()
	}, any](elem, style.PaddingLeft); ok {
		e.SetPaddingLeft(v)
		removers = append(removers, e.UnsetPaddingLeft)
	}
	if e, v, ok := matchMethod[interface {
		SetPaddingRight(any)
		UnsetPaddingRight()
	}, any](elem, style.PaddingRight); ok {
		e.SetPaddingRight(v)
		removers = append(removers, e.UnsetPaddingRight)
	}

	if e, v, ok := matchMethod[interface {
		SetMarginAll(any)
		UnsetMarginAll()
	}, any](elem, style.MarginAll); ok {
		e.SetMarginAll(v)
		removers = append(removers, e.UnsetMarginAll)
	}
	if e, v, ok := matchMethod[interface {
		SetMarginVertical(any)
		UnsetMarginVertical()
	}, any](elem, style.MarginVertical); ok {
		e.SetMarginVertical(v)
		removers = append(removers, e.UnsetMarginVertical)
	}
	if e, v, ok := matchMethod[interface {
		SetMarginHorizontal(any)
		UnsetMarginHorizontal()
	}, any](elem, style.MarginHorizontal); ok {
		e.SetMarginHorizontal(v)
		removers = append(removers, e.UnsetMarginHorizontal)
	}
	if e, v, ok := matchMethod[interface {
		SetMarginTop(any)
		UnsetMarginTop()
	}, any](elem, style.MarginTop); ok {
		e.SetMarginTop(v)
		removers = append(removers, e.UnsetMarginTop)
	}
	if e, v, ok := matchMethod[interface {
		SetMarginBottom(any)
		UnsetMarginBottom()
	}, any](elem, style.MarginBottom); ok {
		e.SetMarginBottom(v)
		removers = append(removers, e.UnsetMarginBottom)
	}
	if e, v, ok := matchMethod[interface {
		SetMarginLeft(any)
		UnsetMarginLeft()
	}, any](elem, style.MarginLeft); ok {
		e.SetMarginLeft(v)
		removers = append(removers, e.UnsetMarginLeft)
	}
	if e, v, ok := matchMethod[interface {
		SetMarginRight(any)
		UnsetMarginRight()
	}, any](elem, style.MarginRight); ok {
		e.SetMarginRight(v)
		removers = append(removers, e.UnsetMarginRight)
	}

	if e, v, ok := matchMethod[interface {
		SetGapAll(any)
		UnsetGapAll()
	}, any](elem, style.GapAll); ok {
		e.SetGapAll(v)
		removers = append(removers, e.UnsetGapAll)
	}
	if e, v, ok := matchMethod[interface {
		SetGapRow(any)
		UnsetGapRow()
	}, any](elem, style.GapRow); ok {
		e.SetGapRow(v)
		removers = append(removers, e.UnsetGapRow)
	}
	if e, v, ok := matchMethod[interface {
		SetGapColumn(any)
		UnsetGapColumn()
	}, any](elem, style.GapColumn); ok {
		e.SetGapColumn(v)
		removers = append(removers, e.UnsetGapColumn)
	}

	return
}

func applyStyleBorder(elem core.Element, style Style) (removers []func()) {
	if e, v, ok := matchMethod[interface {
		SetBorderAll(string)
		UnsetBorderAll()
	}, string](elem, style.BorderAll); ok && v != "" {
		e.SetBorderAll(v)
		removers = append(removers, e.UnsetBorderAll)
	}
	if e, v, ok := matchMethod[interface {
		SetBorderVertical(string)
		UnsetBorderVertical()
	}, string](elem, style.BorderVertical); ok && v != "" {
		e.SetBorderVertical(v)
		removers = append(removers, e.UnsetBorderVertical)
	}
	if e, v, ok := matchMethod[interface {
		SetBorderHorizontal(string)
		UnsetBorderHorizontal()
	}, string](elem, style.BorderHorizontal); ok && v != "" {
		e.SetBorderHorizontal(v)
		removers = append(removers, e.UnsetBorderHorizontal)
	}
	if e, v, ok := matchMethod[interface {
		SetBorderTop(string)
		UnsetBorderTop()
	}, string](elem, style.BorderTop); ok && v != "" {
		e.SetBorderTop(v)
		removers = append(removers, e.UnsetBorderTop)
	}
	if e, v, ok := matchMethod[interface {
		SetBorderBottom(string)
		UnsetBorderBottom()
	}, string](elem, style.BorderBottom); ok && v != "" {
		e.SetBorderBottom(v)
		removers = append(removers, e.UnsetBorderBottom)
	}
	if e, v, ok := matchMethod[interface {
		SetBorderLeft(string)
		UnsetBorderLeft()
	}, string](elem, style.BorderLeft); ok && v != "" {
		e.SetBorderLeft(v)
		removers = append(removers, e.UnsetBorderLeft)
	}
	if e, v, ok := matchMethod[interface {
		SetBorderRight(string)
		UnsetBorderRight()
	}, string](elem, style.BorderRight); ok && v != "" {
		e.SetBorderRight(v)
		removers = append(removers, e.UnsetBorderRight)
	}

	if e, v, ok := matchMethod[interface {
		SetBorderColor(string)
		UnsetBorderColor()
	}, string](elem, style.BorderColor); ok && v != "" {
		e.SetBorderColor(v)
		removers = append(removers, e.UnsetBorderColor)
	}

	return
}

func applyStyleDisplay(elem core.Element, style Style) (removers []func()) {
	if e, v, ok := matchMethod[interface {
		SetDisplay(string)
		UnsetDisplay()
	}, string](elem, style.Display); ok && v != "" {
		e.SetDisplay(v)
		removers = append(removers, e.UnsetDisplay)
	}

	return
}

func applyStyleFlexbox(elem core.Element, style Style) (removers []func()) {
	if e, v, ok := matchMethod[interface {
		SetAlignSelf(string)
		UnsetAlignSelf()
	}, string](elem, style.AlignSelf); ok && v != "" {
		e.SetAlignSelf(v)
		removers = append(removers, e.UnsetAlignSelf)
	}
	if e, v, ok := matchMethod[interface {
		SetAlignItems(string)
		UnsetAlignItems()
	}, string](elem, style.AlignItems); ok && v != "" {
		e.SetAlignItems(v)
		removers = append(removers, e.UnsetAlignItems)
	}
	if e, v, ok := matchMethod[interface {
		SetAlignContent(string)
		UnsetAlignContent()
	}, string](elem, style.AlignContent); ok && v != "" {
		e.SetAlignContent(v)
		removers = append(removers, e.UnsetAlignContent)
	}

	if e, v, ok := matchMethod[interface {
		SetJustifyContent(string)
		UnsetJustifyContent()
	}, string](elem, style.JustifyContent); ok && v != "" {
		e.SetJustifyContent(v)
		removers = append(removers, e.UnsetJustifyContent)
	}

	if e, v, ok := matchMethod[interface {
		SetFlexDirection(string)
		UnsetFlexDirection()
	}, string](elem, style.FlexDirection); ok && v != "" {
		e.SetFlexDirection(v)
		removers = append(removers, e.UnsetFlexDirection)
	}

	if e, v, ok := matchMethod[interface {
		SetFlexWrap(string)
		UnsetFlexWrap()
	}, string](elem, style.FlexWrap); ok && v != "" {
		e.SetFlexWrap(v)
		removers = append(removers, e.UnsetFlexWrap)
	}

	if e, v, ok := matchMethod[interface {
		SetFlexGrow(string)
		UnsetFlexGrow()
	}, string](elem, style.FlexGrow); ok && v != "" {
		e.SetFlexGrow(v)
		removers = append(removers, e.UnsetFlexGrow)
	}
	if e, v, ok := matchMethod[interface {
		SetFlexShrink(string)
		UnsetFlexShrink()
	}, string](elem, style.FlexShrink); ok && v != "" {
		e.SetFlexShrink(v)
		removers = append(removers, e.UnsetFlexShrink)
	}

	return
}

func applyStyleOverflow(elem core.Element, style Style) (removers []func()) {
	if e, v, ok := matchMethod[interface {
		SetOverflow(string)
		UnsetOverflow()
	}, string](elem, style.Overflow); ok && v != "" {
		e.SetOverflow(v)
		removers = append(removers, e.UnsetOverflow)
	}

	return
}

func applyStyleBackground(elem core.Element, style Style) (removers []func()) {
	if style.BackgroundColor != nil {
		if e, v, ok := matchMethod[interface {
			SetBackgroundColor(string)
			UnsetBackgroundColor()
		}, string](elem, style.BackgroundColor); ok && v != "" {
			e.SetBackgroundColor(v)
			removers = append(removers, e.UnsetBackgroundColor)
		} else if e, v, ok := matchMethod[interface {
			SetTextBackground(string)
			UnsetTextBackground()
		}, string](elem, style.BackgroundColor); ok && v != "" && style.DropColor == nil {
			e.SetTextBackground(v)
			removers = append(removers, e.UnsetTextBackground)
		}
	}

	if style.BackgroundOpacity != nil {
		if e, v, ok := matchMethod[interface {
			SetBackgroundOpacity(float32)
			UnsetBackgroundOpacity()
		}, float32](elem, style.BackgroundOpacity); ok {
			e.SetBackgroundOpacity(v)
			removers = append(removers, e.UnsetBackgroundOpacity)
		}
	}

	return
}

func applyStyleText(elem core.Element, style Style) (removers []func()) {
	if e, v, ok := matchMethod[interface {
		SetTextForeground(string)
		UnsetTextForeground()
	}, string](elem, style.Color); ok && v != "" {
		e.SetTextForeground(v)
		removers = append(removers, e.UnsetTextForeground)
	}
	if e, v, ok := matchMethod[interface {
		SetTextBackground(string)
		UnsetTextBackground()
	}, string](elem, style.DropColor); ok && v != "" {
		e.SetTextBackground(v)
		removers = append(removers, e.UnsetTextBackground)
	}

	if e, v, ok := matchMethod[interface {
		SetFontWeight(string)
		UnsetFontWeight()
	}, string](elem, style.FontWeight); ok && v != "" {
		e.SetFontWeight(v)
		removers = append(removers, e.UnsetFontWeight)
	}
	if e, v, ok := matchMethod[interface {
		SetFontStyle(string)
		UnsetFontStyle()
	}, string](elem, style.FontStyle); ok && v != "" {
		e.SetFontStyle(v)
		removers = append(removers, e.UnsetFontStyle)
	}

	if e, v, ok := matchMethod[interface {
		SetTextDecoration(string)
		UnsetTextDecoration()
	}, string](elem, style.TextDecoration); ok && v != "" {
		e.SetTextDecoration(v)
		removers = append(removers, e.UnsetTextDecoration)
	}
	if e, v, ok := matchMethod[interface {
		SetWrap(string)
		UnsetWrap()
	}, string](elem, style.TextWrap); ok && v != "" {
		e.SetWrap(v)
		removers = append(removers, e.UnsetWrap)
	}

	return
}

func applyStylePlaceholder(elem core.Element, style Style) (removers []func()) {
	if e, v, ok := matchMethod[interface {
		SetPlaceholderForeground(string)
		UnsetPlaceholderForeground()
	}, string](elem, style.PlaceholderColor); ok && v != "" {
		e.SetPlaceholderForeground(v)
		removers = append(removers, e.UnsetPlaceholderForeground)
	}
	if e, v, ok := matchMethod[interface {
		SetPlaceholderBackground(string)
		UnsetPlaceholderBackground()
	}, string](elem, style.PlaceholderDropColor); ok && v != "" {
		e.SetPlaceholderBackground(v)
		removers = append(removers, e.UnsetPlaceholderBackground)
	}

	if e, v, ok := matchMethod[interface {
		SetPlaceholderFontWeight(string)
		UnsetPlaceholderFontWeight()
	}, string](elem, style.PlaceholderFontWeight); ok && v != "" {
		e.SetPlaceholderFontWeight(v)
		removers = append(removers, e.UnsetPlaceholderFontWeight)
	}
	if e, v, ok := matchMethod[interface {
		SetPlaceholderFontStyle(string)
		UnsetPlaceholderFontStyle()
	}, string](elem, style.PlaceholderFontStyle); ok && v != "" {
		e.SetPlaceholderFontStyle(v)
		removers = append(removers, e.UnsetPlaceholderFontStyle)
	}

	if e, v, ok := matchMethod[interface {
		SetPlaceholderDecoration(string)
		UnsetPlaceholderDecoration()
	}, string](elem, style.PlaceholderDecoration); ok && v != "" {
		e.SetPlaceholderDecoration(v)
		removers = append(removers, e.UnsetPlaceholderDecoration)
	}

	return
}

func matchMethod[E, V any](elem any, value any) (E, V, bool) {
	var v V
	var vok bool
	var e E
	var eok bool

	v, vok = unwrapAccessor[V](value)
	if !vok {
		return e, v, false
	}

	e, eok = elem.(E)
	if !eok {
		return e, v, false
	}

	return e, v, true
}

// little helper to unwrap a `T | func() T`. mainly used in Appliers
func unwrapAccessor[V any](value any) (v V, vok bool) {
	if sig, ok := value.(func() V); ok {
		value = sig()
	} else if fn, ok := value.(func() any); ok {
		value = fn()
	} else if fn, ok := value.(signals.Accessor[V]); ok {
		value = fn()
	} else if fn, ok := value.(signals.Accessor[any]); ok {
		value = fn()
	} else {
		// fallback to reflect for calling the accessor
		rv := reflect.ValueOf(value)
		isFunc := rv.Kind() == reflect.Func
		if isFunc && rv.Type().NumIn() == 0 && rv.Type().NumOut() == 1 {
			value = rv.Call(nil)[0].Interface()
		}
	}

	if val, ok := value.(V); ok {
		v = val
		vok = true
	} else {
		// fallback to reflect for converting the value
		rv := reflect.ValueOf(value)
		target := reflect.TypeFor[V]()
		if rv.IsValid() && rv.Type().ConvertibleTo(target) {
			v = rv.Convert(target).Interface().(V)
			vok = true
		}
	}

	return
}
