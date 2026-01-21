package types

import "fmt"

type EventResize struct {
	BaseEvent
	Height int
	Width  int
}

func (e EventResize) String() string {
	return fmt.Sprintf("Resize(%dx%d)", e.Width, e.Height)
}
