package elements

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/AnatoleLucet/tess"
)

type BaseElementStyle struct {
	base *BaseElement

	ctx    context.Context
	rdrctx *RenderContext

	node *tess.Node

	translateX float32
	translateY float32
}

func NewBaseElementStyle(ctx context.Context, base *BaseElement) (*BaseElementStyle, error) {
	xyz, err := tess.NewNode()
	if err != nil {
		return nil, fmt.Errorf("failed to create tess node: %w", err)
	}

	return &BaseElementStyle{
		base: base,
		ctx:  ctx,
		node: xyz,
	}, nil
}

func (b *BaseElementStyle) xyz() *tess.Node {
	return b.node
}

func (b *BaseElementStyle) free() {
	b.node.Free()
	b.node = nil
}

func (b *BaseElement) ZIndex() (zindex int) {
	scheduleAccess(b.Self(), func() {
		b.mu.RLock()
		defer b.mu.RUnlock()
		zindex = b.zindexUnsafe()
	})

	return
}

func (b *BaseElement) zindexUnsafe() int {
	return b.zindex
}

func (b *BaseElement) SetZIndex(zIndex int) {
	scheduleUpdate(b.Self(), func() error {
		b.mu.Lock()
		defer b.mu.Unlock()

		if b.parent != nil {
			if err := b.parent.updateZIndex(b.Self(), b.zindex, zIndex); err != nil {
				return fmt.Errorf("%w: %w", ErrFailedToUpdateZIndex, err)
			}
		}

		b.zindex = zIndex
		return nil
	})
}

func (b *BaseElement) updateZIndex(child Element, oldz, newz int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// remove from old z-index
	children := b.children[oldz]
	i := slices.Index(children, child)
	if i == -1 {
		return fmt.Errorf("%w: child not found", ErrFailedToUpdateChild)
	}
	b.children[oldz] = slices.Delete(children, i, i+1)

	// add to new z-index
	b.children[newz] = append(b.children[newz], child)

	return nil
}

func (e *BaseElementStyle) SetWidth(width any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(width)
		if err != nil {
			return fmt.Errorf("%w: invalid width: %v", err, width)
		}

		err = e.xyz().SetWidth(value)
		if err != nil {
			return fmt.Errorf("failed to set width: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetMinWidth(minWidth any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(minWidth)
		if err != nil {
			return fmt.Errorf("%w: invalid min width: %v", err, minWidth)
		}

		err = e.xyz().SetMinWidth(value)
		if err != nil {
			return fmt.Errorf("failed to set min width: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetMaxWidth(maxWidth any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(maxWidth)
		if err != nil {
			return fmt.Errorf("%w: invalid max width: %v", err, maxWidth)
		}

		err = e.xyz().SetMaxWidth(value)
		if err != nil {
			return fmt.Errorf("failed to set max width: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetHeight(height any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(height)
		if err != nil {
			return fmt.Errorf("%w: invalid height: %v", err, height)
		}

		err = e.xyz().SetHeight(value)
		if err != nil {
			return fmt.Errorf("failed to set height: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetMinHeight(minHeight any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(minHeight)
		if err != nil {
			return fmt.Errorf("%w: invalid min height: %v", err, minHeight)
		}

		err = e.xyz().SetMinHeight(value)
		if err != nil {
			return fmt.Errorf("failed to set min height: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetMaxHeight(maxHeight any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(maxHeight)
		if err != nil {
			return fmt.Errorf("%w: invalid max height: %v", err, maxHeight)
		}

		err = e.xyz().SetMaxHeight(value)
		if err != nil {
			return fmt.Errorf("failed to set max height: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetTranslate(x, y float32) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		return e.setTranslateUnsafe(x, y)
	})
}

func (e *BaseElementStyle) setTranslateUnsafe(x, y float32) error {
	e.translateX = x
	e.translateY = y
	return nil
}

func (e *BaseElementStyle) SetTop(top any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(top)
		if err != nil {
			return fmt.Errorf("%w: invalid top: %v", err, top)
		}

		err = e.xyz().SetTop(value)
		if err != nil {
			return fmt.Errorf("failed to set top: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetBottom(bottom any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(bottom)
		if err != nil {
			return fmt.Errorf("%w: invalid bottom: %v", err, bottom)
		}

		err = e.xyz().SetBottom(value)
		if err != nil {
			return fmt.Errorf("failed to set bottom: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetLeft(left any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(left)
		if err != nil {
			return fmt.Errorf("%w: invalid left: %v", err, left)
		}

		err = e.xyz().SetLeft(value)
		if err != nil {
			return fmt.Errorf("failed to set left: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetRight(right any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(right)
		if err != nil {
			return fmt.Errorf("%w: invalid right: %v", err, right)
		}

		err = e.xyz().SetRight(value)
		if err != nil {
			return fmt.Errorf("failed to set right: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetPosition(position string) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		switch position {
		case "static":
			if err := e.xyz().SetPosition(tess.Static); err != nil {
				return fmt.Errorf("failed to set position: %w", err)
			}
		case "relative":
			if err := e.xyz().SetPosition(tess.Relative); err != nil {
				return fmt.Errorf("failed to set position: %w", err)
			}
		case "absolute":
			if err := e.xyz().SetPosition(tess.Absolute); err != nil {
				return fmt.Errorf("failed to set position: %w", err)
			}
		default:
			return fmt.Errorf("%w: failed to set position: '%s' is not recognised", ErrInvalidStyleValue, position)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetPaddingAll(padding any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(padding)
		if err != nil {
			return fmt.Errorf("%w: invalid padding: %v", err, padding)
		}

		err = e.xyz().SetPadding(tess.Edges{All: value})
		if err != nil {
			return fmt.Errorf("failed to set padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetPaddingVertical(padding any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(padding)
		if err != nil {
			return fmt.Errorf("%w: invalid vertical padding: %v", err, padding)
		}

		err = e.xyz().SetPadding(tess.Edges{Vertical: value})
		if err != nil {
			return fmt.Errorf("failed to set vertical padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetPaddingHorizontal(padding any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(padding)
		if err != nil {
			return fmt.Errorf("%w: invalid horizontal padding: %v", err, padding)
		}

		err = e.xyz().SetPadding(tess.Edges{Horizontal: value})
		if err != nil {
			return fmt.Errorf("failed to set horizontal padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetPaddingTop(padding any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(padding)
		if err != nil {
			return fmt.Errorf("%w: invalid top padding: %v", err, padding)
		}

		err = e.xyz().SetPadding(tess.Edges{Top: value})
		if err != nil {
			return fmt.Errorf("failed to set top padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetPaddingBottom(padding any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(padding)
		if err != nil {
			return fmt.Errorf("%w: invalid bottom padding: %v", err, padding)
		}

		err = e.xyz().SetPadding(tess.Edges{Bottom: value})
		if err != nil {
			return fmt.Errorf("failed to set bottom padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetPaddingLeft(padding any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(padding)
		if err != nil {
			return fmt.Errorf("%w: invalid left padding: %v", err, padding)
		}

		err = e.xyz().SetPadding(tess.Edges{Left: value})
		if err != nil {
			return fmt.Errorf("failed to set left padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetPaddingRight(padding any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(padding)
		if err != nil {
			return fmt.Errorf("%w: invalid right padding: %v", err, padding)
		}

		err = e.xyz().SetPadding(tess.Edges{Right: value})
		if err != nil {
			return fmt.Errorf("failed to set right padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetMarginAll(margin any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(margin)
		if err != nil {
			return fmt.Errorf("%w: invalid margin: %v", err, margin)
		}

		err = e.xyz().SetMargin(tess.Edges{All: value})
		if err != nil {
			return fmt.Errorf("failed to set margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetMarginVertical(margin any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(margin)
		if err != nil {
			return fmt.Errorf("%w: invalid vertical margin: %v", err, margin)
		}

		err = e.xyz().SetMargin(tess.Edges{Vertical: value})
		if err != nil {
			return fmt.Errorf("failed to set vertical margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetMarginHorizontal(margin any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(margin)
		if err != nil {
			return fmt.Errorf("%w: invalid horizontal margin: %v", err, margin)
		}

		err = e.xyz().SetMargin(tess.Edges{Horizontal: value})
		if err != nil {
			return fmt.Errorf("failed to set horizontal margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetMarginTop(margin any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(margin)
		if err != nil {
			return fmt.Errorf("%w: invalid top margin: %v", err, margin)
		}

		err = e.xyz().SetMargin(tess.Edges{Top: value})
		if err != nil {
			return fmt.Errorf("failed to set top margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetMarginBottom(margin any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(margin)
		if err != nil {
			return fmt.Errorf("%w: invalid bottom margin: %v", err, margin)
		}

		err = e.xyz().SetMargin(tess.Edges{Bottom: value})
		if err != nil {
			return fmt.Errorf("failed to set bottom margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetMarginLeft(margin any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(margin)
		if err != nil {
			return fmt.Errorf("%w: invalid left margin: %v", err, margin)
		}

		err = e.xyz().SetMargin(tess.Edges{Left: value})
		if err != nil {
			return fmt.Errorf("failed to set left margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetMarginRight(margin any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(margin)
		if err != nil {
			return fmt.Errorf("%w: invalid right margin: %v", err, margin)
		}

		err = e.xyz().SetMargin(tess.Edges{Right: value})
		if err != nil {
			return fmt.Errorf("failed to set right margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetDisplay(display string) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		switch display {
		case "none":
			if err := e.xyz().SetDisplay(tess.None); err != nil {
				return fmt.Errorf("failed to set display: %w", err)
			}
		case "flex":
			if err := e.xyz().SetDisplay(tess.Flex); err != nil {
				return fmt.Errorf("failed to set display: %w", err)
			}
		case "contents":
			if err := e.xyz().SetDisplay(tess.Contents); err != nil {
				return fmt.Errorf("failed to set display: %w", err)
			}
		default:
			return fmt.Errorf("%w: failed to set display: '%s' is not recognised", ErrInvalidStyleValue, display)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetAlignSelf(alignSelf string) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		switch alignSelf {
		case "auto":
			if err := e.xyz().SetAlignSelf(tess.AlignAuto); err != nil {
				return fmt.Errorf("failed to set align self: %w", err)
			}
		case "start":
			if err := e.xyz().SetAlignSelf(tess.AlignStart); err != nil {
				return fmt.Errorf("failed to set align self: %w", err)
			}
		case "end":
			if err := e.xyz().SetAlignSelf(tess.AlignEnd); err != nil {
				return fmt.Errorf("failed to set align self: %w", err)
			}
		case "center":
			if err := e.xyz().SetAlignSelf(tess.AlignCenter); err != nil {
				return fmt.Errorf("failed to set align self: %w", err)
			}
		case "stretch":
			if err := e.xyz().SetAlignSelf(tess.AlignStretch); err != nil {
				return fmt.Errorf("failed to set align self: %w", err)
			}
		case "baseline":
			if err := e.xyz().SetAlignSelf(tess.AlignBaseline); err != nil {
				return fmt.Errorf("failed to set align self: %w", err)
			}
		default:
			return fmt.Errorf("%w: failed to set align self: '%s' is not recognised", ErrInvalidStyleValue, alignSelf)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetAlignItems(alignItems string) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		switch alignItems {
		case "start":
			if err := e.xyz().SetAlignItems(tess.AlignStart); err != nil {
				return fmt.Errorf("failed to set align items: %w", err)
			}
		case "end":
			if err := e.xyz().SetAlignItems(tess.AlignEnd); err != nil {
				return fmt.Errorf("failed to set align items: %w", err)
			}
		case "center":
			if err := e.xyz().SetAlignItems(tess.AlignCenter); err != nil {
				return fmt.Errorf("failed to set align items: %w", err)
			}
		case "stretch":
			if err := e.xyz().SetAlignItems(tess.AlignStretch); err != nil {
				return fmt.Errorf("failed to set align items: %w", err)
			}
		case "baseline":
			if err := e.xyz().SetAlignItems(tess.AlignBaseline); err != nil {
				return fmt.Errorf("failed to set align items: %w", err)
			}
		default:
			return fmt.Errorf("%w: failed to set align items: '%s' is not recognised", ErrInvalidStyleValue, alignItems)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetAlignContent(alignContent string) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		switch alignContent {
		case "start":
			if err := e.xyz().SetAlignContent(tess.AlignStart); err != nil {
				return fmt.Errorf("failed to set align content: %w", err)
			}
		case "end":
			if err := e.xyz().SetAlignContent(tess.AlignEnd); err != nil {
				return fmt.Errorf("failed to set align content: %w", err)
			}
		case "center":
			if err := e.xyz().SetAlignContent(tess.AlignCenter); err != nil {
				return fmt.Errorf("failed to set align content: %w", err)
			}
		case "stretch":
			if err := e.xyz().SetAlignContent(tess.AlignStretch); err != nil {
				return fmt.Errorf("failed to set align content: %w", err)
			}
		case "baseline":
			if err := e.xyz().SetAlignContent(tess.AlignBaseline); err != nil {
				return fmt.Errorf("failed to set align content: %w", err)
			}
		default:
			return fmt.Errorf("%w: failed to set align content: '%s' is not recognised", ErrInvalidStyleValue, alignContent)
		}

		return nil

	})
}

func (e *BaseElementStyle) SetJustifyContent(justifyContent string) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		switch justifyContent {
		case "start":
			if err := e.xyz().SetJustifyContent(tess.JustifyStart); err != nil {
				return fmt.Errorf("failed to set justify content: %w", err)
			}
		case "end":
			if err := e.xyz().SetJustifyContent(tess.JustifyEnd); err != nil {
				return fmt.Errorf("failed to set justify content: %w", err)
			}
		case "center":
			if err := e.xyz().SetJustifyContent(tess.JustifyCenter); err != nil {
				return fmt.Errorf("failed to set justify content: %w", err)
			}
		case "space-between":
			if err := e.xyz().SetJustifyContent(tess.JustifySpaceBetween); err != nil {
				return fmt.Errorf("failed to set justify content: %w", err)
			}
		case "space-around":
			if err := e.xyz().SetJustifyContent(tess.JustifySpaceAround); err != nil {
				return fmt.Errorf("failed to set justify content: %w", err)
			}
		case "space-evenly":
			if err := e.xyz().SetJustifyContent(tess.JustifySpaceEvenly); err != nil {
				return fmt.Errorf("failed to set justify content: %w", err)
			}
		default:
			return fmt.Errorf("%w: failed to set justify content: '%s' is not recognised", ErrInvalidStyleValue, justifyContent)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetFlexDirection(flexDirection string) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		switch flexDirection {
		case "row":
			if err := e.xyz().SetFlexDirection(tess.Row); err != nil {
				return fmt.Errorf("failed to set flex direction: %w", err)
			}
		case "row-reverse":
			if err := e.xyz().SetFlexDirection(tess.RowReverse); err != nil {
				return fmt.Errorf("failed to set flex direction: %w", err)
			}
		case "column":
			if err := e.xyz().SetFlexDirection(tess.Column); err != nil {
				return fmt.Errorf("failed to set flex direction: %w", err)
			}
		case "column-reverse":
			if err := e.xyz().SetFlexDirection(tess.ColumnReverse); err != nil {
				return fmt.Errorf("failed to set flex direction: %w", err)
			}
		default:
			return fmt.Errorf("%w: failed to set flex direction: '%s' is not recognised", ErrInvalidStyleValue, flexDirection)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetFlexWrap(flexWrap string) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		switch flexWrap {
		case "nowrap":
			if err := e.xyz().SetFlexWrap(tess.NoWrap); err != nil {
				return fmt.Errorf("failed to set flex wrap: %w", err)
			}
		case "wrap":
			if err := e.xyz().SetFlexWrap(tess.Wrap); err != nil {
				return fmt.Errorf("failed to set flex wrap: %w", err)
			}
		case "wrap-reverse":
			if err := e.xyz().SetFlexWrap(tess.WrapReverse); err != nil {
				return fmt.Errorf("failed to set flex wrap: %w", err)
			}
		default:
			return fmt.Errorf("%w: failed to set flex wrap: '%s' is not recognised", ErrInvalidStyleValue, flexWrap)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetFlexGrow(flexGrow string) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		if flexGrow == "none" {
			return nil
		}

		grow, err := strconv.ParseFloat(flexGrow, 32)
		if err != nil {
			return fmt.Errorf("%w: invalid flex grow value '%s'", ErrInvalidStyleValue, flexGrow)
		}

		err = e.xyz().SetFlexGrow(float32(grow))
		if err != nil {
			return fmt.Errorf("failed to set flex grow: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetFlexShrink(flexShrink string) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		if flexShrink == "none" {
			return nil
		}

		shrink, err := strconv.ParseFloat(flexShrink, 32)
		if err != nil {
			return fmt.Errorf("%w: invalid flex shrink value '%s'", ErrInvalidStyleValue, flexShrink)
		}

		err = e.xyz().SetFlexShrink(float32(shrink))
		if err != nil {
			return fmt.Errorf("failed to set flex shrink: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetGapAll(gap any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(gap)
		if err != nil {
			return fmt.Errorf("%w: invalid gap: %v", err, gap)
		}

		err = e.xyz().SetGap(tess.Gap{All: value})
		if err != nil {
			return fmt.Errorf("failed to set gap: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetGapRow(gap any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(gap)
		if err != nil {
			return fmt.Errorf("%w: invalid row gap: %v", err, gap)
		}

		err = e.xyz().SetGap(tess.Gap{Row: value})
		if err != nil {
			return fmt.Errorf("failed to set row gap: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetGapColumn(gap any) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		value, err := toTessValue(gap)
		if err != nil {
			return fmt.Errorf("%w: invalid column gap: %v", err, gap)
		}

		err = e.xyz().SetGap(tess.Gap{Column: value})
		if err != nil {
			return fmt.Errorf("failed to set column gap: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) SetOverflow(overflow string) {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		switch overflow {
		case "visible":
			if err := e.xyz().SetOverflow(tess.Visible); err != nil {
				return fmt.Errorf("failed to set overflow: %w", err)
			}
		case "hidden":
			if err := e.xyz().SetOverflow(tess.Hidden); err != nil {
				return fmt.Errorf("failed to set overflow: %w", err)
			}
		default:
			return fmt.Errorf("%w: failed to set overflow: '%s' is not recognised", ErrInvalidStyleValue, overflow)
		}

		return nil
	})
}

func (b *BaseElement) UnsetZIndex() {
	b.SetZIndex(0)
}

func (e *BaseElementStyle) UnsetWidth() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetWidth(tess.Undefined())
		if err != nil {
			return fmt.Errorf("failed to unset width: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetMinWidth() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetMinWidth(tess.Undefined())
		if err != nil {
			return fmt.Errorf("failed to unset min width: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetMaxWidth() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetMaxWidth(tess.Undefined())
		if err != nil {
			return fmt.Errorf("failed to unset max width: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetHeight() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetHeight(tess.Undefined())
		if err != nil {
			return fmt.Errorf("failed to unset height: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetMinHeight() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetMinHeight(tess.Undefined())
		if err != nil {
			return fmt.Errorf("failed to unset min height: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetMaxHeight() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetMaxHeight(tess.Undefined())
		if err != nil {
			return fmt.Errorf("failed to unset max height: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetTranslate() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		e.translateX = 0
		e.translateY = 0
		return nil
	})
}

func (e *BaseElementStyle) UnsetTop() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetTop(tess.Undefined())
		if err != nil {
			return fmt.Errorf("failed to unset top: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetBottom() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetBottom(tess.Undefined())
		if err != nil {
			return fmt.Errorf("failed to unset bottom: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetLeft() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetLeft(tess.Undefined())
		if err != nil {
			return fmt.Errorf("failed to unset left: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetRight() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetRight(tess.Undefined())
		if err != nil {
			return fmt.Errorf("failed to unset right: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetPosition() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetPosition(tess.Static)
		if err != nil {
			return fmt.Errorf("failed to unset position: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetPaddingAll() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetPadding(tess.Edges{All: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetPaddingVertical() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetPadding(tess.Edges{Vertical: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset vertical padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetPaddingHorizontal() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetPadding(tess.Edges{Horizontal: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset horizontal padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetPaddingTop() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetPadding(tess.Edges{Top: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset top padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetPaddingBottom() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetPadding(tess.Edges{Bottom: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset bottom padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetPaddingLeft() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetPadding(tess.Edges{Left: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset left padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetPaddingRight() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetPadding(tess.Edges{Right: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset right padding: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetMarginAll() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetMargin(tess.Edges{All: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetMarginVertical() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetMargin(tess.Edges{Vertical: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset vertical margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetMarginHorizontal() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetMargin(tess.Edges{Horizontal: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset horizontal margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetMarginTop() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetMargin(tess.Edges{Top: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset top margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetMarginBottom() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetMargin(tess.Edges{Bottom: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset bottom margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetMarginLeft() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetMargin(tess.Edges{Left: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset left margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetMarginRight() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetMargin(tess.Edges{Right: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset right margin: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetDisplay() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetDisplay(tess.Flex)
		if err != nil {
			return fmt.Errorf("failed to unset display: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetAlignSelf() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetAlignSelf(tess.AlignAuto)
		if err != nil {
			return fmt.Errorf("failed to unset align self: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetAlignItems() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetAlignItems(tess.AlignStretch)
		if err != nil {
			return fmt.Errorf("failed to unset align items: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetAlignContent() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetAlignContent(tess.AlignStretch)
		if err != nil {
			return fmt.Errorf("failed to unset align content: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetJustifyContent() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetJustifyContent(tess.JustifyStart)
		if err != nil {
			return fmt.Errorf("failed to unset justify content: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetFlexDirection() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetFlexDirection(tess.Column)
		if err != nil {
			return fmt.Errorf("failed to unset flex direction: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetFlexWrap() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetFlexWrap(tess.NoWrap)
		if err != nil {
			return fmt.Errorf("failed to unset flex wrap: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetFlexGrow() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetFlexGrow(0)
		if err != nil {
			return fmt.Errorf("failed to unset flex grow: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetFlexShrink() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetFlexShrink(1)
		if err != nil {
			return fmt.Errorf("failed to unset flex shrink: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetGapAll() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetGap(tess.Gap{All: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset gap: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetGapRow() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetGap(tess.Gap{Row: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset row gap: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetGapColumn() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetGap(tess.Gap{Column: tess.Undefined()})
		if err != nil {
			return fmt.Errorf("failed to unset column gap: %w", err)
		}

		return nil
	})
}

func (e *BaseElementStyle) UnsetOverflow() {
	scheduleUpdate(e.base.Self(), func() error {
		e.base.mu.Lock()
		defer e.base.mu.Unlock()

		err := e.xyz().SetOverflow(tess.Visible)
		if err != nil {
			return fmt.Errorf("failed to unset overflow: %w", err)
		}

		return nil
	})
}

func toTessValue(value any) (tess.Value, error) {
	switch v := value.(type) {
	case nil:
		return tess.Undefined(), nil
	case int, int32, int64, uint, uint32, uint64:
		return tess.Point(float32(v.(int))), nil
	case float32:
		return tess.Point(v), nil
	case float64:
		return tess.Point(float32(v)), nil
	case string:
		switch v {
		case "", "undefined", "none":
			return tess.Undefined(), nil
		case "auto":
			return tess.Auto(), nil
		case "max-content":
			return tess.MaxContent(), nil
		case "fit-content":
			return tess.FitContent(), nil
		case "stretch":
			return tess.Stretch(), nil
		}

		// 100pt or 100.5pt
		if regexp.MustCompile(`^\d+(\.\d+)?pt$`).MatchString(v) {
			pointStr := strings.TrimSuffix(v, "pt")
			points, err := strconv.ParseFloat(pointStr, 32)
			if err != nil {
				return tess.Undefined(), fmt.Errorf("%w: invalid point value '%s'", ErrInvalidStyleValue, v)
			}

			return tess.Point(float32(points)), nil
		}

		// 100% or 50.5%
		if regexp.MustCompile(`^\d+(\.\d+)?%$`).MatchString(v) {
			percentStr := strings.TrimSuffix(v, "%")
			percent, err := strconv.ParseFloat(percentStr, 32)
			if err != nil {
				return tess.Undefined(), fmt.Errorf("%w: invalid percent value '%s'", ErrInvalidStyleValue, v)
			}

			return tess.Percent(float32(percent)), nil
		}
	}

	return tess.Undefined(), fmt.Errorf("%w: '%v' is not recognised", ErrInvalidStyleValue, value)
}
