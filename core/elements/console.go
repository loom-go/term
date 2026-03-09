package elements

import (
	"fmt"
	"math"

	"github.com/loom-go/term/core/debug"
)

const consoleBG = "#111111dd"
const consoleLogsHeaderBG = "#00000055"
const consoleLogsScrollFactorY = 2

var consoleLogsLevelColors = map[debug.LogLevel]string{
	debug.LogLevelDebug:   "#6b7280",
	debug.LogLevelInfo:    "#3b82f6",
	debug.LogLevelWarning: "#f59e0b",
	debug.LogLevelError:   "#ef4444",
}

type ConsoleElement struct {
	*BaseElement
}

func NewConsoleElement() (e *ConsoleElement, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Console: %w: %v", ErrFailedToInitializeElement, err)
		}
	}()

	base, err := NewBaseElement()
	if err != nil {
		return nil, err
	}

	e = &ConsoleElement{
		BaseElement: base,
	}
	base.self = e
	e.SetZIndex(math.MaxInt)

	stats, err := newStatsElement()
	if err != nil {
		return nil, err
	}
	e.AppendChild(stats)

	logs, err := newLogsElement()
	if err != nil {
		return nil, err
	}
	e.AppendChild(logs)

	return e, nil
}
