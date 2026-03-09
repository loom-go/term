package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/loom-go/loom"
	. "github.com/loom-go/loom/components"
	"github.com/loom-go/term"
	. "github.com/loom-go/term/components"
)

var (
	styleContainer = Style{
		Width:           "100%",
		Height:          "100%",
		PaddingVertical: "4%",
		AlignItems:      "center",
		FlexDirection:   "column",
		GapRow:          3,
	}

	styleContent = Style{
		AlignItems:    "center",
		FlexDirection: "column",
		MaxWidth:      80,
	}

	styleActionsContainer = Style{
		AlignItems:    "center",
		FlexDirection: "column",
		GapRow:        1,
	}
	styleActions = Style{
		GapColumn: 2,
	}
	styleBtn = Style{
		PaddingHorizontal: 2,
	}
	styleBtnText = Style{
		Color: "white",
	}

	styleDebug = Style{BackgroundColor: "#6b7280"}
	StyleInfo  = Style{BackgroundColor: "#3b82f6"}
	styleWarn  = Style{BackgroundColor: "#f59e0b"}
	styleError = Style{BackgroundColor: "#ef4444"}
)

func App() loom.Node {
	go func() {
		for t := range time.Tick(time.Second) {
			term.LogInfof("tick at %s", t.Format(time.RFC3339))
		}
	}()

	return Box(
		Console(term.IsDev()), // only enable the console when the program is run with "go run"

		Box(
			B(Text("LoomTerm Console")),
			Text(""),
			Text("Press '`' to open the console."),
			Apply(styleContent),
		),

		Box(
			Text("The console supports multiple log levels."),
			Text("You can log messages from anywhere in your app, and they will appear here."),
			Text("Feel free to scroll up and down to see the history of log messages,"),
			Text("and to resize the console with your mouse to see more!"),
			Apply(styleContent),
		),

		Actions(),

		Apply(styleContainer),
	)
}

func Actions() loom.Node {
	var user bytes.Buffer
	json.Indent(
		&user,
		[]byte(os.ExpandEnv(`{"name":"$USER","age":28,"skills":["Go","Loom","TUI"]}`)),
		"",
		"  ",
	)

	logDebug := func(*term.EventMouse) { term.LogDebug("This is a debug message.") }
	logInfo := func(*term.EventMouse) { term.LogInfof("Some json data: %s", user.String()) }
	logWarn := func(*term.EventMouse) { term.LogWarning("WATCH OUT!") }
	logError := func(*term.EventMouse) { term.LogError("Oh no :((( something went wrong...") }

	return Box(
		Text("Click the buttons:"),
		Box(
			Box(
				Text("DEBUG", Apply(styleBtnText)),
				Apply(On{Click: logDebug}, styleBtn, styleDebug),
			),
			Box(
				Text("INFO", Apply(styleBtnText)),
				Apply(On{Click: logInfo}, styleBtn, StyleInfo),
			),
			Box(
				Text("WARN", Apply(styleBtnText)),
				Apply(On{Click: logWarn}, styleBtn, styleWarn),
			),
			Box(
				Text("ERROR", Apply(styleBtnText)),
				Apply(On{Click: logError}, styleBtn, styleError),
			),

			Apply(styleActions),
		),

		Apply(styleActionsContainer),
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
