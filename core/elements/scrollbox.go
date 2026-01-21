package elements

import (
	"github.com/AnatoleLucet/go-opentui"
	"github.com/AnatoleLucet/loom-term/core/types"
)

type ScrollBoxElement struct {
	*BoxElement

	factor float32

	scrollY float32
	scrollX float32
}

func NewScrollBoxElement(ctx types.RenderContext) (*ScrollBoxElement, error) {
	box, err := NewBoxElement(ctx)
	if err != nil {
		return nil, err
	}

	e := &ScrollBoxElement{
		BoxElement: box,
		factor:     1.0,
	}
	e.SetOverflow("hidden")

	e.OnMouseScroll(func(event *types.EventMouse) {
		e.mu.Lock()
		if err := e.guardUpdate(); err != nil {
			return
		}
		oldScrollY := e.scrollY
		oldScrollX := e.scrollX

		if event.Button == types.MouseButtonWheelUp {
			e.scrollX, e.scrollY = e.clamp(e.scrollX, e.scrollY-1*e.factor)
		}
		if event.Button == types.MouseButtonWheelDown {
			e.scrollX, e.scrollY = e.clamp(e.scrollX, e.scrollY+1*e.factor)
		}
		if event.Button == types.MouseButtonWheelLeft {
			e.scrollX, e.scrollY = e.clamp(e.scrollX-1*e.factor, e.scrollY)
		}
		if event.Button == types.MouseButtonWheelRight {
			e.scrollX, e.scrollY = e.clamp(e.scrollX+1*e.factor, e.scrollY)
		}

		newScrollY := e.scrollY
		newScrollX := e.scrollX
		e.mu.Unlock()

		if oldScrollY != newScrollY || oldScrollX != newScrollX {
			event.StopPropagation()
			ctx.Render()
		}
	})

	return e, nil
}

func (e *ScrollBoxElement) SetScrollFactor(factor float32) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return err
	}

	e.factor = factor
	return nil
}

func (e *ScrollBoxElement) ScrollTo(x, y float32) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.guardUpdate(); err != nil {
		return
	}

	e.scrollX, e.scrollY = e.clamp(x, y)
}

func (e *ScrollBoxElement) Paint(buffer *opentui.Buffer, x, y float32) error {
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

func (e *ScrollBoxElement) paintChildren(buffer *opentui.Buffer, x, y float32) error {
	if len(e.children) == 0 {
		return nil
	}

	// clamp again in case content or viewport size changed since last scroll
	scrollX, scrollY := e.clamp(e.scrollX, e.scrollY)

	return e.withClip(buffer, x, y, func() error {
		for child := range e.childrenUnsafe() {
			l := child.Layout()

			err := child.Paint(buffer, x+l.Left()-scrollX, y+l.Top()-scrollY)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (e *ScrollBoxElement) clamp(x, y float32) (scrollX, scrollY float32) {
	viewportW, viewportH := e.viewportSize()
	contentMinX, contentMinY, contentMaxX, contentMaxY := e.contentBounds()

	minScrollY := min(0, contentMinY)
	maxScrollY := max(0, contentMaxY-viewportH)
	minScrollX := min(0, contentMinX)
	maxScrollX := max(0, contentMaxX-viewportW)

	scrollY = min(max(y, minScrollY), maxScrollY)
	scrollX = min(max(x, minScrollX), maxScrollX)

	return scrollX, scrollY
}

func (e *ScrollBoxElement) viewportSize() (width, height float32) {
	layout := e.xyz.GetLayout()

	viewportW := layout.Width() - layout.Padding().Left() - layout.Padding().Right()
	viewportH := layout.Height() - layout.Padding().Top() - layout.Padding().Bottom()

	return viewportW, viewportH
}

func (e *ScrollBoxElement) contentBounds() (minX, minY, maxX, maxY float32) {
	for child := range e.childrenUnsafe() {
		l := child.Layout()
		minY = min(minY, l.Top())
		maxY = max(maxY, l.Top()+l.Height())
		minX = min(minX, l.Left())
		maxX = max(maxX, l.Left()+l.Width())
	}

	return minX, minY, maxX, maxY
}
