package components

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/AnatoleLucet/loom"
	termerror "github.com/AnatoleLucet/loom-term/error"
	"github.com/AnatoleLucet/loom-term/internal"
	"github.com/AnatoleLucet/loom-term/opentui"
	"github.com/AnatoleLucet/tess"
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

	Display string // "none" | "flex" | "contents"

	AlignItems     string // "flex-start" | "flex-end" | "center" | "stretch" | "baseline"
	JustifyContent string // "flex-start" | "flex-end" | "center" | "space-between" | "space-around" | "space-evenly"
	FlexDirection  string // "row" | "row-reverse" | "column" | "column-reverse"
	FlexWrap       string // "nowrap" | "wrap" | "wrap-reverse"

	BorderStyle string // "single" | "double" | "rounded"  | "heavy"

	Color string // "transparent" | "#RGB" | "#RRGGBBAA"
	// BackgroundColor string // "transparent" | "#RGB" | "#RRGGBBAA"

	BackgroundColor opentui.RGBA
}

func (s *Style) ID() string {
	return "term.Style"
}

func (s *Style) Mount(slot *loom.Slot) error {
	return s.Update(slot)
}

func (s *Style) Update(slot *loom.Slot) error {
	node, ok := slot.Parent().(*internal.Element)
	if !ok {
		return nil
	}

	err := s.apply(node)

	return err
}

func (s *Style) Unmount(slot *loom.Slot) error {
	return nil
}

func (s *Style) apply(node *internal.Element) (err error) {
	err = s.applyDimension(node)
	err = s.applyOffset(node)
	err = s.applySpacing(node)
	err = s.applyDisplay(node)
	err = s.applyFlexbox(node)

	// todo: maybe these should be handled in the components directly?
	// err = s.applyBorder(node)
	// err = s.applyBackground(node)
	// err = s.applyText(node)

	// temporary ugly fix
	node.SetBackgroundColor(s.BackgroundColor)

	return err
}

func (s *Style) applyDimension(node *internal.Element) error {
	layout := node.Layout()

	if s.Width != nil {
		v, err := toTessValue(s.Width)
		if err != nil {
			return fmt.Errorf("invalid width: %w", err)
		}

		layout.SetWidth(v)
	}

	if s.MinWidth != nil {
		v, err := toTessValue(s.MinWidth)
		if err != nil {
			return fmt.Errorf("invalid min width: %w", err)
		}

		layout.SetMinWidth(v)
	}

	if s.MaxWidth != nil {
		v, err := toTessValue(s.MaxWidth)
		if err != nil {
			return fmt.Errorf("invalid max width: %w", err)
		}

		layout.SetMaxWidth(v)
	}

	if s.Height != nil {
		v, err := toTessValue(s.Height)
		if err != nil {
			return fmt.Errorf("invalid height: %w", err)
		}

		layout.SetHeight(v)
	}

	if s.MinHeight != nil {
		v, err := toTessValue(s.MinHeight)
		if err != nil {
			return fmt.Errorf("invalid min height: %w", err)
		}

		layout.SetMinHeight(v)
	}

	if s.MaxHeight != nil {
		v, err := toTessValue(s.MaxHeight)
		if err != nil {
			return fmt.Errorf("invalid max height: %w", err)
		}

		layout.SetMaxHeight(v)
	}

	return nil
}

func (s *Style) applyOffset(node *internal.Element) error {
	layout := node.Layout()

	if s.Top != nil {
		v, err := toTessValue(s.Top)
		if err != nil {
			return fmt.Errorf("invalid top position: %w", err)
		}

		layout.SetTop(v)
	}

	if s.Bottom != nil {
		v, err := toTessValue(s.Bottom)
		if err != nil {
			return fmt.Errorf("invalid bottom position: %w", err)
		}

		layout.SetBottom(v)
	}

	if s.Left != nil {
		v, err := toTessValue(s.Left)
		if err != nil {
			return fmt.Errorf("invalid left position: %w", err)
		}

		layout.SetLeft(v)
	}

	if s.Right != nil {
		v, err := toTessValue(s.Right)
		if err != nil {
			return fmt.Errorf("invalid right position: %w", err)
		}

		layout.SetRight(v)
	}

	if s.Position != "" {
		positionMap := map[string]tess.PositionType{
			"static":   tess.Static,
			"relative": tess.Relative,
			"absolute": tess.Absolute,
		}

		if position, ok := positionMap[s.Position]; ok {
			err := layout.SetPosition(position)
			if err != nil {
				return fmt.Errorf("unable to set position type: %w", err)
			}
		}
	}

	return nil
}

func (s *Style) applySpacing(node *internal.Element) error {
	layout := node.Layout()

	padding, err := toTessEdges(
		s.PaddingAll,
		s.PaddingVertical,
		s.PaddingHorizontal,
		s.PaddingTop,
		s.PaddingBottom,
		s.PaddingLeft,
		s.PaddingRight,
	)
	if err != nil {
		return fmt.Errorf("invalid padding: %w", err)
	}

	margin, err := toTessEdges(
		s.MarginAll,
		s.MarginVertical,
		s.MarginHorizontal,
		s.MarginTop,
		s.MarginBottom,
		s.MarginLeft,
		s.MarginRight,
	)
	if err != nil {
		return fmt.Errorf("invalid margin: %w", err)
	}

	err = layout.SetPadding(padding)
	if err != nil {
		return fmt.Errorf("unable to set padding: %w", err)
	}

	err = layout.SetMargin(margin)
	if err != nil {
		return fmt.Errorf("unable to set margin: %w", err)
	}

	return nil
}

func (s *Style) applyDisplay(node *internal.Element) error {
	layout := node.Layout()

	displayMap := map[string]tess.DisplayType{
		"none":     tess.None,
		"flex":     tess.Flex,
		"contents": tess.Contents,
	}

	if display, ok := displayMap[s.Display]; ok {
		err := layout.SetDisplay(display)

		if err != nil {
			return fmt.Errorf("unable to set display type: %w", err)
		}
	}

	return nil
}

func (s *Style) applyFlexbox(node *internal.Element) error {
	layout := node.Layout()

	if s.AlignItems != "" {
		alignMap := map[string]tess.FlexAlign{
			"flex-start": tess.AlignStart,
			"flex-end":   tess.AlignEnd,
			"center":     tess.AlignCenter,
			"stretch":    tess.AlignStretch,
			"baseline":   tess.AlignBaseline,
		}

		if align, ok := alignMap[s.AlignItems]; ok {
			err := layout.SetAlignItems(align)
			if err != nil {
				return fmt.Errorf("unable to set align items: %w", err)
			}
		}
	}

	if s.JustifyContent != "" {
		justifyMap := map[string]tess.FlexJustify{
			"flex-start":    tess.JustifyStart,
			"flex-end":      tess.JustifyEnd,
			"center":        tess.JustifyCenter,
			"space-between": tess.JustifySpaceBetween,
			"space-around":  tess.JustifySpaceAround,
			"space-evenly":  tess.JustifySpaceEvenly,
		}

		if justify, ok := justifyMap[s.JustifyContent]; ok {
			err := layout.SetJustifyContent(justify)
			if err != nil {
				return fmt.Errorf("unable to set justify content: %w", err)
			}
		}
	}

	if s.FlexDirection != "" {
		directionMap := map[string]tess.FlexDirection{
			"row":            tess.Row,
			"row-reverse":    tess.RowReverse,
			"column":         tess.Column,
			"column-reverse": tess.ColumnReverse,
		}

		if direction, ok := directionMap[s.FlexDirection]; ok {
			err := layout.SetFlexDirection(direction)
			if err != nil {
				return fmt.Errorf("unable to set flex direction: %w", err)
			}
		}
	}

	if s.FlexWrap != "" {
		wrapMap := map[string]tess.FlexWrap{
			"nowrap":       tess.NoWrap,
			"wrap":         tess.Wrap,
			"wrap-reverse": tess.WrapReverse,
		}

		if wrap, ok := wrapMap[s.FlexWrap]; ok {
			err := layout.SetFlexWrap(wrap)
			if err != nil {
				return fmt.Errorf("unable to set flex wrap: %w", err)
			}
		}
	}

	return nil
}

func toTessValue(value any) (tess.Value, error) {
	switch v := value.(type) {
	case nil:
		return tess.Undefined(), nil
	case int, int32, int64:
		return tess.Point(float32(v.(int))), nil
	case float32:
		return tess.Point(v), nil
	case float64:
		return tess.Point(float32(v)), nil
	case string:
		vMap := map[string]tess.Value{
			"":            tess.Undefined(),
			"undefined":   tess.Undefined(),
			"none":        tess.Undefined(),
			"auto":        tess.Auto(),
			"max-content": tess.MaxContent(),
			"fit-content": tess.FitContent(),
			"stretch":     tess.Stretch(),
		}

		if val, ok := vMap[v]; ok {
			return val, nil
		}

		// 100pt or 100.5pt
		if regexp.MustCompile(`^\d+(\.\d+)?pt$`).MatchString(v) {
			pointStr := strings.TrimSuffix(v, "pt")
			points, err := strconv.ParseFloat(pointStr, 32)
			if err != nil {
				return tess.Undefined(), fmt.Errorf("%w: invalid point value '%s'", termerror.ErrInvalidStyleValue, v)
			}

			return tess.Point(float32(points)), nil
		}

		// 100% or 50.5%
		if regexp.MustCompile(`^\d+(\.\d+)?%$`).MatchString(v) {
			percentStr := strings.TrimSuffix(v, "%")
			percent, err := strconv.ParseFloat(percentStr, 32)
			if err != nil {
				return tess.Undefined(), fmt.Errorf("%w: invalid percent value '%s'", termerror.ErrInvalidStyleValue, v)
			}

			return tess.Percent(float32(percent)), nil
		}
	}

	return tess.Undefined(), fmt.Errorf("%w: '%v' is not recognised", termerror.ErrInvalidStyleValue, value)
}

func toTessEdges(all any, vertical any, horizontal any, top any, bottom any, left any, right any) (tess.Edges, error) {
	edges := tess.Edges{}

	if all != nil {
		v, err := toTessValue(all)
		if err != nil {
			return edges, err
		}

		edges.All = v
	}

	if vertical != nil {
		v, err := toTessValue(vertical)
		if err != nil {
			return edges, err
		}

		edges.Vertical = v
	}

	if horizontal != nil {
		v, err := toTessValue(horizontal)
		if err != nil {
			return edges, err
		}

		edges.Horizontal = v
	}

	if top != nil {
		v, err := toTessValue(top)
		if err != nil {
			return edges, err
		}

		edges.Top = v
	}

	if bottom != nil {
		v, err := toTessValue(bottom)
		if err != nil {
			return edges, err
		}

		edges.Bottom = v
	}

	if left != nil {
		v, err := toTessValue(left)
		if err != nil {
			return edges, err
		}

		edges.Left = v
	}

	if right != nil {
		v, err := toTessValue(right)
		if err != nil {
			return edges, err
		}

		edges.Right = v
	}

	return edges, nil
}
