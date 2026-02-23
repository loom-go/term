package events

import (
	"fmt"
	"strconv"
)

type EventMouse struct {
	X, Y   int
	Key    Key         // keyboard modifiers (ctrl, shift, alt, meta)
	Button MouseButton // mouse button (left, right, middle, wheel, etc)
	Action MouseAction // press, release, move, drag, scroll
}

func (e EventMouse) String() string {
	mods := ""
	if e.Key != 0 {
		mods = e.Key.String() + "+"
	}

	btn := e.Button.String()
	if btn == "none" && e.Action == MouseActionMove {
		return fmt.Sprintf("Mouse(%s at %d,%d)", e.Action, e.X, e.Y)
	}

	return fmt.Sprintf("Mouse(%s, %s%s at %d,%d)", e.Action, mods, btn, e.X, e.Y)
}

type MouseAction int

const (
	MouseActionPress MouseAction = iota
	MouseActionRelease
	MouseActionMove
	MouseActionDrag
	MouseActionScroll
)

func (a MouseAction) String() string {
	switch a {
	case MouseActionPress:
		return "press"
	case MouseActionRelease:
		return "release"
	case MouseActionMove:
		return "move"
	case MouseActionDrag:
		return "drag"
	case MouseActionScroll:
		return "scroll"
	default:
		return strconv.Itoa(int(a))
	}
}

type MouseButton int

const (
	MouseNone MouseButton = iota
	MouseLeft
	MouseMiddle
	MouseRight
	MouseWheelUp
	MouseWheelDown
	MouseWheelLeft
	MouseWheelRight
	MouseBackward
	MouseForward
	Mouse10
	Mouse11
)

func (b MouseButton) String() string {
	if name, ok := mouseButtonNames[b]; ok {
		return name
	}

	return fmt.Sprintf("btn(%d)", b)
}

var mouseButtonNames = map[MouseButton]string{
	MouseNone:       "none",
	MouseLeft:       "left",
	MouseMiddle:     "middle",
	MouseRight:      "right",
	MouseWheelUp:    "wheelup",
	MouseWheelDown:  "wheeldown",
	MouseWheelLeft:  "wheelleft",
	MouseWheelRight: "wheelright",
	MouseBackward:   "backward",
	MouseForward:    "forward",
	Mouse10:         "btn10",
	Mouse11:         "btn11",
}
