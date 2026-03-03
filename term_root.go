package term

import (
	"context"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom-term/internal/app"
	. "github.com/AnatoleLucet/loom/components"
)

type Element = core.Element

func Root() Element {
	ctx, err := app.GetContext()
	if err != nil {
		// fine to panic because this can only be called in the reactive system.
		// else we have a reason to *panic*
		panic("term.Root: " + err.Error())
	}

	return ctx.Root()
}

type rootNode struct {
	ctx    context.Context
	appctx *app.AppContext

	fn func() loom.Node
}

func newRootNode(ctx context.Context, appctx *app.AppContext, fn func() loom.Node) (*rootNode, error) {
	return &rootNode{
		ctx:    ctx,
		appctx: appctx,
		fn:     fn,
	}, nil
}

func (n *rootNode) ID() string {
	return "term.Root"
}

func (n *rootNode) Mount(slot *loom.Slot) error {
	n.appctx.PushRenderHold()
	defer n.appctx.PopRenderHold()

	slot.SetSelf(n.appctx.Root())

	return n.Update(slot)
}

func (n *rootNode) Update(slot *loom.Slot) error {
	n.appctx.PushRenderHold()
	defer n.appctx.PopRenderHold()

	return slot.RenderChildren(Provider(app.Context, n.appctx, n.fn))
}

func (n *rootNode) Unmount(slot *loom.Slot) error {
	return nil
}
