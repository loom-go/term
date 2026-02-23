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

	factor float32

	scrollY float32
	scrollX float32

	// we cannot scroll directly when calling a ScrollToX method,
	// because the layout might not be up to date yet.
	// so we store the action, and run it during the record phase.
	pendingAction *scrollAction
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
	container.SetMinWidth("100%")
	container.SetMinHeight("100%")
	container.SetAlignSelf("start")
	container.SetFlexShrink("0")
	container.SetFlexGrow("0")
	box.AppendChild(container)

	// use an inner content element to prevent https://github.com/facebook/yoga/issues/872 (via tess)
	content, err := NewBaseElement()
	if err != nil {
		return nil, err
	}
	container.AppendChild(content)

	scrollb = &ScrollBoxElement{
		BoxElement: box,
		container:  container,
		content:    content,
		factor:     1,
	}
	box.self = scrollb

	remove := scrollb.mouseScrollAction(func(event *EventMouse) {
		scrollb.mu.Lock()
		oldScrollY := scrollb.scrollY
		oldScrollX := scrollb.scrollX

		delta := 1 * scrollb.factor

		if event.Button == events.MouseWheelUp {
			scrollb.scrollX, scrollb.scrollY = scrollb.clamp(scrollb.scrollX, scrollb.scrollY-delta)
		}
		if event.Button == events.MouseWheelDown {
			scrollb.scrollX, scrollb.scrollY = scrollb.clamp(scrollb.scrollX, scrollb.scrollY+delta)
		}
		if event.Button == events.MouseWheelLeft {
			scrollb.scrollX, scrollb.scrollY = scrollb.clamp(scrollb.scrollX-delta, scrollb.scrollY)
		}
		if event.Button == events.MouseWheelRight {
			scrollb.scrollX, scrollb.scrollY = scrollb.clamp(scrollb.scrollX+delta, scrollb.scrollY)
		}

		newScrollY := scrollb.scrollY
		newScrollX := scrollb.scrollX
		scrollb.mu.Unlock()

		if oldScrollY != newScrollY || oldScrollX != newScrollX {
			scrollb.rdrctx.ScheduleRender()
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

func (e *ScrollBoxElement) SetScrollFactor(factor float32) {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		e.factor = factor
		return nil
	})
}

func (e *ScrollBoxElement) ScrollY() float32 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.scrollY
}

func (e *ScrollBoxElement) ScrollX() float32 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.scrollX
}

func (e *ScrollBoxElement) ViewportHeight() float32 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	_, h := e.viewportSize()
	return h
}

func (e *ScrollBoxElement) ViewportWidth() float32 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	w, _ := e.viewportSize()
	return w
}

func (e *ScrollBoxElement) ContentHeight() float32 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	_, h := e.contentSize()
	return h
}

func (e *ScrollBoxElement) ContentWidth() float32 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	w, _ := e.contentSize()
	return w
}

func (e *ScrollBoxElement) ScrollTo(x, y float32) {
	e.scheduleUpdate(func() error {
		e.mu.RLock()
		defer e.mu.RUnlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		e.scrollX, e.scrollY = e.clamp(x, y)
		return nil
	})
}

func (e *ScrollBoxElement) ScrollBy(dx, dy float32) {
	e.scheduleUpdate(func() error {
		e.mu.RLock()
		defer e.mu.RUnlock()

		if err := guardDestroyed(e.ctx); err != nil {
			return err
		}

		e.scrollX, e.scrollY = e.clamp(e.scrollX+dx, e.scrollY+dy)
		return nil
	})
}

func (e *ScrollBoxElement) ScrollToTop() {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		action := scrollActionTop
		e.pendingAction = &action
		e.mu.Unlock()

		return nil
	})
}

func (e *ScrollBoxElement) scrollToTop() {
	e.scrollX, e.scrollY = e.clamp(e.scrollX, 0)
}

func (e *ScrollBoxElement) ScrollToBottom() {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		action := scrollActionBottom
		e.pendingAction = &action
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
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		action := scrollActionLeft
		e.pendingAction = &action
		e.mu.Unlock()

		return nil
	})
}

func (e *ScrollBoxElement) scrollToLeft() error {
	e.scrollX, e.scrollY = e.clamp(0, e.scrollY)
	return nil
}

func (e *ScrollBoxElement) ScrollToRight() {
	e.scheduleUpdate(func() error {
		e.mu.Lock()
		action := scrollActionRight
		e.pendingAction = &action
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
	e.mu.RLock()
	defer e.mu.RUnlock()

	viewportW, viewportH := e.viewportSize()
	contentW, contentH := e.contentSize()

	atTop = e.scrollY <= 0
	atBottom = e.scrollY >= contentH-viewportH
	atLeft = e.scrollX <= 0
	atRight = e.scrollX >= contentW-viewportW

	return
}

func (b *ScrollBoxElement) Record(cb *gfx.CommandBuffer, container gfx.Rect) error {
	if b.pendingAction != nil {
		switch *b.pendingAction {
		case scrollActionTop:
			b.scrollToTop()
		case scrollActionBottom:
			b.scrollToBottom()
		case scrollActionLeft:
			b.scrollToLeft()
		case scrollActionRight:
			b.scrollToRight()
		}
		b.pendingAction = nil
	}

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
	e.scheduleUpdate(func() error {
		return guardDestroyed(e.ctx)
	})
}
func (e *ScrollBoxElement) UnsetOverflow() {
	e.scheduleUpdate(func() error {
		return guardDestroyed(e.ctx)
	})
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
