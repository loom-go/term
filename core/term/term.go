package term

import (
	"fmt"
	"os"

	"github.com/loom-go/term/core/stdio"

	"golang.org/x/term"
)

type State = term.State

func MakeRaw() (*State, error) {
	fd := int(os.Stdin.Fd())

	state, err := term.MakeRaw(fd)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToSetRawMode, err)
	}

	return state, nil
}

func Restore(state *State) error {
	fd := int(os.Stdin.Fd())

	err := term.Restore(fd, state)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToRestoreRawMode, err)
	}

	return nil
}

func Size() (width, height int, err error) {
	fd := int(os.Stdout.Fd())

	width, height, err = term.GetSize(fd)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrFailedToGetTerminalSize, err)
	}

	return width, height, nil
}

func CursorPos() (row, col int, err error) {
	state, err := MakeRaw()
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrFailedToGetCursorPosition, err)
	}
	defer Restore(state)

	stdin := stdio.Stdin().Listen(256)

	// query cursor pos
	fmt.Print("\x1b[6n")

	for bytes := range stdin {
		if len(bytes) == 0 {
			continue
		}

		n, err := fmt.Sscanf(string(bytes), "\x1b[%d;%dR", &row, &col)
		if err != nil || n != 2 {
			return 0, 0, fmt.Errorf("%w: %w", ErrFailedToGetCursorPosition, err)
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
