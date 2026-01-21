package terminal

import (
	"fmt"
	"os"

	"github.com/AnatoleLucet/loom-term/core/stdio"
	"github.com/AnatoleLucet/loom-term/core/types"
	"golang.org/x/term"
)

func Size() (width, height int, err error) {
	fd := int(os.Stdout.Fd())

	width, height, err = term.GetSize(fd)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: failed to get terminal size: %w", types.ErrFailedToGetTerminalSize, err)
	}

	return width, height, nil
}

func CursorPos() (row, col int, err error) {
	fd := int(os.Stdin.Fd())

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: failed to set terminal to raw mode: %w", types.ErrFailedToGetCursorPosition, err)
	}
	defer term.Restore(fd, oldState)

	stdin := stdio.Stdin.Listen(256)

	// query cursor pos
	fmt.Print("\x1b[6n")

	for bytes := range stdin {
		if len(bytes) == 0 {
			continue
		}

		n, err := fmt.Sscanf(string(bytes), "\x1b[%d;%dR", &row, &col)
		if err != nil || n != 2 {
			return 0, 0, fmt.Errorf("%w: failed to parse cursor position response: %w", types.ErrFailedToGetCursorPosition, err)
		}

		break
	}

	return row, col, nil
}

func ScrollUp(lines int) {
	if lines > 0 {
		fmt.Printf("\x1b[%dS", lines)
	}
}

func ScrollDown(lines int) {
	if lines > 0 {
		fmt.Printf("\x1b[%dT", lines)
	}
}
