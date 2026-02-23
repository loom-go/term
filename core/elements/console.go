package elements

import (
	"fmt"
	"math"

	"github.com/AnatoleLucet/loom-term/core/debug"
	"github.com/AnatoleLucet/loom-term/core/term"
)

const consoleBG = "#111111dd"
const consoleLogsHeaderBG = "#00000055"

var consoleLogsLevelColors = map[debug.LogLevel]string{
	debug.LogLevelDebug:   "#6B7280",
	debug.LogLevelInfo:    "#3B82F6",
	debug.LogLevelWarning: "#F59E0B",
	debug.LogLevelError:   "#EF4444",
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

	_, height, err := term.Size()
	if err != nil {
		return nil, err
	}

	if height < 20 {
		return nil, fmt.Errorf("not enough vertical space to render console (height: %d)", height)
	}

	base, err := NewBaseElement()
	if err != nil {
		return nil, err
	}

	e = &ConsoleElement{
		BaseElement: base,
	}
	base.self = e
	e.SetZIndex(math.MaxInt)
	e.SetMinHeight(height / 2)

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
