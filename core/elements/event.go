package elements

import (
	"github.com/loom-go/term/core/elements/events"
	"sync"
)

type EventMouse struct {
	Event
	events.EventMouse
}

type EventKey struct {
	Event
	events.EventKey
}

type EventPaste struct {
	Event
	events.EventPaste
}

type EventFocus struct {
	Event
	Blurred Element
}

type EventBlur struct {
	Event
	Focused Element
}

type EventInput struct {
	Event
	Value string
}

type EventSubmit struct {
	Event
	Value string
}

type EventType string

const (
	EventTypeMousePress   EventType = "mousepress"
	EventTypeMouseRelease EventType = "mouserelease"
	EventTypeMouseEnter   EventType = "mouseenter"
	EventTypeMouseLeave   EventType = "mouseleave"
	EventTypeMouseMove    EventType = "mousemove"
	EventTypeMouseScroll  EventType = "mousescroll"
	EventTypeMouseDrag    EventType = "mousedrag"

	EventTypeKeyPress   EventType = "keypress"
	EventTypeKeyRelease EventType = "keyrelease"

	EventTypePaste EventType = "paste"

	EventTypeFocus EventType = "focus"
	EventTypeBlur  EventType = "blur"

	EventTypeInput  EventType = "input"
	EventTypeSubmit EventType = "submit"

	EventTypeDestroy EventType = "destroy"
)

// EventPhase disinguishes the different phases of event propagation.
// Most events traverse the whole tree from `root to target` (the capture phase) and then from `target to root` (the bubbling phase).
// Some events, like focus and blur, only have a single phase (the default phase) and do not propagate.
// In any cases, the event's target default behavior will happend during the action phase (executed last and only on the target element),
// which can be prevented by calling `PreventDefault` in default, capture or bubble phase.
type EventPhase int

const (
	EventPhaseDefault EventPhase = iota
	EventPhaseCapture
	EventPhaseBubble
	EventPhaseAction
)

type EventOptions struct {
	Capture bool
}

type Event struct {
	mu sync.RWMutex

	prevented bool
	stopped   bool

	phase  EventPhase
	target Element
}

func (e *Event) IsPrevented() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.prevented
}

func (e *Event) PreventDefault() {
	e.mu.Lock()
	e.prevented = true
	e.mu.Unlock()
}

func (e *Event) IsStopped() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.stopped
}

func (e *Event) StopPropagation() {
	e.mu.Lock()
	e.stopped = true
	e.mu.Unlock()
}

func (e *Event) Phase() EventPhase {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.phase
}

func (e *Event) setPhase(phase EventPhase) {
	e.mu.Lock()
	e.phase = phase
	e.mu.Unlock()
}

func (e *Event) Target() Element {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.target
}

func (e *Event) setTarget(target Element) {
	e.mu.Lock()
	e.target = target
	e.mu.Unlock()
}
