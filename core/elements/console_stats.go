package elements

import (
	"fmt"
	"sync"
	"time"

	"github.com/AnatoleLucet/go-opentui"
	"github.com/AnatoleLucet/loom-term/core/debug"
	"github.com/AnatoleLucet/loom-term/core/types"
)

type statsElement struct {
	*BoxElement
}

func newStatsElement(ctx types.RenderContext) (*statsElement, error) {
	box, err := NewBoxElement(ctx)
	if err != nil {
		return nil, err
	}

	e := &statsElement{box}
	e.SetTop("40%")
	e.SetRight(0)
	e.SetGapColumn(2)
	e.SetPaddingVertical(1)
	e.SetPaddingHorizontal(2)
	e.SetOverflow("hidden")
	e.SetPosition("absolute")
	e.SetBackgroundColor(consoleBG)

	labelsElem, err := newStatsLabelsElement(ctx)
	if err != nil {
		return nil, err
	}
	valuesElem, err := newStatsValuesElement(ctx)
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

func newStatsLabelsElement(ctx types.RenderContext) (*statsLabelsElement, error) {
	base, err := NewElement(ctx)
	if err != nil {
		return nil, err
	}

	e := &statsLabelsElement{base}
	e.SetFlexDirection("column")

	labels := []string{
		"frame/s:",
		"frame:",
		"layout:",
		"paint:",
		"render:",
	}
	for _, label := range labels {
		labelElem, err := NewTextElement(ctx)
		if err != nil {
			return nil, err
		}
		labelElem.SetContent(label)
		e.AppendChild(labelElem)
	}

	return e, nil
}

type statsValuesElement struct {
	*BaseElement
}

func newStatsValuesElement(ctx types.RenderContext) (*statsValuesElement, error) {
	base, err := NewElement(ctx)
	if err != nil {
		return nil, err
	}

	e := &statsValuesElement{base}
	e.SetFlexDirection("column")

	fps, cancelFps := debug.FPS()
	frameTime, cancelFrameTime := debug.FrameTime()
	layoutTime, cancelLayoutTime := debug.LayoutTime()
	paintTime, cancelPaintTime := debug.PaintTime()
	renderTime, cancelRenderTime := debug.RenderTime()

	fpsElem, err := newStatElement(ctx, fps, cancelFps, e.fpsToString)
	if err != nil {
		return nil, err
	}
	frameTimeElem, err := newStatElement(ctx, frameTime, cancelFrameTime, e.timingRecordToString)
	if err != nil {
		return nil, err
	}
	layoutTimeElem, err := newStatElement(ctx, layoutTime, cancelLayoutTime, e.timingRecordToString)
	if err != nil {
		return nil, err
	}
	paintTimeElem, err := newStatElement(ctx, paintTime, cancelPaintTime, e.timingRecordToString)
	if err != nil {
		return nil, err
	}
	renderTimeElem, err := newStatElement(ctx, renderTime, cancelRenderTime, e.timingRecordToString)
	if err != nil {
		return nil, err
	}

	e.AppendChild(fpsElem)
	e.AppendChild(frameTimeElem)
	e.AppendChild(layoutTimeElem)
	e.AppendChild(paintTimeElem)
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
	*TextElement

	statMu sync.Mutex

	latest string

	ch     <-chan T
	cancel func()
	format func(T) string
}

func newStatElement[T any](
	ctx types.RenderContext,
	ch <-chan T,
	cancel func(),
	format func(T) string,
) (*statElement[T], error) {
	textElem, err := NewTextElement(ctx)
	if err != nil {
		return nil, err
	}

	e := &statElement[T]{
		TextElement: textElem,
		latest:      "n/a",
		ch:          ch,
		cancel:      cancel,
		format:      format,
	}
	// need to have static dimensions because we can't rely on layout recalculation
	e.SetWidth(33)
	e.SetHeight(1)
	e.XYZ().UnsetMeasureFunc()

	go e.watch()

	return e, nil
}

func (s *statElement[T]) Destroy() error {
	s.mu.Lock()
	if s.destroyed {
		return nil
	}

	s.cancel()
	s.mu.Unlock()
	return s.TextElement.Destroy()
}

func (s *statElement[T]) watch() {
	for v := range s.ch {
		s.statMu.Lock()
		s.latest = s.format(v)
		s.statMu.Unlock()
	}
}

func (s *statElement[T]) Paint(buffer *opentui.Buffer, x, y float32) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.guardPaint(); err != nil {
		return err
	}

	s.statMu.Lock()
	latest := s.latest
	s.statMu.Unlock()

	if latest != s.textBuffer.GetPlainText(0) {
		s.textBuffer.Reset()
		s.textBuffer.Append(latest)
	}

	err := buffer.DrawTextBufferView(s.textBufferView, int32(x), int32(y))
	if err != nil {
		return err
	}

	return nil
}
