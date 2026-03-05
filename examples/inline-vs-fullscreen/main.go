package main

import (
	"fmt"
	"os"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term"
	. "github.com/AnatoleLucet/loom-term/components"
)

var (
	renderTypes = map[string]term.RenderType{
		"--fullscreen": term.RenderFullscreen,
		"--inline":     term.RenderInline,
	}

	styleBox = Style{
		PaddingVertical:   2,
		PaddingHorizontal: 6,
		AlignSelf:         "start",
		AlignItems:        "center",
		FlexDirection:     "column",
		BackgroundColor:   "#374151",
	}
)

func Inline() loom.Node {
	return Box(
		B(Text("Rendering Inline")),
		Text(""),
		Text("The app is rendered right bellow your"),
		Text("terminal prompt. It does not take the entire screen."),

		Apply(styleBox),
	)
}

func Fullscreen() loom.Node {
	return Box(
		B(Text("Rendering Fullscreen")),
		Text(""),
		Text("The app takes the entire screen."),
		Text("You can go back to your terminal prompt"),
		Text("by closing the app."),

		Apply(styleBox),
	)
}

func App() loom.Node {
	ctx := term.Context()

	switch ctx.RenderType() {
	case term.RenderInline:
		return Inline()
	case term.RenderFullscreen:
		return Fullscreen()
	}

	return nil
}

func main() {
	var typ term.RenderType
	var ok bool
	for _, arg := range os.Args {
		if t, found := renderTypes[arg]; found {
			typ = t
			ok = true
			break
		}
	}

	if !ok {
		fmt.Println("You must specify either --inline or --fullscreen")
		os.Exit(1)
	}

	app := term.NewApp()
	for err := range app.Run(typ, App) {
		app.Close()
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
