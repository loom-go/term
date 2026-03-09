package events

import (
	"bytes"
	"context"
	"github.com/loom-go/term/core/stdio"
	"github.com/loom-go/term/core/sync"
	"strconv"
	"strings"
)

// bracketed paste markers
const (
	pasteStart = "\x1b[200~"
	pasteEnd   = "\x1b[201~"
)

type KeyboardListener struct {
	ctx         context.Context
	keyEvents   *sync.Broadcaster[*EventKey]
	pasteEvents *sync.Broadcaster[*EventPaste]
}

func NewKeyboardListener(ctx context.Context) *KeyboardListener {
	listener := &KeyboardListener{
		ctx:         ctx,
		keyEvents:   sync.NewBroadcaster[*EventKey](ctx),
		pasteEvents: sync.NewBroadcaster[*EventPaste](ctx),
	}

	go listener.watch()

	return listener
}

func (l *KeyboardListener) ListenKey(ctx context.Context) <-chan *EventKey {
	return l.keyEvents.Listen(ctx)
}

func (l *KeyboardListener) ListenPaste(ctx context.Context) <-chan *EventPaste {
	return l.pasteEvents.Listen(ctx)
}

func (l *KeyboardListener) watch() {
	stdin := stdio.Stdin.Listen(1024)
	events := stdio.NewBufferedConsumer(func(buf []byte) (consumed int, complete bool) {
		if len(buf) == 0 {
			return 0, false
		}

		pasteEvent, consumed := l.parseBracketedPaste(buf)
		if pasteEvent != nil {
			l.pasteEvents.Broadcast(pasteEvent)
			return consumed, true
		}

		if consumed > 0 {
			return consumed, false
		}

		keyEvent, consumed := l.parseKeyEvent(buf)
		if keyEvent != nil {
			l.keyEvents.Broadcast(keyEvent)
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

func (l *KeyboardListener) parseBracketedPaste(buf []byte) (*EventPaste, int) {
	if !bytes.HasPrefix(buf, []byte(pasteStart)) {
		return nil, 0
	}

	endIdx := bytes.Index(buf[len(pasteStart):], []byte(pasteEnd))
	if endIdx == -1 {
		return nil, 0
	}

	start := len(pasteStart)
	end := start + endIdx
	pastedText := string(buf[start:end])
	consumed := end + len(pasteEnd)

	return &EventPaste{Value: pastedText}, consumed
}

func (l *KeyboardListener) parseKeyEvent(buf []byte) (event *EventKey, consumed int) {
	if len(buf) == 0 {
		return nil, 0
	}

	if isStandaloneModifier(buf[0]) {
		return nil, 1
	}

	if buf[0] == 0x1b {
		if len(buf) < 2 {
			return nil, 0
		}

		switch buf[1] {
		case '[':
			return l.parseCSISequence(buf)
		case 'O':
			return l.parseSS3Sequence(buf)
		case ']':
			return l.skipSequence(buf, 0x07, 0x1b, 0x5c)
		case 'P':
			return l.skipSequence(buf, 0x07, 0x1b, 0x5c)
		case '_':
			return l.skipSequence(buf, 0x07, 0x1b, 0x5c)
		case '^':
			return l.skipSequence(buf, 0x07, 0x1b, 0x5c)
		}

		if buf[1] != 0x1b {
			return l.parseAltKey(buf[1])
		}

		return &EventKey{Key: KeyEscape, Action: KeyActionPress}, 1
	}

	return l.parseSingleByte(buf[0])
}

func (l *KeyboardListener) parseCSISequence(buf []byte) (event *EventKey, consumed int) {
	if len(buf) < 3 {
		return nil, 0
	}

	if buf[2] == '<' {
		return l.skipToTerminator(buf, 3, 'M', 'm')
	}
	if buf[2] == 'M' && len(buf) >= 6 {
		return nil, 6
	}

	content := buf[2:]

	if idx := indexOf(content, 'u'); idx >= 0 {
		return l.parseKittySequence(buf, idx+2)
	}

	if len(content) >= 2 && content[0] == '2' && content[1] == '7' {
		return l.parseModifyOtherKeys(buf)
	}

	if len(content) >= 3 && content[0] == '2' && content[1] == '0' && content[2] == '0' {
		return nil, 0
	}
	if len(content) >= 3 && content[0] == '2' && content[1] == '0' && content[2] == '1' {
		return l.skipTildeTerminator(buf, 5)
	}

	if content[0] == 'I' {
		return nil, 3
	}
	if content[0] == 'O' && len(buf) == 3 {
		return nil, 3
	}

	for i, b := range content {
		switch b {
		case 'R', 'c', 't':
			return nil, i + 3
		}
	}

	for i, b := range content {
		if b == '~' {
			return l.parseTildeSequence(buf[:i+3], i+3)
		}
		if key, ok := csiLetterMap[b]; ok {
			return l.parseLetterSequence(buf[:i+3], i+3, key)
		}
		if b >= 'A' && b <= 'Z' {
			return nil, i + 3
		}
		if b >= 0x40 && b <= 0x7E {
			return nil, i + 3
		}
	}

	if len(buf) > 50 {
		return nil, 2
	}

	return nil, 0
}

func (l *KeyboardListener) parseSS3Sequence(buf []byte) (event *EventKey, consumed int) {
	if len(buf) < 3 {
		return nil, 0
	}

	switch buf[2] {
	case 'P':
		return &EventKey{Key: KeyF1, Action: KeyActionPress}, 3
	case 'Q':
		return &EventKey{Key: KeyF2, Action: KeyActionPress}, 3
	case 'R':
		return &EventKey{Key: KeyF3, Action: KeyActionPress}, 3
	case 'S':
		return &EventKey{Key: KeyF4, Action: KeyActionPress}, 3
	}

	return nil, 3
}

func (l *KeyboardListener) parseKittySequence(buf []byte, uIdx int) (event *EventKey, consumed int) {
	content := string(buf[2:uIdx])
	fields := splitParams(content, ";")

	if len(fields) == 0 {
		return nil, uIdx + 1
	}

	field1Parts := splitParams(fields[0], ":")
	codepointStr := field1Parts[0]

	codepoint, err := strconv.Atoi(codepointStr)
	if err != nil {
		return nil, uIdx + 1
	}

	if isKittyModifier(codepoint) {
		return nil, uIdx + 1
	}

	var key Key
	var action KeyAction = KeyActionPress
	var mods Key

	if codepoint == 1 && len(fields) >= 2 {
		lastField := fields[len(fields)-1]
		if len(lastField) == 1 {
			if k, ok := csiLetterMap[lastField[0]]; ok {
				key = k
				if len(fields) >= 2 {
					mods = l.parseKittyModifiers(fields[1])
					action = l.parseKittyEventType(fields[1])
				}
				return &EventKey{Key: key | mods, Action: action}, uIdx + 1
			}
		}
	}

	if strings.HasSuffix(content, "~") {
		tildeNum, _ := strconv.Atoi(field1Parts[0])
		if k, ok := tildeCodeMap[tildeNum]; ok {
			key = k
			if len(fields) >= 2 {
				mods = l.parseKittyModifiers(fields[1])
				action = l.parseKittyEventType(fields[1])
			}
			return &EventKey{Key: key | mods, Action: action}, uIdx + 1
		}
	}

	shiftedCodepoint := 0
	if len(field1Parts) >= 2 {
		shiftedCodepoint, _ = strconv.Atoi(field1Parts[1])
	}

	if len(fields) >= 2 {
		mods = l.parseKittyModifiers(fields[1])
		action = l.parseKittyEventType(fields[1])
	}

	if k, ok := kittyKeyMap[codepoint]; ok {
		key = k
	} else if codepoint >= 32 && codepoint <= 126 {
		if mods&KeyShift != 0 && shiftedCodepoint > 0 {
			key = KeyRune(rune(shiftedCodepoint))
		} else {
			key = KeyRune(rune(codepoint))
		}
	} else if codepoint >= 1 && codepoint <= 26 {
		key = KeyRune(rune('a'+codepoint-1)) | KeyCtrl
	} else {
		key = KeyRune(rune(codepoint))
	}

	return &EventKey{Key: key | mods, Action: action}, uIdx + 1
}

func (l *KeyboardListener) parseModifyOtherKeys(buf []byte) (event *EventKey, consumed int) {
	for i := 2; i < len(buf); i++ {
		if buf[i] == '~' {
			content := string(buf[2:i])
			parts := splitParams(content, ";")
			if len(parts) < 3 {
				return nil, i + 1
			}

			modifier, _ := strconv.Atoi(parts[1])
			charCode, _ := strconv.Atoi(parts[2])

			mods := l.parseTraditionalModifiers(modifier - 1)
			key := l.charCodeToKey(charCode)

			return &EventKey{Key: key | mods, Action: KeyActionPress}, i + 1
		}
	}

	return nil, 0
}

func (l *KeyboardListener) parseTildeSequence(buf []byte, totalLen int) (event *EventKey, consumed int) {
	content := string(buf[2 : len(buf)-1])
	parts := splitParams(content, ";")

	if len(parts) == 0 {
		return nil, totalLen
	}

	code, _ := strconv.Atoi(parts[0])
	key := l.tildeCodeToKey(code)

	modifier, action := 1, KeyActionPress
	if len(parts) >= 2 {
		modifier, action = parseModField(parts[1])
	}

	return &EventKey{Key: key | l.parseTraditionalModifiers(modifier-1), Action: action}, totalLen
}

func (l *KeyboardListener) parseLetterSequence(buf []byte, totalLen int, key Key) (event *EventKey, consumed int) {
	content := string(buf[2 : len(buf)-1])
	parts := splitParams(content, ";")

	modifier, action := 1, KeyActionPress
	if len(parts) >= 2 {
		modifier, action = parseModField(parts[1])
	}

	return &EventKey{Key: key | l.parseTraditionalModifiers(modifier-1), Action: action}, totalLen
}

func (l *KeyboardListener) parseAltKey(b byte) (event *EventKey, consumed int) {
	key := byteToKey(b)
	return &EventKey{Key: key | KeyAlt, Action: KeyActionPress}, 2
}

func (l *KeyboardListener) parseSingleByte(b byte) (event *EventKey, consumed int) {
	key := byteToControlKey(b)
	if key != 0 {
		return &EventKey{Key: key, Action: KeyActionPress}, 1
	}

	if b == 32 {
		return &EventKey{Key: KeySpace, Action: KeyActionPress}, 1
	}

	if b > 32 && b <= 126 {
		return &EventKey{Key: KeyRune(rune(b)), Action: KeyActionPress}, 1
	}

	return &EventKey{Key: KeyUnknown, Action: KeyActionPress}, 1
}

func (l *KeyboardListener) parseKittyModifiers(field string) Key {
	parts := splitParams(field, ":")
	mod, _ := strconv.Atoi(parts[0])
	return l.parseTraditionalModifiers(mod - 1)
}

func (l *KeyboardListener) parseKittyEventType(field string) KeyAction {
	parts := splitParams(field, ":")
	if len(parts) < 2 {
		return KeyActionPress
	}

	eventType, _ := strconv.Atoi(parts[1])
	switch eventType {
	case 2:
		return KeyActionRepeat
	case 3:
		return KeyActionRelease
	}
	return KeyActionPress
}

func parseModField(field string) (mod int, action KeyAction) {
	parts := splitParams(field, ":")
	mod, _ = strconv.Atoi(parts[0])
	if mod <= 0 {
		mod = 1
	}
	if len(parts) >= 2 {
		eventType, _ := strconv.Atoi(parts[1])
		switch eventType {
		case 2:
			action = KeyActionRepeat
		case 3:
			action = KeyActionRelease
		default:
			action = KeyActionPress
		}
	} else {
		action = KeyActionPress
	}
	return
}

func (l *KeyboardListener) parseTraditionalModifiers(bitmask int) Key {
	var key Key
	if bitmask&1 != 0 {
		key |= KeyShift
	}
	if bitmask&2 != 0 {
		key |= KeyAlt
	}
	if bitmask&4 != 0 {
		key |= KeyCtrl
	}
	if bitmask&8 != 0 {
		key |= KeyMeta
	}
	return key
}

func (l *KeyboardListener) tildeCodeToKey(code int) Key {
	if k, ok := tildeCodeMap[code]; ok {
		return k
	}
	return KeyUnknown
}

func (l *KeyboardListener) charCodeToKey(charCode int) Key {
	switch charCode {
	case 9:
		return KeyTab
	case 13:
		return KeyEnter
	case 27:
		return KeyEscape
	case 127:
		return KeyBackspace
	}

	if charCode >= 1 && charCode <= 26 {
		return KeyRune(rune('a'+charCode-1)) | KeyCtrl
	}
	if charCode >= 32 && charCode <= 126 {
		return KeyRune(rune(charCode))
	}

	return KeyUnknown
}

func (l *KeyboardListener) skipSequence(buf []byte, term1 byte, term2 byte, term3 byte) (event *EventKey, consumed int) {
	if len(buf) < 3 {
		return nil, 0
	}

	for i := 2; i < len(buf); i++ {
		if buf[i] == term1 {
			return nil, i + 1
		}
		if buf[i] == term2 && i+1 < len(buf) && buf[i+1] == term3 {
			return nil, i + 2
		}
	}

	if len(buf) > 200 {
		return nil, 2
	}

	return nil, 0
}

func (l *KeyboardListener) skipToTerminator(buf []byte, start int, term1 byte, term2 byte) (event *EventKey, consumed int) {
	for i := start; i < len(buf); i++ {
		if buf[i] == term1 || buf[i] == term2 {
			return nil, i + 1
		}
	}
	return nil, 0
}

func (l *KeyboardListener) skipTildeTerminator(buf []byte, start int) (event *EventKey, consumed int) {
	for i := start; i < len(buf); i++ {
		if buf[i] == '~' {
			return nil, i + 1
		}
	}
	return nil, 0
}

func isStandaloneModifier(b byte) bool {
	return b >= 0x1C && b <= 0x1F
}

func isKittyModifier(codepoint int) bool {
	return codepoint >= 57441 && codepoint <= 57454
}

func splitParams(s string, sep string) []string {
	var parts []string
	start := 0
	sepLen := len(sep)
	for i := 0; i <= len(s)-sepLen; i++ {
		if s[i:i+sepLen] == sep {
			parts = append(parts, s[start:i])
			start = i + sepLen
			i += sepLen - 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func indexOf(buf []byte, target byte) int {
	for i, b := range buf {
		if b == target {
			return i
		}
	}
	return -1
}

func byteToKey(b byte) Key {
	if b >= 'a' && b <= 'z' {
		return KeyRune(rune(b - 32))
	}
	if b >= 'A' && b <= 'Z' {
		return KeyRune(rune(b))
	}
	return KeyRune(rune(b))
}

func byteToControlKey(b byte) Key {
	switch b {
	case 0x00:
		return KeySpace | KeyCtrl
	case 0x01:
		return KeyRune('a') | KeyCtrl
	case 0x02:
		return KeyRune('b') | KeyCtrl
	case 0x03:
		return KeyRune('c') | KeyCtrl
	case 0x04:
		return KeyRune('d') | KeyCtrl
	case 0x05:
		return KeyRune('e') | KeyCtrl
	case 0x06:
		return KeyRune('f') | KeyCtrl
	case 0x07:
		return KeyRune('g') | KeyCtrl
	case 0x08:
		return KeyBackspace
	case 0x09:
		return KeyTab
	case 0x0A:
		return KeyEnter
	case 0x0B:
		return KeyRune('k') | KeyCtrl
	case 0x0C:
		return KeyRune('l') | KeyCtrl
	case 0x0D:
		return KeyEnter
	case 0x0E:
		return KeyRune('n') | KeyCtrl
	case 0x0F:
		return KeyRune('o') | KeyCtrl
	case 0x10:
		return KeyRune('p') | KeyCtrl
	case 0x11:
		return KeyRune('q') | KeyCtrl
	case 0x12:
		return KeyRune('r') | KeyCtrl
	case 0x13:
		return KeyRune('s') | KeyCtrl
	case 0x14:
		return KeyRune('t') | KeyCtrl
	case 0x15:
		return KeyRune('u') | KeyCtrl
	case 0x16:
		return KeyRune('v') | KeyCtrl
	case 0x17:
		return KeyRune('w') | KeyCtrl
	case 0x18:
		return KeyRune('x') | KeyCtrl
	case 0x19:
		return KeyRune('y') | KeyCtrl
	case 0x1A:
		return KeyRune('z') | KeyCtrl
	case 0x1B:
		return KeyEscape
	case 0x1C:
		return KeyRune('\\') | KeyCtrl
	case 0x1D:
		return KeyRune(']') | KeyCtrl
	case 0x1E:
		return KeyRune('^') | KeyCtrl
	case 0x1F:
		return KeyRune('_') | KeyCtrl
	case 0x7F:
		return KeyBackspace
	}
	return 0
}

var kittyKeyMap = map[int]Key{
	9:     KeyTab,
	13:    KeyEnter,
	27:    KeyEscape,
	32:    KeySpace,
	127:   KeyBackspace,
	57348: KeyInsert,
	57349: KeyDelete,
	57350: KeyLeft,
	57351: KeyRight,
	57352: KeyUp,
	57353: KeyDown,
	57354: KeyPageUp,
	57355: KeyPageDown,
	57356: KeyHome,
	57357: KeyEnd,
	57364: KeyF1,
	57365: KeyF2,
	57366: KeyF3,
	57367: KeyF4,
	57368: KeyF5,
	57369: KeyF6,
	57370: KeyF7,
	57371: KeyF8,
	57372: KeyF9,
	57373: KeyF10,
	57374: KeyF11,
	57375: KeyF12,
}

var csiLetterMap = map[byte]Key{
	'A': KeyUp,
	'B': KeyDown,
	'C': KeyRight,
	'D': KeyLeft,
	'F': KeyEnd,
	'H': KeyHome,
	'P': KeyF1,
	'Q': KeyF2,
	'R': KeyF3,
	'S': KeyF4,
}

var tildeCodeMap = map[int]Key{
	1:  KeyHome,
	2:  KeyInsert,
	3:  KeyDelete,
	4:  KeyEnd,
	5:  KeyPageUp,
	6:  KeyPageDown,
	7:  KeyHome,
	8:  KeyEnd,
	11: KeyF1,
	12: KeyF2,
	13: KeyF3,
	14: KeyF4,
	15: KeyF5,
	17: KeyF6,
	18: KeyF7,
	19: KeyF8,
	20: KeyF9,
	21: KeyF10,
	23: KeyF11,
	24: KeyF12,
}
