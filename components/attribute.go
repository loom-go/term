package components

import (
	"fmt"

	appctx "github.com/AnatoleLucet/loom-term/components/context"
)

// Attr is an alias for Attribute
type Attr = Attribute

// Attribute assigns the given values on the parent.
type Attribute struct {
	Title       any // string | func() string
	Value       any // string | func() string
	Placeholder any // string | func() string
}

func (n Attribute) Apply(parent any) (func() error, error) {
	ctx, err := appctx.Get()
	if err != nil {
		return nil, fmt.Errorf("Attribute: %w", err)
	}

	removers := []func(){}
	remove := func() error {
		return ctx.BatchRender(func() error {
			for _, r := range removers {
				r()
			}

			return nil
		})
	}

	err = ctx.BatchRender(func() error {
		if n.Title != nil {
			if e, v, ok := matchMethod[interface {
				SetTitle(string)
				UnsetTitle()
			}, string](parent, n.Title); ok {
				e.SetTitle(v)
				removers = append(removers, e.UnsetTitle)
			}
		}

		if n.Value != nil {
			if e, v, ok := matchMethod[interface {
				SetValue(string)
				UnsetValue()
			}, string](parent, n.Value); ok {
				e.SetValue(v)
				removers = append(removers, e.UnsetValue)
			}
		}

		if n.Placeholder != nil {
			if e, v, ok := matchMethod[interface {
				SetPlaceholder(string)
				UnsetPlaceholder()
			}, string](parent, n.Placeholder); ok {
				e.SetPlaceholder(v)
				removers = append(removers, e.UnsetPlaceholder)
			}
		}

		return nil
	})

	return remove, err
}
