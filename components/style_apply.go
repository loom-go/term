package components

import (
	"fmt"
	"sync/atomic"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom-term/internal/app"
	. "github.com/AnatoleLucet/loom/components"
)

var nextID atomic.Uint32

func newID() uint32 {
	return nextID.Add(1)
}

type applyNode struct {
	style Style
}

// Apply applies the given Style to a node.
func Apply(style Style) *applyNode {
	return &applyNode{style}
}

func BindApply(fn func() Style) loom.Node {
	return Bind(func() loom.Node {
		return Apply(fn())
	})
}

// ApplyOn applies a Style to a specific event (e.g., "hover", "focus").
func ApplyOn(event string, style Style) *applyNode {
	// todo: impl
	return &applyNode{style}
}

func BindApplyOn(event string, fn func() Style) *applyNode {
	// todo: impl
	return &applyNode{fn()}
}

func (s *applyNode) ID() string {
	return "term.Style"
}

func (s *applyNode) Mount(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Apply (style): %w", err)
	}

	id := newID()
	slot.SetSelf(id)

	parent := slot.Parent().(core.Element)

	stack := getStyleStack(parent)
	stack.Push(id, s.style)

	applyStyle(parent, &s.style)
	ctx.ScheduleRender()

	return nil
}

func (s *applyNode) Update(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Apply (style): %w", err)
	}

	self := slot.Self().(uint32)
	parent := slot.Parent().(core.Element)

	stack := getStyleStack(parent)
	stack.Replace(self, s.style)

	applyStyleStack(parent)
	ctx.ScheduleRender()

	return nil
}

func (s *applyNode) Unmount(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Apply (style): %w", err)
	}

	self := slot.Self().(uint32)
	parent := slot.Parent().(core.Element)

	stack := getStyleStack(parent)
	stack.Pop(self)

	removeStyle(parent, &s.style)
	applyStyleStack(parent)
	ctx.ScheduleRender()

	return nil
}
