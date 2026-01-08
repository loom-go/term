package internal

import (
	"fmt"
	"sync"

	termerror "github.com/AnatoleLucet/loom-term/error"
	"github.com/AnatoleLucet/loom-term/opentui"
	"github.com/AnatoleLucet/tess"
)

type RendererType int

const (
	RendererTypeInline RendererType = iota
	RendererTypeFullscreen
)

type TermRenderer struct {
	mu sync.RWMutex

	typ RendererType

	root     *Element
	queue    *RenderQueue
	renderer *opentui.Renderer

	clock        int
	batchDepth   int
	renderOffset int
}

func NewRenderer(typ RendererType) (tr *TermRenderer, err error) {
	tr = &TermRenderer{typ: typ}

	tr.root, err = NewElement(nil)
	if err != nil {
		return nil, fmt.Errorf("Renderer: %w: %w: %w", termerror.ErrFailedToInitializeRenderer, termerror.ErrFailedToCreateRootNode, err)
	}

	tr.renderer = opentui.NewRenderer(uint32(1), uint32(1))
	if tr.renderer == nil {
		return nil, fmt.Errorf("Renderer: %w", termerror.ErrFailedToInitializeRenderer)
	}

	tr.queue = NewRenderQueue(func(snapshot *Element) error {
		return tr.render(snapshot)
	})

	return tr, nil
}

func (r *TermRenderer) Type() RendererType {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.typ
}

func (r *TermRenderer) Elem() *Element {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.root
}

func (r *TermRenderer) enqueue(snapshot *Element) error {
	return r.queue.Enqueue(snapshot)
}

func (r *TermRenderer) render(snapshot *Element) error {
	prevWidth := snapshot.Layout().GetLayout().Width()
	prevHeight := snapshot.Layout().GetLayout().Height()

	// Compute layout
	err := r.computeLayout(snapshot)
	if err != nil {
		return fmt.Errorf("Renderer: %w", err)
	}

	newWidth := snapshot.Layout().GetLayout().Width()
	newHeight := snapshot.Layout().GetLayout().Height()

	// init on first render
	if r.clock == 0 {
		err := r.init()
		if err != nil {
			return fmt.Errorf("Renderer: %w", err)
		}
	}
	r.clock++

	err = r.updateRendererSize(prevWidth, prevHeight, newWidth, newHeight)
	if err != nil {
		return fmt.Errorf("Renderer: %w", err)
	}

	buffer, err := r.nextBuffer()
	if err != nil {
		return fmt.Errorf("Renderer: %w", err)
	}

	err = snapshot.PaintAll(buffer)
	if err != nil {
		return fmt.Errorf("Renderer: %w", err)
	}

	r.renderer.Render(false)

	return nil
}

func (r *TermRenderer) Close() error {
	if r.renderer == nil {
		return nil
	}

	if r.typ == RendererTypeInline {
		contentHeight := r.root.Layout().GetLayout().Height()
		return r.renderer.CloseWithOptions(false, uint32(contentHeight))
	}

	return r.renderer.Close()
}

// Batch groups multiple updates into a single render on the next frame.
func (r *TermRenderer) Batch(fn func() error) error {
	r.mu.Lock()
	r.batchDepth++
	r.mu.Unlock()

	err := fn()

	r.mu.Lock()
	r.batchDepth--
	isOutermost := r.batchDepth == 0
	r.mu.Unlock()

	if err != nil || !isOutermost {
		return err
	}

	// todo: dont block reactive system. this will block the effect that's running the Batch
	return r.enqueue(r.root.Clone())
}

func (r *TermRenderer) init() error {
	if r.typ == RendererTypeFullscreen {
		return r.initFullscreen()
	}

	if r.typ == RendererTypeInline {
		return r.initInline()
	}

	return nil
}

func (r *TermRenderer) initFullscreen() error {
	r.renderer.SetupTerminal(true)
	r.renderer.ClearTerminal()
	return nil
}

func (r *TermRenderer) initInline() error {
	cursorRow, _, err := CursorPos()
	if err != nil {
		return termerror.ErrFailedToGetCursorPosition
	}

	r.renderer.SetupTerminal(false)
	r.updateRenderOffset(cursorRow - 1)
	return nil
}

func (r *TermRenderer) computeLayout(node *Element) error {
	termWidth, termHeight, err := TerminalSize()
	if err != nil {
		return termerror.ErrFailedToGetTerminalSize
	}

	if r.typ == RendererTypeFullscreen {
		err := node.Layout().ComputeLayout(tess.Container{
			Width:  float32(termWidth),
			Height: float32(termHeight),
		})
		return err
	}

	if r.typ == RendererTypeInline {
		err := node.Layout().ComputeLayout(tess.Container{
			Width: float32(termWidth),
		})
		return err
	}

	return nil
}

func (r *TermRenderer) updateRendererSize(prevWidth, prevHeight, newWidth, newHeight float32) error {
	if prevWidth == newWidth && prevHeight == newHeight {
		return nil
	}

	termWidth, termHeight, err := TerminalSize()
	if err != nil {
		return termerror.ErrFailedToGetTerminalSize
	}

	if r.typ == RendererTypeFullscreen {
		r.renderer.Resize(
			uint32(termWidth),
			uint32(termHeight),
		)
	}

	if r.typ == RendererTypeInline {
		shrunk := newHeight < prevHeight
		if shrunk {
			buffer, err := r.nextBuffer()
			if err != nil {
				return fmt.Errorf("%w: %w", termerror.ErrFailedToGetBuffer, err)
			}

			// clear the area below the content to avoid artifacts
			buffer.FillRect(
				0,
				uint32(newHeight),
				uint32(termWidth),
				uint32(prevHeight-newHeight),
				opentui.Transparent,
			)

			r.renderer.Render(false)
		}

		r.renderer.Resize(
			uint32(termWidth),
			uint32(newHeight),
		)

		// scroll up if content overflows terminal height
		contentEndRow := r.renderOffset + int(newHeight)
		if contentEndRow > termHeight {
			ScrollUp(contentEndRow - termHeight)

			r.updateRenderOffset(max(0, termHeight-int(newHeight)))
		}

	}

	return nil
}

func (r *TermRenderer) nextBuffer() (*opentui.Buffer, error) {
	buffer, err := r.renderer.GetNextBuffer()
	if err != nil {
		return nil, fmt.Errorf("Renderer: %w: %w", termerror.ErrFailedToGetBuffer, err)
	}

	return buffer, nil
}

func (r *TermRenderer) updateRenderOffset(newOffset int) {
	if r.typ != RendererTypeInline {
		return
	}

	r.renderOffset = newOffset
	r.renderer.SetRenderOffset(uint32(r.renderOffset))
}

// 1. x make the renderer render imediatly with a queue if updates are comming faster than can render
// 2. have a global frame pacer that's only used for animations
// 3. have an Animate struct that paces each update with the global frame pacer
// 4. to better optimize the frame pacer, signals.Batch the animations updates, to only trigger one render for every animations of the app
