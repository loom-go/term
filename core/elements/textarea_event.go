package elements

import (
	"github.com/AnatoleLucet/loom-term/core/elements/events"

	"github.com/AnatoleLucet/go-opentui"
)

func (a *TextAreaElement) handlePaste(event *EventPaste) {
	a.InsertText(event.Text)
	a.rdrctx.ScheduleRender()
}

func (a *TextAreaElement) handleKeyPress(event *EventKey) {
	switch event.Key.String() {
	case "enter":
		a.InsertText("\n")
		a.rdrctx.ScheduleRender()
	case "tab":
		a.InsertText("\t")
		a.rdrctx.ScheduleRender()
	case "space":
		a.InsertText(" ")
		a.rdrctx.ScheduleRender()

	case "left", "ctrl+b":
		a.MoveCursorLeft()
		a.rdrctx.ScheduleRender()
	case "ctrl+left":
		a.MoveCursorWordLeft()
		a.rdrctx.ScheduleRender()
	case "right", "ctrl+f":
		a.MoveCursorRight()
		a.rdrctx.ScheduleRender()
	case "ctrl+right":
		a.MoveCursorWordRight()
		a.rdrctx.ScheduleRender()
	case "up":
		a.MoveCursorUp()
		a.rdrctx.ScheduleRender()
	case "down":
		a.MoveCursorDown()
		a.rdrctx.ScheduleRender()
	case "home":
		a.SetCursor(0, 0)
		a.rdrctx.ScheduleRender()
	case "end":
		lineCount := a.editBuffer.GetTextBuffer().GetLineCount()
		a.SetCursor(lineCount-1, 0)
		eol := a.editBuffer.GetEOL()
		a.SetCursor(eol.Row, eol.Col)
		a.rdrctx.ScheduleRender()

	case "ctrl+a":
		row, col := a.GetCursor()
		if col == 0 && row > 0 {
			a.SetCursor(row-1, 0)
			eol := a.editBuffer.GetEOL()
			a.SetCursor(eol.Row, eol.Col)
			a.rdrctx.ScheduleRender()
		} else {
			a.SetCursor(row, 0)
			a.rdrctx.ScheduleRender()
		}
	case "ctrl+e":
		row, col := a.GetCursor()
		eol := a.editBuffer.GetEOL()
		lineCount := a.editBuffer.GetTextBuffer().GetLineCount()
		if col == eol.Col && row < lineCount-1 {
			a.SetCursor(row+1, 0)
			a.rdrctx.ScheduleRender()
		} else {
			a.SetCursor(eol.Row, eol.Col)
			a.rdrctx.ScheduleRender()
		}

	case "backspace":
		a.RemoveLeft()
		a.rdrctx.ScheduleRender()
	case "delete", "ctrl+d":
		a.RemoveRight()
		a.rdrctx.ScheduleRender()
	case "ctrl+w", "ctrl+backspace":
		a.RemoveWordLeft()
		a.rdrctx.ScheduleRender()
	case "meta+d", "ctrl+delete":
		a.RemoveWordRight()
		a.rdrctx.ScheduleRender()
	case "ctrl+k":
		a.RemoveToLineEnd()
		a.rdrctx.ScheduleRender()
	case "ctrl+u":
		a.RemoveToLineStart()
		a.rdrctx.ScheduleRender()
	case "ctrl+shift+d":
		a.RemoveLine()
		a.rdrctx.ScheduleRender()

	case "meta+enter":
		a.Submit()
		a.rdrctx.ScheduleRender()

	default:
		if event.Key.Ctrl() || event.Key.Alt() || event.Key.Meta() || event.Key.IsUnknown() {
			return
		}

		val := event.Key.Value()
		if (val >= events.Key0 && val <= events.Key9) || event.Key.IsRune() {
			a.InsertText(event.Key.String())
			a.rdrctx.ScheduleRender()
		}
	}
}

func (a *TextAreaElement) handleMousePress(event *EventMouse) {
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
	a.rdrctx.ScheduleRender()
}

func (a *TextAreaElement) handleMouseScroll(event *EventMouse) {
	viewportX, viewportY, viewportW, viewportH, ok := a.editBufferView.GetViewport()
	if !ok {
		return
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
	a.rdrctx.ScheduleRender()
}
