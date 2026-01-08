package main

import "github.com/AnatoleLucet/loom-term/internal"

// TerminalSize returns the current size of the terminal (rows, cols).
// If unable to get the size, it returns (0, 0).
func TerminalSize() (rows, cols int) {
	rows, cols, err := internal.TerminalSize()
	if err != nil {
		return 0, 0
	}

	return rows, cols
}

// CursorPosition returns the current position of the terminal cursor.
// If unable to get the position, it returns (0, 0).
func CursorPosition() (row, col int) {
	row, col, err := internal.CursorPos()
	if err != nil {
		return 0, 0
	}

	return row, col
}

// ScrollUp scrolls the terminal content up by the specified number of lines.
func ScrollUp(lines int) {
	internal.ScrollUp(lines)
}

// ScrollDown scrolls the terminal content down by the specified number of lines.
func ScrollDown(lines int) {
	internal.ScrollDown(lines)
}
