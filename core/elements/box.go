package elements

import (
	"github.com/AnatoleLucet/go-opentui"
	"github.com/AnatoleLucet/loom-term/core/types"
)

type BoxElement struct {
	*BaseElement

	bgcolor *Color
}

func NewBoxElement(ctx types.RenderContext) (*BoxElement, error) {
	base, err := NewElement(ctx)
	if err != nil {
		return nil, err
	}

	e := &BoxElement{
		BaseElement: base,
		bgcolor:     &Color{"transparent", opentui.Transparent},
	}

	return e, nil
}

func (e *BoxElement) SetBackgroundColor(color string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	return e.bgcolor.Set(color)
}

func (e *BoxElement) UnsetBackgroundColor() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	return e.bgcolor.Set("transparent")
}

func (e *BoxElement) Paint(buffer *opentui.Buffer, x, y float32) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if err := e.guardPaint(); err != nil {
		return err
	}

	if err := e.paintBox(buffer, x, y); err != nil {
		return err
	}

	if err := e.paint(buffer, x, y); err != nil {
		return err
	}

	return e.paintChildren(buffer, x, y)
}

func (e *BoxElement) paintBox(buffer *opentui.Buffer, x, y float32) error {
	layout := e.Layout()

	bgcolor := e.bgcolor.RGBA()
	if bgcolor.A > 0 {
		buffer.FillRect(
			uint32(x),
			uint32(y),
			uint32(layout.Width()),
			uint32(layout.Height()),
			bgcolor,
		)
	}

	return nil
}
