package main

import (
	"math/rand"
	"strings"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term"
	. "github.com/AnatoleLucet/loom-term/components"
	. "github.com/AnatoleLucet/loom/components"
)

var (
	examples = []string{
		"touch grass",
		"*nothing*",
		"stare at the sun",
		"bake a cake!",
		"quit ai",
		"follow @AnatoleLucet on x",
		"watch foodfight",
		"call your mom",
		"tell my friends about loom",
	}

	styleTodoApp = Style{
		Height:            "100%",
		Width:             "100%",
		MaxWidth:          70,
		PaddingVertical:   1,
		PaddingHorizontal: 4,
		AlignItems:        "center",
		FlexDirection:     "column",
		GapRow:            2,
	}

	styleTodoForm = Style{
		Width:             "100%",
		MaxWidth:          66,
		PaddingVertical:   1,
		PaddingHorizontal: 2,
		FlexShrink:        "0",
		BackgroundColor:   "#374151",
	}
	styleTodoInput = Style{
		PlaceholderColor:     "#9ca3af",
		PlaceholderFontStyle: "italic",
	}

	styleTodoList = Style{
		Width:             "100%",
		MaxWidth:          55,
		FlexDirection:     "column",
		GapRow:            1,
		PaddingHorizontal: 2,
		PaddingVertical:   1,
		BackgroundColor:   "#1f2937",
	}
	styleTodoItem = Style{
		PaddingVertical:   1,
		PaddingHorizontal: 3,
		JustifyContent:    "space-between",
		BackgroundColor:   "#374151",
	}
	styleTodoItemHover = Style{
		BackgroundOpacity: 0.5,
	}
	styleTodoItemBtn = Style{
		PaddingHorizontal: 1,
		BackgroundColor:   "#1f2937",
		BackgroundOpacity: 0.3,
	}
	styleTodoItemBtnHover = Style{
		BackgroundOpacity: 0.8,
	}
	styleTodoItemBtnText = Style{
		Color: "#9ca3af",
	}
)

func TodoApp() loom.Node {
	store := NewTodoStore()
	store.Toggle(store.Add("learn loom"))
	store.Toggle(store.Add("build a todo list app"))
	store.Add("profit")

	return TodoListContext.Provider(store, func() loom.Node {
		return Box(
			TodoForm(),
			TodoList(),

			Apply(styleTodoApp),
		)
	})
}

func TodoForm() loom.Node {
	list := todos()

	var input term.InputElement
	onRef := func(el term.InputElement) {
		input = el
		input.Focus()
	}

	onSubmit := func(e *term.EventSubmit) {
		value := strings.TrimSpace(e.Value)
		if value != "" {
			list.Add(value)
			input.Clear()
		}
	}

	submit := func(*term.EventMouse) {
		input.Submit()
	}

	return Box(
		Input(Apply(
			Ref{Fn: onRef},
			Attr{Placeholder: examples[rand.Intn(len(examples))]},
			styleTodoInput,
		)),
		Box(
			Text(">", Apply(styleTodoItemBtnText)),
			Apply(On{Click: submit}, styleTodoItemBtn),
			ApplyOn("hover", styleTodoItemBtnHover),
		),

		Apply(On{Submit: onSubmit}, styleTodoForm),
	)
}

func TodoList() loom.Node {
	list := todos()

	return ScrollBox(
		For(list.All, func(todo *Todo, index Accessor[int]) loom.Node {
			return TodoItem(todo)
		}),

		Apply(styleTodoList),
	)
}

func TodoItem(todo *Todo) loom.Node {
	list := todos()

	toggle := func(*term.EventMouse) {
		list.Toggle(todo.id)
	}
	remove := func(e *term.EventMouse) {
		list.Remove(todo.id)
	}

	checkbox := func() string {
		if todo.Done() {
			return "[*]"
		}
		return "[ ]"
	}

	return Box(
		P(
			BindText(checkbox, Apply(styleTodoItemBtnText)),
			Text(" "),
			Text(todo.title),
		),

		Box(
			Text("X", Apply(styleTodoItemBtnText)),
			Apply(On{Click: remove}, styleTodoItemBtn),
			ApplyOn("hover", styleTodoItemBtnHover),
		),

		Apply(On{Click: toggle}, styleTodoItem),
		ApplyOn("hover", styleTodoItemHover),
	)
}
