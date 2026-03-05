package core

import (
	"github.com/AnatoleLucet/loom-term/core/elements"
	"github.com/AnatoleLucet/loom-term/core/elements/events"
)

type Event = elements.Event
type EventOptions = elements.EventOptions

type EventType = elements.EventType

const (
	EventTypeMousePress   EventType = elements.EventTypeMousePress
	EventTypeMouseRelease EventType = elements.EventTypeMouseRelease
	EventTypeMouseEnter   EventType = elements.EventTypeMouseEnter
	EventTypeMouseLeave   EventType = elements.EventTypeMouseLeave
	EventTypeMouseMove    EventType = elements.EventTypeMouseMove
	EventTypeMouseScroll  EventType = elements.EventTypeMouseScroll
	EventTypeMouseDrag    EventType = elements.EventTypeMouseDrag

	EventTypeKeyPress   EventType = elements.EventTypeKeyPress
	EventTypeKeyRelease EventType = elements.EventTypeKeyRelease

	EventTypePaste EventType = elements.EventTypePaste

	EventTypeFocus EventType = elements.EventTypeFocus
	EventTypeBlur  EventType = elements.EventTypeBlur

	EventTypeDestroy EventType = elements.EventTypeDestroy
)

type EventFocus = elements.EventFocus

type EventBlur = elements.EventBlur

type EventInput = elements.EventInput

type EventSubmit = elements.EventSubmit

type EventPaste = elements.EventPaste

type EventKey = elements.EventKey

type Key = events.Key
type KeyAction = events.KeyAction

const (
	KeyActionPress   KeyAction = events.KeyActionPress
	KeyActionRelease KeyAction = events.KeyActionRelease
	KeyActionRepeat  KeyAction = events.KeyActionRepeat
)

const (
	KeyUnknown Key = events.KeyUnknown

	KeyCtrl  Key = events.KeyCtrl
	KeyShift Key = events.KeyShift
	KeyAlt   Key = events.KeyAlt
	KeyMeta  Key = events.KeyMeta

	Key0 Key = events.Key0
	Key1 Key = events.Key1
	Key2 Key = events.Key2
	Key3 Key = events.Key3
	Key4 Key = events.Key4
	Key5 Key = events.Key5
	Key6 Key = events.Key6
	Key7 Key = events.Key7
	Key8 Key = events.Key8
	Key9 Key = events.Key9

	KeyF1  Key = events.KeyF1
	KeyF2  Key = events.KeyF2
	KeyF3  Key = events.KeyF3
	KeyF4  Key = events.KeyF4
	KeyF5  Key = events.KeyF5
	KeyF6  Key = events.KeyF6
	KeyF7  Key = events.KeyF7
	KeyF8  Key = events.KeyF8
	KeyF9  Key = events.KeyF9
	KeyF10 Key = events.KeyF10
	KeyF11 Key = events.KeyF11
	KeyF12 Key = events.KeyF12

	KeyUp    Key = events.KeyUp
	KeyDown  Key = events.KeyDown
	KeyLeft  Key = events.KeyLeft
	KeyRight Key = events.KeyRight

	KeyEnter     Key = events.KeyEnter
	KeyTab       Key = events.KeyTab
	KeyBackspace Key = events.KeyBackspace
	KeySpace     Key = events.KeySpace
	KeyEscape    Key = events.KeyEscape
	KeyClear     Key = events.KeyClear
	KeyEnd       Key = events.KeyEnd
	KeyHome      Key = events.KeyHome
	KeyInsert    Key = events.KeyInsert
	KeyDelete    Key = events.KeyDelete
	KeyPageUp    Key = events.KeyPageUp
	KeyPageDown  Key = events.KeyPageDown
)

func KeyRune(r rune) Key { return events.KeyRune(r) }

type EventMouse = elements.EventMouse

type MouseAction = events.MouseAction
type MouseButton = events.MouseButton

const (
	MouseActionPress   MouseAction = events.MouseActionPress
	MouseActionRelease MouseAction = events.MouseActionRelease
	MouseActionMove    MouseAction = events.MouseActionMove
	MouseActionDrag    MouseAction = events.MouseActionDrag
	MouseActionScroll  MouseAction = events.MouseActionScroll
)

const (
	MouseNone       MouseButton = events.MouseNone
	MouseLeft       MouseButton = events.MouseLeft
	MouseMiddle     MouseButton = events.MouseMiddle
	MouseRight      MouseButton = events.MouseRight
	MouseWheelUp    MouseButton = events.MouseWheelUp
	MouseWheelDown  MouseButton = events.MouseWheelDown
	MouseWheelLeft  MouseButton = events.MouseWheelLeft
	MouseWheelRight MouseButton = events.MouseWheelRight
	MouseBackward   MouseButton = events.MouseBackward
	MouseForward    MouseButton = events.MouseForward
	Mouse10         MouseButton = events.Mouse10
	Mouse11         MouseButton = events.Mouse11
)
