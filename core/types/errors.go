package types

import "errors"

var (
	ErrUsingClosedRuntime = errors.New("using closed runtime")

	ErrFailedToRenderFrame = errors.New("failed to render frame")

	ErrRendererFailedToInitialize              = errors.New("failed to initialize renderer")
	ErrRendererFailedToGetBuffer               = errors.New("failed to get render buffer from renderer")
	ErrRendererFailedToGetBufferDimensions     = errors.New("failed to get render buffer dimensions from renderer")
	ErrRendererFailedToProcessTermCapabilities = errors.New("failed to process terminal capabilities from renderer")
	ErrRendererFailedToGetTermCapabilities     = errors.New("failed to get terminal capabilities from renderer")

	ErrUpdatingDestroyedElement = errors.New("cannot update destroyed element")
	ErrPaintingDestroyedElement = errors.New("cannot paint destroyed element")

	ErrFailedToGetTerminalSize   = errors.New("failed to get terminal size")
	ErrFailedToGetCursorPosition = errors.New("failed to get cursor position")

	ErrInvalidStyleValue = errors.New("invalid style value")
)
