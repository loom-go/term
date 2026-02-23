package gfx

import "github.com/AnatoleLucet/go-opentui"

type Renderable interface {
	Render(buffer *opentui.Buffer, rect Rect) (err error)
}

type Recordable interface {
	Record(cb *CommandBuffer, rect Rect) (err error)
}

type CommandType int

const (
	CmdRender CommandType = iota

	CmdPushOverflowScissors
	CmdPopOverflowScissors

	CmdPushHitGridScissors
	CmdPopHitGridScissors

	CmdPushOpacity
	CmdPopOpacity
)

type Command struct {
	Type CommandType

	Element Renderable

	Callback func() error

	// For CmdRender
	Rect Rect

	// For CmdPushOverflowScissors and CmdPushHitGridScissors
	Scissors Rect

	// For CmdPushOpacity
	Opacity float32
}

func NewCommand(typ CommandType, element Renderable) *Command {
	return &Command{
		Type:    typ,
		Element: element,
	}
}

func (c Command) WithCallback(callback func() error) *Command {
	c.Callback = callback
	return &c
}

func (c Command) WithRect(rect Rect) *Command {
	c.Rect = rect
	return &c
}

func (c Command) WithScissors(rect Rect) *Command {
	c.Scissors = rect
	return &c
}

func (c Command) WithOpacity(opacity float32) *Command {
	c.Opacity = opacity
	return &c
}
