package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"

	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom-term/core/types"
	"golang.org/x/term"
)

func main() {
	state, _ := term.MakeRaw(int(os.Stdin.Fd()))
	term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), state)

	// width, height := core.TerminalSize()
	width, _ := core.TerminalSize()

	rt, _ := core.NewRuntime(core.RenderInline)
	// rt, _ := core.NewRuntime(core.RenderFullscreen)
	rt.Root().SetWidth(width)
	rt.Root().SetHeight(40)
	// rt.Root().SetWidth(width)
	// rt.Root().SetHeight(height)

	console, _ := core.NewConsoleElement(rt.RenderContext())
	rt.Root().AppendChild(console)

	c1, _ := core.NewBoxElement(rt.RenderContext())
	c1.SetBackgroundColor("#00f")
	c1.SetWidth("50%")
	c1.SetHeight(4)
	rt.Root().AppendChild(c1)

	c1.OnMouseEnter(func(*types.EventMouse) {
		core.LogDebug("enter")
		c1.SetBackgroundColor("#0ff")
		rt.Render()
	})
	c1.OnMouseLeave(func(*types.EventMouse) {
		core.LogDebug("leave")
		c1.SetBackgroundColor("#00f")
		rt.Render()
	})

	c2, _ := core.NewBoxElement(rt.RenderContext())
	c2.SetBackgroundColor("#f00")
	c2.SetWidth("50%")
	c2.SetHeight(4)
	rt.Root().AppendChild(c2)

	d1 := draggable(rt.RenderContext())
	rt.Root().AppendChild(d1)

	d2 := draggable(rt.RenderContext())
	rt.Root().AppendChild(d2)

	d3 := draggable(rt.RenderContext())
	rt.Root().AppendChild(d3)

	rt.Render()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		evts, _ := rt.Events()

		for event := range evts {
			switch e := event.(type) {
			case *types.EventKey:
				core.LogDebug(fmt.Sprintf("Key event: %+v", e))

				if e.Ctrl && e.Rune == 'c' {
					rt.Close()
					return
				}
			}
		}
	}()

	wg.Wait()
}

func draggable(ctx types.RenderContext) core.Element {
	x := rand.Intn(50) + 10
	y := rand.Intn(20) + 10
	w := 12
	h := 5
	bg := fmt.Sprintf("#%06x", rand.Intn(0xffffff))

	e, _ := core.NewBoxElement(ctx)
	e.SetWidth(w)
	e.SetHeight(h)
	e.SetTop(y)
	e.SetLeft(x)
	e.SetPosition("absolute")
	e.SetBackgroundColor(bg)

	e.OnMousePress(func(evt *types.EventMouse) {
		e.SetBackgroundColor(bg + "88")
		ctx.Render()
	})
	e.OnMouseRelease(func(evt *types.EventMouse) {
		e.SetBackgroundColor(bg)
		ctx.Render()
	})

	e.OnMouseDrag(func(evt *types.EventMouse) {
		x = evt.X - w/2
		y = evt.Y - h/2
		e.SetTop(y)
		e.SetLeft(x)
		ctx.Render()
	})

	return e
}
