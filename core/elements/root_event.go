package elements

import (
	"context"
	"github.com/AnatoleLucet/loom-term/core/elements/events"
)

func (r *RootElement) listenToMouseEvents(ctx context.Context) {
	listenner := events.NewMouseListener(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case event := <-listenner.Listen(ctx):
			offsetX, offsetY := r.rdrctx.RenderOffset()

			evt := &EventMouse{EventMouse: *event}

			// normalize coordinate with render offset before using the event
			evt.Y -= offsetY
			evt.X -= offsetX

			target := r.rdrctx.CheckHit(evt.X, evt.Y)
			if target == nil {
				target = r
			}

			evt.setTarget(target)
			r.rdrctx.SetMousePosition(evt.X, evt.Y)

			switch evt.Action {
			case events.MouseActionPress:
				r.rdrctx.SetPressedElement(target)
				r.rdrctx.DispatchEvent(EventTypeMousePress, target, evt)

			case events.MouseActionRelease:
				r.rdrctx.SetPressedElement(nil)
				r.rdrctx.DispatchEvent(EventTypeMouseRelease, target, evt)

			case events.MouseActionMove:
				hovered := r.rdrctx.HoveredElement()
				if hovered != target {
					r.rdrctx.SetHoveredElement(target)

					if hovered != nil {
						hovered.broadcastEvent(EventTypeMouseLeave, evt)
					}

					target.broadcastEvent(EventTypeMouseEnter, evt)
				}

				r.rdrctx.DispatchEvent(EventTypeMouseMove, target, evt)

			case events.MouseActionScroll:
				r.rdrctx.DispatchEvent(EventTypeMouseScroll, target, evt)

			case events.MouseActionDrag:
				pressed := r.rdrctx.PressedElement()
				if pressed != nil {
					r.rdrctx.DispatchEvent(EventTypeMouseDrag, pressed, evt)
				} else {
					r.rdrctx.DispatchEvent(EventTypeMouseDrag, target, evt)
				}
			}
		}
	}
}

func (r *RootElement) listenToKeyboardEvents(ctx context.Context) {
	listener := events.NewKeyboardListener(ctx)

	// keys
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case event := <-listener.ListenKey(ctx):
				focused := r.rdrctx.FocusedElement()
				if focused == nil {
					focused = r
				}

				evt := &EventKey{EventKey: *event}
				evt.setTarget(focused)

				switch event.Action {
				case events.KeyActionPress:
					r.rdrctx.DispatchEvent(EventTypeKeyPress, focused, evt)
				case events.KeyActionRelease:
					r.rdrctx.DispatchEvent(EventTypeKeyRelease, focused, evt)
				}
			}
		}
	}()

	// paste
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case event := <-listener.ListenPaste(ctx):
				focused := r.rdrctx.FocusedElement()
				if focused == nil {
					focused = r
				}

				evt := &EventPaste{EventPaste: *event}
				evt.setTarget(focused)

				r.rdrctx.DispatchEvent(EventTypePaste, focused, evt)
			}
		}
	}()
}

func (r *RootElement) listenToResizeEvents(ctx context.Context) {
	listenner := events.NewResizeListener(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case event := <-listenner.Listen(ctx):
			r.SetWidth(event.Width)
			if r.rdrctx.RenderType() == RenderTypeFullscreen {
				r.SetHeight(event.Height)
			}

			r.rdrctx.ScheduleRender()
		}
	}
}

func (r *RootElement) listenToCapabilites(ctx context.Context) {
	listenner := events.NewCapabilitiesListener(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case event := <-listenner.Listen(ctx):
			err := r.rdrctx.UpdateCapabilites(event.Raw)
			if err != nil {
				r.rdrctx.emitError(err)
			}
		}
	}
}

func (r *RootElement) listenToExitEvents(ctx context.Context) {
	// handle ctrl+c
	remove := r.keyPressAction(func(event *EventKey) {
		if event.Key.String() == "ctrl+c" {
			r.Destroy()
		}
	})
	go func() {
		<-ctx.Done()
		remove()
	}()

	// handle sigint, sigterm, sigquit etc
	go func() {
		listenner := events.NewExitListener(ctx)

		for {
			select {
			case <-ctx.Done():
				return

			case <-listenner.Listen(ctx):
				r.Destroy()
				return
			}
		}
	}()
}

func (rc *RenderContext) DispatchEvent(typ EventType, target Element, event any) {
	proxy, ok := event.(interface {
		IsStopped() bool
		IsPrevented() bool
		setTarget(Element)
		setPhase(EventPhase)
	})
	if !ok {
		return
	}
	proxy.setTarget(target)

	path := pathToRoot(target)

	// capture phase
	proxy.setPhase(EventPhaseCapture)
	for i := len(path) - 1; i >= 0; i-- {
		if proxy.IsStopped() {
			break
		}

		path[i].broadcastEvent(typ, event)
	}

	// bubble phase
	proxy.setPhase(EventPhaseBubble)
	for i := 0; i < len(path); i++ {
		if proxy.IsStopped() {
			break
		}

		path[i].broadcastEvent(typ, event)
	}

	if proxy.IsPrevented() {
		return
	}

	// action phase
	proxy.setPhase(EventPhaseAction)
	for i := 0; i < len(path); i++ {
		if proxy.IsStopped() {
			break
		}

		path[i].broadcastEvent(typ, event)
	}
}

func (rc *RenderContext) RefreshMouseState() {
	hovered := rc.HoveredElement()
	target := rc.CheckHit(rc.lastMouseX, rc.lastMouseY)
	if target == nil {
		target = rc.root
	}

	if hovered == target {
		return
	}

	evt := &EventMouse{EventMouse: events.EventMouse{
		X:      rc.lastMouseX,
		Y:      rc.lastMouseY,
		Action: events.MouseActionMove,
	}}
	evt.setTarget(target)

	rc.SetHoveredElement(target)

	if hovered != nil {
		hovered.broadcastEvent(EventTypeMouseLeave, evt)
	}
	target.broadcastEvent(EventTypeMouseEnter, evt)
}

func pathToRoot(target Element) []Element {
	var path []Element
	for elem := target; elem != nil; elem = elem.Parent() {
		path = append(path, elem)
	}

	return path
}
