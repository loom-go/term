package events

import (
	"bytes"
	"sort"
	"unicode/utf8"

	"github.com/AnatoleLucet/loom-term/core/types"
)

// NOTE: a lot of this code has been heavily inspired from Charm's BubbleTea. thanks a lot to them for the work they did on that project!
// https://github.com/charmbracelet/bubbletea/blob/f9233d51192293dadda7184a4de347738606c328/key.go
// https://github.com/charmbracelet/bubbletea/blob/f9233d51192293dadda7184a4de347738606c328/key_sequences.go

// sequences maps escape sequences to key events
var sequences = map[string]types.EventKey{
	// Arrow keys
	"\x1b[A":    {Key: types.KeyUp},
	"\x1b[B":    {Key: types.KeyDown},
	"\x1b[C":    {Key: types.KeyRight},
	"\x1b[D":    {Key: types.KeyLeft},
	"\x1b[1;2A": {Key: types.KeyShiftUp},
	"\x1b[1;2B": {Key: types.KeyShiftDown},
	"\x1b[1;2C": {Key: types.KeyShiftRight},
	"\x1b[1;2D": {Key: types.KeyShiftLeft},
	"\x1b[OA":   {Key: types.KeyShiftUp},
	"\x1b[OB":   {Key: types.KeyShiftDown},
	"\x1b[OC":   {Key: types.KeyShiftRight},
	"\x1b[OD":   {Key: types.KeyShiftLeft},
	"\x1b[a":    {Key: types.KeyShiftUp},
	"\x1b[b":    {Key: types.KeyShiftDown},
	"\x1b[c":    {Key: types.KeyShiftRight},
	"\x1b[d":    {Key: types.KeyShiftLeft},
	"\x1b[1;3A": {Key: types.KeyUp, Alt: true},
	"\x1b[1;3B": {Key: types.KeyDown, Alt: true},
	"\x1b[1;3C": {Key: types.KeyRight, Alt: true},
	"\x1b[1;3D": {Key: types.KeyLeft, Alt: true},
	"\x1b[1;4A": {Key: types.KeyShiftUp, Alt: true},
	"\x1b[1;4B": {Key: types.KeyShiftDown, Alt: true},
	"\x1b[1;4C": {Key: types.KeyShiftRight, Alt: true},
	"\x1b[1;4D": {Key: types.KeyShiftLeft, Alt: true},
	"\x1b[1;5A": {Key: types.KeyCtrlUp},
	"\x1b[1;5B": {Key: types.KeyCtrlDown},
	"\x1b[1;5C": {Key: types.KeyCtrlRight},
	"\x1b[1;5D": {Key: types.KeyCtrlLeft},
	"\x1b[Oa":   {Key: types.KeyCtrlUp, Alt: true},
	"\x1b[Ob":   {Key: types.KeyCtrlDown, Alt: true},
	"\x1b[Oc":   {Key: types.KeyCtrlRight, Alt: true},
	"\x1b[Od":   {Key: types.KeyCtrlLeft, Alt: true},
	"\x1b[1;6A": {Key: types.KeyCtrlShiftUp},
	"\x1b[1;6B": {Key: types.KeyCtrlShiftDown},
	"\x1b[1;6C": {Key: types.KeyCtrlShiftRight},
	"\x1b[1;6D": {Key: types.KeyCtrlShiftLeft},
	"\x1b[1;7A": {Key: types.KeyCtrlUp, Alt: true},
	"\x1b[1;7B": {Key: types.KeyCtrlDown, Alt: true},
	"\x1b[1;7C": {Key: types.KeyCtrlRight, Alt: true},
	"\x1b[1;7D": {Key: types.KeyCtrlLeft, Alt: true},
	"\x1b[1;8A": {Key: types.KeyCtrlShiftUp, Alt: true},
	"\x1b[1;8B": {Key: types.KeyCtrlShiftDown, Alt: true},
	"\x1b[1;8C": {Key: types.KeyCtrlShiftRight, Alt: true},
	"\x1b[1;8D": {Key: types.KeyCtrlShiftLeft, Alt: true},

	// Miscellaneous keys
	"\x1b[Z": {Key: types.KeyShiftTab},

	"\x1b[2~":   {Key: types.KeyInsert},
	"\x1b[3;2~": {Key: types.KeyInsert, Alt: true},

	"\x1b[3~":   {Key: types.KeyDelete},
	"\x1b[3;3~": {Key: types.KeyDelete, Alt: true},

	"\x1b[5~":   {Key: types.KeyPgUp},
	"\x1b[5;3~": {Key: types.KeyPgUp, Alt: true},
	"\x1b[5;5~": {Key: types.KeyCtrlPgUp},
	"\x1b[5^":   {Key: types.KeyCtrlPgUp},
	"\x1b[5;7~": {Key: types.KeyCtrlPgUp, Alt: true},

	"\x1b[6~":   {Key: types.KeyPgDown},
	"\x1b[6;3~": {Key: types.KeyPgDown, Alt: true},
	"\x1b[6;5~": {Key: types.KeyCtrlPgDown},
	"\x1b[6^":   {Key: types.KeyCtrlPgDown},
	"\x1b[6;7~": {Key: types.KeyCtrlPgDown, Alt: true},

	"\x1b[1~":   {Key: types.KeyHome},
	"\x1b[H":    {Key: types.KeyHome},
	"\x1b[1;3H": {Key: types.KeyHome, Alt: true},
	"\x1b[1;5H": {Key: types.KeyCtrlHome},
	"\x1b[1;7H": {Key: types.KeyCtrlHome, Alt: true},
	"\x1b[1;2H": {Key: types.KeyShiftHome},
	"\x1b[1;4H": {Key: types.KeyShiftHome, Alt: true},
	"\x1b[1;6H": {Key: types.KeyCtrlShiftHome},
	"\x1b[1;8H": {Key: types.KeyCtrlShiftHome, Alt: true},

	"\x1b[4~":   {Key: types.KeyEnd},
	"\x1b[F":    {Key: types.KeyEnd},
	"\x1b[1;3F": {Key: types.KeyEnd, Alt: true},
	"\x1b[1;5F": {Key: types.KeyCtrlEnd},
	"\x1b[1;7F": {Key: types.KeyCtrlEnd, Alt: true},
	"\x1b[1;2F": {Key: types.KeyShiftEnd},
	"\x1b[1;4F": {Key: types.KeyShiftEnd, Alt: true},
	"\x1b[1;6F": {Key: types.KeyCtrlShiftEnd},
	"\x1b[1;8F": {Key: types.KeyCtrlShiftEnd, Alt: true},

	"\x1b[7~": {Key: types.KeyHome},
	"\x1b[7^": {Key: types.KeyCtrlHome},
	"\x1b[7$": {Key: types.KeyShiftHome},
	"\x1b[7@": {Key: types.KeyCtrlShiftHome},

	"\x1b[8~": {Key: types.KeyEnd},
	"\x1b[8^": {Key: types.KeyCtrlEnd},
	"\x1b[8$": {Key: types.KeyShiftEnd},
	"\x1b[8@": {Key: types.KeyCtrlShiftEnd},

	// Function keys
	"\x1b[[A": {Key: types.KeyF1},
	"\x1b[[B": {Key: types.KeyF2},
	"\x1b[[C": {Key: types.KeyF3},
	"\x1b[[D": {Key: types.KeyF4},
	"\x1b[[E": {Key: types.KeyF5},

	"\x1bOP": {Key: types.KeyF1},
	"\x1bOQ": {Key: types.KeyF2},
	"\x1bOR": {Key: types.KeyF3},
	"\x1bOS": {Key: types.KeyF4},

	"\x1b[1;3P": {Key: types.KeyF1, Alt: true},
	"\x1b[1;3Q": {Key: types.KeyF2, Alt: true},
	"\x1b[1;3R": {Key: types.KeyF3, Alt: true},
	"\x1b[1;3S": {Key: types.KeyF4, Alt: true},

	"\x1b[11~": {Key: types.KeyF1},
	"\x1b[12~": {Key: types.KeyF2},
	"\x1b[13~": {Key: types.KeyF3},
	"\x1b[14~": {Key: types.KeyF4},

	"\x1b[15~":   {Key: types.KeyF5},
	"\x1b[15;3~": {Key: types.KeyF5, Alt: true},

	"\x1b[17~": {Key: types.KeyF6},
	"\x1b[18~": {Key: types.KeyF7},
	"\x1b[19~": {Key: types.KeyF8},
	"\x1b[20~": {Key: types.KeyF9},
	"\x1b[21~": {Key: types.KeyF10},

	"\x1b[17;3~": {Key: types.KeyF6, Alt: true},
	"\x1b[18;3~": {Key: types.KeyF7, Alt: true},
	"\x1b[19;3~": {Key: types.KeyF8, Alt: true},
	"\x1b[20;3~": {Key: types.KeyF9, Alt: true},
	"\x1b[21;3~": {Key: types.KeyF10, Alt: true},

	"\x1b[23~": {Key: types.KeyF11},
	"\x1b[24~": {Key: types.KeyF12},

	"\x1b[23;3~": {Key: types.KeyF11, Alt: true},
	"\x1b[24;3~": {Key: types.KeyF12, Alt: true},

	"\x1b[1;2P": {Key: types.KeyF13},
	"\x1b[1;2Q": {Key: types.KeyF14},

	"\x1b[25~": {Key: types.KeyF13},
	"\x1b[26~": {Key: types.KeyF14},

	"\x1b[25;3~": {Key: types.KeyF13, Alt: true},
	"\x1b[26;3~": {Key: types.KeyF14, Alt: true},

	"\x1b[1;2R": {Key: types.KeyF15},
	"\x1b[1;2S": {Key: types.KeyF16},

	"\x1b[28~": {Key: types.KeyF15},
	"\x1b[29~": {Key: types.KeyF16},

	"\x1b[28;3~": {Key: types.KeyF15, Alt: true},
	"\x1b[29;3~": {Key: types.KeyF16, Alt: true},

	"\x1b[15;2~": {Key: types.KeyF17},
	"\x1b[17;2~": {Key: types.KeyF18},
	"\x1b[18;2~": {Key: types.KeyF19},
	"\x1b[19;2~": {Key: types.KeyF20},

	"\x1b[31~": {Key: types.KeyF17},
	"\x1b[32~": {Key: types.KeyF18},
	"\x1b[33~": {Key: types.KeyF19},
	"\x1b[34~": {Key: types.KeyF20},

	// Powershell sequences
	"\x1bOA": {Key: types.KeyUp},
	"\x1bOB": {Key: types.KeyDown},
	"\x1bOC": {Key: types.KeyRight},
	"\x1bOD": {Key: types.KeyLeft},
}

// extSequences includes all sequences plus Alt variants
var extSequences = func() map[string]*types.EventKey {
	s := make(map[string]*types.EventKey, len(sequences)*2+128)

	// Copy all sequences
	for seq, key := range sequences {
		s[seq] = &key
		// Add Alt variant if not already Alt
		if !key.Alt {
			key.Alt = true
			s["\x1b"+seq] = &key
		}
	}

	// Control character mappings
	// These map control bytes to either special keys or runes with Ctrl=true
	ctrlMappings := map[byte]types.EventKey{
		0x00: {Key: types.KeySpace, Rune: ' ', Ctrl: true}, // Ctrl+Space (same byte as Ctrl+@ / Ctrl+2 in many terminals)
		0x01: {Key: types.KeyRunes, Rune: 'a', Ctrl: true}, // Ctrl+A
		0x02: {Key: types.KeyRunes, Rune: 'b', Ctrl: true}, // Ctrl+B
		0x03: {Key: types.KeyRunes, Rune: 'c', Ctrl: true}, // Ctrl+C
		0x04: {Key: types.KeyRunes, Rune: 'd', Ctrl: true}, // Ctrl+D
		0x05: {Key: types.KeyRunes, Rune: 'e', Ctrl: true}, // Ctrl+E
		0x06: {Key: types.KeyRunes, Rune: 'f', Ctrl: true}, // Ctrl+F
		0x07: {Key: types.KeyRunes, Rune: 'g', Ctrl: true}, // Ctrl+G
		0x08: {Key: types.KeyBackspace, Ctrl: true},        // Ctrl+H (Backspace)
		0x09: {Key: types.KeyTab},                          // Ctrl+I (Tab)
		0x0A: {Key: types.KeyRunes, Rune: 'j', Ctrl: true}, // Ctrl+J
		0x0B: {Key: types.KeyRunes, Rune: 'k', Ctrl: true}, // Ctrl+K
		0x0C: {Key: types.KeyRunes, Rune: 'l', Ctrl: true}, // Ctrl+L
		0x0D: {Key: types.KeyEnter},                        // Ctrl+M (Enter)
		0x0E: {Key: types.KeyRunes, Rune: 'n', Ctrl: true}, // Ctrl+N
		0x0F: {Key: types.KeyRunes, Rune: 'o', Ctrl: true}, // Ctrl+O
		0x10: {Key: types.KeyRunes, Rune: 'p', Ctrl: true}, // Ctrl+P
		0x11: {Key: types.KeyRunes, Rune: 'q', Ctrl: true}, // Ctrl+Q
		0x12: {Key: types.KeyRunes, Rune: 'r', Ctrl: true}, // Ctrl+R
		0x13: {Key: types.KeyRunes, Rune: 's', Ctrl: true}, // Ctrl+S
		0x14: {Key: types.KeyRunes, Rune: 't', Ctrl: true}, // Ctrl+T
		0x15: {Key: types.KeyRunes, Rune: 'u', Ctrl: true}, // Ctrl+U
		0x16: {Key: types.KeyRunes, Rune: 'v', Ctrl: true}, // Ctrl+V
		0x17: {Key: types.KeyRunes, Rune: 'w', Ctrl: true}, // Ctrl+W
		0x18: {Key: types.KeyRunes, Rune: 'x', Ctrl: true}, // Ctrl+X
		0x19: {Key: types.KeyRunes, Rune: 'y', Ctrl: true}, // Ctrl+Y
		0x1A: {Key: types.KeyRunes, Rune: 'z', Ctrl: true}, // Ctrl+Z
		// Note: 0x1B (ESC) is NOT included here to allow Alt+char sequences to work.
		// Lone ESC is handled separately in parsetypes.KeyEvent.
		// Ctrl+3 also sends 0x1B (same as ESC), so it cannot be distinguished in standard terminal mode.
		0x1C: {Key: types.KeyRunes, Rune: '\\', Ctrl: true}, // Ctrl+\
		0x1D: {Key: types.KeyRunes, Rune: ']', Ctrl: true},  // Ctrl+]
		0x1E: {Key: types.KeyRunes, Rune: '^', Ctrl: true},  // Ctrl+^
		// NOTE: 0x1f is "US" (unit separator). Many terminals send this for Ctrl+/ as well
		// as Ctrl+_. We prefer reporting Ctrl+/ since it's the more common shortcut.
		0x1F: {Key: types.KeyRunes, Rune: '/', Ctrl: true}, // Ctrl+/ (ambiguous with Ctrl+_)
		0x7F: {Key: types.KeyBackspace},                    // Backspace (DEL)
	}

	for b, key := range ctrlMappings {
		s[string([]byte{b})] = &key
		altKey := key
		altKey.Alt = true
		s[string([]byte{'\x1b', b})] = &altKey
	}

	// Add space
	s[" "] = &types.EventKey{Key: types.KeySpace, Rune: ' '}
	s["\x1b "] = &types.EventKey{Key: types.KeySpace, Rune: ' ', Alt: true}

	// Double ESC is Alt+Esc
	s["\x1b\x1b"] = &types.EventKey{Key: types.KeyEsc, Alt: true}

	return s
}()

// seqLengths is sorted sequence lengths (largest first)
var seqLengths = func() []int {
	sizes := make(map[int]struct{})
	for seq := range extSequences {
		sizes[len(seq)] = struct{}{}
	}
	lsizes := make([]int, 0, len(sizes))
	for sz := range sizes {
		lsizes = append(lsizes, sz)
	}
	sort.Slice(lsizes, func(i, j int) bool { return lsizes[i] > lsizes[j] })
	return lsizes
}()

func parseKeyEvent(buf []byte) (*types.EventKey, int, bool) {
	if len(buf) == 0 {
		return nil, 0, false
	}

	// Escape-prefixed sequences (Alt, CSI, SS3, OSC, ...).
	// We treat a lone ESC as incomplete to disambiguate between ESC vs Alt+<key>.
	if buf[0] == '\x1b' {
		if len(buf) == 1 {
			return nil, 0, false
		}

		n, complete := escapeSequenceLen(buf)
		if !complete {
			return nil, 0, false
		}
		seq := buf[:n]

		// OSC/DCS/APC are terminal control/response sequences, not keypresses.
		switch buf[1] {
		case ']', 'P', '_':
			return nil, n, true
		}

		// Kitty keyboard protocol: CSI ... u
		if key, ok := parseKittyKeyboard(seq); ok {
			return key, n, true
		}

		// modifyOtherKeys: CSI 27;mod;code~
		if key, ok := parseModifyOtherKeys(seq); ok {
			return key, n, true
		}

		// Known escape sequences
		if key, found := extSequences[string(seq)]; found {
			return key, n, true
		}

		// Terminal responses (DSR/DA/window reports, etc.) should be ignored.
		if isTerminalResponse(seq) {
			return nil, n, true
		}

		// Alt+<control> combos (ESC + control byte)
		if len(seq) == 2 && (seq[1] <= 0x1a) {
			r := rune(seq[1]) + 'a' - 1
			return &types.EventKey{Key: types.KeyRunes, Rune: r, Alt: true, Ctrl: true}, 2, true
		}

		// Alt+<rune> (ESC + UTF-8)
		if len(seq) >= 2 && seq[1] != '[' && seq[1] != 'O' {
			if !utf8.FullRune(seq[1:]) {
				return nil, 0, false
			}
			r, size := utf8.DecodeRune(seq[1:])
			if r != utf8.RuneError && r != 0x7f {
				return &types.EventKey{Key: types.KeyRunes, Rune: r, Alt: true}, 1 + size, true
			}
		}

		return &types.EventKey{Key: types.KeyUnknown}, n, true
	}

	// Single-byte mappings (control bytes, space, etc.)
	if key, found := extSequences[string([]byte{buf[0]})]; found {
		return key, 1, true
	}

	// 8-bit encoded Alt+char (high bit set). Some terminals send Alt+a as 0xe1 instead of ESC a.
	if buf[0] >= 0x80 {
		// Check if clearing high bit gives us a printable character
		c := buf[0] & 0x7f
		if c > ' ' && c < 0x7f {
			r, _ := utf8.DecodeRune([]byte{c})
			if r != utf8.RuneError {
				return &types.EventKey{Key: types.KeyRunes, Rune: r, Alt: true}, 1, true
			}
		}
	}

	// UTF-8 runes
	if !utf8.FullRune(buf) {
		return nil, 0, false
	}
	r, size := utf8.DecodeRune(buf)
	if r == utf8.RuneError {
		return &types.EventKey{Key: types.KeyUnknown}, 1, true
	}
	return &types.EventKey{Key: types.KeyRunes, Rune: r}, size, true
}

func escapeSequenceLen(buf []byte) (int, bool) {
	if len(buf) < 2 || buf[0] != '\x1b' {
		return 0, false
	}

	switch buf[1] {
	case '[':
		// X10 mouse: ESC [ M Cb Cx Cy
		if len(buf) >= 3 && buf[2] == 'M' {
			if len(buf) >= 6 {
				return 6, true
			}
			return 0, false
		}
		// CSI ends with a final byte in 0x40-0x7E.
		for i := 2; i < len(buf); i++ {
			b := buf[i]
			if b >= 0x40 && b <= 0x7e {
				return i + 1, true
			}
		}
		return 0, false
	case ']':
		// OSC ends with BEL or ST (ESC \)
		for i := 2; i < len(buf); i++ {
			if buf[i] == 0x07 {
				return i + 1, true
			}
			if buf[i] == '\x1b' && i+1 < len(buf) && buf[i+1] == '\\' {
				return i + 2, true
			}
		}
		return 0, false
	case 'P', '_':
		// DCS/APC ends with ST (ESC \)
		for i := 2; i < len(buf); i++ {
			if buf[i] == '\x1b' && i+1 < len(buf) && buf[i+1] == '\\' {
				return i + 2, true
			}
		}
		return 0, false
	case 'O':
		// SS3: ESC O <final>
		if len(buf) >= 3 {
			return 3, true
		}
		return 0, false
	default:
		// Meta: ESC + a single (possibly multi-byte) rune
		if !utf8.FullRune(buf[1:]) {
			return 0, false
		}
		_, size := utf8.DecodeRune(buf[1:])
		if size <= 0 {
			return 0, false
		}
		return 1 + size, true
	}
}

func parseModifyOtherKeys(seq []byte) (*types.EventKey, bool) {
	// CSI 27;modifier;code~
	if len(seq) < 7 || seq[0] != '\x1b' || seq[1] != '[' || seq[len(seq)-1] != '~' {
		return nil, false
	}
	params := seq[2 : len(seq)-1]
	if !bytes.HasPrefix(params, []byte("27;")) {
		return nil, false
	}
	params = params[3:]
	parts := bytes.Split(params, []byte{';'})
	if len(parts) != 2 {
		return nil, false
	}
	mod, ok := parseDecimal(parts[0])
	if !ok {
		return nil, false
	}
	code, ok := parseDecimal(parts[1])
	if !ok {
		return nil, false
	}
	mod-- // xterm encodes modifiers starting at 1
	shift := mod&1 != 0
	alt := mod&2 != 0
	ctrl := mod&4 != 0

	switch code {
	case 13:
		return &types.EventKey{Key: types.KeyEnter, Alt: alt, Ctrl: ctrl, Shift: shift}, true
	case 27:
		return &types.EventKey{Key: types.KeyEsc, Alt: alt, Ctrl: ctrl, Shift: shift}, true
	case 9:
		if shift && !alt && !ctrl {
			return &types.EventKey{Key: types.KeyShiftTab}, true
		}
		return &types.EventKey{Key: types.KeyTab, Alt: alt, Ctrl: ctrl, Shift: shift}, true
	case 32:
		return &types.EventKey{Key: types.KeySpace, Rune: ' ', Alt: alt, Ctrl: ctrl}, true
	case 8, 127:
		return &types.EventKey{Key: types.KeyBackspace, Alt: alt, Ctrl: ctrl, Shift: shift}, true
	default:
		if code <= 0 || code > utf8.MaxRune {
			return nil, false
		}
		r := rune(code)
		// For runes we keep the shifted codepoint; Shift is intentionally not set to match legacy behavior.
		return &types.EventKey{Key: types.KeyRunes, Rune: r, Alt: alt, Ctrl: ctrl}, true
	}
}

func parseKittyKeyboard(seq []byte) (*types.EventKey, bool) {
	if len(seq) < 4 || seq[0] != '\x1b' || seq[1] != '[' {
		return nil, false
	}

	// Special key form with event types: CSI <n>;<mods>:<etype><final>
	if key, ok := parseKittySpecialKey(seq); ok {
		return key, true
	}

	// Main Kitty form: CSI ... u
	if seq[len(seq)-1] != 'u' {
		return nil, false
	}
	params := seq[2 : len(seq)-1]
	fields := bytes.Split(params, []byte{';'})
	if len(fields) < 1 {
		return nil, false
	}

	// field 1: codepoint[:shifted[:base]]
	f1 := bytes.Split(fields[0], []byte{':'})
	code, ok := parseDecimal(f1[0])
	if !ok {
		return nil, false
	}
	shifted := 0
	if len(f1) >= 2 {
		if v, ok := parseDecimal(f1[1]); ok {
			shifted = v
		}
	}

	alt, ctrl, shift := false, false, false
	// Event type: 1=press, 2=repeat, 3=release.
	// We ignore releases to avoid double events for normal typing.
	eventType := 1
	if len(fields) >= 2 {
		modsPart := fields[1]
		modsField := modsPart
		if i := bytes.IndexByte(modsPart, ':'); i != -1 {
			modsField = modsPart[:i]
			// Parse optional event type (e.g. ";1:3" => release).
			if et, ok := parseDecimal(modsPart[i+1:]); ok {
				eventType = et
			}
		}
		if modsField != nil && len(modsField) > 0 {
			m, ok := parseDecimal(modsField)
			if ok && m > 0 {
				m-- // kitty encodes modifiers starting at 1
				shift = m&1 != 0
				alt = m&2 != 0
				ctrl = m&4 != 0
			}
		}
	}

	if eventType == 3 {
		return nil, true
	}

	if base, ok := kittyKeyBase(code); ok {
		return normalizeSpecialKey(base, alt, ctrl, shift), true
	}

	if code <= 0 || code > utf8.MaxRune {
		return nil, false
	}
	r := rune(code)
	if shift && shifted > 0 && shifted <= utf8.MaxRune {
		r = rune(shifted)
	}
	return &types.EventKey{Key: types.KeyRunes, Rune: r, Alt: alt, Ctrl: ctrl}, true
}

func parseKittySpecialKey(seq []byte) (*types.EventKey, bool) {
	// CSI <n>;<mods>:<etype><final>
	if len(seq) < 6 || seq[0] != '\x1b' || seq[1] != '[' {
		return nil, false
	}
	final := seq[len(seq)-1]
	if !(final == '~' || (final >= 'A' && final <= 'Z')) {
		return nil, false
	}

	body := seq[2 : len(seq)-1]
	semi := bytes.IndexByte(body, ';')
	if semi == -1 {
		return nil, false
	}
	n, ok := parseDecimal(body[:semi])
	if !ok {
		return nil, false
	}
	modsAndType := body[semi+1:]
	colon := bytes.IndexByte(modsAndType, ':')
	if colon == -1 {
		return nil, false
	}
	mods, ok := parseDecimal(modsAndType[:colon])
	if !ok {
		return nil, false
	}
	etype, ok := parseDecimal(modsAndType[colon+1:])
	if !ok {
		return nil, false
	}
	if etype == 3 {
		return nil, true
	}
	mods--
	shift := mods&1 != 0
	alt := mods&2 != 0
	ctrl := mods&4 != 0

	var base types.KeyType
	switch final {
	case '~':
		base, ok = kittyTildeBase(n)
		if !ok {
			return nil, false
		}
	default:
		base, ok = kittyFunctionalBase(final)
		if !ok {
			return nil, false
		}
	}
	return normalizeSpecialKey(base, alt, ctrl, shift), true
}

func normalizeSpecialKey(base types.KeyType, alt, ctrl, shift bool) *types.EventKey {
	// Keys where ctrl/shift are encoded in KeyType names.
	switch base {
	case types.KeyUp, types.KeyDown, types.KeyLeft, types.KeyRight:
		return &types.EventKey{Key: directionalKey(base, ctrl, shift), Alt: alt}
	case types.KeyHome, types.KeyEnd:
		return &types.EventKey{Key: homeEndKey(base, ctrl, shift), Alt: alt}
	case types.KeyPgUp:
		if ctrl {
			return &types.EventKey{Key: types.KeyCtrlPgUp, Alt: alt}
		}
		return &types.EventKey{Key: types.KeyPgUp, Alt: alt}
	case types.KeyPgDown:
		if ctrl {
			return &types.EventKey{Key: types.KeyCtrlPgDown, Alt: alt}
		}
		return &types.EventKey{Key: types.KeyPgDown, Alt: alt}
	case types.KeyTab:
		if shift && !ctrl {
			return &types.EventKey{Key: types.KeyShiftTab, Alt: alt}
		}
		return &types.EventKey{Key: types.KeyTab, Alt: alt, Ctrl: ctrl}
	default:
		return &types.EventKey{Key: base, Alt: alt, Ctrl: ctrl, Shift: shift}
	}
}

func directionalKey(base types.KeyType, ctrl, shift bool) types.KeyType {
	if ctrl && shift {
		switch base {
		case types.KeyUp:
			return types.KeyCtrlShiftUp
		case types.KeyDown:
			return types.KeyCtrlShiftDown
		case types.KeyLeft:
			return types.KeyCtrlShiftLeft
		case types.KeyRight:
			return types.KeyCtrlShiftRight
		}
	}
	if ctrl {
		switch base {
		case types.KeyUp:
			return types.KeyCtrlUp
		case types.KeyDown:
			return types.KeyCtrlDown
		case types.KeyLeft:
			return types.KeyCtrlLeft
		case types.KeyRight:
			return types.KeyCtrlRight
		}
	}
	if shift {
		switch base {
		case types.KeyUp:
			return types.KeyShiftUp
		case types.KeyDown:
			return types.KeyShiftDown
		case types.KeyLeft:
			return types.KeyShiftLeft
		case types.KeyRight:
			return types.KeyShiftRight
		}
	}
	return base
}

func homeEndKey(base types.KeyType, ctrl, shift bool) types.KeyType {
	if ctrl && shift {
		if base == types.KeyHome {
			return types.KeyCtrlShiftHome
		}
		return types.KeyCtrlShiftEnd
	}
	if ctrl {
		if base == types.KeyHome {
			return types.KeyCtrlHome
		}
		return types.KeyCtrlEnd
	}
	if shift {
		if base == types.KeyHome {
			return types.KeyShiftHome
		}
		return types.KeyShiftEnd
	}
	return base
}

func kittyKeyBase(code int) (types.KeyType, bool) {
	switch code {
	case 27, 57344:
		return types.KeyEsc, true
	case 13, 57345:
		return types.KeyEnter, true
	case 9, 57346:
		return types.KeyTab, true
	case 127, 57347:
		return types.KeyBackspace, true
	case 57348:
		return types.KeyInsert, true
	case 57349:
		return types.KeyDelete, true
	case 57350:
		return types.KeyLeft, true
	case 57351:
		return types.KeyRight, true
	case 57352:
		return types.KeyUp, true
	case 57353:
		return types.KeyDown, true
	case 57354:
		return types.KeyPgUp, true
	case 57355:
		return types.KeyPgDown, true
	case 57356:
		return types.KeyHome, true
	case 57357:
		return types.KeyEnd, true
	case 57364:
		return types.KeyF1, true
	case 57365:
		return types.KeyF2, true
	case 57366:
		return types.KeyF3, true
	case 57367:
		return types.KeyF4, true
	case 57368:
		return types.KeyF5, true
	case 57369:
		return types.KeyF6, true
	case 57370:
		return types.KeyF7, true
	case 57371:
		return types.KeyF8, true
	case 57372:
		return types.KeyF9, true
	case 57373:
		return types.KeyF10, true
	case 57374:
		return types.KeyF11, true
	case 57375:
		return types.KeyF12, true
	case 57376:
		return types.KeyF13, true
	case 57377:
		return types.KeyF14, true
	case 57378:
		return types.KeyF15, true
	case 57379:
		return types.KeyF16, true
	case 57380:
		return types.KeyF17, true
	case 57381:
		return types.KeyF18, true
	case 57382:
		return types.KeyF19, true
	case 57383:
		return types.KeyF20, true
	default:
		return 0, false
	}
}

func kittyFunctionalBase(final byte) (types.KeyType, bool) {
	switch final {
	case 'A':
		return types.KeyUp, true
	case 'B':
		return types.KeyDown, true
	case 'C':
		return types.KeyRight, true
	case 'D':
		return types.KeyLeft, true
	case 'H':
		return types.KeyHome, true
	case 'F':
		return types.KeyEnd, true
	case 'P':
		return types.KeyF1, true
	case 'Q':
		return types.KeyF2, true
	case 'R':
		return types.KeyF3, true
	case 'S':
		return types.KeyF4, true
	default:
		return 0, false
	}
}

func kittyTildeBase(n int) (types.KeyType, bool) {
	switch n {
	case 1, 7:
		return types.KeyHome, true
	case 4, 8:
		return types.KeyEnd, true
	case 2:
		return types.KeyInsert, true
	case 3:
		return types.KeyDelete, true
	case 5:
		return types.KeyPgUp, true
	case 6:
		return types.KeyPgDown, true
	case 11:
		return types.KeyF1, true
	case 12:
		return types.KeyF2, true
	case 13:
		return types.KeyF3, true
	case 14:
		return types.KeyF4, true
	case 15:
		return types.KeyF5, true
	case 17:
		return types.KeyF6, true
	case 18:
		return types.KeyF7, true
	case 19:
		return types.KeyF8, true
	case 20:
		return types.KeyF9, true
	case 21:
		return types.KeyF10, true
	case 23:
		return types.KeyF11, true
	case 24:
		return types.KeyF12, true
	default:
		return 0, false
	}
}

func parseDecimal(b []byte) (int, bool) {
	if len(b) == 0 {
		return 0, false
	}
	n := 0
	for _, c := range b {
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	return n, true
}

func isTerminalResponse(seq []byte) bool {
	if len(seq) < 2 || seq[0] != '\x1b' {
		return false
	}
	// Focus tracking: ESC[I / ESC[O
	if bytes.Equal(seq, []byte("\x1b[I")) || bytes.Equal(seq, []byte("\x1b[O")) {
		return true
	}
	if seq[1] == ']' || seq[1] == 'P' || seq[1] == '_' {
		return true
	}
	if seq[1] != '[' {
		return false
	}
	// DSR cursor position: ESC[row;colR
	if seq[len(seq)-1] == 'R' {
		return true
	}
	// Window/cell size reports: ESC[4;height;widtht or ESC[8;rows;colst
	if seq[len(seq)-1] == 't' {
		return true
	}
	// Device attributes: ESC[?...c
	if seq[len(seq)-1] == 'c' && len(seq) > 3 && seq[2] == '?' {
		return true
	}
	// Mode reports: ESC[?...$y
	if len(seq) >= 4 && seq[len(seq)-2] == '$' && seq[len(seq)-1] == 'y' && seq[2] == '?' {
		return true
	}
	return false
}
