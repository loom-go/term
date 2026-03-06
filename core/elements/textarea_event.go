package elements

import (
	"github.com/AnatoleLucet/loom-term/core/elements/events"

	"github.com/AnatoleLucet/go-opentui"
)

func (a *TextAreaElement) dispatchInputEvent(value string) {
	a.rdrctx.DispatchEvent(EventTypeInput, a.Self(), &EventInput{Value: value})
}

func (a *TextAreaElement) Submit() {
	a.rdrctx.DispatchEvent(EventTypeSubmit, a.Self(), &EventSubmit{Value: a.Value()})
}

func (a *TextAreaElement) handlePaste(event *EventPaste) {
	a.InsertValue(event.Value)
	event.StopPropagation()
	a.dispatchInputEvent(a.Value())
}

func (a *TextAreaElement) handleKeyPress(event *EventKey) {
	switch event.Key.String() {
	case "enter":
		a.InsertValue("\n")
		event.StopPropagation()
		a.dispatchInputEvent(a.Value())
	case "tab":
		a.InsertValue("\t")
		event.StopPropagation()
		a.dispatchInputEvent(a.Value())
	case "space":
		a.InsertValue(" ")
		event.StopPropagation()
		a.dispatchInputEvent(a.Value())

	case "left", "ctrl+b":
		a.MoveCursorLeft()
		event.StopPropagation()
	case "ctrl+left":
		a.MoveCursorWordLeft()
		event.StopPropagation()
	case "right", "ctrl+f":
		a.MoveCursorRight()
		event.StopPropagation()
	case "ctrl+right":
		a.MoveCursorWordRight()
		event.StopPropagation()
	case "up":
		a.MoveCursorUp()
		event.StopPropagation()
	case "down":
		a.MoveCursorDown()
		event.StopPropagation()
	case "home":
		a.SetCursor(0, 0)
		event.StopPropagation()
	case "end":
		lineCount := a.editBuffer.GetTextBuffer().GetLineCount()
		a.SetCursor(lineCount-1, 0)
		eol := a.editBuffer.GetEOL()
		a.SetCursor(eol.Row, eol.Col)
		event.StopPropagation()

	case "ctrl+a":
		row, col := a.GetCursor()
		if col == 0 && row > 0 {
			a.SetCursor(row-1, 0)
			eol := a.editBuffer.GetEOL()
			a.SetCursor(eol.Row, eol.Col)
			event.StopPropagation()
		} else {
			a.SetCursor(row, 0)
			event.StopPropagation()
		}
	case "ctrl+e":
		row, col := a.GetCursor()
		eol := a.editBuffer.GetEOL()
		lineCount := a.editBuffer.GetTextBuffer().GetLineCount()
		if col == eol.Col && row < lineCount-1 {
			a.SetCursor(row+1, 0)
			event.StopPropagation()
		} else {
			a.SetCursor(eol.Row, eol.Col)
			event.StopPropagation()
		}

	case "backspace":
		a.RemoveLeft()
		event.StopPropagation()
		a.dispatchInputEvent(a.Value())
	case "delete", "ctrl+d":
		a.RemoveRight()
		event.StopPropagation()
		a.dispatchInputEvent(a.Value())
	case "ctrl+w", "ctrl+backspace":
		a.RemoveWordLeft()
		event.StopPropagation()
		a.dispatchInputEvent(a.Value())
	case "meta+d", "ctrl+delete":
		a.RemoveWordRight()
		event.StopPropagation()
		a.dispatchInputEvent(a.Value())
	case "ctrl+k":
		a.RemoveToLineEnd()
		event.StopPropagation()
		a.dispatchInputEvent(a.Value())
	case "ctrl+u":
		a.RemoveToLineStart()
		event.StopPropagation()
		a.dispatchInputEvent(a.Value())
	case "ctrl+shift+d":
		a.RemoveLine()
		event.StopPropagation()
		a.dispatchInputEvent(a.Value())

	case "meta+enter":
		a.Submit()
		event.StopPropagation()

	default:
		if event.Key.Ctrl() || event.Key.Alt() || event.Key.Meta() || event.Key.IsUnknown() {
			return
		}

		val := event.Key.Value()
		if (val >= events.Key0 && val <= events.Key9) || event.Key.IsRune() {
			a.InsertValue(event.Key.String())
			event.StopPropagation()
			a.dispatchInputEvent(a.Value())
		}
	}
}

func (a *TextAreaElement) handleMousePress(event *EventMouse) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.RLock()
		layout := a.xyz().GetLayout()
		elemX := layout.AbsoluteLeft()
		elemY := layout.AbsoluteTop()
		a.mu.RUnlock()

		relativeX := max(0, event.X-int(elemX))
		relativeY := max(0, event.Y-int(elemY))

		a.editBufferView.SetLocalSelection(
			int32(relativeX), int32(relativeY),
			int32(relativeX), int32(relativeY), // todo: impl focus on click and drag
			nil, nil,
			true, false,
		)

		return nil
	})
}

func (a *TextAreaElement) handleMouseScroll(event *EventMouse) {
	scheduleUpdate(a.Self(), func() error {
		viewportX, viewportY, viewportW, viewportH, ok := a.editBufferView.GetViewport()
		if !ok {
			return nil
		}

		newViewportX := viewportX
		newViewportY := viewportY
		delta := uint32(1 * a.scrollFactor)

		switch event.Button {
		case events.MouseWheelUp:
			newOffsetY := viewportY - delta
			if delta > viewportY {
				newOffsetY = 0
			}
			newViewportY = newOffsetY

		case events.MouseWheelDown:
			lineCount := a.editBufferView.GetTotalVirtualLineCount()
			maxOffsetY := lineCount - viewportH
			if viewportH > lineCount {
				maxOffsetY = 0
			}
			newOffsetY := min(maxOffsetY, viewportY+delta)
			newViewportY = newOffsetY
		}

		if a.wrapMode == opentui.WrapModeNone {
			switch event.Button {
			case events.MouseWheelLeft:
				newOffsetX := viewportX - delta
				if delta > viewportX {
					newOffsetX = 0
				}
				newViewportX = newOffsetX

			case events.MouseWheelRight:
				newOffsetX := viewportX + delta
				newViewportX = newOffsetX
			}
		}

		a.editBufferView.SetViewport(newViewportX, newViewportY, viewportW, viewportH, true)
		return nil
	})
}
