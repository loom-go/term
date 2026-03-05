package elements

import (
	"fmt"
	"strings"

	"github.com/AnatoleLucet/loom-term/core/elements/events"
)

type InputElement struct {
	*TextAreaElement
}

func NewInputElement() (input *InputElement, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Input: %w: %w", ErrFailedToInitializeElement, err)
		}
	}()

	ta, err := NewTextAreaElement()
	if err != nil {
		return nil, err
	}

	input = &InputElement{
		TextAreaElement: ta,
	}
	ta.self = input

	ta.SetHeight(1)
	ta.SetWrap("none")

	return input, nil
}

func (i *InputElement) SetHeight(height any) {
	scheduleUpdate(i.Self(), func() error { return nil })
}
func (i *InputElement) SetWrap(wrap any) {
	scheduleUpdate(i.Self(), func() error { return nil })
}

func (i *InputElement) SetValue(text string) {
	i.TextAreaElement.SetValue(i.sanitize(text))
}

func (i *InputElement) InsertValue(text string) {
	i.TextAreaElement.InsertValue(i.sanitize(text))
}

func (i *InputElement) handleKeyPress(event *EventKey) {
	key := event.Key.String()
	if key == "enter" {
		i.Submit()
		return
	}

	i.TextAreaElement.handleKeyPress(event)
}

func (i *InputElement) handlePaste(event *EventPaste) {
	evt := &EventPaste{EventPaste: events.EventPaste{
		Value: i.sanitize(event.Value),
	}}
	evt.setTarget(event.Target())

	i.TextAreaElement.handlePaste(evt)
}

func (i *InputElement) sanitize(text string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' {
			return -1
		}

		return r
	}, text)
}
