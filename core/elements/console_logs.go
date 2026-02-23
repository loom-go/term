package elements

import (
	"fmt"
	"sync"

	"github.com/AnatoleLucet/loom-term/core/debug"
	"github.com/AnatoleLucet/loom-term/core/elements/events"
)

var (
	consoleLogsStateMu sync.Mutex

	consoleLogsDefaultHeight = "35%"
	consoleLogsMinHeight     = 7
	consoleLogsLastHeight    = 0
)

type logsElement struct {
	*BoxElement

	cancel func()
}

func newLogsElement() (*logsElement, error) {
	box, err := NewBoxElement()
	if err != nil {
		return nil, err
	}

	e := &logsElement{
		BoxElement: box,
	}
	box.self = e
	e.SetWidth("100%")
	e.SetLeft(0)
	e.SetBottom(0)
	e.SetPosition("absolute")
	e.SetFlexDirection("column")
	e.SetBackgroundColor(consoleBG)

	if consoleLogsLastHeight == 0 {
		e.SetHeight(consoleLogsDefaultHeight)
	} else {
		e.SetHeight(consoleLogsLastHeight)
	}

	header, err := newLogsHeaderElement()
	list, err := newLogsListElement()
	attachButton, err := newLogsAttachButtonElement()
	if err != nil {
		return nil, err
	}

	e.AppendChild(header)
	e.AppendChild(list)

	attachButton.OnMousePress(func(event *EventMouse) {
		if event.Button == events.MouseLeft && list.IsDetached() {
			list.Attach()
			e.RemoveChild(attachButton)
		}
	})

	list.OnMouseScroll(func(event *EventMouse) {
		viewportH := list.ViewportHeight()
		contentH := list.ContentHeight()
		maxScrollY := contentH - viewportH

		if !list.IsDetached() && event.Button == events.MouseWheelUp {
			list.Detach()
			e.AppendChild(attachButton)
		}

		isNearBottom := list.ScrollY() >= maxScrollY-2
		if isNearBottom && list.IsDetached() && event.Button == events.MouseWheelDown {
			list.Attach()
			e.RemoveChild(attachButton)
		}
	})

	header.mouseDragAction(func(event *EventMouse) {
		if e.rdrctx == nil {
			return
		}

		e.mu.RLock()
		currentY := int(e.xyz().GetLayout().AbsoluteTop())
		currentHeight := int(e.xyz().GetLayout().Height())
		newHeight := max(consoleLogsMinHeight, currentHeight-(event.Y-currentY))
		e.mu.RUnlock()

		consoleLogsStateMu.Lock()
		consoleLogsLastHeight = newHeight
		consoleLogsStateMu.Unlock()

		e.SetHeight(newHeight)
		e.rdrctx.ScheduleRender()
	})

	return e, nil
}

type logsHeaderElement struct {
	*BoxElement
}

func newLogsHeaderElement() (*logsHeaderElement, error) {
	box, err := NewBoxElement()
	if err != nil {
		return nil, err
	}

	e := &logsHeaderElement{
		BoxElement: box,
	}
	box.self = e
	e.SetWidth("100%")
	e.SetJustifyContent("center")
	e.SetBackgroundColor(consoleLogsHeaderBG)

	text, err := NewTextElement()
	if err != nil {
		return nil, err
	}
	text.SetText("Console")

	e.AppendChild(text)

	return e, nil
}

type logsListElement struct {
	*ScrollBoxElement

	detached bool
}

func newLogsListElement() (*logsListElement, error) {
	box, err := NewScrollBoxElement()
	if err != nil {
		return nil, err
	}

	e := &logsListElement{
		ScrollBoxElement: box,
	}
	box.self = e
	e.SetMaxHeight("100%")
	e.SetFlexDirection("column")
	e.SetJustifyContent("end")

	logs, cancel := debug.Logs()
	e.OnDestroy(cancel)

	go e.watchLogs(logs)

	return e, nil
}

func (s *logsListElement) watchLogs(logs <-chan *debug.LogEntry) {
	for log := range logs {
		s.addLog(log)
	}
}

func (s *logsListElement) addLog(log *debug.LogEntry) error {
	container, err := NewTextElement()

	level, err := NewTextElement()
	level.SetText(fmt.Sprintf("[%s] ", log.Level))
	level.SetTextForeground(consoleLogsLevelColors[log.Level])

	date, err := NewTextElement()
	date.SetText(fmt.Sprintf("[%s] ", log.Time.Format("15:04:05")))

	message, err := NewTextElement()
	message.SetText(log.Message)

	container.AppendChild(level)
	container.AppendChild(date)
	container.AppendChild(message)
	s.AppendChild(container)

	if err != nil {
		return err
	}

	if !s.detached {
		s.ScrollToBottom()
	}

	if s.rdrctx != nil {
		s.rdrctx.ScheduleRender()
	}
	return nil
}

func (s *logsListElement) IsDetached() bool {
	return s.detached
}

func (s *logsListElement) Detach() {
	s.detached = true
}

func (s *logsListElement) Attach() {
	s.detached = false
	s.ScrollToBottom()
}

type logsAttachButtonElement struct {
	*BoxElement
}

func newLogsAttachButtonElement() (*logsAttachButtonElement, error) {
	box, err := NewBoxElement()
	if err != nil {
		return nil, err
	}

	e := &logsAttachButtonElement{BoxElement: box}
	box.self = e
	box.SetRight(4)
	box.SetBottom(0)
	box.SetPosition("absolute")
	box.SetPaddingHorizontal(1)
	box.SetBackgroundColor(consoleLogsHeaderBG)

	text, err := NewTextElement()
	if err != nil {
		return nil, err
	}
	text.SetText(" ⮟ ")
	e.AppendChild(text)

	return e, nil
}
