package components

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom-term/internal/app"
	. "github.com/AnatoleLucet/loom/components"
)

func newId() (string, error) {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", fmt.Errorf("unable to generate random ID: %w", err) // todo: error
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
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

	id, err := newId()
	if err != nil {
		return err
	}
	slot.SetSelf(id)

	parent := slot.Parent().(core.Element)

	stack := getStyleStack(parent)
	stack.Push(id, s.style)

	err = ctx.DoSafely(func() error {
		err = applyStyle(parent, &s.style)
		if err != nil {
			return err
		}

		return ctx.RequestRender()
	})

	if err != nil {
		return fmt.Errorf("Apply (style): %w", err)
	}

	return nil
}

func (s *applyNode) Update(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Apply (style): %w", err)
	}

	self := slot.Self().(string)
	parent := slot.Parent().(core.Element)

	stack := getStyleStack(parent)
	stack.Replace(self, s.style)

	err = ctx.DoSafely(func() error {
		err = applyStyleStack(parent)
		if err != nil {
			return err
		}

		return ctx.RequestRender()
	})

	if err != nil {
		return fmt.Errorf("Apply (style): %w", err)
	}

	return nil
}

func (s *applyNode) Unmount(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Apply (style): %w", err)
	}

	self := slot.Self().(string)
	parent := slot.Parent().(core.Element)

	stack := getStyleStack(parent)
	stack.Pop(self)

	err = ctx.DoSafely(func() error {
		err = removeStyle(parent, &s.style)
		if err != nil {
			return err
		}

		err = applyStyleStack(parent)
		if err != nil {
			return err
		}

		return ctx.RequestRender()
	})

	if err != nil {
		return fmt.Errorf("Apply (style): %w", err)
	}

	return nil
}
