package main

import (
	"fmt"
	"os"

	"github.com/AnatoleLucet/loom"
	. "github.com/AnatoleLucet/loom-term/components"
	"github.com/AnatoleLucet/loom-term/internal"
	"github.com/AnatoleLucet/loom-term/opentui"
	. "github.com/AnatoleLucet/loom/components"
	. "github.com/AnatoleLucet/loom/signals"
)

// triangle returns a value that bounces between 0 and max
func triangle(t, max int) int {
	period := 2 * max
	v := t % period
	if v > max {
		return period - v
	}
	return v
}

func MyApp() loom.Node {
	termWidth, _, _ := internal.TerminalSize()
	maxPos := termWidth - 2 // Leave room for the 2-wide bar

	frame, setFrame := Signal(0)
	pos, setPos := Signal(0)

	return InlineRenderer(func() loom.Node {
		// return FullscreenRenderer(func() loom.Node {

		go func() {
			internal.Animate(internal.Animation{
				Tick: func(progress float64) {
					setFrame(frame() + 1)
					setPos(triangle(frame(), maxPos))
				},
			})
		}()

		return Box(
			// Console(),

			Box(
				BindText(func() string {
					return fmt.Sprintf("Frame: %d, Position: %d, MaxPos: %d", frame(), pos(), maxPos)
				}),

				&Style{
					Height: 1,
					Width:  "100%",
				},
			),

			// Progress bar container
			Box(
				// Moving bar - position based on frame count
				Bind(func() loom.Node {
					return Box(
						&Style{
							MarginLeft:      pos(), // Integer position - no rounding needed
							Width:           2,
							Height:          2,
							BackgroundColor: opentui.NewRGBA(0.2, 0.8, 0.4, 1), // green
						},
					)
				}),

				&Style{
					Width:           "100%",
					Height:          2,
					BackgroundColor: opentui.NewRGBA(0.2, 0.2, 0.2, 1), // dark gray
					FlexDirection:   "row",
					AlignItems:      "center",
				},
			),

			&Style{
				Width:         "100%",
				FlexDirection: "column",
				MarginTop:     2,
			},
		)
	})
}

func main() {
	app := NewApp()

	for err := range app.Run(MyApp) {
		fmt.Printf("Error: %v\n", err)
		app.Close()
		os.Exit(1)
	}
}
