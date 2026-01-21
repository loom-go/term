package types

import (
	"fmt"
	"strconv"
	"strings"
)

type EventMouse struct {
	BaseEvent
	X      int
	Y      int
	Alt    bool
	Ctrl   bool
	Shift  bool
	Action MouseAction
	Button MouseButton
}

func (e EventMouse) String() string {
	var modifiers []string
	if e.Shift {
		modifiers = append(modifiers, "shift")
	}
	if e.Alt {
		modifiers = append(modifiers, "alt")
	}
	if e.Ctrl {
		modifiers = append(modifiers, "ctrl")
	}

	modStr := ""
	if len(modifiers) > 0 {
		modStr = fmt.Sprintf(" + %s", strings.Join(modifiers, " + "))
	}

	return fmt.Sprintf("Mouse(%s at %d,%d%s)", e.Action, e.X, e.Y, modStr)
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
		return "Press"
	case MouseActionRelease:
		return "Release"
	case MouseActionMove:
		return "Move"
	case MouseActionDrag:
		return "Drag"
	case MouseActionScroll:
		return "Scroll"
	default:
		return strconv.Itoa(int(a))
	}
}

type MouseButton int

func (b MouseButton) String() string {
	switch b {
	case MouseButtonNone:
		return "None"
	case MouseButtonLeft:
		return "Left"
	case MouseButtonMiddle:
		return "Middle"
	case MouseButtonRight:
		return "Right"
	case MouseButtonWheelUp:
		return "WheelUp"
	case MouseButtonWheelDown:
		return "WheelDown"
	case MouseButtonWheelLeft:
		return "WheelLeft"
	case MouseButtonWheelRight:
		return "WheelRight"
	case MouseButtonBackward:
		return "Backward"
	case MouseButtonForward:
		return "Forward"
	default:
		return strconv.Itoa(int(b))
	}
}

// http://xahlee.info/linux/linux_x11_mouse_button_number.html
const (
	MouseButtonNone MouseButton = iota
	MouseButtonLeft
	MouseButtonMiddle
	MouseButtonRight
	MouseButtonWheelUp
	MouseButtonWheelDown
	MouseButtonWheelLeft
	MouseButtonWheelRight
	MouseButtonBackward
	MouseButtonForward
	MouseButton10
	MouseButton11
)
