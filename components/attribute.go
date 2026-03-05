package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/core"
	. "github.com/AnatoleLucet/loom/components"
)

type AttributeName string

const (
	AttrTitle       AttributeName = "title"
	AttrValue       AttributeName = "value"
	AttrPlaceholder AttributeName = "placeholder"
)

func Attr[T any](name AttributeName, value T) loom.Node {
	return &attrNode[T]{
		name:  name,
		value: value,
	}
}

func BindAttr[T any](name AttributeName, fn func() T) loom.Node {
	return Bind(func() loom.Node {
		return Attr(name, fn())
	})
}

type attrNode[T any] struct {
	name  AttributeName
	value T
}

func (n *attrNode[T]) ID() string {
	return "term.Attr"
}

func (n *attrNode[T]) Mount(slot *loom.Slot) error {
	return n.Update(slot)
}

func (n *attrNode[T]) Update(slot *loom.Slot) error {
	parent := slot.Parent().(core.Element)

	var remove func()

	switch n.name {
	case AttrTitle:
		if e, ok := parent.(interface {
			SetTitle(string)
			UnsetTitle()
		}); ok {
			e.SetTitle(fmt.Sprintf("%v", n.value))
			remove = e.UnsetTitle
		}

	case AttrValue:
		if e, ok := parent.(interface {
			SetValue(string)
			Clear()
		}); ok {
			e.SetValue(fmt.Sprintf("%v", n.value))
			remove = e.Clear
		}

	case AttrPlaceholder:
		if e, ok := parent.(interface {
			SetPlaceholder(string)
			UnsetPlaceholder()
		}); ok {
			e.SetPlaceholder(fmt.Sprintf("%v", n.value))
			remove = e.UnsetPlaceholder
		}

	default:
		return fmt.Errorf("Attr: unsupported attribute name: %s", n.name)
	}

	slot.SetSelf(remove)

	return nil
}

func (n *attrNode[T]) Unmount(slot *loom.Slot) error {
	remove := slot.Self()
	if remove != nil {
		remove.(func())()
		slot.SetSelf(nil)
	}

	return nil
}
