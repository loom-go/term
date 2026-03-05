package term

import (
	"os"
	"strings"
)

// source: https://github.com/pocketbase/pocketbase/blob/93e3eb3a35e35948a2edeaa67c356ac63a192ad6/tools/osutils/run.go#L11

// IsDev returns true if the program was started with "go run".
func IsDev() bool {
	for _, dir := range runDirs {
		if dir != "" && strings.HasPrefix(os.Args[0], dir) {
			return true
		}
	}

	return false
}

var runDirs = []string{os.TempDir(), cacheDir()}

func cacheDir() string {
	dir := os.Getenv("GOCACHE")
	if dir == "off" {
		return ""
	}

	if dir == "" {
		dir, _ = os.UserCacheDir()
	}

	return dir
}
