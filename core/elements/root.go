package elements

import (
	"fmt"

	"github.com/loom-go/term/core/gfx"

	"github.com/AnatoleLucet/go-opentui"
)

type RootElement struct {
	*BaseElement

	initialized bool

	cb  *gfx.CommandBuffer
	rdr *opentui.Renderer
}

func NewRootElement(typ RenderType) (root *RootElement, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Root: %w: %w", ErrFailedToInitializeRoot, err)
		}
	}()

	base, err := NewBaseElement()
	if err != nil {
		return nil, err
	}

	root = &RootElement{BaseElement: base}
	root.cb = gfx.NewCommandBuffer(root.ctx)
	root.rdr = opentui.NewRenderer(1, 1)

	rc, err := NewRenderContext(root.ctx, typ, root)
	if err != nil {
		return nil, err
	}
	root.SetRenderContext(rc)

	go root.listenToMouseEvents(root.ctx)
	go root.listenToKeyboardEvents(root.ctx)
	go root.listenToResizeEvents(root.ctx)
	go root.listenToCapabilites(root.ctx)
	go root.listenToExitEvents(root.ctx)

	root.OnDestroy(func() {
		root.rdr.Close()
	})

	return root, nil
}
