package core

import "github.com/loom-go/term/core/term"

// TerminalSize returns the current size of the terminal (rows, cols).
// If unable to get the size, it returns (0, 0).
func TerminalSize() (width, height int) {
	width, height, err := term.Size()
	if err != nil {
		return 0, 0
	}

	return width, height
}

// CursorPosition returns the current position of the terminal cursor.
// If unable to get the position, it returns (0, 0).
func CursorPosition() (row, col int) {
	row, col, err := term.CursorPos()
	if err != nil {
		return 0, 0
	}

	return row, col
}

// ScrollUp scrolls the terminal content up by the specified number of lines.
func ScrollUp(lines int) {
	term.ScrollUp(lines)
}

// ScrollDown scrolls the terminal content down by the specified number of lines.
func ScrollDown(lines int) {
	term.ScrollDown(lines)
}
