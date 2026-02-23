package elements

import "errors"

var (
	ErrPanicDuringRender     = errors.New("panic during render")
	ErrFailedToRenderFrame   = errors.New("failed to render frame")
	ErrFailedToComputeLayout = errors.New("failed to compute layout")
	ErrFailedToRecordFrame   = errors.New("failed to record frame")
	ErrFailedToDrawFrame     = errors.New("failed to draw frame")

	ErrFailedToInitializeRoot    = errors.New("failed to initialize root element")
	ErrFailedToInitializeElement = errors.New("failed to initialize element")

	ErrFailedToGetBuffer                   = errors.New("failed to get render buffer")
	ErrFailedToGetTerminalCapabilities     = errors.New("failed to get terminal capabilities from renderer")
	ErrFailedToProcessTerminalCapabilities = errors.New("failed to process terminal capabilities from renderer")

	ErrUsingDestroyedElement = errors.New("cannot use an element that has been destroyed")

	ErrFailedToRemoveChild = errors.New("failed to remove child element")
	ErrFailedToAppendChild = errors.New("failed to append child element")
	ErrFailedToUpdateChild = errors.New("failed to update child element")

	ErrInvalidStyleValue    = errors.New("invalid style value")
	ErrFailedToUpdateZIndex = errors.New("failed to update z-index")
)
