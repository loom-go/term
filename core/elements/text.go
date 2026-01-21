package elements

import (
	"errors"
	"math"

	"github.com/AnatoleLucet/go-opentui"
	"github.com/AnatoleLucet/loom-term/core/types"
	"github.com/AnatoleLucet/tess"
)

type TextElement struct {
	*BaseElement

	textBuffer     *opentui.TextBuffer
	textBufferView *opentui.TextBufferView
}

func NewTextElement(ctx types.RenderContext) (*TextElement, error) {
	base, err := NewElement(ctx)
	if err != nil {
		return nil, err
	}

	tb := opentui.NewTextBuffer(0)
	if tb == nil {
		return nil, errors.New("failed to create text buffer") // todo: better error
	}
	tb.SetDefaultFg(opentui.White)
	tb.SetDefaultBg(opentui.Transparent)

	view := opentui.NewTextBufferView(tb)
	if view == nil {
		tb.Close()
		return nil, errors.New("failed to create text buffer view") // todo: better error
	}
	view.SetWrapMode(opentui.WrapModeWord)

	e := &TextElement{
		BaseElement: base,

		textBuffer:     tb,
		textBufferView: view,
	}

	e.XYZ().SetMeasureFunc(e.measureFunc)
	e.SetFlexShrink("0")
	e.SetFlexGrow("0")

	return e, nil
}

func (t *TextElement) SetContent(content string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.guardUpdate(); err != nil {
		return err
	}

	t.textBuffer.Reset()
	t.textBuffer.Append(content)
	t.xyz.MarkDirty()

	return nil
}

func (t *TextElement) SetForegroundColor(color string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.guardUpdate(); err != nil {
		return err
	}

	c, err := toOpenTUIColor(color)
	if err != nil {
		return err
	}

	t.textBuffer.SetDefaultBg(c)
	return nil
}

func (t *TextElement) UnsetForegroundColor() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.guardUpdate(); err != nil {
		return err
	}

	t.textBuffer.SetDefaultFg(opentui.White)
	return nil
}

func (t *TextElement) SetBackgroundColor(color string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.guardUpdate(); err != nil {
		return err
	}

	c, err := toOpenTUIColor(color)
	if err != nil {
		return err
	}

	t.textBuffer.SetDefaultBg(c)
	return nil
}

func (t *TextElement) UnsetBackgroundColor() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.guardUpdate(); err != nil {
		return err
	}

	t.textBuffer.SetDefaultBg(opentui.Transparent)
	return nil
}

func (e *TextElement) Paint(buffer *opentui.Buffer, x, y float32) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if err := e.guardPaint(); err != nil {
		return err
	}

	layout := e.Layout()
	width := layout.Width()
	height := layout.Height()

	e.textBufferView.SetViewportSize(uint32(width), uint32(height))

	err := buffer.DrawTextBufferView(e.textBufferView, int32(x), int32(y))
	if err != nil {
		return err
	}

	if err := e.paint(buffer, x, y); err != nil {
		return err
	}

	return e.paintChildren(buffer, x, y)
}

func (t *TextElement) Destroy() error {
	t.mu.Lock()

	if t.destroyed {
		t.mu.Unlock()
		return nil
	}

	// beffer view must be closed before the buffer
	// https://github.com/anomalyco/opentui/blob/5958ce8060af43c0d4300cfbddeaf32d67bfb94c/packages/core/src/zig/text-buffer-view.zig#L208
	bv := t.textBufferView
	if bv != nil {
		bv.Close()
		t.textBufferView = nil
	}

	tb := t.textBuffer
	if tb != nil {
		tb.Close()
		t.textBuffer = nil
	}

	t.mu.Unlock()
	return t.BaseElement.Destroy()
}

// source: https://github.com/anomalyco/opentui/blob/9092e7c366ee04ceec208dddc74bd49efc632d2f/packages/core/src/renderables/TextBufferRenderable.ts#L376-L416
func (t *TextElement) measureFunc(node *tess.Node, width float32, widthMode tess.MeasureMode, height float32, heightMode tess.MeasureMode) tess.Size {
	var effectiveWidth uint32
	if widthMode == tess.MeasureModeUndefined || math.IsNaN(float64(width)) {
		effectiveWidth = 0
	} else {
		effectiveWidth = uint32(width)
	}

	var effectiveHeight uint32
	if math.IsNaN(float64(height)) {
		effectiveHeight = 1
	} else {
		effectiveHeight = uint32(height)
	}

	outWidth, outHeight, _ := t.textBufferView.MeasureForDimensions(effectiveWidth, effectiveHeight)

	measuredWidth := max(1, float32(outWidth))
	measuredHeight := max(1, float32(outHeight))

	if widthMode == tess.MeasureModeAtMost && node.GetPosition() == tess.Absolute {
		return tess.Size{
			Width:  min(float32(effectiveWidth), measuredWidth),
			Height: min(float32(effectiveHeight), measuredHeight),
		}
	}

	return tess.Size{
		Width:  measuredWidth,
		Height: measuredHeight,
	}
}
