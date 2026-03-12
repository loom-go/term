//go:build !windows

package term

func Init() (func() error, error) {
	state, err := MakeRaw()
	if err != nil {
		return nil, err
	}

	return func() error {
		return Restore(state)
	}, nil
}
