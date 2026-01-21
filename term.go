package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom-term/core/runtime"
	"github.com/AnatoleLucet/loom-term/internal/app"
	"github.com/AnatoleLucet/loom-term/internal/errs"
	. "github.com/AnatoleLucet/loom/components"
	"github.com/AnatoleLucet/loom/signals"
)

type RenderType = core.RenderType

const RenderInline RenderType = core.RenderInline
const RenderFullscreen RenderType = core.RenderFullscreen

type App struct {
	running bool

	slot  *loom.Slot
	owner *signals.Owner

	errorChan chan any
}

func NewApp() *App {
	return &App{
		running:   false,
		owner:     signals.NewOwner(),
		errorChan: make(chan any, 1),
	}
}

func (a *App) Run(typ RenderType, fn func() loom.Node) <-chan any {
	if a.running {
		a.errorChan <- errors.New("app is already running")
		return a.errorChan
	}

	a.owner.OnError(func(err any) {
		a.errorChan <- err
	})

	err := a.Render(typ, fn)
	if err != nil {
		a.errorChan <- err
		return a.errorChan
	}

	a.running = true
	return a.errorChan
}

func (a *App) Render(typ RenderType, fn func() loom.Node) error {
	err := a.owner.Run(func() error {
		rt, err := runtime.NewRuntime(typ)
		if err != nil {
			return fmt.Errorf("App: %w: %w", errs.ErrAppFailedToInitialize, err)
		}

		ctx := app.NewAppContext(rt)
		go a.forwardErrors(ctx.Errors())

		node := &rootNode{
			typ: typ,
			fn:  fn,
			ctx: ctx,
		}

		slot, err := loom.Render(nil, node)
		a.slot = slot

		return err
	})

	if err != nil {
		a.owner.Dispose()
		return err
	}

	return nil
}

func (a *App) Close() {
	if !a.running {
		return
	}
	a.running = false

	if a.slot != nil {
		a.slot.Unmount()
	}

	// todo: move in renderer.Close()
	os.Stdout.WriteString("\x1b[?25h") // ensure cursor is visible

	a.owner.Dispose()
}

func (a *App) forwardErrors(errs <-chan error) {
	a.errorChan <- <-errs
}

type rootNode struct {
	typ RenderType
	ctx *app.AppContext
	fn  func() loom.Node
}

func (n *rootNode) ID() string {
	return "term.Root"
}

func (n *rootNode) Mount(slot *loom.Slot) error {
	n.ctx.PushRenderHold()
	defer n.ctx.PopRenderHold()

	slot.SetSelf(n.ctx.Root())

	n.ctx.Root().SetPosition("relative")

	width, height := TerminalSize()
	if width > 0 {
		n.ctx.Root().SetWidth(width)
	}
	if height > 0 {
		if n.typ == RenderFullscreen {
			n.ctx.Root().SetHeight(height)
		}

		if n.typ == RenderInline {
			n.ctx.Root().SetMaxHeight(height)
		}
	}

	return n.Update(slot)
}

func (n *rootNode) Update(slot *loom.Slot) error {
	n.ctx.PushRenderHold()
	defer n.ctx.PopRenderHold()

	return slot.RenderChildren(Provider(app.Context, n.ctx, n.fn))
}

func (n *rootNode) Unmount(slot *loom.Slot) error {
	if n.ctx != nil {
		n.ctx.Close()
	}

	return nil
}
