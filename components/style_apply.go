package components

import (
	"fmt"
	"sync/atomic"

	"github.com/AnatoleLucet/loom"
	appctx "github.com/AnatoleLucet/loom-term/components/context"
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom/signals"
)

var nextID atomic.Uint32

func newID() uint32 {
	return nextID.Add(1)
}

type applyNode struct {
	event  string
	styles []Style
}

// Apply applies the given Styles to a node.
func Apply(styles ...Style) *applyNode {
	return &applyNode{
		styles: styles,
	}
}

// ApplyOn applies the Styles on a specific event ("hover", "focus", "active").
func ApplyOn(event string, styles ...Style) *applyNode {
	return &applyNode{
		event:  event,
		styles: styles,
	}
}

func (s *applyNode) ID() string {
	return "term.Apply"
}

func (s *applyNode) Mount(slot *loom.Slot) error {
	id := newID()
	layer := &styleLayer{
		id:      id,
		event:   s.event,
		styles:  s.styles,
		visible: s.event == "",
	}
	slot.SetSelf(layer)

	return s.run(slot, true)
}

func (s *applyNode) Update(slot *loom.Slot) error {
	return s.run(slot, false)
}

func (s *applyNode) Unmount(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("Apply (style): %w", err)
	}

	self := slot.Self().(*styleLayer)
	parent := slot.Parent().(core.Element)

	return ctx.BatchRender(func() error {
		// remove our layer from the stack
		stack := getStyleStack(parent)
		stack.Pop(self.id)

		// unset our layer styles
		self.visible = false
		s.applyStyleLayer(parent, self)

		// reapply the properties that we might have unset
		s.applyStyleStack(slot)

		return nil
	})
}

func (s *applyNode) applyStyleStack(slot *loom.Slot) {
	self := slot.Self().(*styleLayer)
	parent := slot.Parent().(core.Element)

	stack := getStyleStack(parent)
	for _, layer := range stack.layers {
		if layer == self {
			s.applyStyleLayer(parent, layer)
		} else {
			// untrack style layers that's not ours. each apply tracks its own layer
			signals.Untrack(func() any {
				s.applyStyleLayer(parent, layer)
				return nil
			})
		}
	}
}

func (s *applyNode) applyStyleLayer(parent core.Element, layer *styleLayer) bool {
	for _, style := range layer.styles {
		if layer.visible {
			applyStyle(parent, style)
		} else {
			removeStyle(parent, style)
		}
	}

	return true
}

func (s *applyNode) run(slot *loom.Slot, initial bool) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("Apply (style): %w", err)
	}

	if s.event == "" {
		s.watch(slot, ctx, initial)
	} else {
		s.registerEvents(slot, ctx, initial)
	}

	return nil
}

func (s *applyNode) watch(slot *loom.Slot, ctx *appctx.Context, initial bool) {
	self := slot.Self().(*styleLayer)
	parent := slot.Parent().(core.Element)

	stack := getStyleStack(parent)

	signals.RenderEffect(func() {
		ctx.BatchRender(func() error {
			// fmt.Printf("applying style layer %d (visible: %t)\n\r", self.id, self.visible)
			if initial {
				// if we're in the initial phase, we can just apply our own layer
				stack.Push(self)
				s.applyStyleLayer(parent, self)
			} else {
				// else we must update our layer (re-prioritizing it) then reapply the whole stack
				// to make sure removed properties in the new layer gets proper fallback
				stack.Replace(self.id, s.styles)
				s.applyStyleStack(slot)
			}
			// fmt.Printf("applied style layer %d (visible: %t)\n\r", self.id, self.visible)

			return nil
		})
	})
}

func (s *applyNode) registerEvents(slot *loom.Slot, ctx *appctx.Context, initial bool) {
	// use a custom owner to dispose the RenderEffect in watch() when the event is removed.
	owner := signals.NewOwner()

	self := slot.Self().(*styleLayer)
	parent := slot.Parent().(core.Element)

	// todo: have visible be a signal, and always start a watch.
	// on event, just set visible to true or false, and the watch will react.
	add := func() {
		owner.Run(func() error {
			self.visible = true
			s.watch(slot, ctx, initial)
			return nil
		})
	}
	remove := func() {
		ctx.BatchRender(func() error {
			self.visible = false
			owner.Dispose()
			s.applyStyleStack(slot)
			return nil
		})
	}

	if s.event == "hover" {
		signals.OnCleanup(parent.OnMouseEnter(func(*core.EventMouse) { add() }))
		signals.OnCleanup(parent.OnMouseLeave(func(*core.EventMouse) { remove() }))
	}
	if s.event == "focus" {
		signals.OnCleanup(parent.OnFocus(func(*core.EventFocus) { add() }))
		signals.OnCleanup(parent.OnBlur(func(*core.EventBlur) { remove() }))
	}
	if s.event == "active" {
		signals.OnCleanup(parent.OnMousePress(func(*core.EventMouse) { add() }))
		signals.OnCleanup(parent.OnMouseRelease(func(*core.EventMouse) { remove() }))
	}
}
