package errs

import "errors"

var (
	ErrAppNotInitialized     = errors.New("it seems the app is not initialized properly. Are you trying to use a component outside of a LoomTerm app?")
	ErrAppFailedToInitialize = errors.New("failed to initialize LoomTerm app")
)
