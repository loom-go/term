package core

import (
	"github.com/AnatoleLucet/loom-term/core/debug"
	"github.com/AnatoleLucet/loom-term/core/elements"
	"github.com/AnatoleLucet/loom-term/core/runtime"
	"github.com/AnatoleLucet/loom-term/core/terminal"
	"github.com/AnatoleLucet/loom-term/core/types"
)

type Element = types.Element
type Runtime = types.Runtime
type RenderType = types.RenderType
type RenderContext = types.RenderContext

const (
	RenderInline     RenderType = types.RenderTypeInline
	RenderFullscreen RenderType = types.RenderTypeFullscreen
)

func NewRuntime(typ RenderType) (Runtime, error) {
	return runtime.NewRuntime(typ)
}

func NewElement(ctx RenderContext) (Element, error) {
	return elements.NewElement(ctx)
}

func NewTextElement(ctx RenderContext) (*elements.TextElement, error) {
	return elements.NewTextElement(ctx)
}

func NewBoxElement(ctx RenderContext) (*elements.BoxElement, error) {
	return elements.NewBoxElement(ctx)
}

func NewScrollBoxElement(ctx RenderContext) (*elements.ScrollBoxElement, error) {
	return elements.NewScrollBoxElement(ctx)
}

func NewConsoleElement(ctx RenderContext) (*elements.ConsoleElement, error) {
	return elements.NewConsoleElement(ctx)
}

func TerminalSize() (width, height int) {
	width, height, err := terminal.Size()
	if err != nil {
		return 0, 0
	}

	return width, height
}

func CursorPosition() (row, col int) {
	row, col, err := terminal.CursorPos()
	if err != nil {
		return 0, 0
	}

	return row, col
}

func ScrollUp(lines int) {
	terminal.ScrollUp(lines)
}

func ScrollDown(lines int) {
	terminal.ScrollDown(lines)
}

func LogDebug(msg string) {
	debug.LogDebug(msg)
}

func LogInfo(msg string) {
	debug.LogInfo(msg)
}

func LogWarning(msg string) {
	debug.LogWarning(msg)
}

func LogError(msg string) {
	debug.LogError(msg)
}
