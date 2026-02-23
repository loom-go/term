package term

import "errors"

var (
	ErrFailedToSetRawMode        = errors.New("failed to set terminal to raw mode")
	ErrFailedToRestoreRawMode    = errors.New("failed to restore terminal raw mode")
	ErrFailedToGetTerminalSize   = errors.New("failed to get terminal size")
	ErrFailedToGetCursorPosition = errors.New("failed to get cursor position")
)
