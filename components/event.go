package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/core"
)

// On registers a callback for a specific event.
//
// Supported events:
//
//	"hover", "click", "mouseenter", "mouseleave" "mousepress", "mouserelease", "mousemove"
//	"scroll", "drag", "keypress", "keyrelease", "paste", "focus", "blur", "submit".
func On[T any](event string, fn func(T)) loom.Node {
	return &eventNode[T]{
		event: event,
		fn:    fn,
	}
}

type eventNode[T any] struct {
	event string
	fn    func(T)
}

func (n *eventNode[T]) ID() string {
	return "term.On"
}

func (n *eventNode[T]) Mount(slot *loom.Slot) error {
	return n.Update(slot)
}

func (n *eventNode[T]) Update(slot *loom.Slot) error {
	n.removeListener(slot)

	parent := slot.Parent().(core.Element)

	var remove func()
	switch n.event {
	case "mousemove":
		remove = parent.OnMouseMove(func(event *core.EventMouse) { n.dispatch(event) })
	case "mouseenter", "hover":
		remove = parent.OnMouseEnter(func(event *core.EventMouse) { n.dispatch(event) })
	case "mouseleave":
		remove = parent.OnMouseLeave(func(event *core.EventMouse) { n.dispatch(event) })
	case "mousepress", "click":
		remove = parent.OnMousePress(func(event *core.EventMouse) { n.dispatch(event) })
	case "mouserelease":
		remove = parent.OnMouseRelease(func(event *core.EventMouse) { n.dispatch(event) })
	case "scroll":
		remove = parent.OnMouseScroll(func(event *core.EventMouse) { n.dispatch(event) })
	case "drag":
		remove = parent.OnMouseDrag(func(event *core.EventMouse) { n.dispatch(event) })
	case "keypress":
		remove = parent.OnKeyPress(func(event *core.EventKey) { n.dispatch(event) })
	case "keyrelease":
		remove = parent.OnKeyRelease(func(event *core.EventKey) { n.dispatch(event) })
	case "paste":
		remove = parent.OnPaste(func(event *core.EventPaste) { n.dispatch(event) })
	case "focus":
		remove = parent.OnFocus(func(event *core.EventFocus) { n.dispatch(event) })
	case "blur":
		remove = parent.OnBlur(func(event *core.EventBlur) { n.dispatch(event) })
	case "submit":
		remove = parent.OnSubmit(func(event *core.EventSubmit) { n.dispatch(event) })
	default:
		return fmt.Errorf("On: unsupported event type %q", n.event)
	}

	slot.SetSelf(remove)

	return nil
}

func (n *eventNode[T]) Unmount(slot *loom.Slot) error {
	n.removeListener(slot)
	return nil
}

func (n *eventNode[T]) dispatch(event any) {
	if ect, ok := event.(T); ok {
		n.fn(ect)
	} else {
		// todo: pipe the ctx error chan instead
		core.LogErrorf("On: expected event callback with type %T, got %T.", event, *(new(T)))
	}
}

func (n *eventNode[T]) removeListener(slot *loom.Slot) {
	remove := slot.Self()
	if remove != nil {
		remove.(func())()
		slot.SetSelf(nil)
	}
}
