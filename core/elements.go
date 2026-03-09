package core

import "github.com/loom-go/term/core/elements"

type Element = elements.Element
type TextElement = *elements.TextElement

type BoxElement = *elements.BoxElement
type ScrollBoxElement = *elements.ScrollBoxElement

type InputElement = *elements.InputElement
type TextAreaElement = *elements.TextAreaElement

type ConsoleElement = *elements.ConsoleElement

func NewElement() (Element, error)         { return elements.NewBaseElement() }
func NewTextElement() (TextElement, error) { return elements.NewTextElement() }

func NewBoxElement() (BoxElement, error)             { return elements.NewBoxElement() }
func NewScrollBoxElement() (ScrollBoxElement, error) { return elements.NewScrollBoxElement() }

func NewInputElement() (InputElement, error)       { return elements.NewInputElement() }
func NewTextAreaElement() (TextAreaElement, error) { return elements.NewTextAreaElement() }

func NewConsoleElement() (ConsoleElement, error) { return elements.NewConsoleElement() }
