package term

import (
	"context"

	"github.com/AnatoleLucet/loom"
	appctx "github.com/AnatoleLucet/loom-term/components/context"
	components "github.com/AnatoleLucet/loom/components"
)

type rootNode struct {
	ctx    context.Context
	appctx *AppContext

	fn func() loom.Node
}

func newRootNode(ctx context.Context, appctx *AppContext, fn func() loom.Node) (*rootNode, error) {
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
	slot.SetSelf(n.appctx.Root())

	return n.Update(slot)
}

func (n *rootNode) Update(slot *loom.Slot) error {
	return n.appctx.BatchRender(func() error {
		return slot.RenderChildren(components.Provider(appctx.AppContext, n.appctx, n.fn))
	})
}

func (n *rootNode) Unmount(slot *loom.Slot) error {
	return nil
}
