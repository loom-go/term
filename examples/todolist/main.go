package main

import (
	"fmt"
	"os"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term"
	. "github.com/AnatoleLucet/loom-term/components"
	. "github.com/AnatoleLucet/loom/components"
)

var art = `
 ███████████              █████           █████████ 
░█░░░███░░░█             ░░███           ███░░░░░███
░   ░███  ░   ██████   ███████   ██████ ░███    ░░░ 
    ░███     ███░░███ ███░░███  ███░░███░░█████████ 
    ░███    ░███ ░███░███ ░███ ░███ ░███ ░░░░░░░░███
    ░███    ░███ ░███░███ ░███ ░███ ░███ ███    ░███
    █████   ░░██████ ░░████████░░██████ ░░█████████ 
   ░░░░░     ░░░░░░   ░░░░░░░░  ░░░░░░   ░░░░░░░░░ `

var (
	containerStyle = Style{
		Width:           "100%",
		Height:          "100%",
		AlignItems:      "center",
		FlexDirection:   "column",
		PaddingVertical: "2%",
		GapRow:          4,
	}
)

func App() loom.Node {
	return Box(
		Console(term.IsDev()),

		Text(art, Apply(Style{Color: "#248564"})),
		TodoApp(),

		Apply(containerStyle),
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
