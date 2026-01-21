package main

import "sync"

// rendering
type RenderContext struct {
	renderMu sync.Mutex
}

type Renderer struct {
	ctx *RenderContext
}

func NewRenderer() *Renderer {
	return &Renderer{&RenderContext{}}
}

func (r *Renderer) Render(eleme *TextElement) {
	r.ctx.renderMu.Lock()
	defer r.ctx.renderMu.Unlock()

	// ...
}

// scheduling and values
type BufferedValue[T any] struct {
	ctx   *RenderContext
	value T
}

func (s *BufferedValue[T]) Update(fn func(T) T) {
	s.ctx.renderMu.Lock()
	s.value = fn(s.value)
	s.ctx.renderMu.Unlock()
	app.RequestRender()
}

// elements
type TextElement struct {
	content BufferedValue[string]
}

func NewTextElement(ctx *RenderContext) *TextElement {
	return &TextElement{
		content: BufferedValue[string]{ctx: ctx},
	}
}

func (e *TextElement) SetContent(content string) {
	e.content.Update(func(_ string) string {
		return content
	})
}

// app orchestration
var app = &App{}

type App struct {
	renderer *Renderer
	root     *TextElement
}

func (a *App) RequestRender() {
	a.renderer.Render(a.root)
}
