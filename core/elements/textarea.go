package elements

import (
	"fmt"

	"github.com/AnatoleLucet/loom-term/core/gfx"

	"github.com/AnatoleLucet/go-opentui"
)

type TextAreaElement struct {
	*BoxElement

	wrapMode     uint8
	scrollFactor float32

	placeholder *TextChunk

	foregroundColor *Color
	backgroundColor *Color

	editBuffer     *opentui.EditBuffer
	editBufferView *opentui.EditorView

	input map[EventPhase]*eventBroadcaster[*EventInput]
}

func NewTextAreaElement() (area *TextAreaElement, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("TextArea: %w: %v", ErrFailedToInitializeElement, err)
		}
	}()

	box, err := NewBoxElement()
	if err != nil {
		return nil, err
	}

	area = &TextAreaElement{
		BoxElement: box,

		placeholder: NewTextChunk(),

		foregroundColor: NewColor(""),
		backgroundColor: NewColor(""),

		scrollFactor: 1,

		input: newPhasedEventBroadcasters[*EventInput](),
	}
	box.self = area
	area.editBuffer = opentui.NewEditBuffer(0)
	area.editBufferView = opentui.NewEditorView(area.editBuffer, 0, 0)

	area.setFocusable()
	area.SetWrap("word")

	// use input.Self() to be able to override the event handlers in derived elems like InputElement
	removePaste := area.pasteAction(func(event *EventPaste) {
		if i, ok := area.Self().(interface{ handlePaste(*EventPaste) }); ok {
			i.handlePaste(event)
		}
	})
	removeKeyPress := area.keyPressAction(func(event *EventKey) {
		if i, ok := area.Self().(interface{ handleKeyPress(*EventKey) }); ok {
			i.handleKeyPress(event)
		}
	})
	removeMousePress := area.mousePressAction(func(event *EventMouse) {
		if i, ok := area.Self().(interface{ handleMousePress(*EventMouse) }); ok {
			i.handleMousePress(event)
		}
	})
	removeMouseScroll := area.mouseScrollAction(func(event *EventMouse) {
		if i, ok := area.Self().(interface{ handleMouseScroll(*EventMouse) }); ok {
			i.handleMouseScroll(event)
		}
	})

	area.OnDestroy(func() {
		removePaste()
		removeKeyPress()
		removeMousePress()
		removeMouseScroll()
	})

	return area, nil
}

func (a *TextAreaElement) Value() (text string) {
	scheduleAccess(a.Self(), func() {
		a.mu.RLock()
		defer a.mu.RUnlock()

		length := a.editBuffer.GetTextBuffer().GetLength()
		text = a.editBuffer.GetText(int(length))
	})

	return
}

func (a *TextAreaElement) SetValue(text string) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.editBuffer.SetText(text)
		return nil
	})
}

func (a *TextAreaElement) InsertValue(text string) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.editBuffer.InsertText(text)
		return nil
	})
}

func (a *TextAreaElement) Clear() {
	a.SetValue("")
}

func (a *TextAreaElement) SetScrollFactor(factor float32) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.scrollFactor = factor
		return nil
	})
}

func (a *TextAreaElement) RemoveLeft() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.editBuffer.DeleteCharBackward()
		return nil
	})
}

func (a *TextAreaElement) RemoveRight() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.editBuffer.DeleteChar()
		return nil
	})
}

func (a *TextAreaElement) RemoveWordLeft() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		cursor := a.editBuffer.GetCursorPosition()
		prevWord := a.editBuffer.GetPrevWordBoundary()
		if prevWord.Offset < cursor.Offset {
			a.editBuffer.DeleteRange(prevWord.Row, prevWord.Col, cursor.Row, cursor.Col)
		}

		return nil
	})
}

func (a *TextAreaElement) RemoveWordRight() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		cursor := a.editBuffer.GetCursorPosition()
		nextWord := a.editBuffer.GetNextWordBoundary()
		if nextWord.Offset > cursor.Offset {
			a.editBuffer.DeleteRange(cursor.Row, cursor.Col, nextWord.Row, nextWord.Col)
		}

		return nil
	})
}

func (a *TextAreaElement) RemoveLine() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.editBuffer.DeleteLine()
		return nil
	})
}

func (a *TextAreaElement) RemoveToLineStart() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		cursor := a.editBuffer.GetCursorPosition()
		if cursor.Col > 0 {
			a.editBuffer.DeleteRange(cursor.Row, 0, cursor.Row, cursor.Col)
		}

		return nil
	})
}

func (a *TextAreaElement) RemoveToLineEnd() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		cursor := a.editBuffer.GetCursorPosition()
		eol := a.editBuffer.GetEOL()
		if cursor.Col < eol.Col {
			a.editBuffer.DeleteRange(cursor.Row, cursor.Col, cursor.Row, eol.Col)
		}

		return nil
	})
}

func (a *TextAreaElement) MoveCursor(offset int) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.editBuffer.SetCursorByOffset(uint32(offset))
		return nil
	})
}

func (a *TextAreaElement) GetCursor() (row, col uint32) {
	scheduleAccess(a.Self(), func() {
		a.mu.RLock()
		defer a.mu.RUnlock()

		row, col = a.editBuffer.GetCursor()
	})

	return
}

func (a *TextAreaElement) SetCursor(row, col uint32) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.editBuffer.SetCursor(row, col)
		return nil
	})
}

func (a *TextAreaElement) MoveCursorUp() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.editBuffer.MoveCursorUp()
		return nil
	})
}

func (a *TextAreaElement) MoveCursorDown() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.editBuffer.MoveCursorDown()
		return nil
	})
}

func (a *TextAreaElement) MoveCursorLeft() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.editBuffer.MoveCursorLeft()
		return nil
	})
}

func (a *TextAreaElement) MoveCursorRight() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.editBuffer.MoveCursorRight()
		return nil
	})
}

func (a *TextAreaElement) MoveCursorWordLeft() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		cursor := a.editBuffer.GetCursorPosition()
		prevWord := a.editBuffer.GetPrevWordBoundary()
		if prevWord.Offset < cursor.Offset {
			a.editBuffer.SetCursor(prevWord.Row, prevWord.Col)
		}

		return nil
	})
}

func (a *TextAreaElement) MoveCursorWordRight() {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		cursor := a.editBuffer.GetCursorPosition()
		nextWord := a.editBuffer.GetNextWordBoundary()
		if nextWord.Offset > cursor.Offset {
			a.editBuffer.SetCursor(nextWord.Row, nextWord.Col)
		}

		return nil
	})
}

// none | word | all
func (a *TextAreaElement) SetWrap(mode string) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		wrapMode, ok := wrapModes[mode]
		if !ok {
			return nil
		}

		a.wrapMode = wrapMode
		a.editBufferView.SetWrapMode(wrapMode)

		return nil
	})
}

func (a *TextAreaElement) UnsetWrap() {
	a.SetWrap("none")
}

func (a *TextAreaElement) SetTextForeground(color string) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.foregroundColor.Set(color)
		rgba := a.foregroundColor.RGBA()
		a.editBuffer.GetTextBuffer().SetDefaultFg(rgba)
		return nil
	})
}

func (a *TextAreaElement) UnsetTextForeground() {
	a.SetTextForeground("")
}

func (a *TextAreaElement) SetTextBackground(color string) {
	scheduleUpdate(a.Self(), func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.backgroundColor.Set(color)
		rgba := a.backgroundColor.RGBA()
		a.editBuffer.GetTextBuffer().SetDefaultBg(rgba)
		return nil
	})
}

func (a *TextAreaElement) UnsetTextBackground() {
	a.SetTextBackground("")
}

func (a *TextAreaElement) Render(buffer *opentui.Buffer, rect gfx.Rect) error {
	err := a.BoxElement.Render(buffer, rect)
	if err != nil {
		return err
	}

	a.editBufferView.SetViewportSize(uint32(rect.Width), uint32(rect.Height))
	err = buffer.DrawEditorView(a.editBufferView, int32(rect.X), int32(rect.Y))
	if err != nil {
		return err
	}

	if a.focused {
		vc := a.editBufferView.GetVisualCursor()
		cursorX := rect.X + int(vc.VisualCol) + 1
		cursorY := rect.Y + int(vc.VisualRow) + 1

		maxX := rect.X + rect.Width
		maxY := rect.Y + rect.Height

		a.rdrctx.SetCursorPosition(min(cursorX, maxX), min(cursorY, maxY), true)
	}

	return nil
}
