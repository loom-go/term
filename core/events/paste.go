package events

import (
	"bytes"

	"github.com/AnatoleLucet/loom-term/core/types"
)

// Bracketed paste markers
const (
	pasteStart = "\x1b[200~"
	pasteEnd   = "\x1b[201~"
)

func parseBracketedPaste(buf []byte) (*types.EventPaste, int, bool) {
	if !bytes.HasPrefix(buf, []byte(pasteStart)) {
		return nil, 0, false
	}

	// find the end marker
	endIdx := bytes.Index(buf[len(pasteStart):], []byte(pasteEnd))
	if endIdx == -1 {
		// need more data
		return nil, 0, false
	}

	start := len(pasteStart)
	end := start + endIdx
	pastedText := string(buf[start:end])
	consumed := end + len(pasteEnd)

	return &types.EventPaste{Text: pastedText}, consumed, true
}
