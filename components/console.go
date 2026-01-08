package components

import (
	"fmt"
	"time"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/opentui"
	. "github.com/AnatoleLucet/loom/components"
	. "github.com/AnatoleLucet/loom/signals"
)

func Console() loom.Node {
	return Box(
		fps(),

		&Style{
			Width:  "100%",
			Height: "100%",

			Position: "absolute",
		},
	)
}

func fps() loom.Node {
	// ctx, err := getAppContext()
	// if err != nil {
	// 	return nil
	// }

	fps, setFps := Signal(0.0)
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()

		for range ticker.C {
			// setFps(renderer.Debug().GetFPS())
			setFps(0)
		}
	}()

	return Box(
		Text("fps: "),
		Bind(func() loom.Node {
			return Text(fmt.Sprintf("%.2f", fps()))
		}),

		&Style{
			Width:  10,
			Height: 1,

			Right: 0,
			Top:   0,

			Position: "absolute",

			BackgroundColor: opentui.NewRGBA(0, 0, 0, 0.5),
		},
	)

}
