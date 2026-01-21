package components

import (
	"fmt"
	"slices"

	"github.com/AnatoleLucet/loom-term/core"
)

type styleLayer struct {
	id    string
	style Style
}

type styleStack struct {
	layers []styleLayer
}

func (s *styleStack) Push(id string, style Style) {
	s.layers = append(s.layers, styleLayer{id, style})
}

func (s *styleStack) Replace(id string, style Style) {
	for i, layer := range s.layers {
		if layer.id == id {
			s.layers[i].style = style
			return
		}
	}
}

func (s *styleStack) Pop(id string) {
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

func applyStyleStack(elem core.Element) error {
	stack := getStyleStack(elem)

	for _, layer := range stack.layers {
		err := applyStyle(elem, &layer.style)
		if err != nil {
			return err
		}
	}

	return nil
}

func applyStyle(elem core.Element, style *Style) (err error) {
	err = applyStyleDimension(elem, style)
	err = applyStylePositionning(elem, style)
	err = applyStyleSpacing(elem, style)
	err = applyStyleDisplay(elem, style)
	err = applyStyleFlexbox(elem, style)
	err = applyStyleOverflow(elem, style)
	err = applyStyleBackground(elem, style)
	err = applyStyleText(elem, style)

	// err = s.applyStyleBorder(elem, style)

	return err
}

func applyStyleDimension(elem core.Element, style *Style) (err error) {
	if style.Width != nil {
		err = elem.SetWidth(style.Width)
	}
	if style.MinWidth != nil {
		err = elem.SetMinWidth(style.MinWidth)
	}
	if style.MaxWidth != nil {
		err = elem.SetMaxWidth(style.MaxWidth)
	}

	if style.Height != nil {
		err = elem.SetHeight(style.Height)
	}
	if style.MinHeight != nil {
		err = elem.SetMinHeight(style.MinHeight)
	}
	if style.MaxHeight != nil {
		err = elem.SetMaxHeight(style.MaxHeight)
	}

	return err
}

func applyStylePositionning(elem core.Element, style *Style) (err error) {
	if style.Top != nil {
		err = elem.SetTop(style.Top)
	}
	if style.Bottom != nil {
		err = elem.SetBottom(style.Bottom)
	}
	if style.Left != nil {
		err = elem.SetLeft(style.Left)
	}
	if style.Right != nil {
		err = elem.SetRight(style.Right)
	}

	if elem.ZIndex() != style.ZIndex {
		err = elem.SetZIndex(style.ZIndex)
	}

	if style.Position != "" {
		err = elem.SetPosition(style.Position)
	}

	return err
}

func applyStyleSpacing(elem core.Element, style *Style) (err error) {
	if style.PaddingAll != nil {
		err = elem.SetPaddingAll(style.PaddingAll)
	}
	if style.PaddingVertical != nil {
		err = elem.SetPaddingVertical(style.PaddingVertical)
	}
	if style.PaddingHorizontal != nil {
		err = elem.SetPaddingHorizontal(style.PaddingHorizontal)
	}
	if style.PaddingTop != nil {
		err = elem.SetPaddingTop(style.PaddingTop)
	}
	if style.PaddingBottom != nil {
		err = elem.SetPaddingBottom(style.PaddingBottom)
	}
	if style.PaddingLeft != nil {
		err = elem.SetPaddingLeft(style.PaddingLeft)
	}
	if style.PaddingRight != nil {
		err = elem.SetPaddingRight(style.PaddingRight)
	}

	if style.MarginAll != nil {
		err = elem.SetMarginAll(style.MarginAll)
	}
	if style.MarginVertical != nil {
		err = elem.SetMarginVertical(style.MarginVertical)
	}
	if style.MarginHorizontal != nil {
		err = elem.SetMarginHorizontal(style.MarginHorizontal)
	}
	if style.MarginTop != nil {
		err = elem.SetMarginTop(style.MarginTop)
	}
	if style.MarginBottom != nil {
		err = elem.SetMarginBottom(style.MarginBottom)
	}
	if style.MarginLeft != nil {
		err = elem.SetMarginLeft(style.MarginLeft)
	}
	if style.MarginRight != nil {
		err = elem.SetMarginRight(style.MarginRight)
	}

	if style.GapAll != nil {
		err = elem.SetGapAll(style.GapAll)
	}
	if style.GapRow != nil {
		err = elem.SetGapRow(style.GapRow)
	}
	if style.GapColumn != nil {
		err = elem.SetGapColumn(style.GapColumn)
	}

	return err
}

func applyStyleDisplay(elem core.Element, style *Style) error {
	if style.Display != "" {
		err := elem.SetDisplay(style.Display)
		if err != nil {
			return fmt.Errorf("unable to set display: %w", err)
		}
	}

	return nil
}

func applyStyleFlexbox(elem core.Element, style *Style) (err error) {
	if style.AlignSelf != "" {
		err = elem.SetAlignSelf(style.AlignSelf)
	}
	if style.AlignItems != "" {
		err = elem.SetAlignItems(style.AlignItems)
	}
	if style.AlignContent != "" {
		err = elem.SetAlignContent(style.AlignContent)
	}

	if style.JustifyContent != "" {
		err = elem.SetJustifyContent(style.JustifyContent)
	}

	if style.FlexDirection != "" {
		err = elem.SetFlexDirection(style.FlexDirection)
	}

	if style.FlexWrap != "" {
		err = elem.SetFlexWrap(style.FlexWrap)
	}

	if style.FlexGrow != "" {
		err = elem.SetFlexGrow(style.FlexGrow)
	}
	if style.FlexShrink != "" {
		err = elem.SetFlexShrink(style.FlexShrink)
	}

	return err
}

func applyStyleOverflow(elem core.Element, style *Style) error {
	if style.Overflow != "" {
		err := elem.SetOverflow(style.Overflow)
		if err != nil {
			return fmt.Errorf("unable to set overflow: %w", err)
		}
	}

	return nil
}

func applyStyleBackground(elem core.Element, style *Style) error {
	if style.BackgroundColor != "" {
		if n, ok := elem.(interface{ SetBackgroundColor(string) error }); ok {
			err := n.SetBackgroundColor(style.BackgroundColor)
			if err != nil {
				return fmt.Errorf("unable to set background color: %w", err)
			}
		}
	}

	return nil
}

func applyStyleText(elem core.Element, style *Style) error {
	if style.Color != "" {
		if n, ok := elem.(interface{ SetForegroundColor(string) error }); ok {
			err := n.SetForegroundColor(style.Color)
			if err != nil {
				return fmt.Errorf("unable to set text color: %w", err)
			}
		}
	}

	return nil
}

func removeStyle(elem core.Element, style *Style) (err error) {
	err = removeStyleDimension(elem, style)
	err = removeStylePositionning(elem, style)
	err = removeStyleSpacing(elem, style)
	err = removeStyleDisplay(elem, style)
	err = removeStyleFlexbox(elem, style)
	err = removeStyleOverflow(elem, style)
	err = removeStyleBackground(elem, style)
	err = removeStyleText(elem, style)

	// err = s.removeStyleBorder(elem, style)

	return err
}

func removeStyleDimension(elem core.Element, style *Style) (err error) {
	if style.Width != nil {
		err = elem.UnsetWidth()
	}
	if style.MinWidth != nil {
		err = elem.UnsetMinWidth()
	}
	if style.MaxWidth != nil {
		err = elem.UnsetMaxWidth()
	}

	if style.Height != nil {
		err = elem.UnsetHeight()
	}
	if style.MinHeight != nil {
		err = elem.UnsetMinHeight()
	}
	if style.MaxHeight != nil {
		err = elem.UnsetMaxHeight()
	}

	return err
}

func removeStylePositionning(elem core.Element, style *Style) (err error) {
	if style.Top != nil {
		err = elem.UnsetTop()
	}
	if style.Bottom != nil {
		err = elem.UnsetBottom()
	}
	if style.Left != nil {
		err = elem.UnsetLeft()
	}
	if style.Right != nil {
		err = elem.UnsetRight()
	}

	if style.ZIndex != 0 {
		err = elem.UnsetZIndex()
	}

	if style.Position != "" {
		err = elem.UnsetPosition()
	}

	return err
}

func removeStyleSpacing(elem core.Element, style *Style) (err error) {
	if style.PaddingAll != nil {
		err = elem.UnsetPaddingAll()
	}
	if style.PaddingVertical != nil {
		err = elem.UnsetPaddingVertical()
	}
	if style.PaddingHorizontal != nil {
		err = elem.UnsetPaddingHorizontal()
	}
	if style.PaddingTop != nil {
		err = elem.UnsetPaddingTop()
	}
	if style.PaddingBottom != nil {
		err = elem.UnsetPaddingBottom()
	}
	if style.PaddingLeft != nil {
		err = elem.UnsetPaddingLeft()
	}
	if style.PaddingRight != nil {
		err = elem.UnsetPaddingRight()
	}

	if style.MarginAll != nil {
		err = elem.UnsetMarginAll()
	}
	if style.MarginVertical != nil {
		err = elem.UnsetMarginVertical()
	}
	if style.MarginHorizontal != nil {
		err = elem.UnsetMarginHorizontal()
	}
	if style.MarginTop != nil {
		err = elem.UnsetMarginTop()
	}
	if style.MarginBottom != nil {
		err = elem.UnsetMarginBottom()
	}
	if style.MarginLeft != nil {
		err = elem.UnsetMarginLeft()
	}
	if style.MarginRight != nil {
		err = elem.UnsetMarginRight()
	}

	if style.GapAll != nil {
		err = elem.UnsetGapAll()
	}
	if style.GapRow != nil {
		err = elem.UnsetGapRow()
	}
	if style.GapColumn != nil {
		err = elem.UnsetGapColumn()
	}

	return err
}

func removeStyleDisplay(elem core.Element, style *Style) error {
	if style.Display != "" {
		err := elem.UnsetDisplay()
		if err != nil {
			return fmt.Errorf("unable to unset display: %w", err)
		}
	}

	return nil
}

func removeStyleFlexbox(elem core.Element, style *Style) (err error) {
	if style.AlignSelf != "" {
		err = elem.UnsetAlignSelf()
	}
	if style.AlignItems != "" {
		err = elem.UnsetAlignItems()
	}
	if style.AlignContent != "" {
		err = elem.UnsetAlignContent()
	}

	if style.JustifyContent != "" {
		err = elem.UnsetJustifyContent()
	}
	if style.FlexDirection != "" {
		err = elem.UnsetFlexDirection()
	}
	if style.FlexWrap != "" {
		err = elem.UnsetFlexWrap()
	}
	if style.FlexGrow != "" {
		err = elem.UnsetFlexGrow()
	}
	if style.FlexShrink != "" {
		err = elem.UnsetFlexShrink()
	}

	return err
}

func removeStyleOverflow(elem core.Element, style *Style) error {
	if style.Overflow != "" {
		err := elem.UnsetOverflow()
		if err != nil {
			return fmt.Errorf("unable to unset overflow: %w", err)
		}
	}

	return nil
}

func removeStyleBackground(elem core.Element, style *Style) error {
	if style.BackgroundColor != "" {
		if n, ok := elem.(interface{ UnsetBackgroundColor() error }); ok {
			err := n.UnsetBackgroundColor()
			if err != nil {
				return fmt.Errorf("unable to unset background color: %w", err)
			}
		}
	}

	return nil
}

func removeStyleText(elem core.Element, style *Style) error {
	if style.Color != "" {
		if n, ok := elem.(interface{ UnsetForegroundColor() error }); ok {
			err := n.UnsetForegroundColor()
			if err != nil {
				return fmt.Errorf("unable to unset text color: %w", err)
			}
		}
	}

	return nil
}
