package events

import (
	"regexp"
	"strconv"

	"github.com/AnatoleLucet/loom-term/core/types"
)

// NOTE: a lot of this code has been heavily inspired from Charm's BubbleTea. thanks a lot to them for the work they did on that project!
// https://github.com/charmbracelet/bubbletea/blob/f9233d51192293dadda7184a4de347738606c328/mouse.go

const x10MouseByteOffset = 32

const (
	bitShift  = 0b0000_0100
	bitAlt    = 0b0000_1000
	bitCtrl   = 0b0001_0000
	bitMotion = 0b0010_0000
	bitWheel  = 0b0100_0000
	bitAdd    = 0b1000_0000

	bitsMask = 0b0000_0011
)

var mouseSGRRegex = regexp.MustCompile(`(\d+);(\d+);(\d+)([Mm])`)

// decodeButton extracts button and modifier flags from a normalized button code (0-31).
// The button code should have the motion bit (32) already stripped if present.
func decodeButton(b int) (btn types.MouseButton, shift, alt, ctrl bool) {
	buttonBits := b & bitsMask

	switch {
	case b&bitAdd != 0:
		btn = types.MouseButtonBackward + types.MouseButton(buttonBits)
	case b&bitWheel != 0:
		btn = types.MouseButtonWheelUp + types.MouseButton(buttonBits)
	case buttonBits == bitsMask:
		btn = types.MouseButtonNone
	default:
		btn = types.MouseButtonLeft + types.MouseButton(buttonBits)
	}

	shift = b&bitShift != 0
	alt = b&bitAlt != 0
	ctrl = b&bitCtrl != 0
	return
}

// determineAction returns the appropriate action based on button, motion state, and release state.
// This centralizes the action logic shared between X10 and SGR protocols.
func determineAction(btn types.MouseButton, isMotion, isRelease bool) types.MouseAction {
	// Wheel buttons are always scroll events (they have no release or motion)
	if btn >= types.MouseButtonWheelUp && btn <= types.MouseButtonWheelRight {
		return types.MouseActionScroll
	}

	if isMotion {
		if btn == types.MouseButtonNone {
			return types.MouseActionMove
		}
		return types.MouseActionDrag
	}

	if isRelease {
		return types.MouseActionRelease
	}

	return types.MouseActionPress
}

// Parse SGR-encoded mouse events; SGR extended mouse events. SGR mouse events
// look like:
//
//	ESC [ < Cb ; Cx ; Cy (M or m)
//
// where:
//
//	Cb is the encoded button code
//	Cx is the x-coordinate of the mouse
//	Cy is the y-coordinate of the mouse
//	M is for button press, m is for button release
//
// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Extended-coordinates
func parseSGRMouseEvent(buf []byte) *types.EventMouse {
	str := string(buf[3:])
	matches := mouseSGRRegex.FindStringSubmatch(str)
	if len(matches) != 5 {
		panic("invalid mouse event") // todo: return error
	}

	rawBtn, _ := strconv.Atoi(matches[1])
	x, _ := strconv.Atoi(matches[2])
	y, _ := strconv.Atoi(matches[3])
	isRelease := matches[4] == "m"

	isWheel := rawBtn&bitWheel != 0
	// SGR encodes motion as values >= 32 (same bit pattern as X10 motion)
	// But wheel events also have values >= 32, so exclude them
	isMotion := !isWheel && rawBtn >= x10MouseByteOffset
	btnCode := rawBtn
	if isMotion {
		btnCode = rawBtn - x10MouseByteOffset
	}

	button, shift, alt, ctrl := decodeButton(btnCode)
	action := determineAction(button, isMotion, isRelease)

	return &types.EventMouse{
		X:      x - 1, // (1,1) is upper left, normalize to (0,0)
		Y:      y - 1,
		Button: button,
		Action: action,
		Shift:  shift,
		Alt:    alt,
		Ctrl:   ctrl,
	}
}

// Parse X10-encoded mouse events; the simplest kind. The last release of X10
// was December 1986, by the way. The original X10 mouse protocol limits the Cx
// and Cy coordinates to 223 (=255-032).
//
// X10 mouse events look like:
//
//	ESC [M Cb Cx Cy
//
// See: http://www.xfree86.org/current/ctlseqs.html#Mouse%20Tracking
func parseX10MouseEvent(buf []byte) *types.EventMouse {
	v := buf[3:6]
	rawBtn := int(v[0])

	isWheel := rawBtn&bitWheel != 0
	// X10 encodes motion as byte value >= 32 (bit 5 set)
	// But wheel events also have values >= 32, so exclude them
	isMotion := !isWheel && rawBtn >= x10MouseByteOffset
	btnCode := rawBtn
	if isMotion {
		btnCode = rawBtn - x10MouseByteOffset
	}

	button, shift, alt, ctrl := decodeButton(btnCode)

	// X10 encodes release as button code 3 (all button bits set)
	// This is only for press events - motion events keep their button info
	isRelease := !isMotion && !isWheel && (btnCode&bitsMask) == bitsMask

	action := determineAction(button, isMotion, isRelease)

	return &types.EventMouse{
		X:      int(v[1]) - x10MouseByteOffset - 1, // (1,1) is upper left, normalize to (0,0)
		Y:      int(v[2]) - x10MouseByteOffset - 1,
		Button: button,
		Action: action,
		Shift:  shift,
		Alt:    alt,
		Ctrl:   ctrl,
	}
}
