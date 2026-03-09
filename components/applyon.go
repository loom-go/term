package components

import (
	"github.com/loom-go/loom"
	. "github.com/loom-go/loom/components"
	term "github.com/loom-go/term"
)

// ApplyOn applies the Styles on a specific event ("hover", "focus", "active").
func ApplyOn(event string, appliers ...loom.Applier) loom.Node {
	shown, setShown := Signal(false)

	var events loom.Node

	switch event {
	case "hover":
		events = Apply(On{
			MouseEnter: func(*term.EventMouse) { setShown(true) },
			MouseLeave: func(*term.EventMouse) { setShown(false) },
		})
	case "focus":
		events = Apply(On{
			Focus: func(*term.EventFocus) { setShown(true) },
			Blur:  func(*term.EventBlur) { setShown(false) },
		})
	case "active":
		events = Apply(On{
			MousePress:   func(*term.EventMouse) { setShown(true) },
			MouseRelease: func(*term.EventMouse) { setShown(false) },
		})
	}

	return Fragment(
		events,
		Show(shown, func() loom.Node {
			return Apply(appliers...)
		}),
	)
}
