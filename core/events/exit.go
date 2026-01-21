package events

import (
	"os"
	"syscall"

	"github.com/AnatoleLucet/loom-term/core/types"
)

func toExitSignal(sig os.Signal) (types.ExitSignal, bool) {
	switch sig {
	case syscall.SIGINT:
		return types.ExitSigInt, true
	case syscall.SIGTERM:
		return types.ExitSigTerm, true
	case syscall.SIGQUIT:
		return types.ExitSigQuit, true
	case syscall.SIGHUP:
		return types.ExitSigHup, true
	default:
		return 0, false
	}
}
