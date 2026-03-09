package main

import (
	"fmt"
	"os"
	"time"

	"github.com/loom-go/loom"
	. "github.com/loom-go/loom/components"
	"github.com/loom-go/term"
	. "github.com/loom-go/term/components"
)

func Counter() loom.Node {
	frame, setFrame := Signal(0)

	go func(self loom.Component) {
		for !self.IsDisposed() {
			time.Sleep(time.Second / 120)
			setFrame(frame() + 1)
		}
	}(Self())

	return Box(Text("Count: "), BindText(frame))
}

func main() {
	app := term.NewApp()

	for err := range app.Run(term.RenderInline, Counter) {
		app.Close()
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
