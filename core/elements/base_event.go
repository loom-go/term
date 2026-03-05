package elements

import (
	"context"
	"sync"
)

type BaseElementEvent struct {
	base *BaseElement

	ctx    context.Context
	rdrctx *RenderContext

	mousePress   map[EventPhase]*eventBroadcaster[*EventMouse]
	mouseRelease map[EventPhase]*eventBroadcaster[*EventMouse]
	mouseMove    map[EventPhase]*eventBroadcaster[*EventMouse]
	mouseScroll  map[EventPhase]*eventBroadcaster[*EventMouse]
	mouseDrag    map[EventPhase]*eventBroadcaster[*EventMouse]
	mouseEnter   map[EventPhase]*eventBroadcaster[*EventMouse]
	mouseLeave   map[EventPhase]*eventBroadcaster[*EventMouse]

	keyPress   map[EventPhase]*eventBroadcaster[*EventKey]
	keyRelease map[EventPhase]*eventBroadcaster[*EventKey]

	paste map[EventPhase]*eventBroadcaster[*EventPaste]

	focus map[EventPhase]*eventBroadcaster[*EventFocus]
	blur  map[EventPhase]*eventBroadcaster[*EventBlur]

	input  map[EventPhase]*eventBroadcaster[*EventInput]
	submit map[EventPhase]*eventBroadcaster[*EventSubmit]

	destroy map[EventPhase]*eventBroadcaster[struct{}]
}

func NewBaseElementEvent(ctx context.Context, base *BaseElement) (*BaseElementEvent, error) {
	e := &BaseElementEvent{base: base, ctx: ctx}

	e.mousePress = newPhasedEventBroadcasters[*EventMouse]()
	e.mouseRelease = newPhasedEventBroadcasters[*EventMouse]()
	e.mouseMove = newPhasedEventBroadcasters[*EventMouse]()
	e.mouseScroll = newPhasedEventBroadcasters[*EventMouse]()
	e.mouseDrag = newPhasedEventBroadcasters[*EventMouse]()
	e.mouseEnter = newPhasedEventBroadcasters[*EventMouse]()
	e.mouseLeave = newPhasedEventBroadcasters[*EventMouse]()

	e.keyPress = newPhasedEventBroadcasters[*EventKey]()
	e.keyRelease = newPhasedEventBroadcasters[*EventKey]()

	e.paste = newPhasedEventBroadcasters[*EventPaste]()

	e.focus = newPhasedEventBroadcasters[*EventFocus]()
	e.blur = newPhasedEventBroadcasters[*EventBlur]()

	e.input = newPhasedEventBroadcasters[*EventInput]()
	e.submit = newPhasedEventBroadcasters[*EventSubmit]()

	e.destroy = newPhasedEventBroadcasters[struct{}]()

	return e, nil
}

func (e *BaseElementEvent) OnMousePress(handler func(*EventMouse), options ...EventOptions) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mousePress[eventPhaseFromOptions(options...)], handler)
}

func (e *BaseElementEvent) mousePressAction(handler func(*EventMouse)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mousePress[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnMouseRelease(handler func(*EventMouse), options ...EventOptions) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mouseRelease[eventPhaseFromOptions(options...)], handler)
}

func (e *BaseElementEvent) mouseReleaseAction(handler func(*EventMouse)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mouseRelease[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnMouseEnter(handler func(*EventMouse)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mouseEnter[EventPhaseDefault], handler)
}

func (e *BaseElementEvent) mouseEnterAction(handler func(*EventMouse)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mouseEnter[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnMouseLeave(handler func(*EventMouse)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mouseLeave[EventPhaseDefault], handler)
}

func (e *BaseElementEvent) mouseLeaveAction(handler func(*EventMouse)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mouseLeave[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnMouseMove(handler func(*EventMouse), options ...EventOptions) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mouseMove[eventPhaseFromOptions(options...)], handler)
}

func (e *BaseElementEvent) mouseMoveAction(handler func(*EventMouse)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mouseMove[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnMouseScroll(handler func(*EventMouse), options ...EventOptions) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mouseScroll[eventPhaseFromOptions(options...)], handler)
}

func (e *BaseElementEvent) mouseScrollAction(handler func(*EventMouse)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mouseScroll[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnMouseDrag(handler func(*EventMouse), options ...EventOptions) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mouseDrag[eventPhaseFromOptions(options...)], handler)
}

func (e *BaseElementEvent) mouseDragAction(handler func(*EventMouse)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.mouseDrag[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnKeyPress(handler func(*EventKey), options ...EventOptions) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.keyPress[eventPhaseFromOptions(options...)], handler)
}

func (e *BaseElementEvent) keyPressAction(handler func(*EventKey)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.keyPress[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnKeyRelease(handler func(*EventKey), options ...EventOptions) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.keyRelease[eventPhaseFromOptions(options...)], handler)
}

func (e *BaseElementEvent) keyReleaseAction(handler func(*EventKey)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.keyRelease[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnPaste(handler func(*EventPaste), options ...EventOptions) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.paste[eventPhaseFromOptions(options...)], handler)
}

func (e *BaseElementEvent) pasteAction(handler func(*EventPaste)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.paste[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnFocus(handler func(*EventFocus), options ...EventOptions) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.focus[EventPhaseDefault], handler)
}

func (e *BaseElementEvent) focusAction(handler func(*EventFocus)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.focus[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnBlur(handler func(*EventBlur), options ...EventOptions) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.blur[EventPhaseDefault], handler)
}

func (e *BaseElementEvent) blurAction(handler func(*EventBlur)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.blur[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnInput(handler func(*EventInput), options ...EventOptions) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.input[eventPhaseFromOptions(options...)], handler)
}

func (e *BaseElementEvent) inputAction(handler func(*EventInput)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.input[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnSubmit(handler func(*EventSubmit), options ...EventOptions) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.submit[eventPhaseFromOptions(options...)], handler)
}

func (e *BaseElementEvent) submitAction(handler func(*EventSubmit)) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.submit[EventPhaseAction], handler)
}

func (e *BaseElementEvent) OnDestroy(handler func()) (remove func()) {
	if err := guardDestroyed(e.ctx); err != nil {
		return func() {}
	}

	return newEventHandler(e.destroy[EventPhaseDefault], func(struct{}) {
		handler()
	})
}

func (e *BaseElementEvent) broadcastEvent(typ EventType, event any) {
	var phase EventPhase
	if evt, ok := event.(interface{ Phase() EventPhase }); ok {
		phase = evt.Phase()
	}

	switch typ {
	case EventTypeMousePress:
		e.mousePress[phase].Broadcast(event.(*EventMouse))
	case EventTypeMouseRelease:
		e.mouseRelease[phase].Broadcast(event.(*EventMouse))
	case EventTypeMouseMove:
		e.mouseMove[phase].Broadcast(event.(*EventMouse))
	case EventTypeMouseScroll:
		e.mouseScroll[phase].Broadcast(event.(*EventMouse))
	case EventTypeMouseDrag:
		e.mouseDrag[phase].Broadcast(event.(*EventMouse))
	case EventTypeMouseEnter:
		e.mouseEnter[phase].Broadcast(event.(*EventMouse))
	case EventTypeMouseLeave:
		e.mouseLeave[phase].Broadcast(event.(*EventMouse))

	case EventTypeKeyPress:
		e.keyPress[phase].Broadcast(event.(*EventKey))
	case EventTypeKeyRelease:
		e.keyRelease[phase].Broadcast(event.(*EventKey))

	case EventTypePaste:
		e.paste[phase].Broadcast(event.(*EventPaste))

	case EventTypeFocus:
		e.focus[phase].Broadcast(event.(*EventFocus))
	case EventTypeBlur:
		e.blur[phase].Broadcast(event.(*EventBlur))

	case EventTypeInput:
		e.input[phase].Broadcast(event.(*EventInput))
	case EventTypeSubmit:
		e.submit[phase].Broadcast(event.(*EventSubmit))

	case EventTypeDestroy:
		e.destroy[phase].Broadcast(struct{}{})
	}
}

func (e *BaseElementEvent) free() {
	for _, b := range e.mousePress {
		b.Clear()
	}
	for _, b := range e.mouseRelease {
		b.Clear()
	}
	for _, b := range e.mouseMove {
		b.Clear()
	}
	for _, b := range e.mouseScroll {
		b.Clear()
	}
	for _, b := range e.mouseDrag {
		b.Clear()
	}
	for _, b := range e.mouseEnter {
		b.Clear()
	}
	for _, b := range e.mouseLeave {
		b.Clear()
	}

	for _, b := range e.keyPress {
		b.Clear()
	}
	for _, b := range e.keyRelease {
		b.Clear()
	}

	for _, b := range e.paste {
		b.Clear()
	}

	for _, b := range e.focus {
		b.Clear()
	}
	for _, b := range e.blur {
		b.Clear()
	}

	for _, b := range e.input {
		b.Clear()
	}
	for _, b := range e.submit {
		b.Clear()
	}
}

func newPhasedEventBroadcasters[T any]() map[EventPhase]*eventBroadcaster[T] {
	return map[EventPhase]*eventBroadcaster[T]{
		EventPhaseDefault: newEventBroadcaster[T](),
		EventPhaseCapture: newEventBroadcaster[T](),
		EventPhaseBubble:  newEventBroadcaster[T](),
		EventPhaseAction:  newEventBroadcaster[T](),
	}
}

func newEventHandler[T any](broadcaster *eventBroadcaster[T], handler func(T)) (remove func()) {
	id := broadcaster.Add(handler)

	return func() {
		broadcaster.Remove(id)
	}
}

func eventPhaseFromOptions(options ...EventOptions) EventPhase {
	var opts EventOptions
	if len(options) > 0 {
		opts = options[0]
	}

	if opts.Capture {
		return EventPhaseCapture
	}

	return EventPhaseBubble
}

type eventEntry[T any] struct {
	id      uint32
	handler func(T)
}

type eventBroadcaster[T any] struct {
	mu sync.RWMutex

	entries []*eventEntry[T]
}

func newEventBroadcaster[T any]() *eventBroadcaster[T] {
	return &eventBroadcaster[T]{}
}

func (b *eventBroadcaster[T]) Broadcast(event T) {
	for _, entry := range b.entries {
		entry.handler(event)
	}
}

func (b *eventBroadcaster[T]) Add(handler func(T)) (id uint32) {
	b.mu.Lock()
	id = newID()

	b.entries = append(b.entries, &eventEntry[T]{
		id:      id,
		handler: handler,
	})
	b.mu.Unlock()

	return id
}

func (b *eventBroadcaster[T]) Remove(id uint32) {
	b.mu.Lock()
	for i, entry := range b.entries {
		if entry.id == id {
			b.entries = append(b.entries[:i], b.entries[i+1:]...)
			break
		}
	}
	b.mu.Unlock()
}

func (b *eventBroadcaster[T]) Clear() {
	b.mu.Lock()
	b.entries = nil
	b.mu.Unlock()
}
