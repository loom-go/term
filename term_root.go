package term

import (
	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/components/context"
)

type rootNode struct {
	appctx *AppContext

	fn func() loom.Node
}

func newRootNode(appctx *AppContext, fn func() loom.Node) (*rootNode, error) {
	return &rootNode{
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
		return slot.RenderChildren(appctx.Provider(n.appctx, func() loom.Node {
			return n.fn()
		}))
	})
}

func (n *rootNode) Unmount(slot *loom.Slot) error {
	return nil
}
