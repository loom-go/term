package events

import "fmt"

type EventResize struct {
	Height int
	Width  int
}

func (e EventResize) String() string {
	return fmt.Sprintf("Resize(%dx%d)", e.Width, e.Height)
}
