package render

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/AnatoleLucet/go-opentui"
	"github.com/AnatoleLucet/loom-term/core/debug"
	"github.com/AnatoleLucet/loom-term/core/terminal"
	"github.com/AnatoleLucet/loom-term/core/types"
	"github.com/AnatoleLucet/tess"
	"golang.org/x/term"
)

type TermRenderer struct {
	mu sync.RWMutex

	typ          types.RenderType
	clock        int
	renderOffset int

	ctx      types.RenderContext
	oldState *term.State
	renderer *opentui.Renderer
	elements map[uint32]types.Element
}

func NewRenderer(typ types.RenderType, ctx types.RenderContext) (tr *TermRenderer, err error) {
	tr = &TermRenderer{
		typ:      typ,
		ctx:      ctx,
		elements: make(map[uint32]types.Element),
	}

	tr.renderer = opentui.NewRenderer(uint32(1), uint32(1))
	if tr.renderer == nil {
		return nil, fmt.Errorf("Renderer: %w", types.ErrRendererFailedToInitialize)
	}

	return tr, nil
}

func (r *TermRenderer) Type() types.RenderType {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.typ
}

func (r *TermRenderer) RenderOffset() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.renderOffset
}

func (r *TermRenderer) Render(elem types.Element) error {
	start := time.Now()
	defer func() {
		end := time.Now()
		debug.EmitFrameTime(end.Sub(start))
	}()

	prevWidth := elem.Layout().Width()
	prevHeight := elem.Layout().Height()

	err := r.computeLayout(elem)
	if err != nil {
		return fmt.Errorf("Renderer: %w", err)
	}

	newWidth := elem.Layout().Width()
	newHeight := elem.Layout().Height()

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

	err = r.paint(elem, buffer)
	if err != nil {
		return fmt.Errorf("Renderer: %w", err)
	}

	err = r.render(false)
	if err != nil {
		return fmt.Errorf("Renderer: %w", err)
	}

	return nil
}

func (r *TermRenderer) CheckHit(x, y int) types.Element {
	id, err := r.renderer.CheckHit(uint32(x), uint32(y))
	if err != nil {
		return nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.elements[id]
}

func (r *TermRenderer) AddToHitGrid(element types.Element, x, y, width, height int) {
	randomId, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	if err != nil {
		return
	}

	id := uint32(randomId.Int64())

	r.mu.Lock()
	r.elements[id] = element
	r.renderer.AddToHitGrid(int32(x), int32(y), uint32(width), uint32(height), id)
	r.mu.Unlock()
}

func (r *TermRenderer) PushHitGridScissors(x, y, width, height int) {
	r.renderer.HitGridPushScissorRect(int32(x), int32(y), uint32(width), uint32(height))
}

func (r *TermRenderer) PopHitGridScissors() {
	r.renderer.HitGridPopScissorRect()
}

func (r *TermRenderer) UpdateCapabilities(buf []byte) error {
	if !terminal.IsCapabilityResponse(buf) {
		return nil
	}

	err := r.renderer.ProcessCapabilityResponse(buf)
	if err != nil {
		return fmt.Errorf("Renderer: %w: %w", types.ErrRendererFailedToProcessTermCapabilities, err)
	}

	caps, err := r.renderer.GetTerminalCapabilities()
	if err != nil {
		return fmt.Errorf("Renderer: %w: %w", types.ErrRendererFailedToGetTermCapabilities, err)
	}

	if caps.KittyKeyboard {
		r.renderer.EnableKittyKeyboard(7)
	}

	if caps.BracketedPaste {
		r.renderer.EnableBracketedPaste()
	}

	r.renderer.EnableMouse(true)

	return nil
}

func (r *TermRenderer) Close() error {
	r.renderer.Close()
	term.Restore(int(os.Stdin.Fd()), r.oldState)
	return nil
}

func (r *TermRenderer) init() error {
	fd := int(os.Stdin.Fd())

	initialState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	r.oldState = initialState

	if r.typ == types.RenderTypeFullscreen {
		err := r.initFullscreen()
		if err != nil {
			return err
		}
	}

	if r.typ == types.RenderTypeInline {
		err := r.initInline()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TermRenderer) initFullscreen() error {
	r.renderer.SetupTerminal(true)

	r.renderer.ClearTerminal()
	return nil
}

func (r *TermRenderer) initInline() error {
	cursorRow, _, err := terminal.CursorPos()
	if err != nil {
		return err
	}

	r.renderer.SetupTerminal(false)

	r.updateRenderOffset(cursorRow - 1)
	return nil
}

func (r *TermRenderer) updateRendererSize(prevWidth, prevHeight, newWidth, newHeight float32) error {
	if prevWidth == newWidth && prevHeight == newHeight {
		return nil
	}

	termWidth, termHeight, err := terminal.Size()
	if err != nil {
		return err
	}

	if r.typ == types.RenderTypeFullscreen {
		r.renderer.Resize(
			uint32(termWidth),
			uint32(termHeight),
		)
	}

	if r.typ == types.RenderTypeInline {
		shrunk := newHeight < prevHeight
		if shrunk {
			buffer, err := r.nextBuffer()
			if err != nil {
				return err
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
			terminal.ScrollUp(contentEndRow - termHeight)

			r.updateRenderOffset(max(0, termHeight-int(newHeight)))
		}

	}

	return nil
}

func (r *TermRenderer) nextBuffer() (*opentui.Buffer, error) {
	buffer, err := r.renderer.GetNextBuffer()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", types.ErrRendererFailedToGetBuffer, err)
	}

	return buffer, nil
}

func (r *TermRenderer) updateRenderOffset(newOffset int) {
	if r.typ != types.RenderTypeInline {
		return
	}

	r.renderOffset = newOffset
	r.renderer.SetRenderOffset(uint32(r.renderOffset))
}

func (r *TermRenderer) computeLayout(elem types.Element) error {
	start := time.Now()
	defer func() {
		end := time.Now()
		debug.EmitLayoutTime(end.Sub(start))
	}()

	termWidth, termHeight, err := terminal.Size()
	if err != nil {
		return err
	}

	if r.typ == types.RenderTypeFullscreen {
		return elem.XYZ().ComputeLayout(tess.Container{
			Width:  float32(termWidth),
			Height: float32(termHeight),
		})
	}

	if r.typ == types.RenderTypeInline {
		return elem.XYZ().ComputeLayout(tess.Container{
			Width: float32(termWidth),
		})
	}

	return nil
}

func (r *TermRenderer) paint(elem types.Element, buffer *opentui.Buffer) error {
	start := time.Now()
	defer func() {
		end := time.Now()
		debug.EmitPaintTime(end.Sub(start))
	}()

	return elem.Paint(buffer, 0, 0)
}

func (r *TermRenderer) render(force bool) error {
	start := time.Now()
	defer func() {
		end := time.Now()
		debug.EmitRenderTime(end.Sub(start))
	}()

	return r.renderer.Render(force)
}
