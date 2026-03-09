package elements

import (
	"context"
	"fmt"
	"iter"
	"maps"
	"runtime/debug"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/loom-go/term/core/gfx"

	"github.com/AnatoleLucet/go-opentui"
	"github.com/AnatoleLucet/tess"
)

// todo: this file need some work. BaseElement is becoming quite big and it's are to reuse only parts of it

var nextID atomic.Uint32

func newID() uint32 {
	return nextID.Add(1)
}

type BaseElement struct {
	*BaseElementStyle
	*BaseElementEvent

	id uint32

	mu sync.RWMutex

	ctx    context.Context
	rdrctx *RenderContext

	self     Element // used to access the concrete element type in base methods
	parent   Element
	children map[int][]Element // map[zindex]children

	zindex int

	focused   bool
	focusable bool

	pendingUpdates []func() error
}

func NewBaseElement() (base *BaseElement, err error) {
	ctx, cancel := context.WithCancel(context.Background())

	base = &BaseElement{
		id:       newID(),
		ctx:      ctx,
		children: make(map[int][]Element),
	}

	base.BaseElementStyle, err = NewBaseElementStyle(ctx, base)
	if err != nil {
		cancel()
		return nil, err
	}
	base.BaseElementEvent, err = NewBaseElementEvent(ctx, base)
	if err != nil {
		cancel()
		return nil, err
	}

	base.OnDestroy(func() {
		cancel()
		base.BaseElementStyle.free()
		base.BaseElementEvent.free()
	})

	return base, nil
}

func (b *BaseElement) lock() {
	b.mu.Lock()
}

func (b *BaseElement) unlock() {
	b.mu.Unlock()
}

func (b *BaseElement) context() context.Context {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.ctx
}

func (b *BaseElement) RenderContext() (ctx *RenderContext) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.rdrctx
}

func (b *BaseElement) ID() uint32 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.id
}

func (b *BaseElement) Self() Element {
	if b.self == nil {
		return b
	}

	return b.self
}

func (b *BaseElement) Parent() (parent Element) {
	scheduleAccess(b.Self(), func() {
		b.mu.RLock()
		defer b.mu.RUnlock()
		parent = b.parent
	})

	return
}

func (b *BaseElement) parentUnsafe() Element {
	return b.parent
}

func (b *BaseElement) SetParent(parent Element) {
	scheduleUpdate(b.Self(), func() error {
		b.mu.Lock()
		defer b.mu.Unlock()

		return b.Self().setParentUnsafe(parent)
	})
}

func (b *BaseElement) setParentUnsafe(parent Element) error {
	b.parent = parent
	return nil
}

func (b *BaseElement) Children() (children iter.Seq[Element]) {
	scheduleAccess(b.Self(), func() {
		b.mu.RLock()
		defer b.mu.RUnlock()

		children = b.Self().childrenUnsafe()
	})

	return
}

func (b *BaseElement) childrenUnsafe() iter.Seq[Element] {
	return func(yield func(Element) bool) {
		zindexes := slices.Sorted(maps.Keys(b.children))
		for _, zindex := range zindexes {
			for _, child := range b.children[zindex] {
				if !yield(child) {
					return
				}
			}
		}
	}
}

func (b *BaseElement) AppendChild(child Element) {
	scheduleUpdate(b.Self(), func() error {
		b.mu.Lock()
		defer b.mu.Unlock()

		return b.Self().appendChildUnsafe(child)
	})
}

func (b *BaseElement) appendChildUnsafe(child Element) error {
	if child.parentUnsafe() != nil {
		return fmt.Errorf("%w: child already has a parent", ErrFailedToAppendChild)
	}

	zindex := child.zindexUnsafe()

	b.xyz().AppendChild(child.xyz())
	b.children[zindex] = append(b.children[zindex], child)

	child.setParentUnsafe(b.Self())
	child.SetRenderContext(b.rdrctx)
	return nil
}

func (b *BaseElement) RemoveChild(child Element) {
	scheduleUpdate(b.Self(), func() error {
		b.mu.Lock()
		defer b.mu.Unlock()

		return b.Self().removeChildUnsafe(child)
	})
}

func (b *BaseElement) removeChildUnsafe(child Element) error {
	zindex := child.zindexUnsafe()
	children, ok := b.children[zindex]
	if !ok {
		return fmt.Errorf("%w: child not found", ErrFailedToRemoveChild)
	}

	i := slices.Index(children, child)
	if i == -1 {
		return fmt.Errorf("%w: child not found", ErrFailedToRemoveChild)
	}

	b.xyz().RemoveChild(child.xyz())
	b.children[zindex] = slices.Delete(children, i, i+1)

	child.setParentUnsafe(nil)
	child.SetRenderContext(nil)
	return nil
}

func (b *BaseElement) Focus() {
	scheduleUpdate(b.Self(), func() error {
		b.mu.Lock()
		if b.focused || !b.focusable {
			b.mu.Unlock()
			return nil
		}
		b.Self().setFocused(true)

		prev := b.rdrctx.FocusedElement()
		if prev != nil {
			prev.setFocused(false)
		}

		b.rdrctx.SetFocusedElement(b.Self())
		b.mu.Unlock()

		focusEvt := &EventFocus{Blurred: prev}
		focusEvt.setTarget(b.Self())
		b.Self().broadcastEvent(EventTypeFocus, focusEvt)
		if prev != nil {
			blurEvt := &EventBlur{Focused: b.Self()}
			blurEvt.setTarget(prev)
			prev.broadcastEvent(EventTypeBlur, blurEvt)
		}

		return nil
	})
}

func (b *BaseElement) Blur() {
	scheduleUpdate(b.Self(), func() error {
		b.mu.Lock()
		if !b.focused {
			b.mu.Unlock()
			return nil
		}
		b.setFocused(false)

		b.rdrctx.SetFocusedElement(nil)
		b.mu.Unlock()

		evt := &EventBlur{Focused: b.Self()}
		evt.setTarget(b.Self())
		b.Self().broadcastEvent(EventTypeBlur, evt)

		return nil
	})
}

func (b *BaseElement) Record(cb *gfx.CommandBuffer, container gfx.Rect) error {
	self := b.Self()
	rect := b.rect(container)

	render := gfx.NewCommand(gfx.CmdRender, self).WithRect(rect).WithCallback(func() error {
		return b.rdrctx.AddToHitGrid(self, rect)
	})
	cb.Add(render)

	if b.xyz().GetOverflow() == tess.Hidden {
		cb.Add(gfx.NewCommand(gfx.CmdPushHitGridScissors, self).WithScissors(rect))
		cb.Add(gfx.NewCommand(gfx.CmdPushOverflowScissors, self).WithScissors(rect))

		defer cb.Add(gfx.NewCommand(gfx.CmdPopHitGridScissors, self))
		defer cb.Add(gfx.NewCommand(gfx.CmdPopOverflowScissors, self))
	}

	for child := range b.Self().childrenUnsafe() {
		err := child.Record(cb, rect)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BaseElement) Render(buf *opentui.Buffer, rect gfx.Rect) error {
	return nil
}

func (b *BaseElement) Destroy() {
	scheduleUpdate(b.Self(), func() error {
		if b.parent != nil {
			b.parent.removeChildUnsafe(b.Self())
		}

		b.Self().destroyUnsafe()
		return nil
	})
}

func (b *BaseElement) destroyUnsafe() {
	for child := range b.Self().childrenUnsafe() {
		child.destroyUnsafe()
	}

	b.Self().broadcastEvent(EventTypeDestroy, nil)
}

func (b *BaseElement) rect(container gfx.Rect) gfx.Rect {
	layout := b.xyz().GetLayout()
	return gfx.Rect{
		X:      container.X + int(layout.Left()) + int(b.translateX),
		Y:      container.Y + int(layout.Top()) + int(b.translateY),
		Width:  int(layout.Width()),
		Height: int(layout.Height()),
	}
}

func (b *BaseElement) contentRect(container gfx.Rect) gfx.Rect {
	layout := b.xyz().GetLayout()
	p := layout.Padding()
	return gfx.Rect{
		X:      container.X + int(layout.Left()+p.Left()) + int(b.translateX),
		Y:      container.Y + int(layout.Top()+p.Top()) + int(b.translateY),
		Width:  int(layout.Width() - p.Left() - p.Right()),
		Height: int(layout.Height() - p.Top() - p.Bottom()),
	}
}

func (b *BaseElement) relativeRect() gfx.Rect {
	layout := b.xyz().GetLayout()
	return gfx.Rect{
		X:      int(layout.Left()),
		Y:      int(layout.Top()),
		Width:  int(layout.Width()),
		Height: int(layout.Height()),
	}
}

func (b *BaseElement) setFocused(focused bool) {
	b.focused = focused
}

func (b *BaseElement) setFocusable() (unset func()) {
	b.mu.Lock()
	if b.focusable {
		b.mu.Unlock()
		return func() {}
	}
	b.focusable = true
	b.mu.Unlock()

	remove := b.mousePressAction(func(e *EventMouse) {
		b.Self().Focus()
	})

	return func() {
		b.mu.Lock()
		b.focusable = false
		remove()
		b.mu.Unlock()
	}
}

func (b *BaseElement) addPendingUpdate(fn func() error) {
	b.mu.Lock()
	b.pendingUpdates = append(b.pendingUpdates, fn)
	b.mu.Unlock()
}

func (b *BaseElement) SetRenderContext(ctx *RenderContext) {
	batchUpdate(b.Self(), func() {
		b.rdrctx = ctx
		b.BaseElementStyle.rdrctx = ctx
		b.BaseElementEvent.rdrctx = ctx

		for child := range b.Self().childrenUnsafe() {
			child.SetRenderContext(ctx)
		}

		if ctx != nil {
			for _, pending := range b.pendingUpdates {
				scheduleUpdate(b.Self(), pending)
			}
			b.pendingUpdates = nil
		}
	})
}

func scheduleAccess(e Element, fn func()) {
	if e.RenderContext() == nil {
		fn()
	} else {
		e.RenderContext().scheduleAccess(fn)
	}
}

func batchUpdate(e Element, fn func()) {
	if e.RenderContext() == nil {
		fn()
	} else {
		e.RenderContext().Batch(fn)
	}
}

func scheduleUpdate(e Element, fn func() error) {
	if e.RenderContext() == nil {
		e.Self().addPendingUpdate(fn)
	} else {
		e.RenderContext().scheduleUpdate(func() error {
			if err := guardDestroyed(e); err != nil {
				return err
			}

			return fn()
		})
	}
}
func scheduleUpdateSync(e Element, fn func() error) {
	if e.RenderContext() == nil {
		e.Self().addPendingUpdate(fn)
	} else {
		e.RenderContext().scheduleUpdateSync(func() error {
			if err := guardDestroyed(e); err != nil {
				return err
			}

			return fn()
		})
	}
}

func guardDestroyed(e Element) error {
	ctx := e.context()

	select {
	case <-ctx.Done():
		return fmt.Errorf("%T: %w\n%s", e, ErrUsingDestroyedElement, debug.Stack())
	default:
		return nil
	}
}
