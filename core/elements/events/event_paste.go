package events

import "fmt"

type EventPaste struct {
	Text string
}

func (e EventPaste) String() string {
	return fmt.Sprintf("Paste(%q)", e.Text)
}
