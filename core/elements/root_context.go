package elements

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/loom-go/term/core/gfx"
	"github.com/loom-go/term/core/term"
)

type RenderType int

const (
	RenderTypeInline RenderType = iota
	RenderTypeFullscreen
)

type RenderContext struct {
	mu  sync.RWMutex
	ctx context.Context

	typ RenderType

	root *RootElement

	hitGrid map[uint32]Element

	scheduler *Scheduler

	focused Element
	pressed Element
	hovered []Element // path from root to hovered element

	lastMouseX       int
	lastMouseY       int
	hasMousePosition bool

	renderOffsetY int
	renderOffsetX int

	errors chan error
}

func NewRenderContext(ctx context.Context, typ RenderType, root *RootElement) (*RenderContext, error) {
	rc := &RenderContext{
		ctx:     ctx,
		typ:     typ,
		root:    root,
		hitGrid: make(map[uint32]Element),
		errors:  make(chan error, 1),
	}

	rc.scheduler = NewScheduler(SchedulerOptions{
		Context:    ctx,
		Errors:     rc.errors,
		Render:     rc.render,
		PostRender: rc.postRender,
	})

	return rc, nil
}

func (rc *RenderContext) render() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: %v:\n%s", ErrPanicDuringRender, r, debug.Stack())
		}
	}()

	return rc.root.Render(nil, gfx.Rect{})
}

func (rc *RenderContext) postRender() (err error) {
	rc.root.refreshMouseState()
	return nil
}

func (rc *RenderContext) Errors() <-chan error {
	return rc.errors
}

func (rc *RenderContext) PushHold() {
	rc.scheduler.PushHold()
}

func (rc *RenderContext) PopHold() {
	rc.scheduler.PopHold()
}

func (rc *RenderContext) Batch(fn func()) {
	rc.scheduler.Batch(fn)
}

func (rc *RenderContext) Settle() {
	rc.scheduler.Settle()
}

func (rc *RenderContext) scheduleAccess(access func()) {
	rc.scheduler.Access(access)
}

func (rc *RenderContext) scheduleUpdate(update func() error) {
	rc.scheduler.Update(update)
}

func (rc *RenderContext) scheduleUpdateSync(update func() error) {
	rc.scheduler.UpdateSync(update)
}

func (rc *RenderContext) RenderType() RenderType {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.typ
}

func (rc *RenderContext) RenderOffset() (x, y int) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.renderOffsetX, rc.renderOffsetY
}

func (rc *RenderContext) SetRenderOffset(x, y int) {
	rc.mu.Lock()
	rc.renderOffsetX = x
	rc.renderOffsetY = y
	rc.mu.Unlock()
}

func (rc *RenderContext) HoverChain() []Element {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.hovered
}

func (rc *RenderContext) SetHoverChain(chain []Element) {
	rc.mu.Lock()
	rc.hovered = chain
	rc.mu.Unlock()
}

func (rc *RenderContext) FocusedElement() Element {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.focused
}

func (rc *RenderContext) SetFocusedElement(element Element) {
	rc.mu.Lock()
	rc.focused = element
	rc.mu.Unlock()
}

func (rc *RenderContext) PressedElement() Element {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.pressed
}

func (rc *RenderContext) SetPressedElement(element Element) {
	rc.mu.Lock()
	rc.pressed = element
	rc.mu.Unlock()
}

func (rc *RenderContext) CheckHit(x, y int) (elem Element) {
	rc.scheduler.Access(func() {
		rc.mu.RLock()
		defer rc.mu.RUnlock()

		id, err := rc.root.rdr.CheckHit(uint32(x), uint32(y))
		if err != nil {
			return
		}

		elem = rc.hitGrid[id]
	})

	return
}

func (rc *RenderContext) AddToHitGrid(element Element, rect gfx.Rect) error {
	rc.mu.Lock()
	id := newID()
	rc.hitGrid[id] = element
	rc.mu.Unlock()

	return rc.root.rdr.AddToHitGrid(int32(rect.X), int32(rect.Y), uint32(rect.Width), uint32(rect.Height), id)
}

func (rc *RenderContext) ClearHitGrid() error {
	rc.mu.Lock()
	rc.hitGrid = make(map[uint32]Element)
	rc.mu.Unlock()

	return rc.root.rdr.ClearCurrentHitGrid()
}

func (rc *RenderContext) SetMousePosition(x, y int) {
	rc.mu.Lock()
	rc.lastMouseX = x
	rc.lastMouseY = y
	rc.hasMousePosition = true
	rc.mu.Unlock()
}

func (rc *RenderContext) SetCursorPosition(x, y int, visible bool) error {
	return rc.root.rdr.SetCursorPosition(int32(x), int32(y), visible)
}

func (rc *RenderContext) ClearCursor() error {
	return rc.root.rdr.SetCursorPosition(0, 0, false)
}

func (rc *RenderContext) UpdateCapabilites(buf []byte) error {
	if !term.IsCapabilityResponse(buf) {
		return nil
	}

	err := rc.root.rdr.ProcessCapabilityResponse(buf)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToProcessTerminalCapabilities, err)
	}

	caps, err := rc.root.rdr.GetTerminalCapabilities()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToProcessTerminalCapabilities, err)
	}

	if caps.KittyKeyboard {
		rc.root.rdr.EnableKittyKeyboard(7)
	}

	if caps.BracketedPaste {
		rc.root.rdr.EnableBracketedPaste()
	}

	rc.root.rdr.EnableMouse(true)

	return nil
}

func (rc *RenderContext) emitError(err error) {
	select {
	case rc.errors <- err:
	default:
	}
}
