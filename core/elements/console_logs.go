package elements

import (
	"github.com/AnatoleLucet/loom-term/core/debug"
	"github.com/AnatoleLucet/loom-term/core/types"
)

type logsElement struct {
	*BoxElement

	cancel func()
}

func newLogsElement(ctx types.RenderContext) (*logsElement, error) {
	box, err := NewBoxElement(ctx)
	if err != nil {
		return nil, err
	}

	e := &logsElement{
		BoxElement: box,
	}
	e.SetWidth("100%")
	e.SetHeight(consoleLogsHeight)
	e.SetLeft(0)
	e.SetBottom(0)
	e.SetPosition("absolute")
	e.SetFlexDirection("column")
	e.SetBackgroundColor(consoleBG)

	header, err := newLogsHeaderElement(ctx)
	if err != nil {
		return nil, err
	}

	list, err := newLogsListElement(ctx)
	if err != nil {
		return nil, err
	}

	e.AppendChild(header)
	e.AppendChild(list)

	return e, nil
}

type logsHeaderElement struct {
	*BoxElement
}

func newLogsHeaderElement(ctx types.RenderContext) (*logsHeaderElement, error) {
	box, err := NewBoxElement(ctx)
	if err != nil {
		return nil, err
	}

	e := &logsHeaderElement{
		BoxElement: box,
	}
	e.SetWidth("100%")
	e.SetJustifyContent("center")
	e.SetBackgroundColor(consoleLogsHeaderBG)

	text, err := NewTextElement(ctx)
	if err != nil {
		return nil, err
	}
	text.SetContent("Console")

	e.AppendChild(text)

	return e, nil
}

type logsListElement struct {
	*ScrollBoxElement

	cancel func()
}

func newLogsListElement(ctx types.RenderContext) (*logsListElement, error) {
	box, err := NewScrollBoxElement(ctx)
	if err != nil {
		return nil, err
	}

	e := &logsListElement{
		ScrollBoxElement: box,
	}
	e.SetMaxHeight("100%")
	e.SetFlexDirection("column")
	e.SetJustifyContent("end")

	logs, cancel := debug.Logs()
	e.cancel = cancel
	go e.watchLogs(logs)

	return e, nil
}

func (s *logsListElement) watchLogs(logs <-chan string) {
	for log := range logs {
		s.addLog(log)
	}
}

func (s *logsListElement) addLog(log string) error {
	elem, err := NewTextElement(s.ctx)
	if err != nil {
		return err
	}
	elem.SetContent(log)
	s.AppendChild(elem)

	go s.ctx.Render()
	return nil
}

func (s *logsListElement) Destroy() error {
	s.mu.Lock()
	if s.destroyed {
		return nil
	}

	s.cancel()
	s.mu.Unlock()
	return s.BoxElement.Destroy()
}
