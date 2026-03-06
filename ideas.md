New agnostic loom Apply:

- ApplyOn would jsut be a renderer thing and wrap Apply in a Show (e.g. `Show(cond, Apply(...))`)

```go
func App() loom.Node {
  var input term.InputElement

  onSubmit := func(evt *term.EventSubmit) {
    fmt.Println("Submited:", evt.Value)
    input.Clear()
  }

  return Box(
    Text("This is an input:"),
    Input(Apply(
        Ref{Ptr: &input},
        Attr{Placeholder: "an input..."},
        Style{PlaceholderFontStyle: "italic"},
    )),

    Apply(
      On{Submit: onSubmit},
      Attr{Title: "a box"},
      Style{Width: 10, Height: 10, BackgroundColor: "red"},
    ),
  )
}
```

```go
func App() loom.Node {
  var input term.InputElement

  onSubmit := func(evt *term.EventSubmit) {
    fmt.Println("Submited:", evt.Value)
    input.Clear()
  }

  return Box(
    Text("This is an input:"),
    Input(
      Ref{Ptr: &input},
      Attr{Placeholder: "an input..."},
      Style{PlaceholderFontStyle: "italic"},
    ),

    On{Submit: onSubmit},
    Attr{Title: "a box"},
    Style{Width: 10, Height: 10, BackgroundColor: "red"},
  )
}
```

```go
func App() loom.Node {
  var input term.InputElement

  onSubmit := func(evt *term.EventSubmit) {
    fmt.Println("Submited:", evt.Value)
    input.Clear()
  }

  return Box(
    Text("This is an input:"),
    Input(
      Ref(&input),
      OnRef(onRef),
      Attr("placeholder", "an input..."),
      Apply(Style{PlaceholderFontStyle: "italic"}),
    ),

    Attr("title", "a box"),
    On("submit", onSubmit),
    Apply(Style{Width: 10, Height: 10, BackgroundColor: "red"}),
  )
}
```

```go
type Ref struct {
	Ptr any // *term.Element
	Fn  any // func(term.Element)
}

func (r Ref) Apply(e any) error {
	// should be able to do all that without reflection and with compile-time type safety
	// once this is added: https://github.com/golang/go/issues/61731

	eType := reflect.TypeOf(e)

	if r.Ptr != nil {
		ptrType := reflect.TypeOf(r.Ptr)
		// Ptr must be a pointer whose elem matches e's type
		if ptrType.Kind() != reflect.Pointer || ptrType.Elem() != eType {
			return fmt.Errorf("ref ptr type mismatch: expected *%s, got %s", eType, ptrType)
		}
		reflect.ValueOf(r.Ptr).Elem().Set(reflect.ValueOf(e))
	}

	if r.Fn != nil {
		fnType := reflect.TypeOf(r.Fn)
		// Fn must be a func with exactly one param matching e's type
		if fnType.Kind() != reflect.Func || fnType.NumIn() != 1 || fnType.In(0) != eType {
			return fmt.Errorf("ref fn type mismatch: expected func(%s), got %s", eType, fnType)
		}
		reflect.ValueOf(r.Fn).Call([]reflect.Value{reflect.ValueOf(e)})
	}

	return nil
}

type BoxElement struct{ v any }
type TextElement struct{ v any }

func main() {
	var box BoxElement
	fmt.Println(Ref{Ptr: &box}.Apply(BoxElement{"works"}))
	fmt.Println(Ref{Ptr: &box}.Apply(TextElement{"doesnt"}))

	var text TextElement
	onRef := func(t TextElement) { text = t }
	fmt.Println(Ref{Fn: onRef}.Apply(TextElement{"works"}))
	fmt.Println(Ref{Fn: onRef}.Apply(BoxElement{"doesnt"}))

	fmt.Printf("box: %+v\n", box)
	fmt.Printf("text: %+v\n", text)
}
```
