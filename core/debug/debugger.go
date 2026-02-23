package debug

import (
	"context"
	"fmt"
	"time"
)

var debugger = NewDebugger()

type Debugger struct {
	logs *Emitter[*LogEntry]

	fps *RateMetric

	frameTime  *TimingMetric
	layoutTime *TimingMetric
	recordTime *TimingMetric
	drawTime   *TimingMetric
	renderTime *TimingMetric
}

func NewDebugger() *Debugger {
	return &Debugger{
		logs: NewEmitter[*LogEntry](),

		fps: NewRateMetric(time.Second),

		frameTime:  NewTimingMetric(),
		renderTime: NewTimingMetric(),
		recordTime: NewTimingMetric(),
		layoutTime: NewTimingMetric(),
		drawTime:   NewTimingMetric(),
	}
}

func (d *Debugger) emitLog(log *LogEntry) {
	go d.logs.Emit(log)
}

func LogDebugf(format string, args ...any) {
	LogDebug(fmt.Sprintf(format, args...))
}

func LogDebug(args ...any) {
	log := &LogEntry{
		Level:   LogLevelDebug,
		Message: fmt.Sprint(args...),
		Time:    time.Now(),
	}

	debugger.emitLog(log)
}

func LogInfof(format string, args ...any) {
	LogInfo(fmt.Sprintf(format, args...))
}

func LogInfo(args ...any) {
	log := &LogEntry{
		Level:   LogLevelInfo,
		Message: fmt.Sprint(args...),
		Time:    time.Now(),
	}

	debugger.emitLog(log)
}

func LogWarningf(format string, args ...any) {
	LogWarning(fmt.Sprintf(format, args...))
}

func LogWarning(args ...any) {
	log := &LogEntry{
		Level:   LogLevelWarning,
		Message: fmt.Sprint(args...),
		Time:    time.Now(),
	}

	debugger.emitLog(log)
}

func LogErrorf(format string, args ...any) {
	LogError(fmt.Sprintf(format, args...))
}

func LogError(args ...any) {
	log := &LogEntry{
		Level:   LogLevelError,
		Message: fmt.Sprint(args...),
		Time:    time.Now(),
	}

	debugger.emitLog(log)
}

func Logs() (logs <-chan *LogEntry, cancel func()) {
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

func EmitRecordTime(duration time.Duration) {
	go debugger.recordTime.Emit(duration)
}

func RecordTime() (durations <-chan *TimingRecord, cancel func()) {
	return debugger.recordTime.Subscribe(10)
}

func EmitDrawTime(duration time.Duration) {
	go debugger.drawTime.Emit(duration)
}

func DrawTime() (durations <-chan *TimingRecord, cancel func()) {
	return debugger.drawTime.Subscribe(10)
}

func EmitRenderTime(duration time.Duration) {
	go debugger.renderTime.Emit(duration)
}

func RenderTime() (durations <-chan *TimingRecord, cancel func()) {
	return debugger.renderTime.Subscribe(10)
}
