package debug

import (
	"context"
	"fmt"
	"time"
)

var debugger = NewDebugger()

type Debugger struct {
	logs *Emitter[string]

	fps *RateMetric

	frameTime  *TimingMetric
	layoutTime *TimingMetric
	paintTime  *TimingMetric
	renderTime *TimingMetric
}

func NewDebugger() *Debugger {
	return &Debugger{
		logs: NewEmitter[string](),

		fps: NewRateMetric(time.Second),

		frameTime:  NewTimingMetric(),
		renderTime: NewTimingMetric(),
		layoutTime: NewTimingMetric(),
		paintTime:  NewTimingMetric(),
	}
}

func (d *Debugger) emitLog(message string) {
	go d.logs.Emit(message)
}

func LogDebug(message string) {
	debugger.emitLog(fmt.Sprintf("[DEBUG] [%s] %s", time.Now().Format(time.TimeOnly), message))
}

func LogInfo(message string) {
	debugger.emitLog(fmt.Sprintf("[INFO] [%s] %s", time.Now().Format(time.TimeOnly), message))
}

func LogWarning(message string) {
	debugger.emitLog(fmt.Sprintf("[WARNING] [%s] %s", time.Now().Format(time.TimeOnly), message))
}

func LogError(message string) {
	debugger.emitLog(fmt.Sprintf("[ERROR] [%s] %s", time.Now().Format(time.TimeOnly), message))
}

func Logs() (logs <-chan string, cancel func()) {
	ctx, cancel := context.WithCancel(context.Background())
	return debugger.logs.Subscribe(ctx, 100), cancel
}

func FPS() (fps <-chan float64, cancel func()) {
	return debugger.fps.Subscribe(10)
}

func EmitFrameTime(duration time.Duration) {
	go debugger.fps.Emit()
	go debugger.frameTime.Emit(duration)
}

func FrameTime() (durations <-chan *TimingRecord, cancel func()) {
	return debugger.frameTime.Subscribe(10)
}

func EmitLayoutTime(duration time.Duration) {
	go debugger.layoutTime.Emit(duration)
}

func LayoutTime() (durations <-chan *TimingRecord, cancel func()) {
	return debugger.layoutTime.Subscribe(10)
}

func EmitPaintTime(duration time.Duration) {
	go debugger.paintTime.Emit(duration)
}

func PaintTime() (durations <-chan *TimingRecord, cancel func()) {
	return debugger.paintTime.Subscribe(10)
}

func EmitRenderTime(duration time.Duration) {
	go debugger.renderTime.Emit(duration)
}

func RenderTime() (durations <-chan *TimingRecord, cancel func()) {
	return debugger.renderTime.Subscribe(10)
}
