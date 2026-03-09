package elements

import (
	"fmt"
	"github.com/loom-go/term/core/debug"
	"github.com/loom-go/term/core/gfx"
	"github.com/loom-go/term/core/term"
	"time"

	"github.com/AnatoleLucet/go-opentui"
	"github.com/AnatoleLucet/tess"
)

func (r *RootElement) Render(*opentui.Buffer, gfx.Rect) (err error) {
	if err := guardDestroyed(r); err != nil {
		return fmt.Errorf("Renderer: %w", err)
	}

	start := time.Now()
	defer func() {
		end := time.Now()
		debug.EmitFrameTime(end.Sub(start))

		if err != nil {
			err = fmt.Errorf("Renderer: %w", err)
		}
	}()

	prevWidth := r.xyz().GetLayout().Width()
	prevHeight := r.xyz().GetLayout().Height()

	err = r.computeLayout(r)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToComputeLayout, err)
	}

	newWidth := r.xyz().GetLayout().Width()
	newHeight := r.xyz().GetLayout().Height()

	if !r.initialized {
		err := r.init()
		if err != nil {
			return err
		}
	}
	r.initialized = true

	resized, err := r.updateRendererSize(prevWidth, prevHeight, newWidth, newHeight)
	if err != nil {
		return err
	}

	err = r.record()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRecordFrame, err)
	}

	buffer, err := r.nextBuffer()
	if err != nil {
		return err
	}

	err = r.rdrctx.ClearCursor()
	if err != nil {
		return err
	}

	err = r.rdrctx.ClearHitGrid()
	if err != nil {
		return err
	}

	err = r.draw(buffer)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToDrawFrame, err)
	}

	err = r.render(resized)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRenderFrame, err)
	}

	return nil
}

func (r *RootElement) record() error {
	start := time.Now()
	defer func() {
		end := time.Now()
		debug.EmitRecordTime(end.Sub(start))
	}()

	return r.Record(r.cb, r.relativeRect())
}

func (r *RootElement) draw(buffer *opentui.Buffer) error {
	start := time.Now()
	defer func() {
		end := time.Now()
		debug.EmitDrawTime(end.Sub(start))
	}()

	return r.cb.Execute(r.rdr, buffer)
}

func (r *RootElement) render(force bool) error {
	start := time.Now()
	defer func() {
		end := time.Now()
		debug.EmitRenderTime(end.Sub(start))
	}()

	return r.rdr.Render(force)
}

func (r *RootElement) init() error {
	initialState, err := term.MakeRaw()
	if err != nil {
		return err
	}
	r.OnDestroy(func() {
		term.Restore(initialState)
	})

	if r.rdrctx.RenderType() == RenderTypeFullscreen {
		err = r.initFullscreen()
	}
	if r.rdrctx.RenderType() == RenderTypeInline {
		err = r.initInline()
	}

	return err
}

func (r *RootElement) initFullscreen() error {
	r.rdr.SetupTerminal(true)
	r.rdr.ClearTerminal()
	return nil
}

func (r *RootElement) initInline() error {
	cursorRow, _, err := term.CursorPos()
	if err != nil {
		return err
	}

	r.rdr.SetupTerminal(false)

	r.updateRenderOffset(0, cursorRow-1)
	return nil
}

func (r *RootElement) updateRendererSize(prevWidth, prevHeight, newWidth, newHeight float32) (bool, error) {
	if prevWidth == newWidth && prevHeight == newHeight {
		return false, nil
	}

	termWidth, termHeight, err := term.Size()
	if err != nil {
		return false, err
	}

	if r.rdrctx.RenderType() == RenderTypeFullscreen {
		r.rdr.Resize(
			uint32(termWidth),
			uint32(termHeight),
		)

		return true, nil
	}

	if r.rdrctx.RenderType() == RenderTypeInline {
		shrunk := newHeight < prevHeight
		if shrunk {
			// if the content shrunk, render a frame before resizing the renderer to
			// clear the bottom area that's about to exit the renderer's viewport
			// to avoid leaving behind any visual artifacts
			r.rdr.Render(false)
		}

		r.rdr.Resize(
			uint32(termWidth),
			uint32(newHeight),
		)

		// scroll if content overflows terminal height
		_, offsetY := r.rdrctx.RenderOffset()
		contentEndRow := offsetY + int(newHeight)
		if contentEndRow > termHeight {
			term.ScrollUp(contentEndRow - termHeight)
			r.updateRenderOffset(0, max(0, termHeight-int(newHeight)))
		}

		return true, nil
	}

	return false, nil
}

func (r *RootElement) nextBuffer() (*opentui.Buffer, error) {
	buffer, err := r.rdr.GetNextBuffer()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetBuffer, err)
	}

	return buffer, nil
}

func (r *RootElement) updateRenderOffset(offsetX, offsetY int) {
	if r.rdrctx.RenderType() != RenderTypeInline {
		return
	}

	r.rdrctx.SetRenderOffset(offsetX, offsetY)
	r.rdr.SetRenderOffset(uint32(offsetY))
}

func (r *RootElement) computeLayout(elem Element) error {
	start := time.Now()
	defer func() {
		end := time.Now()
		debug.EmitLayoutTime(end.Sub(start))
	}()

	termWidth, termHeight, err := term.Size()
	if err != nil {
		return err
	}

	if r.rdrctx.RenderType() == RenderTypeFullscreen {
		return elem.xyz().ComputeLayout(tess.Container{
			Width:  float32(termWidth),
			Height: float32(termHeight),
		})
	}

	if r.rdrctx.RenderType() == RenderTypeInline {
		return elem.xyz().ComputeLayout(tess.Container{
			Width: float32(termWidth),
		})
	}

	return nil
}
