//go:build windows

package term

import "golang.org/x/sys/windows"

func Init() (func() error, error) {
	prevCP, err := windows.GetConsoleCP()
	prevOutputCP, err := windows.GetConsoleOutputCP()
	if err != nil {
		return nil, err
	}

	// set console code page as utf-8
	windows.SetConsoleCP(65001)
	windows.SetConsoleOutputCP(65001)

	prevState, err := MakeRaw()
	if err != nil {
		return nil, err
	}

	return func() error {
		windows.SetConsoleCP(prevCP)
		windows.SetConsoleOutputCP(prevOutputCP)
		return Restore(prevState)
	}, nil
}
