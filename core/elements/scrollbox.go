package elements

import (
	"fmt"
	"iter"

	"github.com/AnatoleLucet/loom-term/core/elements/events"
	"github.com/AnatoleLucet/loom-term/core/gfx"
)

type scrollAction int

const (
	scrollActionTop scrollAction = iota
	scrollActionBottom
	scrollActionLeft
	scrollActionRight
)

type ScrollBoxElement struct {
	*BoxElement

	container *BaseElement
	content   *BaseElement

	factorX float32
	factorY float32

	scrollY float32
	scrollX float32

	// we cannot scroll directly when calling a ScrollToX method,
	// because the layout might not be up to date yet.
	// so we store the action, and run it during the record phase (after the layout has been computed)
	pendingActions []scrollAction
}

func NewScrollBoxElement() (scrollb *ScrollBoxElement, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("ScrollBox: %w: %v", ErrFailedToInitializeElement, err)
		}
	}()

	box, err := NewBoxElement()
	if err != nil {
		return nil, err
	}
	box.SetOverflow("hidden")

	container, err := NewBaseElement()
	if err != nil {
		return nil, err
	}
	container.SetAlignSelf("start")
	container.SetFlexShrink("0")
	container.SetFlexGrow("1")
	box.AppendChild(container)

	// use an inner content element to prevent https://github.com/facebook/yoga/issues/872 (via tess)
	content, err := NewBaseElement()
	if err != nil {
		return nil, err
	}
	content.SetFlexGrow("1")
	container.AppendChild(content)

	scrollb = &ScrollBoxElement{
		BoxElement: box,
		container:  container,
		content:    content,
		factorX:    2,
		factorY:    1,
	}
	box.self = scrollb

	remove := scrollb.mouseScrollAction(func(event *EventMouse) {
		scrollb.mu.Lock()
		oldScrollY := scrollb.scrollY
		oldScrollX := scrollb.scrollX

		if event.Button == events.MouseWheelUp {
			delta := 1 * scrollb.factorY
			scrollb.scrollX, scrollb.scrollY = scrollb.clamp(scrollb.scrollX, scrollb.scrollY-delta)
		}
		if event.Button == events.MouseWheelDown {
			delta := 1 * scrollb.factorY
			scrollb.scrollX, scrollb.scrollY = scrollb.clamp(scrollb.scrollX, scrollb.scrollY+delta)
		}
		if event.Button == events.MouseWheelLeft {
			delta := 1 * scrollb.factorX
			scrollb.scrollX, scrollb.scrollY = scrollb.clamp(scrollb.scrollX-delta, scrollb.scrollY)
		}
		if event.Button == events.MouseWheelRight {
			delta := 1 * scrollb.factorX
			scrollb.scrollX, scrollb.scrollY = scrollb.clamp(scrollb.scrollX+delta, scrollb.scrollY)
		}

		newScrollY := scrollb.scrollY
		newScrollX := scrollb.scrollX
		scrollb.mu.Unlock()

		if oldScrollY != newScrollY || oldScrollX != newScrollX {
			scheduleUpdate(scrollb.Self(), func() error { return nil }) // schedule a nil update to trigger a render
		}
	})
	scrollb.OnDestroy(remove)

	return scrollb, nil
}

func (e *ScrollBoxElement) Children() iter.Seq[Element] {
	return e.content.Children()
}

func (e *ScrollBoxElement) AppendChild(child Element) {
	e.content.AppendChild(child)
}

func (e *ScrollBoxElement) RemoveChild(child Element) {
	e.content.RemoveChild(child)
}

func (e *ScrollBoxElement) SetScrollFactorX(factor float32) {
	scheduleUpdate(e.Self(), func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		e.factorX = factor
		return nil
	})
}

func (e *ScrollBoxElement) SetScrollFactorY(factor float32) {
	scheduleUpdate(e.Self(), func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		e.factorY = factor
		return nil
	})
}

func (e *ScrollBoxElement) ScrollY() (scrollY float32) {
	scheduleAccess(e.Self(), func() {
		e.mu.RLock()
		defer e.mu.RUnlock()

		scrollY = e.scrollY
	})

	return
}

func (e *ScrollBoxElement) ScrollX() (scrollX float32) {
	scheduleAccess(e.Self(), func() {
		e.mu.RLock()
		defer e.mu.RUnlock()

		scrollX = e.scrollX
	})

	return
}

func (e *ScrollBoxElement) ViewportHeight() (height float32) {
	scheduleAccess(e.Self(), func() {
		e.mu.RLock()
		defer e.mu.RUnlock()

		_, height = e.viewportSize()
	})

	return
}

func (e *ScrollBoxElement) ViewportWidth() (width float32) {
	scheduleAccess(e.Self(), func() {
		e.mu.RLock()
		defer e.mu.RUnlock()

		width, _ = e.viewportSize()
	})

	return
}

func (e *ScrollBoxElement) ContentHeight() (height float32) {
	scheduleAccess(e.Self(), func() {
		e.mu.RLock()
		defer e.mu.RUnlock()

		_, height = e.contentSize()
	})

	return
}

func (e *ScrollBoxElement) ContentWidth() (width float32) {
	scheduleAccess(e.Self(), func() {
		e.mu.RLock()
		defer e.mu.RUnlock()

		width, _ = e.contentSize()
	})

	return
}

func (e *ScrollBoxElement) ScrollTo(x, y float32) {
	scheduleUpdate(e.Self(), func() error {
		e.mu.RLock()
		defer e.mu.RUnlock()

		e.scrollX, e.scrollY = e.clamp(x, y)
		return nil
	})
}

func (e *ScrollBoxElement) ScrollBy(dx, dy float32) {
	scheduleUpdate(e.Self(), func() error {
		e.mu.RLock()
		defer e.mu.RUnlock()

		e.scrollX, e.scrollY = e.clamp(e.scrollX+dx, e.scrollY+dy)
		return nil
	})
}

func (e *ScrollBoxElement) ScrollToTop() {
	scheduleUpdate(e.Self(), func() error {
		e.mu.Lock()
		e.pendingActions = append(e.pendingActions, scrollActionTop)
		e.mu.Unlock()

		return nil
	})
}

func (e *ScrollBoxElement) scrollToTop() {
	e.scrollX, e.scrollY = e.clamp(e.scrollX, 0)
}

func (e *ScrollBoxElement) ScrollToBottom() {
	scheduleUpdate(e.Self(), func() error {
		e.mu.Lock()
		e.pendingActions = append(e.pendingActions, scrollActionBottom)
		e.mu.Unlock()

		return nil
	})
}

func (e *ScrollBoxElement) scrollToBottom() {
	_, viewportH := e.viewportSize()
	_, contentH := e.contentSize()
	e.scrollX, e.scrollY = e.clamp(e.scrollX, contentH-viewportH)
}

func (e *ScrollBoxElement) ScrollToLeft() {
	scheduleUpdate(e.Self(), func() error {
		e.mu.Lock()
		e.pendingActions = append(e.pendingActions, scrollActionLeft)
		e.mu.Unlock()

		return nil
	})
}

func (e *ScrollBoxElement) scrollToLeft() error {
	e.scrollX, e.scrollY = e.clamp(0, e.scrollY)
	return nil
}

func (e *ScrollBoxElement) ScrollToRight() {
	scheduleUpdate(e.Self(), func() error {
		e.mu.Lock()
		e.pendingActions = append(e.pendingActions, scrollActionRight)
		e.mu.Unlock()

		return nil
	})
}

func (e *ScrollBoxElement) scrollToRight() error {
	viewportW, _ := e.viewportSize()
	contentW, _ := e.contentSize()
	e.scrollX, e.scrollY = e.clamp(contentW-viewportW, e.scrollY)
	return nil
}

func (e *ScrollBoxElement) IsAtScrollEdge() (atTop, atBottom, atLeft, atRight bool) {
	scheduleAccess(e.Self(), func() {
		e.mu.RLock()
		defer e.mu.RUnlock()

		viewportW, viewportH := e.viewportSize()
		contentW, contentH := e.contentSize()

		atTop = e.scrollY <= 0
		atBottom = e.scrollY >= contentH-viewportH
		atLeft = e.scrollX <= 0
		atRight = e.scrollX >= contentW-viewportW
	})

	return
}

func (b *ScrollBoxElement) Record(cb *gfx.CommandBuffer, container gfx.Rect) error {
	for _, action := range b.pendingActions {
		switch action {
		case scrollActionTop:
			b.scrollToTop()
		case scrollActionBottom:
			b.scrollToBottom()
		case scrollActionLeft:
			b.scrollToLeft()
		case scrollActionRight:
			b.scrollToRight()
		}
	}
	b.pendingActions = nil

	// clamp again in case the element was resized since the last scroll action
	x, y := b.clamp(b.scrollX, b.scrollY)

	b.container.setTranslateUnsafe(-x, -y)
	return b.BoxElement.Record(cb, container)
}

func (e *ScrollBoxElement) clamp(x, y float32) (scrollX, scrollY float32) {
	viewportW, viewportH := e.viewportSize()
	contentW, contentH := e.contentSize()

	maxScrollY := max(0, contentH-viewportH)
	maxScrollX := max(0, contentW-viewportW)

	scrollY = min(max(y, 0), maxScrollY)
	scrollX = min(max(x, 0), maxScrollX)

	return scrollX, scrollY
}

func (e *ScrollBoxElement) viewportSize() (width, height float32) {
	l := e.xyz().GetLayout()
	return l.Width(), l.Height()
}

func (e *ScrollBoxElement) contentSize() (width, height float32) {
	l := e.content.xyz().GetLayout()
	return l.Width(), l.Height()
}

func (e *ScrollBoxElement) SetOverflow(overflow string) {
	scheduleUpdate(e.Self(), func() error { return nil })
}
func (e *ScrollBoxElement) UnsetOverflow() {
	scheduleUpdate(e.Self(), func() error { return nil })
}

func (e *ScrollBoxElement) SetPaddingAll(padding any) {
	e.content.SetPaddingAll(padding)
}

func (e *ScrollBoxElement) UnsetPaddingAll() {
	e.content.UnsetPaddingAll()
}

func (e *ScrollBoxElement) SetPaddingVertical(padding any) {
	e.content.SetPaddingVertical(padding)
}

func (e *ScrollBoxElement) UnsetPaddingVertical() {
	e.content.UnsetPaddingVertical()
}

func (e *ScrollBoxElement) SetPaddingHorizontal(padding any) {
	e.content.SetPaddingHorizontal(padding)
}

func (e *ScrollBoxElement) SetPaddingTop(padding any) {
	e.content.SetPaddingTop(padding)
}

func (e *ScrollBoxElement) UnsetPaddingTop() {
	e.content.UnsetPaddingTop()
}

func (e *ScrollBoxElement) SetPaddingRight(padding any) {
	e.content.SetPaddingRight(padding)
}

func (e *ScrollBoxElement) UnsetPaddingRight() {
	e.content.UnsetPaddingRight()
}

func (e *ScrollBoxElement) SetPaddingBottom(padding any) {
	e.content.SetPaddingBottom(padding)
}

func (e *ScrollBoxElement) UnsetPaddingBottom() {
	e.content.UnsetPaddingBottom()
}

func (e *ScrollBoxElement) SetPaddingLeft(padding any) {
	e.content.SetPaddingLeft(padding)
}

func (e *ScrollBoxElement) UnsetPaddingLeft() {
	e.content.UnsetPaddingLeft()
}

func (e *ScrollBoxElement) SetFlexDirection(direction string) {
	e.content.SetFlexDirection(direction)
}

func (e *ScrollBoxElement) UnsetFlexDirection() {
	e.content.UnsetFlexDirection()
}

func (e *ScrollBoxElement) SetFlexWrap(wrap string) {
	e.content.SetFlexWrap(wrap)
}

func (e *ScrollBoxElement) UnsetFlexWrap() {
	e.content.UnsetFlexWrap()
}

func (e *ScrollBoxElement) SetGapAll(gap any) {
	e.content.SetGapAll(gap)
}

func (e *ScrollBoxElement) UnsetGapAll() {
	e.content.UnsetGapAll()
}

func (e *ScrollBoxElement) SetGapRow(gap any) {
	e.content.SetGapRow(gap)
}

func (e *ScrollBoxElement) UnsetGapRow() {
	e.content.UnsetGapRow()
}

func (e *ScrollBoxElement) SetGapColumn(gap any) {
	e.content.SetGapColumn(gap)
}

func (e *ScrollBoxElement) UnsetGapColumn() {
	e.content.UnsetGapColumn()
}

func (e *ScrollBoxElement) SetAlignItems(alignItems string) {
	e.content.SetAlignItems(alignItems)
}

func (e *ScrollBoxElement) UnsetAlignItems() {
	e.content.UnsetAlignItems()
}

func (e *ScrollBoxElement) SetAlignContent(alignContent string) {
	e.content.SetAlignContent(alignContent)
}

func (e *ScrollBoxElement) UnsetAlignContent() {
	e.content.UnsetAlignContent()
}

func (e *ScrollBoxElement) SetJustifyContent(justifyContent string) {
	e.content.SetJustifyContent(justifyContent)
}

func (e *ScrollBoxElement) UnsetJustifyContent() {
	e.content.UnsetJustifyContent()
}
