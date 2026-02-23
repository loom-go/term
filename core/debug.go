package core

import "github.com/AnatoleLucet/loom-term/core/debug"

func LogDebug(args ...any)                 { debug.LogDebug(args...) }
func LogDebugf(format string, args ...any) { debug.LogDebugf(format, args...) }

func LogInfo(args ...any)                 { debug.LogInfo(args...) }
func LogInfof(format string, args ...any) { debug.LogInfof(format, args...) }

func LogWarning(args ...any)                 { debug.LogWarning(args...) }
func LogWarningf(format string, args ...any) { debug.LogWarningf(format, args...) }

func LogError(args ...any)                 { debug.LogError(args...) }
func LogErrorf(format string, args ...any) { debug.LogErrorf(format, args...) }
