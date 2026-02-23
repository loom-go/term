package components

import (
	"slices"

	"github.com/AnatoleLucet/loom-term/core"
)

type styleLayer struct {
	id    uint32
	style Style
}

type styleStack struct {
	layers []styleLayer
}

func (s *styleStack) Push(id uint32, style Style) {
	s.layers = append(s.layers, styleLayer{id, style})
}

func (s *styleStack) Replace(id uint32, style Style) {
	for i, layer := range s.layers {
		if layer.id == id {
			s.layers[i].style = style
			return
		}
	}
}

func (s *styleStack) Pop(id uint32) {
	s.layers = slices.DeleteFunc(s.layers, func(layer styleLayer) bool {
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

func applyStyleStack(elem core.Element) {
	stack := getStyleStack(elem)

	for _, layer := range stack.layers {
		applyStyle(elem, &layer.style)
	}
}

func applyStyle(elem core.Element, style *Style) {
	applyStyleDimension(elem, style)
	applyStylePositionning(elem, style)
	applyStyleSpacing(elem, style)
	applyStyleBorder(elem, style)
	applyStyleDisplay(elem, style)
	applyStyleFlexbox(elem, style)
	applyStyleOverflow(elem, style)
	applyStyleBackground(elem, style)
	applyStyleText(elem, style)
}

func applyStyleDimension(elem core.Element, style *Style) {
	if style.Width != nil {
		elem.SetWidth(style.Width)
	}
	if style.MinWidth != nil {
		elem.SetMinWidth(style.MinWidth)
	}
	if style.MaxWidth != nil {
		elem.SetMaxWidth(style.MaxWidth)
	}

	if style.Height != nil {
		elem.SetHeight(style.Height)
	}
	if style.MinHeight != nil {
		elem.SetMinHeight(style.MinHeight)
	}
	if style.MaxHeight != nil {
		elem.SetMaxHeight(style.MaxHeight)
	}
}

func applyStylePositionning(elem core.Element, style *Style) {
	if style.Top != nil {
		elem.SetTop(style.Top)
	}
	if style.Bottom != nil {
		elem.SetBottom(style.Bottom)
	}
	if style.Left != nil {
		elem.SetLeft(style.Left)
	}
	if style.Right != nil {
		elem.SetRight(style.Right)
	}

	if elem.ZIndex() != style.ZIndex {
		elem.SetZIndex(style.ZIndex)
	}

	if style.Position != "" {
		elem.SetPosition(style.Position)
	}
}

func applyStyleSpacing(elem core.Element, style *Style) {
	if style.PaddingAll != nil {
		elem.SetPaddingAll(style.PaddingAll)
	}
	if style.PaddingVertical != nil {
		elem.SetPaddingVertical(style.PaddingVertical)
	}
	if style.PaddingHorizontal != nil {
		elem.SetPaddingHorizontal(style.PaddingHorizontal)
	}
	if style.PaddingTop != nil {
		elem.SetPaddingTop(style.PaddingTop)
	}
	if style.PaddingBottom != nil {
		elem.SetPaddingBottom(style.PaddingBottom)
	}
	if style.PaddingLeft != nil {
		elem.SetPaddingLeft(style.PaddingLeft)
	}
	if style.PaddingRight != nil {
		elem.SetPaddingRight(style.PaddingRight)
	}

	if style.MarginAll != nil {
		elem.SetMarginAll(style.MarginAll)
	}
	if style.MarginVertical != nil {
		elem.SetMarginVertical(style.MarginVertical)
	}
	if style.MarginHorizontal != nil {
		elem.SetMarginHorizontal(style.MarginHorizontal)
	}
	if style.MarginTop != nil {
		elem.SetMarginTop(style.MarginTop)
	}
	if style.MarginBottom != nil {
		elem.SetMarginBottom(style.MarginBottom)
	}
	if style.MarginLeft != nil {
		elem.SetMarginLeft(style.MarginLeft)
	}
	if style.MarginRight != nil {
		elem.SetMarginRight(style.MarginRight)
	}

	if style.GapAll != nil {
		elem.SetGapAll(style.GapAll)
	}
	if style.GapRow != nil {
		elem.SetGapRow(style.GapRow)
	}
	if style.GapColumn != nil {
		elem.SetGapColumn(style.GapColumn)
	}
}

func applyStyleBorder(elem core.Element, style *Style) {
	if style.BorderAll != "" {
		if n, ok := elem.(interface{ SetBorderAll(string) }); ok {
			n.SetBorderAll(style.BorderAll)
		}
	}
	if style.BorderVertical != "" {
		if n, ok := elem.(interface{ SetBorderVertical(string) }); ok {
			n.SetBorderVertical(style.BorderVertical)
		}
	}
	if style.BorderHorizontal != "" {
		if n, ok := elem.(interface{ SetBorderHorizontal(string) }); ok {
			n.SetBorderHorizontal(style.BorderHorizontal)
		}
	}
	if style.BorderTop != "" {
		if n, ok := elem.(interface{ SetBorderTop(string) }); ok {
			n.SetBorderTop(style.BorderTop)
		}
	}
	if style.BorderBottom != "" {
		if n, ok := elem.(interface{ SetBorderBottom(string) }); ok {
			n.SetBorderBottom(style.BorderBottom)
		}
	}
	if style.BorderLeft != "" {
		if n, ok := elem.(interface{ SetBorderLeft(string) }); ok {
			n.SetBorderLeft(style.BorderLeft)
		}
	}
	if style.BorderRight != "" {
		if n, ok := elem.(interface{ SetBorderRight(string) }); ok {
			n.SetBorderRight(style.BorderRight)
		}
	}
}

func applyStyleDisplay(elem core.Element, style *Style) {
	if style.Display != "" {
		elem.SetDisplay(style.Display)
	}
}

func applyStyleFlexbox(elem core.Element, style *Style) {
	if style.AlignSelf != "" {
		elem.SetAlignSelf(style.AlignSelf)
	}
	if style.AlignItems != "" {
		elem.SetAlignItems(style.AlignItems)
	}
	if style.AlignContent != "" {
		elem.SetAlignContent(style.AlignContent)
	}

	if style.JustifyContent != "" {
		elem.SetJustifyContent(style.JustifyContent)
	}

	if style.FlexDirection != "" {
		elem.SetFlexDirection(style.FlexDirection)
	}

	if style.FlexWrap != "" {
		elem.SetFlexWrap(style.FlexWrap)
	}

	if style.FlexGrow != "" {
		elem.SetFlexGrow(style.FlexGrow)
	}
	if style.FlexShrink != "" {
		elem.SetFlexShrink(style.FlexShrink)
	}
}

func applyStyleOverflow(elem core.Element, style *Style) {
	if style.Overflow != "" {
		elem.SetOverflow(style.Overflow)
	}
}

func applyStyleBackground(elem core.Element, style *Style) {
	if style.BackgroundColor != "" {
		if n, ok := elem.(interface{ SetBackgroundColor(string) }); ok {
			n.SetBackgroundColor(style.BackgroundColor)
		} else if n, ok := elem.(interface{ SetTextBackground(string) }); ok && style.DropColor == "" {
			n.SetTextBackground(style.BackgroundColor)
		}
	}
}

func applyStyleText(elem core.Element, style *Style) {
	if style.Color != "" {
		if n, ok := elem.(interface{ SetTextForeground(string) }); ok {
			n.SetTextForeground(style.Color)
		}
	}
	if style.DropColor != "" {
		if n, ok := elem.(interface{ SetTextBackground(string) }); ok {
			n.SetTextBackground(style.DropColor)
		}
	}

	if style.FontWeight != "" {
		if n, ok := elem.(interface{ SetFontWeight(string) }); ok {
			n.SetFontWeight(style.FontWeight)
		}
	}
	if style.FontStyle != "" {
		if n, ok := elem.(interface{ SetFontStyle(string) }); ok {
			n.SetFontStyle(style.FontStyle)
		}
	}

	if style.TextDecoration != "" {
		if n, ok := elem.(interface{ SetTextDecoration(string) }); ok {
			n.SetTextDecoration(style.TextDecoration)
		}
	}
	if style.TextWrap != "" {
		if n, ok := elem.(interface{ SetWrap(string) }); ok {
			n.SetWrap(style.TextWrap)
		}
	}
}

func removeStyle(elem core.Element, style *Style) {
	removeStyleDimension(elem, style)
	removeStylePositionning(elem, style)
	removeStyleSpacing(elem, style)
	removeStyleBorder(elem, style)
	removeStyleDisplay(elem, style)
	removeStyleFlexbox(elem, style)
	removeStyleOverflow(elem, style)
	removeStyleBackground(elem, style)
	removeStyleText(elem, style)
}

func removeStyleDimension(elem core.Element, style *Style) {
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

func removeStylePositionning(elem core.Element, style *Style) {
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

	if style.ZIndex != 0 {
		elem.UnsetZIndex()
	}

	if style.Position != "" {
		elem.UnsetPosition()
	}
}

func removeStyleSpacing(elem core.Element, style *Style) {
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

func removeStyleBorder(elem core.Element, style *Style) {
	if style.BorderAll != "" {
		if n, ok := elem.(interface{ UnsetBorderAll() }); ok {
			n.UnsetBorderAll()
		}
	}
	if style.BorderVertical != "" {
		if n, ok := elem.(interface{ UnsetBorderVertical() }); ok {
			n.UnsetBorderVertical()
		}
	}
	if style.BorderHorizontal != "" {
		if n, ok := elem.(interface{ UnsetBorderHorizontal() }); ok {
			n.UnsetBorderHorizontal()
		}
	}
	if style.BorderTop != "" {
		if n, ok := elem.(interface{ UnsetBorderTop() }); ok {
			n.UnsetBorderTop()
		}
	}
	if style.BorderBottom != "" {
		if n, ok := elem.(interface{ UnsetBorderBottom() }); ok {
			n.UnsetBorderBottom()
		}
	}
	if style.BorderLeft != "" {
		if n, ok := elem.(interface{ UnsetBorderLeft() }); ok {
			n.UnsetBorderLeft()
		}
	}
	if style.BorderRight != "" {
		if n, ok := elem.(interface{ UnsetBorderRight() }); ok {
			n.UnsetBorderRight()
		}
	}
}

func removeStyleDisplay(elem core.Element, style *Style) {
	if style.Display != "" {
		elem.UnsetDisplay()
	}
}

func removeStyleFlexbox(elem core.Element, style *Style) {
	if style.AlignSelf != "" {
		elem.UnsetAlignSelf()
	}
	if style.AlignItems != "" {
		elem.UnsetAlignItems()
	}
	if style.AlignContent != "" {
		elem.UnsetAlignContent()
	}

	if style.JustifyContent != "" {
		elem.UnsetJustifyContent()
	}
	if style.FlexDirection != "" {
		elem.UnsetFlexDirection()
	}
	if style.FlexWrap != "" {
		elem.UnsetFlexWrap()
	}
	if style.FlexGrow != "" {
		elem.UnsetFlexGrow()
	}
	if style.FlexShrink != "" {
		elem.UnsetFlexShrink()
	}
}

func removeStyleOverflow(elem core.Element, style *Style) {
	if style.Overflow != "" {
		elem.UnsetOverflow()
	}
}

func removeStyleBackground(elem core.Element, style *Style) {
	if style.BackgroundColor != "" {
		if n, ok := elem.(interface{ UnsetBackgroundColor() }); ok {
			n.UnsetBackgroundColor()
		} else if n, ok := elem.(interface{ UnsetTextBackground() }); ok && style.DropColor == "" {
			n.UnsetTextBackground()
		}
	}
}

func removeStyleText(elem core.Element, style *Style) {
	if style.Color != "" {
		if n, ok := elem.(interface{ UnsetForegroundColor() }); ok {
			n.UnsetForegroundColor()
		}
	}
	if style.DropColor != "" {
		if n, ok := elem.(interface{ UnsetTextBackground() }); ok {
			n.UnsetTextBackground()
		}
	}

	if style.FontWeight != "" {
		if n, ok := elem.(interface{ UnsetFontWeight() }); ok {
			n.UnsetFontWeight()
		}
	}
	if style.FontStyle != "" {
		if n, ok := elem.(interface{ UnsetFontStyle() }); ok {
			n.UnsetFontStyle()
		}
	}

	if style.TextDecoration != "" {
		if n, ok := elem.(interface{ UnsetTextDecoration() }); ok {
			n.UnsetTextDecoration()
		}
	}
	if style.TextWrap != "" {
		if n, ok := elem.(interface{ UnsetWrap() }); ok {
			n.UnsetWrap()
		}
	}
}
