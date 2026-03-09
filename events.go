package term

import "github.com/loom-go/term/core"

// todo: see if we could generate these barrel files. its prone to error and anoying to maintain

type Event = core.Event
type EventOptions = core.EventOptions

type EventType = core.EventType

const (
	EventTypeMousePress   EventType = core.EventTypeMousePress
	EventTypeMouseRelease EventType = core.EventTypeMouseRelease
	EventTypeMouseEnter   EventType = core.EventTypeMouseEnter
	EventTypeMouseLeave   EventType = core.EventTypeMouseLeave
	EventTypeMouseMove    EventType = core.EventTypeMouseMove
	EventTypeMouseScroll  EventType = core.EventTypeMouseScroll
	EventTypeMouseDrag    EventType = core.EventTypeMouseDrag

	EventTypeKeyPress   EventType = core.EventTypeKeyPress
	EventTypeKeyRelease EventType = core.EventTypeKeyRelease

	EventTypePaste EventType = core.EventTypePaste

	EventTypeFocus EventType = core.EventTypeFocus
	EventTypeBlur  EventType = core.EventTypeBlur

	EventTypeDestroy EventType = core.EventTypeDestroy
)

type EventFocus = core.EventFocus

type EventBlur = core.EventBlur

type EventInput = core.EventInput

type EventSubmit = core.EventSubmit

type EventPaste = core.EventPaste

type EventKey = core.EventKey

type Key = core.Key
type KeyAction = core.KeyAction

const (
	KeyActionPress   KeyAction = core.KeyActionPress
	KeyActionRelease KeyAction = core.KeyActionRelease
	KeyActionRepeat  KeyAction = core.KeyActionRepeat
)

const (
	KeyUnknown Key = core.KeyUnknown

	KeyCtrl  Key = core.KeyCtrl
	KeyShift Key = core.KeyShift
	KeyAlt   Key = core.KeyAlt
	KeyMeta  Key = core.KeyMeta

	Key0 Key = core.Key0
	Key1 Key = core.Key1
	Key2 Key = core.Key2
	Key3 Key = core.Key3
	Key4 Key = core.Key4
	Key5 Key = core.Key5
	Key6 Key = core.Key6
	Key7 Key = core.Key7
	Key8 Key = core.Key8
	Key9 Key = core.Key9

	KeyF1  Key = core.KeyF1
	KeyF2  Key = core.KeyF2
	KeyF3  Key = core.KeyF3
	KeyF4  Key = core.KeyF4
	KeyF5  Key = core.KeyF5
	KeyF6  Key = core.KeyF6
	KeyF7  Key = core.KeyF7
	KeyF8  Key = core.KeyF8
	KeyF9  Key = core.KeyF9
	KeyF10 Key = core.KeyF10
	KeyF11 Key = core.KeyF11
	KeyF12 Key = core.KeyF12

	KeyUp    Key = core.KeyUp
	KeyDown  Key = core.KeyDown
	KeyLeft  Key = core.KeyLeft
	KeyRight Key = core.KeyRight

	KeyEnter     Key = core.KeyEnter
	KeyTab       Key = core.KeyTab
	KeyBackspace Key = core.KeyBackspace
	KeySpace     Key = core.KeySpace
	KeyEscape    Key = core.KeyEscape
	KeyClear     Key = core.KeyClear
	KeyEnd       Key = core.KeyEnd
	KeyHome      Key = core.KeyHome
	KeyInsert    Key = core.KeyInsert
	KeyDelete    Key = core.KeyDelete
	KeyPageUp    Key = core.KeyPageUp
	KeyPageDown  Key = core.KeyPageDown
)

func KeyRune(r rune) Key { return core.KeyRune(r) }

type EventMouse = core.EventMouse

type MouseAction = core.MouseAction
type MouseButton = core.MouseButton

const (
	MouseActionPress   MouseAction = core.MouseActionPress
	MouseActionRelease MouseAction = core.MouseActionRelease
	MouseActionMove    MouseAction = core.MouseActionMove
	MouseActionDrag    MouseAction = core.MouseActionDrag
	MouseActionScroll  MouseAction = core.MouseActionScroll
)

const (
	MouseNone       MouseButton = core.MouseNone
	MouseLeft       MouseButton = core.MouseLeft
	MouseMiddle     MouseButton = core.MouseMiddle
	MouseRight      MouseButton = core.MouseRight
	MouseWheelUp    MouseButton = core.MouseWheelUp
	MouseWheelDown  MouseButton = core.MouseWheelDown
	MouseWheelLeft  MouseButton = core.MouseWheelLeft
	MouseWheelRight MouseButton = core.MouseWheelRight
	MouseBackward   MouseButton = core.MouseBackward
	MouseForward    MouseButton = core.MouseForward
	Mouse10         MouseButton = core.Mouse10
	Mouse11         MouseButton = core.Mouse11
)
