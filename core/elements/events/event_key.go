package events

import "fmt"

type EventKey struct {
	Key    Key
	Action KeyAction
}

func (e EventKey) Pressed() bool {
	return e.Action == KeyActionPress
}

func (e EventKey) Released() bool {
	return e.Action == KeyActionRelease
}

func (e EventKey) Repeated() bool {
	return e.Action == KeyActionRepeat
}

func (e EventKey) String() string {
	return "Key(" + e.Key.String() + ", " + e.Action.String() + ")"
}

type KeyAction int

const (
	KeyActionPress KeyAction = iota
	KeyActionRelease
	KeyActionRepeat
)

func (a KeyAction) String() string {
	switch a {
	case KeyActionPress:
		return "press"
	case KeyActionRelease:
		return "release"
	case KeyActionRepeat:
		return "repeat"
	default:
		return "unknown"
	}
}

// Key encodes the key type, value, and modifiers in a uint32.
//
// Layout:
//
//	bits 0-19: key value (key code or unicode rune)
//	bits 20-21: key type (normal, rune, unknown)
//	bits 22-28: modifiers (ctrl, shift, alt, meta)
//	bits 29-31: reserved/unused
type Key uint32

const (
	keyValueMask Key = 0xFFFFF     // bits 0-19
	keyTypeMask  Key = 0b11 << 20  // bits 20-21
	keyModMask   Key = 0x1FF << 22 // bits 22-28

	keyTypeNamed   Key = 0b00 << 20
	keyTypeRune    Key = 0b01 << 20
	keyTypeUnknown Key = 0b10 << 20
)

const (
	KeyUnknown Key = keyTypeUnknown
)

const (
	KeyCtrl Key = 1 << (22 + iota)
	KeyShift
	KeyAlt
	KeyMeta
)

const (
	KeyNull Key = keyTypeNamed | iota

	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9

	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12

	KeyUp
	KeyDown
	KeyLeft
	KeyRight

	KeyEnter
	KeyTab
	KeyBackspace
	KeySpace
	KeyEscape
	KeyClear
	KeyEnd
	KeyHome
	KeyInsert
	KeyDelete
	KeyPageUp
	KeyPageDown
)

func KeyRune(r rune) Key {
	if r > rune(keyValueMask) {
		return KeyUnknown
	}

	return keyTypeRune | Key(r)
}

func (k Key) IsRune() bool    { return k&keyTypeMask == keyTypeRune }
func (k Key) IsUnknown() bool { return k&keyTypeMask == keyTypeUnknown }

func (k Key) Ctrl() bool  { return k&KeyCtrl != 0 }
func (k Key) Shift() bool { return k&KeyShift != 0 }
func (k Key) Alt() bool   { return k&KeyAlt != 0 }
func (k Key) Meta() bool  { return k&KeyMeta != 0 }

func (k Key) Value() Key { return k & keyValueMask }
func (k Key) Mods() Key  { return k & keyModMask }
func (k Key) Rune() rune { return rune(k & keyValueMask) }

func (k Key) String() string {
	var str string

	if k.Ctrl() {
		str += "ctrl+"
	}
	if k.Shift() && !k.IsRune() {
		str += "shift+"
	}
	if k.Alt() {
		str += "alt+"
	}
	if k.Meta() {
		str += "meta+"
	}

	switch k & keyTypeMask {
	case keyTypeNamed:
		key := k.Value()
		if name, ok := keyNames[key]; ok {
			str += name
		} else {
			str += fmt.Sprintf("unknown(%d)", key)
		}
	case keyTypeRune:
		str += string(k.Rune())
	default:
		str += "unknown"
	}

	return str
}

var keyNames = map[Key]string{
	KeyNull: "",

	Key0: "0",
	Key1: "1",
	Key2: "2",
	Key3: "3",
	Key4: "4",
	Key5: "5",
	Key6: "6",
	Key7: "7",
	Key8: "8",
	Key9: "9",

	KeyF1:  "F1",
	KeyF2:  "F2",
	KeyF3:  "F3",
	KeyF4:  "F4",
	KeyF5:  "F5",
	KeyF6:  "F6",
	KeyF7:  "F7",
	KeyF8:  "F8",
	KeyF9:  "F9",
	KeyF10: "F10",
	KeyF11: "F11",
	KeyF12: "F12",

	KeyUp:    "up",
	KeyDown:  "down",
	KeyLeft:  "left",
	KeyRight: "right",

	KeyEnter:     "enter",
	KeyTab:       "tab",
	KeyBackspace: "backspace",
	KeySpace:     "space",
	KeyEscape:    "escape",
	KeyClear:     "clear",
	KeyEnd:       "end",
	KeyHome:      "home",
	KeyInsert:    "insert",
	KeyDelete:    "delete",
	KeyPageUp:    "pageup",
	KeyPageDown:  "pagedown",
}
