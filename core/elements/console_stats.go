package elements

import (
	"fmt"
	"github.com/loom-go/term/core/debug"
	"github.com/loom-go/term/core/gfx"
	"time"

	"github.com/AnatoleLucet/go-opentui"
)

type statsElement struct {
	*BoxElement
}

func newStatsElement() (*statsElement, error) {
	box, err := NewBoxElement()
	if err != nil {
		return nil, err
	}

	e := &statsElement{box}
	box.self = e
	e.SetTop("40%")
	e.SetRight(0)
	e.SetGapColumn(2)
	e.SetPaddingVertical(1)
	e.SetPaddingHorizontal(2)
	e.SetOverflow("hidden")
	e.SetPosition("absolute")
	e.SetBackgroundColor(consoleBG)

	labelsElem, err := newStatsLabelsElement()
	if err != nil {
		return nil, err
	}
	valuesElem, err := newStatsValuesElement()
	if err != nil {
		return nil, err
	}

	e.AppendChild(labelsElem)
	e.AppendChild(valuesElem)

	return e, nil
}

type statsLabelsElement struct {
	*BaseElement
}

func newStatsLabelsElement() (*statsLabelsElement, error) {
	base, err := NewBaseElement()
	if err != nil {
		return nil, err
	}

	e := &statsLabelsElement{base}
	base.self = e
	e.SetFlexDirection("column")

	labels := []string{
		"frame/s:",
		"frame:",
		"layout:",
		"record:",
		"draw:",
		"render:",
	}
	for _, label := range labels {
		labelElem, err := NewTextElement()
		if err != nil {
			return nil, err
		}
		labelElem.SetText(label)
		e.AppendChild(labelElem)
	}

	return e, nil
}

type statsValuesElement struct {
	*BaseElement
}

func newStatsValuesElement() (*statsValuesElement, error) {
	base, err := NewBaseElement()
	if err != nil {
		return nil, err
	}

	e := &statsValuesElement{base}
	base.self = e
	e.SetFlexDirection("column")

	fps, cancelFps := debug.FPS()
	frameTime, cancelFrameTime := debug.FrameTime()
	layoutTime, cancelLayoutTime := debug.LayoutTime()
	recordTime, cancelRecordTime := debug.RecordTime()
	drawTime, cancelDrawTime := debug.DrawTime()
	renderTime, cancelRenderTime := debug.RenderTime()

	fpsElem, err := newStatElement(fps, cancelFps, e.fpsToString)
	frameTimeElem, err := newStatElement(frameTime, cancelFrameTime, e.timingRecordToString)
	layoutTimeElem, err := newStatElement(layoutTime, cancelLayoutTime, e.timingRecordToString)
	recordTimeElem, err := newStatElement(recordTime, cancelRecordTime, e.timingRecordToString)
	drawTimeElem, err := newStatElement(drawTime, cancelDrawTime, e.timingRecordToString)
	renderTimeElem, err := newStatElement(renderTime, cancelRenderTime, e.timingRecordToString)
	if err != nil {
		return nil, err
	}

	e.AppendChild(fpsElem)
	e.AppendChild(frameTimeElem)
	e.AppendChild(layoutTimeElem)
	e.AppendChild(recordTimeElem)
	e.AppendChild(drawTimeElem)
	e.AppendChild(renderTimeElem)

	return e, nil
}

func (s *statsValuesElement) fpsToString(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

func (s *statsValuesElement) timingRecordToString(v *debug.TimingRecord) string {
	if v == nil {
		return "n/a"
	}

	return fmt.Sprintf(
		"%0.2fms (avg: %0.2fms, max: %0.2fms)",
		float32(v.Last)/float32(time.Millisecond),
		float32(v.Avg)/float32(time.Millisecond),
		float32(v.Max)/float32(time.Millisecond),
	)
}

type statElement[T any] struct {
	*BaseElement

	ch     <-chan T
	format func(T) string

	latest string

	textBuffer     *opentui.TextBuffer
	textBufferView *opentui.TextBufferView
}

func newStatElement[T any](
	ch <-chan T,
	cancel func(),
	format func(T) string,
) (*statElement[T], error) {
	base, err := NewBaseElement()
	if err != nil {
		return nil, err
	}

	e := &statElement[T]{
		BaseElement: base,
		ch:          ch,
		format:      format,
	}
	base.self = e
	e.SetWidth(33)
	e.SetHeight(1)

	e.textBuffer = opentui.NewTextBuffer(0)
	e.textBufferView = opentui.NewTextBufferView(e.textBuffer)
	e.textBuffer.Append("n/a")

	e.OnDestroy(cancel)

	go e.watch()

	return e, nil
}

func (s *statElement[T]) watch() {
	for v := range s.ch {
		// use a locked field instead of updating the buffer directly
		// to make sure we're not going to update or reset the buffer mid-render
		s.mu.Lock()
		s.latest = s.format(v)
		s.mu.Unlock()
	}
}

func (s *statElement[T]) Render(buffer *opentui.Buffer, rect gfx.Rect) error {
	s.mu.Lock()
	latest := s.latest
	s.mu.Unlock()

	s.textBuffer.Reset()
	s.textBuffer.Append(latest)

	return buffer.DrawTextBufferView(s.textBufferView, int32(rect.X), int32(rect.Y))
}
