package events

import (
	"context"
	"github.com/loom-go/term/core/stdio"
	"github.com/loom-go/term/core/sync"
	"strconv"
)

type MouseListener struct {
	ctx           context.Context
	events        *sync.Broadcaster[*EventMouse]
	pressedButton MouseButton
}

func NewMouseListener(ctx context.Context) *MouseListener {
	listener := &MouseListener{
		ctx:           ctx,
		events:        sync.NewBroadcaster[*EventMouse](ctx),
		pressedButton: MouseNone,
	}

	go listener.watch()

	return listener
}

func (l *MouseListener) Listen(ctx context.Context) <-chan *EventMouse {
	return l.events.Listen(ctx)
}

func (l *MouseListener) watch() {
	stdin := stdio.Stdin.Listen(1024)
	events := stdio.NewBufferedConsumer(func(buf []byte) (consumed int, complete bool) {
		event, consumed := l.parseMouseEvent(buf)
		if event != nil {
			l.events.Broadcast(event)
			return consumed, true
		}

		return consumed, consumed > 0
	})

	for {
		select {
		case <-l.ctx.Done():
			return
		case buf := <-stdin:
			events.Feed(buf)
		}
	}
}

func (l *MouseListener) parseMouseEvent(buf []byte) (event *EventMouse, consumed int) {
	if len(buf) < 3 {
		return nil, 0
	}

	// If first byte isn't ESC, this isn't a mouse event - skip it
	if buf[0] != 0x1b {
		return nil, 1
	}

	// If second byte isn't '[', this might be an OSC sequence (ESC ] ...) or other
	// Skip just the ESC byte to let other parsers handle it
	if buf[1] != '[' {
		return nil, 1
	}

	// Now we have ESC [, check if it's a mouse event
	// Mouse events are: ESC [ < (SGR) or ESC [ M (X10)
	if len(buf) >= 4 && buf[2] == '<' {
		return l.parseSGRMouse(buf)
	}

	if len(buf) >= 6 && buf[2] == 'M' {
		return l.parseX10Mouse(buf)
	}

	// This is some other CSI sequence (like arrow keys: ESC [ A, ESC [ B, etc.)
	// Skip the whole CSI sequence to avoid breaking the stream
	// CSI sequences end with a byte in the range 0x40-0x7E (@ to ~)
	for i := 2; i < len(buf); i++ {
		if buf[i] >= 0x40 && buf[i] <= 0x7E {
			// Found the end of CSI sequence
			return nil, i + 1
		}
	}

	// Incomplete CSI sequence, don't consume anything yet
	return nil, 0
}

// parseX10Mouse parses X10 mouse encoding: ESC [ M Cb Cx Cy
// Cb, Cx, Cy are encoded values (button, x, y each offset by 32)
func (l *MouseListener) parseX10Mouse(buf []byte) (event *EventMouse, consumed int) {
	if len(buf) < 6 {
		return nil, 0
	}

	// Verify format: ESC [ M
	if buf[0] != 0x1b || buf[1] != '[' || buf[2] != 'M' {
		return nil, 0
	}

	cb := buf[3] - 32         // button encoding
	x := int(buf[4]) - 32 - 1 // x coordinate (1-based in protocol, convert to 0-based)
	y := int(buf[5]) - 32 - 1 // y coordinate (1-based in protocol, convert to 0-based)

	button, key := decodeButton(cb)
	action := l.determineAction(button, cb, false, cb&32 != 0)

	return &EventMouse{
		X:      x,
		Y:      y,
		Key:    key,
		Button: button,
		Action: action,
	}, 6
}

// parseSGRMouse parses SGR (1006) mouse encoding:
// Press: ESC [ < Cb ; Cx ; Cy M
// Release: ESC [ < Cb ; Cx ; Cy m
func (l *MouseListener) parseSGRMouse(buf []byte) (event *EventMouse, consumed int) {
	if len(buf) < 9 { // Minimum: ESC [ < 0 ; 0 ; 0 M (9 chars)
		return nil, 0
	}

	// Verify format starts with ESC [ <
	if buf[0] != 0x1b || buf[1] != '[' || buf[2] != '<' {
		return nil, 0
	}

	// Find the terminating M (press) or m (release)
	var endIdx int
	var isRelease bool
	for i := 3; i < len(buf); i++ {
		if buf[i] == 'm' {
			isRelease = true
			endIdx = i
			break
		}
		if buf[i] == 'M' {
			endIdx = i
			break
		}
	}

	if endIdx == 0 {
		return nil, 0 // No terminator found
	}

	// Parse the content between < and M/m: Cb ; Cx ; Cy
	content := string(buf[3:endIdx])
	parts := splitSGRParams(content)
	if len(parts) < 3 {
		return nil, 0
	}

	cb, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, 0
	}

	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, 0
	}

	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, 0
	}

	button, key := decodeButton(byte(cb))
	action := l.determineAction(button, byte(cb), isRelease, cb&32 != 0)

	// SGR uses 1-based coordinates, convert to 0-based
	return &EventMouse{
		X:      x - 1,
		Y:      y - 1,
		Key:    key,
		Button: button,
		Action: action,
	}, endIdx + 1
}

// determineAction determines the mouse action based on button type, button code, and event state
func (l *MouseListener) determineAction(button MouseButton, cb byte, isRelease, isMotionBitSet bool) MouseAction {
	// Wheel events are always scroll actions
	if isWheel(button) {
		return MouseActionScroll
	}

	// Motion event (bit 5 set in button code)
	if isMotionBitSet {
		if l.pressedButton != MouseNone {
			// Motion while button is held = drag
			return MouseActionDrag
		}
		return MouseActionMove
	}

	// Release event
	if isRelease || cb&3 == 3 {
		l.pressedButton = MouseNone
		return MouseActionRelease
	}

	// Press event - remember which button
	l.pressedButton = button
	return MouseActionPress
}

// isWheel checks if the button is a wheel button
func isWheel(button MouseButton) bool {
	return button == MouseWheelUp || button == MouseWheelDown ||
		button == MouseWheelLeft || button == MouseWheelRight
}

// splitSGRParams splits SGR mouse parameter string by semicolons
func splitSGRParams(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ';' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

// decodeButton decodes the button byte into button type and modifier key
// Bit layout: 0-1 = button (0=left, 1=middle, 2=right, 3=release/motion)
//
//	2   = shift
//	3   = meta/alt
//	4   = ctrl
//	5   = motion (32)
//	6   = wheel indicator (64)
//
// Wheel values (with bit 6 set):
//
//	64 = wheel up, 65 = wheel down
//	66 = wheel left, 67 = wheel right
func decodeButton(cb byte) (MouseButton, Key) {
	var key Key

	// Extract modifiers from bits 2, 3, 4
	if cb&4 != 0 {
		key |= KeyShift
	}
	if cb&8 != 0 {
		key |= KeyAlt
	}
	if cb&16 != 0 {
		key |= KeyCtrl
	}

	// Check for mouse wheel events (bit 6 set)
	if cb&64 != 0 {
		if cb == 66 {
			return MouseWheelLeft, key
		}
		if cb == 67 {
			return MouseWheelRight, key
		}
		if cb == 64 {
			return MouseWheelUp, key
		}
		if cb == 65 {
			return MouseWheelDown, key
		}
	}

	// Regular button detection from bits 0-1
	btn := cb & 3
	switch btn {
	case 0:
		return MouseLeft, key
	case 1:
		return MouseMiddle, key
	case 2:
		return MouseRight, key
	case 3:
		// Button 3 indicates button release or motion without button
		return MouseNone, key
	}

	return MouseNone, key
}
