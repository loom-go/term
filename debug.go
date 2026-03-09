package term

import "github.com/loom-go/term/core"

func LogDebug(args ...any)                 { core.LogDebug(args...) }
func LogDebugf(format string, args ...any) { core.LogDebugf(format, args...) }

func LogInfo(args ...any)                 { core.LogInfo(args...) }
func LogInfof(format string, args ...any) { core.LogInfof(format, args...) }

func LogWarning(args ...any)                 { core.LogWarning(args...) }
func LogWarningf(format string, args ...any) { core.LogWarningf(format, args...) }

func LogError(args ...any)                 { core.LogError(args...) }
func LogErrorf(format string, args ...any) { core.LogErrorf(format, args...) }
