package types

import "fmt"

type EventExit struct {
	BaseEvent
	Signal ExitSignal
}

func (e EventExit) String() string {
	return fmt.Sprintf("Exit(%s)", e.Signal)
}

type ExitSignal int

const (
	ExitSigInt ExitSignal = iota
	ExitSigTerm
	ExitSigQuit
	ExitSigHup
)

func (s ExitSignal) String() string {
	switch s {
	case ExitSigInt:
		return "SIGINT"
	case ExitSigTerm:
		return "SIGTERM"
	case ExitSigQuit:
		return "SIGQUIT"
	case ExitSigHup:
		return "SIGHUP"
	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}
