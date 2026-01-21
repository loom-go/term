package elements

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/AnatoleLucet/go-opentui"
	"github.com/AnatoleLucet/loom-term/core/types"
	"github.com/AnatoleLucet/tess"
)

type BaseElement struct {
	mu sync.RWMutex

	ctx types.RenderContext

	xyz *tess.Node

	parent   types.Element
	children map[int][]types.Element // map[zindex]children

	zindex int

	destroyed bool

	nextListenerID        int
	mouseMoveListeners    map[int]func(*types.EventMouse)
	mouseEnterListeners   map[int]func(*types.EventMouse)
	mouseLeaveListeners   map[int]func(*types.EventMouse)
	mousePressListeners   map[int]func(*types.EventMouse)
	mouseReleaseListeners map[int]func(*types.EventMouse)
	mouseDragListeners    map[int]func(*types.EventMouse)
	mouseScrollListeners  map[int]func(*types.EventMouse)
}

func NewElement(ctx types.RenderContext) (*BaseElement, error) {
	xyz, err := tess.NewNode()
	if err != nil {
		return nil, err
	}

	return &BaseElement{
		ctx:                   ctx,
		xyz:                   xyz,
		parent:                nil,
		children:              make(map[int][]types.Element),
		zindex:                0,
		destroyed:             false,
		nextListenerID:        0,
		mouseMoveListeners:    make(map[int]func(*types.EventMouse)),
		mouseEnterListeners:   make(map[int]func(*types.EventMouse)),
		mouseLeaveListeners:   make(map[int]func(*types.EventMouse)),
		mousePressListeners:   make(map[int]func(*types.EventMouse)),
		mouseReleaseListeners: make(map[int]func(*types.EventMouse)),
		mouseScrollListeners:  make(map[int]func(*types.EventMouse)),
		mouseDragListeners:    make(map[int]func(*types.EventMouse)),
	}, nil
}

func (e *BaseElement) XYZ() *tess.Node {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.xyz
}

func (e *BaseElement) Layout() *tess.Layout {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.xyz.GetLayout()
}

func (e *BaseElement) ZIndex() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.zindex
}

func (e *BaseElement) Parent() types.Element {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.parent
}

func (e *BaseElement) SetParent(parent types.Element) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	e.parent = parent
	return nil
}

func (e *BaseElement) Children() iter.Seq[types.Element] {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.childrenUnsafe()
}

func (e *BaseElement) childrenUnsafe() iter.Seq[types.Element] {
	return func(yield func(types.Element) bool) {
		zindexes := slices.Sorted(maps.Keys(e.children))
		for _, zindex := range zindexes {
			for _, child := range e.children[zindex] {
				if !yield(child) {
					return
				}
			}
		}
	}
}

func (e *BaseElement) AppendChild(child types.Element) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	return e.appendChild(child)
}

func (e *BaseElement) appendChild(child types.Element) error {
	e.xyz.AppendChild(child.XYZ())

	zindex := child.ZIndex()
	e.children[zindex] = append(e.children[zindex], child)

	child.SetParent(e)
	return nil
}

func (e *BaseElement) RemoveChild(child types.Element) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	return e.removeChild(child)
}

func (e *BaseElement) removeChild(child types.Element) error {
	zindex := child.ZIndex()
	children, ok := e.children[zindex]
	if !ok {
		return errors.New("child not found") // todo: better error
	}

	i := slices.Index(children, child)
	if i == -1 {
		return errors.New("child not found") // todo: better error
	}

	e.xyz.RemoveChild(child.XYZ())
	e.children[zindex] = slices.Delete(children, i, i+1)

	child.SetParent(nil)
	return nil
}

func (e *BaseElement) OnMouseMove(handler func(*types.EventMouse)) (remove func()) {
	e.mu.Lock()
	id := e.nextListenerID
	e.nextListenerID++

	e.mouseMoveListeners[id] = handler
	e.mu.Unlock()

	return func() {
		e.mu.Lock()
		delete(e.mouseMoveListeners, id)
		e.mu.Unlock()
	}
}

func (e *BaseElement) OnMouseEnter(handler func(*types.EventMouse)) (remove func()) {
	e.mu.Lock()
	id := e.nextListenerID
	e.nextListenerID++

	e.mouseEnterListeners[id] = handler
	e.mu.Unlock()

	return func() {
		e.mu.Lock()
		delete(e.mouseEnterListeners, id)
		e.mu.Unlock()
	}
}

func (e *BaseElement) OnMouseLeave(handler func(*types.EventMouse)) (remove func()) {
	e.mu.Lock()
	id := e.nextListenerID
	e.nextListenerID++

	e.mouseLeaveListeners[id] = handler
	e.mu.Unlock()

	return func() {
		e.mu.Lock()
		delete(e.mouseLeaveListeners, id)
		e.mu.Unlock()
	}
}

func (e *BaseElement) OnMousePress(handler func(*types.EventMouse)) (remove func()) {
	e.mu.Lock()
	id := e.nextListenerID
	e.nextListenerID++

	e.mousePressListeners[id] = handler
	e.mu.Unlock()

	return func() {
		e.mu.Lock()
		delete(e.mousePressListeners, id)
		e.mu.Unlock()
	}
}

func (e *BaseElement) OnMouseRelease(handler func(*types.EventMouse)) (remove func()) {
	e.mu.Lock()
	id := e.nextListenerID
	e.nextListenerID++

	e.mouseReleaseListeners[id] = handler
	e.mu.Unlock()

	return func() {
		e.mu.Lock()
		delete(e.mouseReleaseListeners, id)
		e.mu.Unlock()
	}
}

func (e *BaseElement) OnMouseScroll(handler func(*types.EventMouse)) (remove func()) {
	e.mu.Lock()
	id := e.nextListenerID
	e.nextListenerID++

	e.mouseScrollListeners[id] = handler
	e.mu.Unlock()

	return func() {
		e.mu.Lock()
		delete(e.mouseScrollListeners, id)
		e.mu.Unlock()
	}
}

func (e *BaseElement) OnMouseDrag(handler func(*types.EventMouse)) (remove func()) {
	e.mu.Lock()
	id := e.nextListenerID
	e.nextListenerID++

	e.mouseDragListeners[id] = handler
	e.mu.Unlock()

	return func() {
		e.mu.Lock()
		delete(e.mouseDragListeners, id)
		e.mu.Unlock()
	}
}

func (e *BaseElement) HandleMouseMove(event *types.EventMouse) error {
	for _, handler := range e.mouseMoveListeners {
		handler(event)
	}

	if e.parent != nil && event.ShouldPropagate() {
		return e.parent.HandleMouseMove(event)
	}

	return nil
}

func (e *BaseElement) HandleMouseEnter(event *types.EventMouse) error {
	for _, handler := range e.mouseEnterListeners {
		handler(event)
	}

	if e.parent != nil && event.ShouldPropagate() {
		return e.parent.HandleMouseEnter(event)
	}

	return nil
}

func (e *BaseElement) HandleMouseLeave(event *types.EventMouse) error {
	for _, handler := range e.mouseLeaveListeners {
		handler(event)
	}

	if e.parent != nil && event.ShouldPropagate() {
		return e.parent.HandleMouseLeave(event)
	}

	return nil
}

func (e *BaseElement) HandleMousePress(event *types.EventMouse) error {
	for _, handler := range e.mousePressListeners {
		handler(event)
	}

	if e.parent != nil && event.ShouldPropagate() {
		return e.parent.HandleMousePress(event)
	}

	return nil
}

func (e *BaseElement) HandleMouseRelease(event *types.EventMouse) error {
	for _, handler := range e.mouseReleaseListeners {
		handler(event)
	}

	if e.parent != nil && event.ShouldPropagate() {
		return e.parent.HandleMouseRelease(event)
	}

	return nil
}

func (e *BaseElement) HandleMouseScroll(event *types.EventMouse) error {
	for _, handler := range e.mouseScrollListeners {
		handler(event)
	}

	if e.parent != nil && event.ShouldPropagate() {
		return e.parent.HandleMouseScroll(event)
	}

	return nil
}

func (e *BaseElement) HandleMouseDrag(event *types.EventMouse) error {
	for _, handler := range e.mouseDragListeners {
		handler(event)
	}

	if e.parent != nil && event.ShouldPropagate() {
		return e.parent.HandleMouseDrag(event)
	}

	return nil
}

func (e *BaseElement) Paint(buffer *opentui.Buffer, x, y float32) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if err := e.guardPaint(); err != nil {
		return err
	}

	if err := e.paint(buffer, x, y); err != nil {
		return err
	}

	return e.paintChildren(buffer, x, y)
}

func (e *BaseElement) paint(_ *opentui.Buffer, x, y float32) error {
	layout := e.xyz.GetLayout()
	e.ctx.AddToHitGrid(e, int(x), int(y), int(layout.Width()), int(layout.Height()))
	return nil
}

func (e *BaseElement) paintChildren(buffer *opentui.Buffer, x, y float32) error {
	if len(e.children) == 0 {
		return nil
	}

	return e.withClip(buffer, x, y, func() error {
		for child := range e.childrenUnsafe() {
			l := child.Layout()

			err := child.Paint(buffer, x+l.Left(), y+l.Top())
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (e *BaseElement) withClip(buffer *opentui.Buffer, x, y float32, fn func() error) error {
	layout := e.xyz.GetLayout()
	overflow := e.xyz.GetOverflow()

	if overflow == tess.Hidden {
		e.ctx.PushHitGridScissors(
			int(x),
			int(y),
			int(layout.Width()),
			int(layout.Height()),
		)

		buffer.PushScissorRect(
			int32(x+layout.Padding().Left()),
			int32(y+layout.Padding().Top()),
			uint32(layout.Width()-layout.Padding().Left()-layout.Padding().Right()),
			uint32(layout.Height()-layout.Padding().Top()-layout.Padding().Bottom()),
		)
	}

	err := fn()

	if overflow == tess.Hidden {
		e.ctx.PopHitGridScissors()
		buffer.PopScissorRect()
	}

	return err
}

func (e *BaseElement) Destroy() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.destroyed {
		return nil
	}
	e.destroyed = true

	for child := range e.childrenUnsafe() {
		err := child.Destroy()
		if err != nil {
			return err
		}
	}

	e.xyz.Free()
	return nil
}

func (e *BaseElement) SetZIndex(zIndex int) error {
	if err := e.guardUpdate(); err != nil {
		return err
	}

	if e.parent != nil {
		err := e.parent.RemoveChild(e)
		if err != nil {
			return fmt.Errorf("failed to update z-index: %w", err) // todo: error
		}
	}

	e.zindex = zIndex

	if e.parent != nil {
		err := e.parent.AppendChild(e)
		if err != nil {
			return fmt.Errorf("failed to update z-index: %w", err) // todo: error
		}
	}

	return nil
}

func (e *BaseElement) SetWidth(width any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(width)
	if err != nil {
		return fmt.Errorf("%w: invalid width: %v", err, width)
	}

	err = e.xyz.SetWidth(value)
	if err != nil {
		return fmt.Errorf("failed to set width: %w", err)
	}

	return nil
}

func (e *BaseElement) SetMinWidth(minWidth any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(minWidth)
	if err != nil {
		return fmt.Errorf("%w: invalid min width: %v", err, minWidth)
	}

	err = e.xyz.SetMinWidth(value)
	if err != nil {
		return fmt.Errorf("failed to set min width: %w", err)
	}

	return nil
}

func (e *BaseElement) SetMaxWidth(maxWidth any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(maxWidth)
	if err != nil {
		return fmt.Errorf("%w: invalid max width: %v", err, maxWidth)
	}

	err = e.xyz.SetMaxWidth(value)
	if err != nil {
		return fmt.Errorf("failed to set max width: %w", err)
	}

	return nil
}

func (e *BaseElement) SetHeight(height any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(height)
	if err != nil {
		return fmt.Errorf("%w: invalid height: %v", err, height)
	}

	err = e.xyz.SetHeight(value)
	if err != nil {
		return fmt.Errorf("failed to set height: %w", err)
	}

	return nil
}

func (e *BaseElement) SetMinHeight(minHeight any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(minHeight)
	if err != nil {
		return fmt.Errorf("%w: invalid min height: %v", err, minHeight)
	}

	err = e.xyz.SetMinHeight(value)
	if err != nil {
		return fmt.Errorf("failed to set min height: %w", err)
	}

	return nil
}

func (e *BaseElement) SetMaxHeight(maxHeight any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(maxHeight)
	if err != nil {
		return fmt.Errorf("%w: invalid max height: %v", err, maxHeight)
	}

	err = e.xyz.SetMaxHeight(value)
	if err != nil {
		return fmt.Errorf("failed to set max height: %w", err)
	}

	return nil
}

func (e *BaseElement) SetTop(top any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(top)
	if err != nil {
		return fmt.Errorf("%w: invalid top: %v", err, top)
	}

	err = e.xyz.SetTop(value)
	if err != nil {
		return fmt.Errorf("failed to set top: %w", err)
	}

	return nil
}

func (e *BaseElement) SetBottom(bottom any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(bottom)
	if err != nil {
		return fmt.Errorf("%w: invalid bottom: %v", err, bottom)
	}

	err = e.xyz.SetBottom(value)
	if err != nil {
		return fmt.Errorf("failed to set bottom: %w", err)
	}

	return nil
}

func (e *BaseElement) SetLeft(left any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(left)
	if err != nil {
		return fmt.Errorf("%w: invalid left: %v", err, left)
	}

	err = e.xyz.SetLeft(value)
	if err != nil {
		return fmt.Errorf("failed to set left: %w", err)
	}

	return nil
}

func (e *BaseElement) SetRight(right any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(right)
	if err != nil {
		return fmt.Errorf("%w: invalid right: %v", err, right)
	}

	err = e.xyz.SetRight(value)
	if err != nil {
		return fmt.Errorf("failed to set right: %w", err)
	}

	return nil
}

func (e *BaseElement) SetPosition(position string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	switch position {
	case "static":
		if err := e.xyz.SetPosition(tess.Static); err != nil {
			return fmt.Errorf("failed to set position: %w", err)
		}
	case "relative":
		if err := e.xyz.SetPosition(tess.Relative); err != nil {
			return fmt.Errorf("failed to set position: %w", err)
		}
	case "absolute":
		if err := e.xyz.SetPosition(tess.Absolute); err != nil {
			return fmt.Errorf("failed to set position: %w", err)
		}
	}

	return fmt.Errorf("%w: '%s' is not recognised", types.ErrInvalidStyleValue, position)
}

func (e *BaseElement) SetPaddingAll(padding any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(padding)
	if err != nil {
		return fmt.Errorf("%w: invalid padding: %v", err, padding)
	}

	err = e.xyz.SetPadding(tess.Edges{All: value})
	if err != nil {
		return fmt.Errorf("failed to set padding: %w", err)
	}

	return nil
}

func (e *BaseElement) SetPaddingVertical(padding any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(padding)
	if err != nil {
		return fmt.Errorf("%w: invalid vertical padding: %v", err, padding)
	}

	err = e.xyz.SetPadding(tess.Edges{Vertical: value})
	if err != nil {
		return fmt.Errorf("failed to set vertical padding: %w", err)
	}

	return nil
}

func (e *BaseElement) SetPaddingHorizontal(padding any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(padding)
	if err != nil {
		return fmt.Errorf("%w: invalid horizontal padding: %v", err, padding)
	}

	err = e.xyz.SetPadding(tess.Edges{Horizontal: value})
	if err != nil {
		return fmt.Errorf("failed to set horizontal padding: %w", err)
	}

	return nil
}

func (e *BaseElement) SetPaddingTop(padding any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(padding)
	if err != nil {
		return fmt.Errorf("%w: invalid top padding: %v", err, padding)
	}

	err = e.xyz.SetPadding(tess.Edges{Top: value})
	if err != nil {
		return fmt.Errorf("failed to set top padding: %w", err)
	}

	return nil
}

func (e *BaseElement) SetPaddingBottom(padding any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(padding)
	if err != nil {
		return fmt.Errorf("%w: invalid bottom padding: %v", err, padding)
	}

	err = e.xyz.SetPadding(tess.Edges{Bottom: value})
	if err != nil {
		return fmt.Errorf("failed to set bottom padding: %w", err)
	}

	return nil
}

func (e *BaseElement) SetPaddingLeft(padding any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(padding)
	if err != nil {
		return fmt.Errorf("%w: invalid left padding: %v", err, padding)
	}

	err = e.xyz.SetPadding(tess.Edges{Left: value})
	if err != nil {
		return fmt.Errorf("failed to set left padding: %w", err)
	}

	return nil
}

func (e *BaseElement) SetPaddingRight(padding any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(padding)
	if err != nil {
		return fmt.Errorf("%w: invalid right padding: %v", err, padding)
	}

	err = e.xyz.SetPadding(tess.Edges{Right: value})
	if err != nil {
		return fmt.Errorf("failed to set right padding: %w", err)
	}

	return nil
}

func (e *BaseElement) SetMarginAll(margin any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(margin)
	if err != nil {
		return fmt.Errorf("%w: invalid margin: %v", err, margin)
	}

	err = e.xyz.SetMargin(tess.Edges{All: value})
	if err != nil {
		return fmt.Errorf("failed to set margin: %w", err)
	}

	return nil
}

func (e *BaseElement) SetMarginVertical(margin any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(margin)
	if err != nil {
		return fmt.Errorf("%w: invalid vertical margin: %v", err, margin)
	}

	err = e.xyz.SetMargin(tess.Edges{Vertical: value})
	if err != nil {
		return fmt.Errorf("failed to set vertical margin: %w", err)
	}

	return nil
}

func (e *BaseElement) SetMarginHorizontal(margin any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(margin)
	if err != nil {
		return fmt.Errorf("%w: invalid horizontal margin: %v", err, margin)
	}

	err = e.xyz.SetMargin(tess.Edges{Horizontal: value})
	if err != nil {
		return fmt.Errorf("failed to set horizontal margin: %w", err)
	}

	return nil
}

func (e *BaseElement) SetMarginTop(margin any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(margin)
	if err != nil {
		return fmt.Errorf("%w: invalid top margin: %v", err, margin)
	}

	err = e.xyz.SetMargin(tess.Edges{Top: value})
	if err != nil {
		return fmt.Errorf("failed to set top margin: %w", err)
	}

	return nil
}

func (e *BaseElement) SetMarginBottom(margin any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(margin)
	if err != nil {
		return fmt.Errorf("%w: invalid bottom margin: %v", err, margin)
	}

	err = e.xyz.SetMargin(tess.Edges{Bottom: value})
	if err != nil {
		return fmt.Errorf("failed to set bottom margin: %w", err)
	}

	return nil
}

func (e *BaseElement) SetMarginLeft(margin any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(margin)
	if err != nil {
		return fmt.Errorf("%w: invalid left margin: %v", err, margin)
	}

	err = e.xyz.SetMargin(tess.Edges{Left: value})
	if err != nil {
		return fmt.Errorf("failed to set left margin: %w", err)
	}

	return nil
}

func (e *BaseElement) SetMarginRight(margin any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(margin)
	if err != nil {
		return fmt.Errorf("%w: invalid right margin: %v", err, margin)
	}

	err = e.xyz.SetMargin(tess.Edges{Right: value})
	if err != nil {
		return fmt.Errorf("failed to set right margin: %w", err)
	}

	return nil
}

func (e *BaseElement) SetDisplay(display string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	switch display {
	case "none":
		if err := e.xyz.SetDisplay(tess.None); err != nil {
			return fmt.Errorf("failed to set display: %w", err)
		}
	case "flex":
		if err := e.xyz.SetDisplay(tess.Flex); err != nil {
			return fmt.Errorf("failed to set display: %w", err)
		}
	case "contents":
		if err := e.xyz.SetDisplay(tess.Contents); err != nil {
			return fmt.Errorf("failed to set display: %w", err)
		}
	}

	return fmt.Errorf("%w: '%s' is not recognised", types.ErrInvalidStyleValue, display)
}

func (e *BaseElement) SetAlignSelf(alignSelf string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	switch alignSelf {
	case "auto":
		if err := e.xyz.SetAlignSelf(tess.AlignAuto); err != nil {
			return fmt.Errorf("failed to set align self: %w", err)
		}
	case "start":
		if err := e.xyz.SetAlignSelf(tess.AlignStart); err != nil {
			return fmt.Errorf("failed to set align self: %w", err)
		}
	case "end":
		if err := e.xyz.SetAlignSelf(tess.AlignEnd); err != nil {
			return fmt.Errorf("failed to set align self: %w", err)
		}
	case "center":
		if err := e.xyz.SetAlignSelf(tess.AlignCenter); err != nil {
			return fmt.Errorf("failed to set align self: %w", err)
		}
	case "stretch":
		if err := e.xyz.SetAlignSelf(tess.AlignStretch); err != nil {
			return fmt.Errorf("failed to set align self: %w", err)
		}
	case "baseline":
		if err := e.xyz.SetAlignSelf(tess.AlignBaseline); err != nil {
			return fmt.Errorf("failed to set align self: %w", err)
		}
	}

	return fmt.Errorf("%w: '%s' is not recognised", types.ErrInvalidStyleValue, alignSelf)
}

func (e *BaseElement) SetAlignItems(alignItems string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	switch alignItems {
	case "start":
		if err := e.xyz.SetAlignItems(tess.AlignStart); err != nil {
			return fmt.Errorf("failed to set align items: %w", err)
		}
	case "end":
		if err := e.xyz.SetAlignItems(tess.AlignEnd); err != nil {
			return fmt.Errorf("failed to set align items: %w", err)
		}
	case "center":
		if err := e.xyz.SetAlignItems(tess.AlignCenter); err != nil {
			return fmt.Errorf("failed to set align items: %w", err)
		}
	case "stretch":
		if err := e.xyz.SetAlignItems(tess.AlignStretch); err != nil {
			return fmt.Errorf("failed to set align items: %w", err)
		}
	case "baseline":
		if err := e.xyz.SetAlignItems(tess.AlignBaseline); err != nil {
			return fmt.Errorf("failed to set align items: %w", err)
		}
	}

	return fmt.Errorf("%w: '%s' is not recognised", types.ErrInvalidStyleValue, alignItems)
}

func (e *BaseElement) SetAlignContent(alignContent string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	switch alignContent {
	case "start":
		if err := e.xyz.SetAlignContent(tess.AlignStart); err != nil {
			return fmt.Errorf("failed to set align content: %w", err)
		}
	case "end":
		if err := e.xyz.SetAlignContent(tess.AlignEnd); err != nil {
			return fmt.Errorf("failed to set align content: %w", err)
		}
	case "center":
		if err := e.xyz.SetAlignContent(tess.AlignCenter); err != nil {
			return fmt.Errorf("failed to set align content: %w", err)
		}
	case "stretch":
		if err := e.xyz.SetAlignContent(tess.AlignStretch); err != nil {
			return fmt.Errorf("failed to set align content: %w", err)
		}
	case "baseline":
		if err := e.xyz.SetAlignContent(tess.AlignBaseline); err != nil {
			return fmt.Errorf("failed to set align content: %w", err)
		}
	}

	return fmt.Errorf("%w: '%s' is not recognised", types.ErrInvalidStyleValue, alignContent)
}

func (e *BaseElement) SetJustifyContent(justifyContent string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	switch justifyContent {
	case "start":
		if err := e.xyz.SetJustifyContent(tess.JustifyStart); err != nil {
			return fmt.Errorf("failed to set justify content: %w", err)
		}
	case "end":
		if err := e.xyz.SetJustifyContent(tess.JustifyEnd); err != nil {
			return fmt.Errorf("failed to set justify content: %w", err)
		}
	case "center":
		if err := e.xyz.SetJustifyContent(tess.JustifyCenter); err != nil {
			return fmt.Errorf("failed to set justify content: %w", err)
		}
	case "space-between":
		if err := e.xyz.SetJustifyContent(tess.JustifySpaceBetween); err != nil {
			return fmt.Errorf("failed to set justify content: %w", err)
		}
	case "space-around":
		if err := e.xyz.SetJustifyContent(tess.JustifySpaceAround); err != nil {
			return fmt.Errorf("failed to set justify content: %w", err)
		}
	case "space-evenly":
		if err := e.xyz.SetJustifyContent(tess.JustifySpaceEvenly); err != nil {
			return fmt.Errorf("failed to set justify content: %w", err)
		}
	}

	return fmt.Errorf("%w: '%s' is not recognised", types.ErrInvalidStyleValue, justifyContent)
}

func (e *BaseElement) SetFlexDirection(flexDirection string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	switch flexDirection {
	case "row":
		if err := e.xyz.SetFlexDirection(tess.Row); err != nil {
			return fmt.Errorf("failed to set flex direction: %w", err)
		}
	case "row-reverse":
		if err := e.xyz.SetFlexDirection(tess.RowReverse); err != nil {
			return fmt.Errorf("failed to set flex direction: %w", err)
		}
	case "column":
		if err := e.xyz.SetFlexDirection(tess.Column); err != nil {
			return fmt.Errorf("failed to set flex direction: %w", err)
		}
	case "column-reverse":
		if err := e.xyz.SetFlexDirection(tess.ColumnReverse); err != nil {
			return fmt.Errorf("failed to set flex direction: %w", err)
		}
	}

	return fmt.Errorf("%w: '%s' is not recognised", types.ErrInvalidStyleValue, flexDirection)
}

func (e *BaseElement) SetFlexWrap(flexWrap string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	switch flexWrap {
	case "nowrap":
		if err := e.xyz.SetFlexWrap(tess.NoWrap); err != nil {
			return fmt.Errorf("failed to set flex wrap: %w", err)
		}
	case "wrap":
		if err := e.xyz.SetFlexWrap(tess.Wrap); err != nil {
			return fmt.Errorf("failed to set flex wrap: %w", err)
		}
	case "wrap-reverse":
		if err := e.xyz.SetFlexWrap(tess.WrapReverse); err != nil {
			return fmt.Errorf("failed to set flex wrap: %w", err)
		}
	}

	return fmt.Errorf("%w: '%s' is not recognised", types.ErrInvalidStyleValue, flexWrap)
}

func (e *BaseElement) SetFlexGrow(flexGrow string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	if flexGrow == "none" {
		return nil
	}

	grow, err := strconv.ParseFloat(flexGrow, 32)
	if err != nil {
		return fmt.Errorf("%w: invalid flex grow value '%s'", types.ErrInvalidStyleValue, flexGrow)
	}

	err = e.xyz.SetFlexGrow(float32(grow))
	if err != nil {
		return fmt.Errorf("failed to set flex grow: %w", err)
	}

	return nil
}

func (e *BaseElement) SetFlexShrink(flexShrink string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	if flexShrink == "none" {
		return nil
	}

	shrink, err := strconv.ParseFloat(flexShrink, 32)
	if err != nil {
		return fmt.Errorf("%w: invalid flex shrink value '%s'", types.ErrInvalidStyleValue, flexShrink)
	}

	err = e.xyz.SetFlexShrink(float32(shrink))
	if err != nil {
		return fmt.Errorf("failed to set flex shrink: %w", err)
	}

	return nil
}

func (e *BaseElement) SetGapAll(gap any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(gap)
	if err != nil {
		return fmt.Errorf("%w: invalid gap: %v", err, gap)
	}

	err = e.xyz.SetGap(tess.Gap{All: value})
	if err != nil {
		return fmt.Errorf("failed to set gap: %w", err)
	}

	return nil
}

func (e *BaseElement) SetGapRow(gap any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(gap)
	if err != nil {
		return fmt.Errorf("%w: invalid row gap: %v", err, gap)
	}

	err = e.xyz.SetGap(tess.Gap{Row: value})
	if err != nil {
		return fmt.Errorf("failed to set row gap: %w", err)
	}

	return nil
}

func (e *BaseElement) SetGapColumn(gap any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	value, err := toTessValue(gap)
	if err != nil {
		return fmt.Errorf("%w: invalid column gap: %v", err, gap)
	}

	err = e.xyz.SetGap(tess.Gap{Column: value})
	if err != nil {
		return fmt.Errorf("failed to set column gap: %w", err)
	}

	return nil
}

func (e *BaseElement) SetOverflow(overflow string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	switch overflow {
	case "visible":
		if err := e.xyz.SetOverflow(tess.Visible); err != nil {
			return fmt.Errorf("failed to set overflow: %w", err)
		}
	case "hidden":
		if err := e.xyz.SetOverflow(tess.Hidden); err != nil {
			return fmt.Errorf("failed to set overflow: %w", err)
		}
	}

	return fmt.Errorf("%w: '%s' is not recognised", types.ErrInvalidStyleValue, overflow)
}

func (e *BaseElement) UnsetZIndex() error {
	return e.SetZIndex(0)
}

func (e *BaseElement) UnsetWidth() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetWidth(tess.Undefined())
	if err != nil {
		return fmt.Errorf("failed to unset width: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetMinWidth() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetMinWidth(tess.Undefined())
	if err != nil {
		return fmt.Errorf("failed to unset min width: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetMaxWidth() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetMaxWidth(tess.Undefined())
	if err != nil {
		return fmt.Errorf("failed to unset max width: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetHeight() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetHeight(tess.Undefined())
	if err != nil {
		return fmt.Errorf("failed to unset height: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetMinHeight() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetMinHeight(tess.Undefined())
	if err != nil {
		return fmt.Errorf("failed to unset min height: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetMaxHeight() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetMaxHeight(tess.Undefined())
	if err != nil {
		return fmt.Errorf("failed to unset max height: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetTop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetTop(tess.Undefined())
	if err != nil {
		return fmt.Errorf("failed to unset top: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetBottom() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetBottom(tess.Undefined())
	if err != nil {
		return fmt.Errorf("failed to unset bottom: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetLeft() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetLeft(tess.Undefined())
	if err != nil {
		return fmt.Errorf("failed to unset left: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetRight() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetRight(tess.Undefined())
	if err != nil {
		return fmt.Errorf("failed to unset right: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetPosition() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetPosition(tess.Static)
	if err != nil {
		return fmt.Errorf("failed to unset position: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetPaddingAll() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetPadding(tess.Edges{All: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset padding: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetPaddingVertical() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetPadding(tess.Edges{Vertical: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset vertical padding: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetPaddingHorizontal() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetPadding(tess.Edges{Horizontal: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset horizontal padding: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetPaddingTop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetPadding(tess.Edges{Top: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset top padding: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetPaddingBottom() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetPadding(tess.Edges{Bottom: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset bottom padding: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetPaddingLeft() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetPadding(tess.Edges{Left: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset left padding: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetPaddingRight() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetPadding(tess.Edges{Right: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset right padding: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetMarginAll() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetMargin(tess.Edges{All: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset margin: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetMarginVertical() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetMargin(tess.Edges{Vertical: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset vertical margin: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetMarginHorizontal() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetMargin(tess.Edges{Horizontal: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset horizontal margin: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetMarginTop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetMargin(tess.Edges{Top: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset top margin: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetMarginBottom() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetMargin(tess.Edges{Bottom: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset bottom margin: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetMarginLeft() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetMargin(tess.Edges{Left: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset left margin: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetMarginRight() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetMargin(tess.Edges{Right: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset right margin: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetDisplay() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetDisplay(tess.Flex)
	if err != nil {
		return fmt.Errorf("failed to unset display: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetAlignSelf() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetAlignSelf(tess.AlignAuto)
	if err != nil {
		return fmt.Errorf("failed to unset align self: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetAlignItems() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetAlignItems(tess.AlignStretch)
	if err != nil {
		return fmt.Errorf("failed to unset align items: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetAlignContent() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetAlignContent(tess.AlignStretch)
	if err != nil {
		return fmt.Errorf("failed to unset align content: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetJustifyContent() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetJustifyContent(tess.JustifyStart)
	if err != nil {
		return fmt.Errorf("failed to unset justify content: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetFlexDirection() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetFlexDirection(tess.Column)
	if err != nil {
		return fmt.Errorf("failed to unset flex direction: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetFlexWrap() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetFlexWrap(tess.NoWrap)
	if err != nil {
		return fmt.Errorf("failed to unset flex wrap: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetFlexGrow() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetFlexGrow(0)
	if err != nil {
		return fmt.Errorf("failed to unset flex grow: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetFlexShrink() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetFlexShrink(1)
	if err != nil {
		return fmt.Errorf("failed to unset flex shrink: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetGapAll() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetGap(tess.Gap{All: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset gap: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetGapRow() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetGap(tess.Gap{Row: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset row gap: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetGapColumn() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetGap(tess.Gap{Column: tess.Undefined()})
	if err != nil {
		return fmt.Errorf("failed to unset column gap: %w", err)
	}

	return nil
}

func (e *BaseElement) UnsetOverflow() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	err := e.xyz.SetOverflow(tess.Visible)
	if err != nil {
		return fmt.Errorf("failed to unset overflow: %w", err)
	}

	return nil
}

func (e *BaseElement) IsDestroyed() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.destroyed
}

func (e *BaseElement) guardUpdate() error {
	if e.destroyed {
		return types.ErrUpdatingDestroyedElement
	}
	return nil
}

func (e *BaseElement) guardPaint() error {
	if e.destroyed {
		return types.ErrPaintingDestroyedElement
	}
	return nil
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
				return tess.Undefined(), fmt.Errorf("%w: invalid point value '%s'", types.ErrInvalidStyleValue, v)
			}

			return tess.Point(float32(points)), nil
		}

		// 100% or 50.5%
		if regexp.MustCompile(`^\d+(\.\d+)?%$`).MatchString(v) {
			percentStr := strings.TrimSuffix(v, "%")
			percent, err := strconv.ParseFloat(percentStr, 32)
			if err != nil {
				return tess.Undefined(), fmt.Errorf("%w: invalid percent value '%s'", types.ErrInvalidStyleValue, v)
			}

			return tess.Percent(float32(percent)), nil
		}
	}

	return tess.Undefined(), fmt.Errorf("%w: '%v' is not recognised", types.ErrInvalidStyleValue, value)
}
