package elements

import (
	"context"
	"fmt"
	"iter"
	"maps"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/AnatoleLucet/loom-term/core/gfx"

	"github.com/AnatoleLucet/go-opentui"
	"github.com/AnatoleLucet/tess"
)

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

func (b *BaseElement) scheduleUpdate(fn func() error) {
	if b.rdrctx == nil {
		b.mu.Lock()
		b.pendingUpdates = append(b.pendingUpdates, fn)
		b.mu.Unlock()
	} else {
		b.rdrctx.ScheduleUpdate(fn)
	}
}

func (b *BaseElement) ID() uint32 {
	return b.id
}

func (b *BaseElement) Self() Element {
	if b.self == nil {
		return b
	}

	return b.self
}

func (b *BaseElement) Parent() Element {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.parent
}

func (b *BaseElement) parentUnsafe() Element {
	return b.parent
}

func (b *BaseElement) SetParent(parent Element) {
	b.scheduleUpdate(func() error {
		b.mu.Lock()
		defer b.mu.Unlock()

		return b.setParentUnsafe(parent)
	})
}

func (b *BaseElement) setParentUnsafe(parent Element) error {
	if err := guardDestroyed(b.ctx); err != nil {
		return err
	}

	b.parent = parent
	return nil
}

func (b *BaseElement) Children() iter.Seq[Element] {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.childrenUnsafe()
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
	b.scheduleUpdate(func() error {
		b.mu.Lock()
		defer b.mu.Unlock()

		return b.appendChildUnsafe(child)
	})
}

func (b *BaseElement) appendChildUnsafe(child Element) error {
	if err := guardDestroyed(b.ctx); err != nil {
		return err
	}

	if child.parentUnsafe() != nil {
		return fmt.Errorf("%w: child already has a parent", ErrFailedToAppendChild)
	}

	b.xyz().AppendChild(child.xyz())

	zindex := child.zindexUnsafe()
	b.children[zindex] = append(b.children[zindex], child)

	child.setParentUnsafe(b)
	return child.setContextUnsafe(b.rdrctx)
}

func (b *BaseElement) RemoveChild(child Element) {
	b.scheduleUpdate(func() error {
		b.mu.Lock()
		defer b.mu.Unlock()

		return b.removeChildUnsafe(child)
	})
}

func (b *BaseElement) removeChildUnsafe(child Element) error {
	if err := guardDestroyed(b.ctx); err != nil {
		return err
	}

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
	return child.setContextUnsafe(nil)
}

func (b *BaseElement) Focus() {
	b.scheduleUpdate(func() error {
		b.mu.Lock()

		if err := guardDestroyed(b.ctx); err != nil {
			b.mu.Unlock()
			return err
		}

		if b.focused || !b.focusable {
			b.mu.Unlock()
			return nil
		}
		b.setFocused(true)

		prev := b.rdrctx.FocusedElement()
		if prev != nil {
			prev.setFocused(false)
		}

		b.rdrctx.SetFocusedElement(b.Self())
		b.mu.Unlock()

		focusEvt := &EventFocus{Blurred: prev}
		focusEvt.setTarget(b.Self())
		b.broadcastEvent(EventTypeFocus, focusEvt)
		if prev != nil {
			blurEvt := &EventBlur{Focused: b.Self()}
			blurEvt.setTarget(prev)
			prev.broadcastEvent(EventTypeBlur, blurEvt)
		}

		return nil
	})
}

func (b *BaseElement) Blur() {
	b.scheduleUpdate(func() error {
		b.mu.Lock()

		if err := guardDestroyed(b.ctx); err != nil {
			b.mu.Unlock()
			return nil
		}

		if !b.focused {
			b.mu.Unlock()
			return nil
		}
		b.setFocused(false)

		b.rdrctx.SetFocusedElement(nil)
		b.mu.Unlock()

		evt := &EventBlur{Focused: b.Self()}
		evt.setTarget(b.Self())
		b.broadcastEvent(EventTypeBlur, evt)

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

	for child := range b.childrenUnsafe() {
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
	b.scheduleUpdate(func() error {
		b.mu.Lock()
		if err := guardDestroyed(b.ctx); err != nil {
			b.mu.Unlock()
			return nil
		}

		if b.parent != nil {
			b.parent.lock()
			b.parent.removeChildUnsafe(b.Self())
			b.parent.unlock()
		}
		b.mu.Unlock()

		b.destroyUnsafe()
		return nil
	})
}

func (b *BaseElement) destroyUnsafe() {
	for child := range b.childrenUnsafe() {
		child.destroyUnsafe()
	}

	b.broadcastEvent(EventTypeDestroy, nil)
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
		b.Focus()
		b.rdrctx.ScheduleRender()
	})

	return func() {
		b.mu.Lock()
		b.focusable = false
		remove()
		b.mu.Unlock()
	}
}

func (b *BaseElement) setContextUnsafe(ctx *RenderContext) error {
	b.rdrctx = ctx
	b.BaseElementStyle.rdrctx = ctx
	b.BaseElementEvent.rdrctx = ctx

	for child := range b.childrenUnsafe() {
		child.lock()
		child.setContextUnsafe(ctx)
		child.unlock()
	}

	for _, pending := range b.pendingUpdates {
		if err := pending(); err != nil {
			return err
		}
	}
	b.pendingUpdates = nil

	return nil
}

func guardDestroyed(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ErrUsingDestroyedElement
	default:
		return nil
	}
}
