package events

import "fmt"

type EventPaste struct {
	Value string
}

func (e EventPaste) String() string {
	return fmt.Sprintf("Paste(%q)", e.Value)
}
