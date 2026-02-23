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

func (a *TextAreaElement) Text() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.editBuffer.GetText(0)
}

func (a *TextAreaElement) SetText(text string) {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.editBuffer.SetText(text)
		return nil
	})
}

func (a *TextAreaElement) InsertText(text string) {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.editBuffer.InsertText(text)
		return nil
	})
}

func (a *TextAreaElement) SetScrollFactor(factor float32) {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.scrollFactor = factor
		return nil
	})
}

func (a *TextAreaElement) RemoveLeft() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.editBuffer.DeleteCharBackward()
		return nil
	})
}

func (a *TextAreaElement) RemoveRight() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.editBuffer.DeleteChar()
		return nil
	})
}

func (a *TextAreaElement) RemoveWordLeft() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		cursor := a.editBuffer.GetCursorPosition()
		prevWord := a.editBuffer.GetPrevWordBoundary()
		if prevWord.Offset < cursor.Offset {
			a.editBuffer.DeleteRange(prevWord.Row, prevWord.Col, cursor.Row, cursor.Col)
		}

		return nil
	})
}

func (a *TextAreaElement) RemoveWordRight() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		cursor := a.editBuffer.GetCursorPosition()
		nextWord := a.editBuffer.GetNextWordBoundary()
		if nextWord.Offset > cursor.Offset {
			a.editBuffer.DeleteRange(cursor.Row, cursor.Col, nextWord.Row, nextWord.Col)
		}

		return nil
	})
}

func (a *TextAreaElement) RemoveLine() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.editBuffer.DeleteLine()
		return nil
	})
}

func (a *TextAreaElement) RemoveToLineStart() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		cursor := a.editBuffer.GetCursorPosition()
		if cursor.Col > 0 {
			a.editBuffer.DeleteRange(cursor.Row, 0, cursor.Row, cursor.Col)
		}

		return nil
	})
}

func (a *TextAreaElement) RemoveToLineEnd() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		cursor := a.editBuffer.GetCursorPosition()
		eol := a.editBuffer.GetEOL()
		if cursor.Col < eol.Col {
			a.editBuffer.DeleteRange(cursor.Row, cursor.Col, cursor.Row, eol.Col)
		}

		return nil
	})
}

func (a *TextAreaElement) MoveCursor(offset int) {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.editBuffer.SetCursorByOffset(uint32(offset))
		return nil
	})
}

func (a *TextAreaElement) GetCursor() (row, col uint32) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.editBuffer.GetCursor()
}

func (a *TextAreaElement) SetCursor(row, col uint32) {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.editBuffer.SetCursor(row, col)
		return nil
	})
}

func (a *TextAreaElement) MoveCursorUp() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.editBuffer.MoveCursorUp()
		return nil
	})
}

func (a *TextAreaElement) MoveCursorDown() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.editBuffer.MoveCursorDown()
		return nil
	})
}

func (a *TextAreaElement) MoveCursorLeft() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.editBuffer.MoveCursorLeft()
		return nil
	})
}

func (a *TextAreaElement) MoveCursorRight() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.editBuffer.MoveCursorRight()
		return nil
	})
}

func (a *TextAreaElement) MoveCursorWordLeft() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		cursor := a.editBuffer.GetCursorPosition()
		prevWord := a.editBuffer.GetPrevWordBoundary()
		if prevWord.Offset < cursor.Offset {
			a.editBuffer.SetCursor(prevWord.Row, prevWord.Col)
		}

		return nil
	})
}

func (a *TextAreaElement) MoveCursorWordRight() {
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

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
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

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
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

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
	a.scheduleUpdate(func() error {
		a.mu.Lock()
		defer a.mu.Unlock()

		if err := guardDestroyed(a.ctx); err != nil {
			return err
		}

		a.backgroundColor.Set(color)
		rgba := a.backgroundColor.RGBA()
		a.editBuffer.GetTextBuffer().SetDefaultBg(rgba)
		return nil
	})
}

func (a *TextAreaElement) UnsetTextBackground() {
	a.SetTextBackground("")
}

func (a *TextAreaElement) Submit() error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if err := guardDestroyed(a.ctx); err != nil {
		return err
	}

	evt := &EventSubmit{Text: a.Text()}
	evt.setTarget(a.Self())

	a.rdrctx.DispatchEvent(EventTypeSubmit, a, evt)
	return nil
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
