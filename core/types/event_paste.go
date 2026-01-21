package types

import "fmt"

type EventPaste struct {
	BaseEvent
	Text string
}

func (e EventPaste) String() string {
	return fmt.Sprintf("Paste(%q)", e.Text)
}
