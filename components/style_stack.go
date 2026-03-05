package components

import (
	"reflect"
	"slices"

	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom/signals"
)

// todo: that's not concurency safe. prob will need to be

type styleLayer struct {
	id      uint32
	event   string
	styles  []Style
	visible bool
}

type styleStack struct {
	layers []*styleLayer
}

func (s *styleStack) Push(layer *styleLayer) {
	s.layers = append(s.layers, layer)
}

func (s *styleStack) Replace(id uint32, styles []Style) {
	i := slices.IndexFunc(s.layers, func(layer *styleLayer) bool {
		return layer.id == id
	})
	if i == -1 {
		return
	}

	// remove old layer and add new one at the end to update its priority
	s.layers = append(slices.Delete(s.layers, i, i+1), &styleLayer{
		id:      id,
		styles:  styles,
		visible: true,
	})
}

func (s *styleStack) Pop(id uint32) {
	s.layers = slices.DeleteFunc(s.layers, func(layer *styleLayer) bool {
		return layer.id == id
	})
}

var styleStacks = map[core.Element]*styleStack{}

func getStyleStack(elem core.Element) *styleStack {
	stack, exists := styleStacks[elem]
	if !exists {
		stack = &styleStack{}
		styleStacks[elem] = stack
	}

	return stack
}

func applyStyle(elem core.Element, style Style) {
	applyStyleDimension(elem, style)
	applyStylePositionning(elem, style)
	applyStyleSpacing(elem, style)
	applyStyleBorder(elem, style)
	applyStyleDisplay(elem, style)
	applyStyleFlexbox(elem, style)
	applyStyleOverflow(elem, style)
	applyStyleBackground(elem, style)
	applyStyleText(elem, style)
	applyStylePlaceholder(elem, style)
}

func applyStyleDimension(elem core.Element, style Style) {
	if e, v, ok := assertSetter[interface{ SetWidth(any) }, any](elem, style.Width); ok {
		e.SetWidth(v)
	}
	if e, v, ok := assertSetter[interface{ SetMinWidth(any) }, any](elem, style.MinWidth); ok {
		e.SetMinWidth(v)
	}
	if e, v, ok := assertSetter[interface{ SetMaxWidth(any) }, any](elem, style.MaxWidth); ok {
		e.SetMaxWidth(v)
	}

	if e, v, ok := assertSetter[interface{ SetHeight(any) }, any](elem, style.Height); ok {
		e.SetHeight(v)
	}
	if e, v, ok := assertSetter[interface{ SetMinHeight(any) }, any](elem, style.MinHeight); ok {
		e.SetMinHeight(v)
	}
	if e, v, ok := assertSetter[interface{ SetMaxHeight(any) }, any](elem, style.MaxHeight); ok {
		e.SetMaxHeight(v)
	}
}

func applyStylePositionning(elem core.Element, style Style) {
	if e, v, ok := assertSetter[interface{ SetTop(any) }, any](elem, style.Top); ok {
		e.SetTop(v)
	}
	if e, v, ok := assertSetter[interface{ SetBottom(any) }, any](elem, style.Bottom); ok {
		e.SetBottom(v)
	}
	if e, v, ok := assertSetter[interface{ SetLeft(any) }, any](elem, style.Left); ok {
		e.SetLeft(v)
	}
	if e, v, ok := assertSetter[interface{ SetRight(any) }, any](elem, style.Right); ok {
		e.SetRight(v)
	}

	if e, v, ok := assertSetter[interface{ SetZIndex(int) }, int](elem, style.ZIndex); ok {
		e.SetZIndex(v)
	}

	if e, v, ok := assertSetter[interface{ SetPosition(string) }, string](elem, style.Position); ok && v != "" {
		e.SetPosition(v)
	}
}

func applyStyleSpacing(elem core.Element, style Style) {
	if e, v, ok := assertSetter[interface{ SetPaddingAll(any) }, any](elem, style.PaddingAll); ok {
		e.SetPaddingAll(v)
	}
	if e, v, ok := assertSetter[interface{ SetPaddingVertical(any) }, any](elem, style.PaddingVertical); ok {
		e.SetPaddingVertical(v)
	}
	if e, v, ok := assertSetter[interface{ SetPaddingHorizontal(any) }, any](elem, style.PaddingHorizontal); ok {
		e.SetPaddingHorizontal(v)
	}
	if e, v, ok := assertSetter[interface{ SetPaddingTop(any) }, any](elem, style.PaddingTop); ok {
		e.SetPaddingTop(v)
	}
	if e, v, ok := assertSetter[interface{ SetPaddingBottom(any) }, any](elem, style.PaddingBottom); ok {
		e.SetPaddingBottom(v)
	}
	if e, v, ok := assertSetter[interface{ SetPaddingLeft(any) }, any](elem, style.PaddingLeft); ok {
		e.SetPaddingLeft(v)
	}
	if e, v, ok := assertSetter[interface{ SetPaddingRight(any) }, any](elem, style.PaddingRight); ok {
		e.SetPaddingRight(v)
	}

	if e, v, ok := assertSetter[interface{ SetMarginAll(any) }, any](elem, style.MarginAll); ok {
		e.SetMarginAll(v)
	}
	if e, v, ok := assertSetter[interface{ SetMarginVertical(any) }, any](elem, style.MarginVertical); ok {
		e.SetMarginVertical(v)
	}
	if e, v, ok := assertSetter[interface{ SetMarginHorizontal(any) }, any](elem, style.MarginHorizontal); ok {
		e.SetMarginHorizontal(v)
	}
	if e, v, ok := assertSetter[interface{ SetMarginTop(any) }, any](elem, style.MarginTop); ok {
		e.SetMarginTop(v)
	}
	if e, v, ok := assertSetter[interface{ SetMarginBottom(any) }, any](elem, style.MarginBottom); ok {
		e.SetMarginBottom(v)
	}
	if e, v, ok := assertSetter[interface{ SetMarginLeft(any) }, any](elem, style.MarginLeft); ok {
		e.SetMarginLeft(v)
	}
	if e, v, ok := assertSetter[interface{ SetMarginRight(any) }, any](elem, style.MarginRight); ok {
		e.SetMarginRight(v)
	}

	if e, v, ok := assertSetter[interface{ SetGapAll(any) }, any](elem, style.GapAll); ok {
		e.SetGapAll(v)
	}
	if e, v, ok := assertSetter[interface{ SetGapRow(any) }, any](elem, style.GapRow); ok {
		e.SetGapRow(v)
	}
	if e, v, ok := assertSetter[interface{ SetGapColumn(any) }, any](elem, style.GapColumn); ok {
		e.SetGapColumn(v)
	}
}

func applyStyleBorder(elem core.Element, style Style) {
	if e, v, ok := assertSetter[interface{ SetBorderAll(string) }, string](elem, style.BorderAll); ok && v != "" {
		e.SetBorderAll(v)
	}
	if e, v, ok := assertSetter[interface{ SetBorderVertical(string) }, string](elem, style.BorderVertical); ok && v != "" {
		e.SetBorderVertical(v)
	}
	if e, v, ok := assertSetter[interface{ SetBorderHorizontal(string) }, string](elem, style.BorderHorizontal); ok && v != "" {
		e.SetBorderHorizontal(v)
	}
	if e, v, ok := assertSetter[interface{ SetBorderTop(string) }, string](elem, style.BorderTop); ok && v != "" {
		e.SetBorderTop(v)
	}
	if e, v, ok := assertSetter[interface{ SetBorderBottom(string) }, string](elem, style.BorderBottom); ok && v != "" {
		e.SetBorderBottom(v)
	}
	if e, v, ok := assertSetter[interface{ SetBorderLeft(string) }, string](elem, style.BorderLeft); ok && v != "" {
		e.SetBorderLeft(v)
	}
	if e, v, ok := assertSetter[interface{ SetBorderRight(string) }, string](elem, style.BorderRight); ok && v != "" {
		e.SetBorderRight(v)
	}

	if e, v, ok := assertSetter[interface{ SetBorderColor(string) }, string](elem, style.BorderColor); ok && v != "" {
		e.SetBorderColor(v)
	}
}

func applyStyleDisplay(elem core.Element, style Style) {
	if e, v, ok := assertSetter[interface{ SetDisplay(string) }, string](elem, style.Display); ok && v != "" {
		e.SetDisplay(v)
	}
}

func applyStyleFlexbox(elem core.Element, style Style) {
	if e, v, ok := assertSetter[interface{ SetAlignSelf(string) }, string](elem, style.AlignSelf); ok && v != "" {
		e.SetAlignSelf(v)
	}
	if e, v, ok := assertSetter[interface{ SetAlignItems(string) }, string](elem, style.AlignItems); ok && v != "" {
		e.SetAlignItems(v)
	}
	if e, v, ok := assertSetter[interface{ SetAlignContent(string) }, string](elem, style.AlignContent); ok && v != "" {
		e.SetAlignContent(v)
	}

	if e, v, ok := assertSetter[interface{ SetJustifyContent(string) }, string](elem, style.JustifyContent); ok && v != "" {
		e.SetJustifyContent(v)
	}

	if e, v, ok := assertSetter[interface{ SetFlexDirection(string) }, string](elem, style.FlexDirection); ok && v != "" {
		e.SetFlexDirection(v)
	}

	if e, v, ok := assertSetter[interface{ SetFlexWrap(string) }, string](elem, style.FlexWrap); ok && v != "" {
		e.SetFlexWrap(v)
	}

	if e, v, ok := assertSetter[interface{ SetFlexGrow(string) }, string](elem, style.FlexGrow); ok && v != "" {
		e.SetFlexGrow(v)
	}
	if e, v, ok := assertSetter[interface{ SetFlexShrink(string) }, string](elem, style.FlexShrink); ok && v != "" {
		e.SetFlexShrink(v)
	}
}

func applyStyleOverflow(elem core.Element, style Style) {
	if e, v, ok := assertSetter[interface{ SetOverflow(string) }, string](elem, style.Overflow); ok && v != "" {
		e.SetOverflow(v)
	}
}

func applyStyleBackground(elem core.Element, style Style) {
	if style.BackgroundColor != nil {
		if e, v, ok := assertSetter[interface{ SetBackgroundColor(string) }, string](elem, style.BackgroundColor); ok && v != "" {
			e.SetBackgroundColor(v)
		} else if e, v, ok := assertSetter[interface{ SetTextBackground(string) }, string](elem, style.BackgroundColor); ok && v != "" && style.DropColor == nil {
			e.SetTextBackground(v)
		}
	}

	if style.BackgroundOpacity != nil {
		if e, v, ok := assertSetter[interface{ SetBackgroundOpacity(float32) }, float32](elem, style.BackgroundOpacity); ok {
			e.SetBackgroundOpacity(v)
		}
	}
}

func applyStyleText(elem core.Element, style Style) {
	if e, v, ok := assertSetter[interface{ SetTextForeground(string) }, string](elem, style.Color); ok && v != "" {
		e.SetTextForeground(v)
	}
	if e, v, ok := assertSetter[interface{ SetTextBackground(string) }, string](elem, style.DropColor); ok && v != "" {
		e.SetTextBackground(v)
	}

	if e, v, ok := assertSetter[interface{ SetFontWeight(string) }, string](elem, style.FontWeight); ok && v != "" {
		e.SetFontWeight(v)
	}
	if e, v, ok := assertSetter[interface{ SetFontStyle(string) }, string](elem, style.FontStyle); ok && v != "" {
		e.SetFontStyle(v)
	}

	if e, v, ok := assertSetter[interface{ SetTextDecoration(string) }, string](elem, style.TextDecoration); ok && v != "" {
		e.SetTextDecoration(v)
	}
	if e, v, ok := assertSetter[interface{ SetWrap(string) }, string](elem, style.TextWrap); ok && v != "" {
		e.SetWrap(v)
	}
}

func applyStylePlaceholder(elem core.Element, style Style) {
	if e, v, ok := assertSetter[interface{ SetPlaceholderForeground(string) }, string](elem, style.PlaceholderColor); ok && v != "" {
		e.SetPlaceholderForeground(v)
	}
	if e, v, ok := assertSetter[interface{ SetPlaceholderBackground(string) }, string](elem, style.PlaceholderDropColor); ok && v != "" {
		e.SetPlaceholderBackground(v)
	}

	if e, v, ok := assertSetter[interface{ SetPlaceholderFontWeight(string) }, string](elem, style.PlaceholderFontWeight); ok && v != "" {
		e.SetPlaceholderFontWeight(v)
	}
	if e, v, ok := assertSetter[interface{ SetPlaceholderFontStyle(string) }, string](elem, style.PlaceholderFontStyle); ok && v != "" {
		e.SetPlaceholderFontStyle(v)
	}

	if e, v, ok := assertSetter[interface{ SetPlaceholderDecoration(string) }, string](elem, style.PlaceholderDecoration); ok && v != "" {
		e.SetPlaceholderDecoration(v)
	}
}

func removeStyle(elem core.Element, style Style) {
	removeStyleDimension(elem, style)
	removeStylePositionning(elem, style)
	removeStyleSpacing(elem, style)
	removeStyleBorder(elem, style)
	removeStyleDisplay(elem, style)
	removeStyleFlexbox(elem, style)
	removeStyleOverflow(elem, style)
	removeStyleBackground(elem, style)
	removeStyleText(elem, style)
	removeStylePlaceholder(elem, style)
}

func removeStyleDimension(elem core.Element, style Style) {
	if style.Width != nil {
		elem.UnsetWidth()
	}
	if style.MinWidth != nil {
		elem.UnsetMinWidth()
	}
	if style.MaxWidth != nil {
		elem.UnsetMaxWidth()
	}

	if style.Height != nil {
		elem.UnsetHeight()
	}
	if style.MinHeight != nil {
		elem.UnsetMinHeight()
	}
	if style.MaxHeight != nil {
		elem.UnsetMaxHeight()
	}

}

func removeStylePositionning(elem core.Element, style Style) {
	if style.Top != nil {
		elem.UnsetTop()
	}
	if style.Bottom != nil {
		elem.UnsetBottom()
	}
	if style.Left != nil {
		elem.UnsetLeft()
	}
	if style.Right != nil {
		elem.UnsetRight()
	}

	if style.ZIndex != nil {
		elem.UnsetZIndex()
	}

	if style.Position != nil {
		elem.UnsetPosition()
	}
}

func removeStyleSpacing(elem core.Element, style Style) {
	if style.PaddingAll != nil {
		elem.UnsetPaddingAll()
	}
	if style.PaddingVertical != nil {
		elem.UnsetPaddingVertical()
	}
	if style.PaddingHorizontal != nil {
		elem.UnsetPaddingHorizontal()
	}
	if style.PaddingTop != nil {
		elem.UnsetPaddingTop()
	}
	if style.PaddingBottom != nil {
		elem.UnsetPaddingBottom()
	}
	if style.PaddingLeft != nil {
		elem.UnsetPaddingLeft()
	}
	if style.PaddingRight != nil {
		elem.UnsetPaddingRight()
	}

	if style.MarginAll != nil {
		elem.UnsetMarginAll()
	}
	if style.MarginVertical != nil {
		elem.UnsetMarginVertical()
	}
	if style.MarginHorizontal != nil {
		elem.UnsetMarginHorizontal()
	}
	if style.MarginTop != nil {
		elem.UnsetMarginTop()
	}
	if style.MarginBottom != nil {
		elem.UnsetMarginBottom()
	}
	if style.MarginLeft != nil {
		elem.UnsetMarginLeft()
	}
	if style.MarginRight != nil {
		elem.UnsetMarginRight()
	}

	if style.GapAll != nil {
		elem.UnsetGapAll()
	}
	if style.GapRow != nil {
		elem.UnsetGapRow()
	}
	if style.GapColumn != nil {
		elem.UnsetGapColumn()
	}
}

func removeStyleBorder(elem core.Element, style Style) {
	if style.BorderAll != nil {
		if n, ok := elem.(interface{ UnsetBorderAll() }); ok {
			n.UnsetBorderAll()
		}
	}
	if style.BorderVertical != nil {
		if n, ok := elem.(interface{ UnsetBorderVertical() }); ok {
			n.UnsetBorderVertical()
		}
	}
	if style.BorderHorizontal != nil {
		if n, ok := elem.(interface{ UnsetBorderHorizontal() }); ok {
			n.UnsetBorderHorizontal()
		}
	}
	if style.BorderTop != nil {
		if n, ok := elem.(interface{ UnsetBorderTop() }); ok {
			n.UnsetBorderTop()
		}
	}
	if style.BorderBottom != nil {
		if n, ok := elem.(interface{ UnsetBorderBottom() }); ok {
			n.UnsetBorderBottom()
		}
	}
	if style.BorderLeft != nil {
		if n, ok := elem.(interface{ UnsetBorderLeft() }); ok {
			n.UnsetBorderLeft()
		}
	}
	if style.BorderRight != nil {
		if n, ok := elem.(interface{ UnsetBorderRight() }); ok {
			n.UnsetBorderRight()
		}
	}

	if style.BorderColor != nil {
		if n, ok := elem.(interface{ UnsetBorderColor() }); ok {
			n.UnsetBorderColor()
		}
	}
}

func removeStyleDisplay(elem core.Element, style Style) {
	if style.Display != nil {
		elem.UnsetDisplay()
	}
}

func removeStyleFlexbox(elem core.Element, style Style) {
	if style.AlignSelf != nil {
		elem.UnsetAlignSelf()
	}
	if style.AlignItems != nil {
		elem.UnsetAlignItems()
	}
	if style.AlignContent != nil {
		elem.UnsetAlignContent()
	}

	if style.JustifyContent != nil {
		elem.UnsetJustifyContent()
	}
	if style.FlexDirection != nil {
		elem.UnsetFlexDirection()
	}
	if style.FlexWrap != nil {
		elem.UnsetFlexWrap()
	}
	if style.FlexGrow != nil {
		elem.UnsetFlexGrow()
	}
	if style.FlexShrink != nil {
		elem.UnsetFlexShrink()
	}
}

func removeStyleOverflow(elem core.Element, style Style) {
	if style.Overflow != nil {
		elem.UnsetOverflow()
	}
}

func removeStyleBackground(elem core.Element, style Style) {
	if style.BackgroundColor != nil {
		if n, ok := elem.(interface{ UnsetBackgroundColor() }); ok {
			n.UnsetBackgroundColor()
		} else if n, ok := elem.(interface{ UnsetTextBackground() }); ok && style.DropColor == "" {
			n.UnsetTextBackground()
		}
	}

	if style.BackgroundOpacity != nil {
		if n, ok := elem.(interface{ UnsetBackgroundOpacity() }); ok {
			n.UnsetBackgroundOpacity()
		}
	}
}

func removeStyleText(elem core.Element, style Style) {
	if style.Color != nil {
		if n, ok := elem.(interface{ UnsetForegroundColor() }); ok {
			n.UnsetForegroundColor()
		}
	}
	if style.DropColor != nil {
		if n, ok := elem.(interface{ UnsetTextBackground() }); ok {
			n.UnsetTextBackground()
		}
	}

	if style.FontWeight != nil {
		if n, ok := elem.(interface{ UnsetFontWeight() }); ok {
			n.UnsetFontWeight()
		}
	}
	if style.FontStyle != nil {
		if n, ok := elem.(interface{ UnsetFontStyle() }); ok {
			n.UnsetFontStyle()
		}
	}

	if style.TextDecoration != nil {
		if n, ok := elem.(interface{ UnsetTextDecoration() }); ok {
			n.UnsetTextDecoration()
		}
	}
	if style.TextWrap != nil {
		if n, ok := elem.(interface{ UnsetWrap() }); ok {
			n.UnsetWrap()
		}
	}
}

func removeStylePlaceholder(elem core.Element, style Style) {
	if style.PlaceholderColor != nil {
		if n, ok := elem.(interface{ UnsetPlaceholderForeground() }); ok {
			n.UnsetPlaceholderForeground()
		}
	}
	if style.PlaceholderDropColor != nil {
		if n, ok := elem.(interface{ UnsetPlaceholderBackground() }); ok {
			n.UnsetPlaceholderBackground()
		}
	}

	if style.PlaceholderFontWeight != nil {
		if n, ok := elem.(interface{ UnsetPlaceholderFontWeight() }); ok {
			n.UnsetPlaceholderFontWeight()
		}
	}
	if style.PlaceholderFontStyle != nil {
		if n, ok := elem.(interface{ UnsetPlaceholderFontStyle() }); ok {
			n.UnsetPlaceholderFontStyle()
		}
	}

	if style.PlaceholderDecoration != nil {
		if n, ok := elem.(interface{ UnsetPlaceholderDecoration() }); ok {
			n.UnsetPlaceholderDecoration()
		}
	}
}

func assertSetter[E, V any](elem core.Element, value any) (E, V, bool) {
	var v V
	var vok bool
	var e E
	var eok bool

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

	if elem, ok := elem.(E); ok {
		e = elem
		eok = true
	}

	return e, v, eok && vok
}
