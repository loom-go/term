package main

import (
	"fmt"
	"os"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term"
	. "github.com/AnatoleLucet/loom-term/components"
	"github.com/AnatoleLucet/loom-term/core"
	. "github.com/AnatoleLucet/loom/components"
)

var (
	styleContainer = Style{
		Width:          "100%",
		AlignItems:     "center",
		JustifyContent: "center",
		FlexDirection:  "column",
		GapRow:         2,
	}

	styleScrollBox = Style{
		Width:     80,
		Height:    40,
		MaxWidth:  "100%",
		MaxHeight: "100%",

		BackgroundColor: "#1f2937",
	}

	styleBoxesList = Style{
		Width: 218,

		AlignSelf: "center",
		FlexWrap:  "wrap",

		GapRow:            1,
		GapColumn:         2,
		PaddingVertical:   1,
		PaddingHorizontal: 2,
	}
	styleBox = Style{
		Width:  10,
		Height: 5,

		AlignItems:     "center",
		JustifyContent: "center",

		BackgroundColor: "#374151",
	}
	styleBoxHover = Style{
		BackgroundOpacity: 0.3,
	}

	actionsContainer = Style{
		AlignItems: "center",
		GapColumn:  3,
	}
	actionsMiddle = Style{
		FlexDirection:  "column",
		JustifyContent: "space-between",
		GapRow:         3,
	}
	btnStyle = Style{
		PaddingHorizontal: 2,
		BackgroundColor:   "#374151",
	}
)

func App() loom.Node {
	var ref core.ScrollBoxElement

	scrollTop := func(*term.EventMouse) { ref.ScrollToTop() }
	scrollBottom := func(*term.EventMouse) { ref.ScrollToBottom() }
	scrollRight := func(*term.EventMouse) { ref.ScrollToRight() }
	scrollLeft := func(*term.EventMouse) { ref.ScrollToLeft() }

	var boxes []loom.Node
	for i := range 500 {
		box := Box(
			Text(fmt.Sprintf("Box %d", i+1)),
			Apply(styleBox),
			ApplyOn("hover", styleBoxHover),
		)
		boxes = append(boxes, box)
	}

	return Box(
		Console(term.IsDev()),

		B(Text("Try scrolling!")),

		ScrollBox(
			Box(
				Fragment(boxes...),
				Apply(styleBoxesList),
			),

			Ref(&ref),
			Apply(styleScrollBox),
		),

		Box(
			Box(Text("LEFT"), On("click", scrollLeft), Apply(btnStyle)),
			Box(
				Box(Text("TOP"), On("click", scrollTop), Apply(btnStyle)),
				Box(Text("BOT"), On("click", scrollBottom), Apply(btnStyle)),
				Apply(actionsMiddle),
			),
			Box(Text("RIGHT"), On("click", scrollRight), Apply(btnStyle)),

			Apply(actionsContainer),
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
