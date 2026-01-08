package termerror

import "errors"

var (
	ErrAlreadyInFullscreenRenderer = errors.New("already in a fullscreen renderer. You probably have nested FullscreenRenderer/InlineRenderer components")
	ErrAlreadyInInlineRenderer     = errors.New("already in an inline renderer. You probably have nested FullscreenRenderer/InlineRenderer components")
	ErrNoRendererInContext         = errors.New("failed to find renderer. You must use FullscreenRenderer or InlineRenderer as a root component")
	ErrFailedToInitializeRenderer  = errors.New("failed to initialize renderer")
	ErrFailedToGetBuffer           = errors.New("failed to get buffer from renderer")
	ErrFailedToCreateRootNode      = errors.New("failed to create root terminal node")

	ErrFailedToGetTerminalSize   = errors.New("failed to get terminal size")
	ErrFailedToGetCursorPosition = errors.New("failed to get cursor position")

	ErrInvalidStyleValue = errors.New("invalid style value")
)
