package main

import (
	"fmt"
	"os"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term"
	. "github.com/AnatoleLucet/loom-term/components"
	. "github.com/AnatoleLucet/loom/components"
)

var (
	styleContainer = Style{
		Width:          "100%",
		Height:         "100%",
		AlignItems:     "center",
		JustifyContent: "center",
		FlexDirection:  "column",
		GapRow:         3,
	}

	styleForm = Style{
		PaddingVertical:   1,
		PaddingHorizontal: 4,
		BackgroundColor:   "#1f2937",
		FlexDirection:     "column",
		GapRow:            1,
	}

	styleTitle = Style{
		AlignSelf:    "center",
		MarginBottom: 1,
	}
	styleInput = Style{
		BackgroundColor:      "#374151",
		PlaceholderFontStyle: "italic",
	}
	styleTextArea = Style{
		Height:               8,
		BackgroundColor:      "#374151",
		PlaceholderFontStyle: "italic",
	}
	styleBtn = Style{
		PaddingHorizontal: 2,
		AlignSelf:         "end",
		BackgroundColor:   "#374151",
	}
)

func InputForm() loom.Node {
	var ref term.InputElement

	result, setResult := Signal("")

	onSubmit := func(event *term.EventSubmit) {
		setResult(event.Value)
		ref.Clear()
	}
	submit := func(*term.EventMouse) {
		ref.Submit()
	}

	hasResult := Memo(func() bool {
		return result() != ""
	})

	return Box(
		B(Text("Input"), Apply(styleTitle)),

		Input(Apply(
			Ref{Ptr: &ref},
			Attr{Placeholder: "Type something..."},
			On{Submit: onSubmit},
			styleInput,
		)),
		Box(
			Text("SUBMIT"),
			Apply(On{Click: submit}, styleBtn),
		),

		Show(hasResult, func() loom.Node {
			return P(Text("Submitted: "), BindText(result))
		}),

		Apply(styleForm, Style{Width: 40}),
	)
}

func TextAreaForm() loom.Node {
	var ref term.TextAreaElement

	result, setResult := Signal("")

	onSubmit := func(event *term.EventSubmit) {
		setResult(event.Value)
		ref.Clear()
	}
	submit := func(*term.EventMouse) {
		ref.Submit()
	}

	hasResult := Memo(func() bool {
		return result() != ""
	})

	return Box(
		B(Text("TextArea"), Apply(styleTitle)),

		TextArea(Apply(
			Ref{Ptr: &ref},
			Attr{Placeholder: "Type something..."},
			On{Submit: onSubmit},
			styleTextArea,
		)),
		Box(
			Text("SUBMIT"),
			Apply(On{Click: submit}, styleBtn),
		),

		Show(hasResult, func() loom.Node {
			return P(Text("Submitted: "), BindText(result))
		}),

		Apply(styleForm, Style{Width: 60}),
	)
}

func App() loom.Node {
	return Box(
		Console(term.IsDev()),

		InputForm(),
		TextAreaForm(),

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
