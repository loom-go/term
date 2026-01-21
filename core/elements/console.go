package elements

import (
	"math"

	"github.com/AnatoleLucet/loom-term/core/types"
)

var consoleBG = "#111111dd"
var consoleLogsHeaderBG = "#00000055"
var consoleLogsHeight = "35%"

type ConsoleElement struct {
	*BaseElement
}

func NewConsoleElement(ctx types.RenderContext) (*ConsoleElement, error) {
	base, err := NewElement(ctx)
	if err != nil {
		return nil, err
	}

	e := &ConsoleElement{
		BaseElement: base,
	}
	e.SetZIndex(math.MaxInt)

	stats, err := newStatsElement(ctx)
	if err != nil {
		return nil, err
	}
	e.AppendChild(stats)

	logs, err := newLogsElement(ctx)
	if err != nil {
		return nil, err
	}
	e.AppendChild(logs)

	return e, nil
}
