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
	i.scheduleUpdate(func() error {
		return guardDestroyed(i.ctx)
	})
}
func (i *InputElement) SetWrap(wrap any) {
	i.scheduleUpdate(func() error {
		return guardDestroyed(i.ctx)
	})
}

func (i *InputElement) SetText(text string) {
	i.TextAreaElement.SetText(i.sanitize(text))
}

func (i *InputElement) InsertText(text string) {
	i.TextAreaElement.InsertText(i.sanitize(text))
}

func (i *InputElement) handleKeyPress(event *EventKey) {
	key := event.Key.String()
	if key == "enter" {
		i.Submit()
		i.rdrctx.ScheduleRender()
		return
	}

	i.TextAreaElement.handleKeyPress(event)
}

func (i *InputElement) handlePaste(event *EventPaste) {
	evt := &EventPaste{EventPaste: events.EventPaste{
		Text: i.sanitize(event.Text),
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
