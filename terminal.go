package term

import "github.com/AnatoleLucet/loom-term/core"

// TerminalSize returns the current size of the terminal (rows, cols).
// If unable to get the size, it returns (0, 0).
func TerminalSize() (width, height int) {
	return core.TerminalSize()
}

// CursorPosition returns the current position of the terminal cursor.
// If unable to get the position, it returns (0, 0).
func CursorPosition() (row, col int) {
	return core.CursorPosition()
}

// ScrollUp scrolls the terminal content up by the specified number of lines.
func ScrollUp(lines int) {
	core.ScrollUp(lines)
}

// ScrollDown scrolls the terminal content down by the specified number of lines.
func ScrollDown(lines int) {
	core.ScrollDown(lines)
}
