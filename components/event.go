package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom-term/core"
)

// On registers a callback for a specific event.
type On struct {
	Hover        func(*core.EventMouse)
	Click        func(*core.EventMouse)
	MouseEnter   func(*core.EventMouse)
	MouseLeave   func(*core.EventMouse)
	MousePress   func(*core.EventMouse)
	MouseRelease func(*core.EventMouse)
	MouseMove    func(*core.EventMouse)
	MouseScroll  func(*core.EventMouse)
	MouseDrag    func(*core.EventMouse)

	KeyPress   func(*core.EventKey)
	KeyRelease func(*core.EventKey)
	Paste      func(*core.EventPaste)

	Focus func(*core.EventFocus)
	Blur  func(*core.EventBlur)

	Input  func(*core.EventInput)
	Submit func(*core.EventSubmit)
}

func (n On) Apply(parent any) (func() error, error) {
	elem, ok := parent.(core.Element)
	if !ok {
		return nil, fmt.Errorf("On: parent node is not an Element")
	}

	var removers []func()

	if n.Click != nil {
		removers = append(removers, elem.OnMousePress(n.Click))
	}
	if n.Hover != nil {
		removers = append(removers, elem.OnMouseEnter(n.Hover))
	}
	if n.MouseEnter != nil {
		removers = append(removers, elem.OnMouseEnter(n.MouseEnter))
	}
	if n.MouseLeave != nil {
		removers = append(removers, elem.OnMouseLeave(n.MouseLeave))
	}
	if n.MousePress != nil {
		removers = append(removers, elem.OnMousePress(n.MousePress))
	}
	if n.MouseRelease != nil {
		removers = append(removers, elem.OnMouseRelease(n.MouseRelease))
	}
	if n.MouseMove != nil {
		removers = append(removers, elem.OnMouseMove(n.MouseMove))
	}
	if n.MouseScroll != nil {
		removers = append(removers, elem.OnMouseScroll(n.MouseScroll))
	}
	if n.MouseDrag != nil {
		removers = append(removers, elem.OnMouseDrag(n.MouseDrag))
	}

	if n.KeyPress != nil {
		removers = append(removers, elem.OnKeyPress(n.KeyPress))
	}
	if n.KeyRelease != nil {
		removers = append(removers, elem.OnKeyRelease(n.KeyRelease))
	}
	if n.Paste != nil {
		removers = append(removers, elem.OnPaste(n.Paste))
	}

	if n.Focus != nil {
		removers = append(removers, elem.OnFocus(n.Focus))
	}
	if n.Blur != nil {
		removers = append(removers, elem.OnBlur(n.Blur))
	}

	if n.Input != nil {
		removers = append(removers, elem.OnInput(n.Input))
	}
	if n.Submit != nil {
		removers = append(removers, elem.OnSubmit(n.Submit))
	}

	remove := func() error {
		for _, r := range removers {
			r()
		}

		return nil
	}

	return remove, nil
}
