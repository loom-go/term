package elements

import (
	"fmt"
	"iter"
	"math"
	"slices"
	"strings"

	"github.com/AnatoleLucet/loom-term/core/gfx"

	"github.com/AnatoleLucet/go-opentui"
	"github.com/AnatoleLucet/tess"
)

var wrapModes = map[string]uint8{
	"":       opentui.WrapModeNone,
	"none":   opentui.WrapModeNone,
	"norwap": opentui.WrapModeNone,

	"word":   opentui.WrapModeWord,
	"normal": opentui.WrapModeWord,

	"all":  opentui.WrapModeChar,
	"char": opentui.WrapModeChar,
}

type TextElement struct {
	// todo: dont use base element. it's carying a lot more than needed
	*BaseElement

	rootText *TextElement

	chunk *TextChunk

	textStyle      *opentui.SyntaxStyle
	textBuffer     *opentui.TextBuffer
	textBufferView *opentui.TextBufferView

	children []*TextElement
}

func NewTextElement() (text *TextElement, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Text: %w: %v", ErrFailedToInitializeElement, err)
		}
	}()

	base, err := NewBaseElement()
	if err != nil {
		return nil, err
	}

	text = &TextElement{
		BaseElement: base,
		chunk:       NewTextChunk(),
	}
	base.self = text

	text.textBuffer = opentui.NewTextBuffer(0)
	text.textBufferView = opentui.NewTextBufferView(text.textBuffer)
	text.textStyle = opentui.NewSyntaxStyle()
	text.textBuffer.SetSyntaxStyle(text.textStyle)
	text.textBufferView.SetWrapMode(opentui.WrapModeWord)

	text.xyz().SetMeasureFunc(text.measureFunc)
	text.SetFlexGrow("0")
	text.SetFlexShrink("0")

	text.OnDestroy(func() {
		// note: text view must be closed before the buffer
		// https://github.com/anomalyco/opentui/blob/5958ce8060af43c0d4300cfbddeaf32d67bfb94c/packages/core/src/zig/text-buffer-view.zig#L208
		text.textBufferView.Close()
		text.textBuffer.Close()
		text.textStyle.Close()
	})

	return text, nil
}

func (t *TextElement) Children() (children iter.Seq[Element]) {
	scheduleAccess(t.Self(), func() {
		t.mu.RLock()
		defer t.mu.RUnlock()

		children = t.childrenUnsafe()
	})

	return
}

func (t *TextElement) childrenUnsafe() iter.Seq[Element] {
	return func(yield func(Element) bool) {
		for _, child := range t.children {
			if !yield(child) {
				return
			}
		}
	}
}

func (t *TextElement) AppendChild(child Element) {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		err := t.appendChildUnsafe(child)
		if err != nil {
			return err
		}

		t.update()
		return nil
	})
}

func (t *TextElement) appendChildUnsafe(child Element) error {
	c, ok := child.(*TextElement)
	if !ok {
		return nil
	}

	if c.parentUnsafe() != nil {
		return fmt.Errorf("%w: child already has a parent", ErrFailedToAppendChild)
	}

	t.children = append(t.children, c)

	c.rootText = t.rootTextElement()
	c.setParentUnsafe(t.Self())
	c.SetRenderContext(t.rdrctx)
	return nil
}

func (t *TextElement) RemoveChild(child Element) {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		err := t.removeChildUnsafe(child)
		if err != nil {
			return err
		}

		t.update()
		return nil
	})
}

func (t *TextElement) removeChildUnsafe(child Element) error {
	c, ok := child.(*TextElement)
	if !ok {
		return nil
	}

	i := slices.Index(t.children, c)
	if i == -1 {
		return fmt.Errorf("%w: child not found", ErrFailedToRemoveChild)
	}

	t.children = slices.Delete(t.children, i, i+1)

	c.rootText = nil
	c.setParentUnsafe(nil)
	c.SetRenderContext(nil)
	return nil
}

func (t *TextElement) Record(cb *gfx.CommandBuffer, container gfx.Rect) error {
	self := t.Self()
	rect := t.rect(container)

	render := gfx.NewCommand(gfx.CmdRender, self).WithRect(rect).WithCallback(func() error {
		return t.rdrctx.AddToHitGrid(self, rect)
	})
	cb.Add(render)

	return nil
}

func (t *TextElement) Render(buffer *opentui.Buffer, rect gfx.Rect) error {
	t.textBufferView.SetWrapWidth(uint32(rect.Width))
	return buffer.DrawTextBufferView(t.textBufferView, int32(rect.X), int32(rect.Y))
}

func (t *TextElement) Text() (text string) {
	scheduleAccess(t.Self(), func() {
		t.mu.RLock()
		defer t.mu.RUnlock()

		text = t.text()
	})

	return
}

func (t *TextElement) SetText(text string) {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.chunk.SetText(text)
		t.update()
		return nil
	})
}

func (t *TextElement) UnsetText() {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.chunk.UnsetText()
		t.update()
		return nil
	})
}

// normal | bold
func (t *TextElement) SetFontWeight(weight string) {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.chunk.SetFontWeight(weight)
		t.update()
		return nil
	})
}

func (t *TextElement) UnsetFontWeight() {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.chunk.UnsetFontWeight()
		t.update()
		return nil
	})
}

// normal | italic
func (t *TextElement) SetFontStyle(style string) {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.chunk.SetFontStyle(style)
		t.update()
		return nil
	})
}

func (t *TextElement) UnsetFontStyle() {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.chunk.UnsetFontStyle()
		t.update()
		return nil
	})
}

// none | underline | strikethrough
func (t *TextElement) SetTextDecoration(decoration string) {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.chunk.SetTextDecoration(decoration)
		t.update()
		return nil
	})
}

func (t *TextElement) UnsetTextDecoration() {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.chunk.UnsetTextDecoration()
		t.update()
		return nil
	})
}

// none | word | all
func (t *TextElement) SetWrap(mode string) {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		wrapMode, ok := wrapModes[mode]
		if !ok {
			return nil
		}

		t.textBufferView.SetWrapMode(wrapMode)
		t.update()
		return nil
	})
}

func (t *TextElement) UnsetWrap() {
	t.SetWrap("none")
}

func (t *TextElement) SetTextForeground(color string) {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.chunk.SetTextForeground(color)
		t.update()
		return nil
	})
}

func (t *TextElement) UnsetTextForeground() {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.chunk.UnsetTextForeground()
		t.update()
		return nil
	})
}

func (t *TextElement) SetTextBackground(color string) {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.chunk.SetTextBackground(color)
		t.update()
		return nil
	})
}

func (t *TextElement) UnsetTextBackground() {
	scheduleUpdate(t.Self(), func() error {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.chunk.UnsetTextBackground()
		t.update()
		return nil
	})
}

func (t *TextElement) rootTextElement() *TextElement {
	if t.rootText != nil {
		return t.rootText
	}

	return t
}

func (t *TextElement) text() string {
	var sb strings.Builder
	sb.WriteString(t.chunk.Text())

	for _, child := range t.children {
		sb.WriteString(child.text())
	}

	return sb.String()
}

func (t *TextElement) update() {
	root := t.rootTextElement()

	root.xyz().MarkDirty()
	root.textBuffer.Reset()
	root.textBuffer.SetStyledText(root.chunks(nil))
}

func (t *TextElement) chunks(parent *TextElement) []opentui.StyledChunk {
	chunks := []opentui.StyledChunk{}
	if t.chunk.Text() != "" {
		chunks = append(chunks, t.chunk.StyledChunk(t.chunk))
	}

	for _, child := range t.children {
		for _, chunk := range child.chunks(t) {
			if chunk.Text != "" {
				chunks = append(chunks, chunk)
			}
		}
	}

	return chunks
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
