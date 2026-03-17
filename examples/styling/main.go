package main

import (
	"fmt"
	"os"

	"github.com/loom-go/loom"
	. "github.com/loom-go/loom/components"
	"github.com/loom-go/term"
	"github.com/loom-go/term/animate"
	. "github.com/loom-go/term/components"
)

var (
	styleContainer = Style{
		Width:           "100%",
		Height:          "100%",
		AlignItems:      "center",
		JustifyContent:  "center",
		GapColumn:       4,
		BackgroundColor: RGB(0, 0, 40),
	}

	styleBox = Style{
		Width:          20,
		Height:         10,
		AlignItems:     "center",
		JustifyContent: "center",
		Position:       "relative",
	}
	styleBoxHover = Style{
		BackgroundOpacity: 0.6,
	}

	styleText = Style{
		Color:           RGB(0, 0, 20),
		BackgroundColor: RGBA(255, 255, 255, 0.4),
	}
)

func App() loom.Node {
	n, setN := Signal(0)

	box1Color := Memo(func() string { return HSL(float64(n()), 100, 50) })
	box2Color := Memo(func() string { return HSL(float64(n()+120), 100, 50) })
	box3Color := Memo(func() string { return HSL(float64(n()+240), 100, 50) })

	// run an animation that increments n every tick (60fps with default)
	// see examples/term/animation/ for more info
	go animate.Run(animate.A{
		Context: Self().Context(),
		Tick: func(float64) {
			setN(n() + 1)
		},
	})

	return Box(
		Console(term.IsDev()),

		Box(
			BindText(box1Color, Apply(styleText)),
			Apply(styleBox, Style{BackgroundColor: box1Color}),
			ApplyOn("hover", styleBoxHover),
		),
		Box(
			BindText(box2Color, Apply(styleText)),
			Apply(styleBox, Style{BackgroundColor: box2Color}),
			ApplyOn("hover", styleBoxHover),
		),
		Box(
			BindText(box3Color, Apply(styleText)),
			Apply(styleBox, Style{BackgroundColor: box3Color}),
			ApplyOn("hover", styleBoxHover),
		),

		Apply(styleContainer),
	)
}

func main() {
	app := term.NewApp()

	for err := range app.Run(term.RenderFullscreen, App) {
		app.Close()
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
