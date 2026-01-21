package types

type RenderType int

const (
	RenderTypeInline RenderType = iota
	RenderTypeFullscreen
)

type Renderer interface {
	Type() RenderType
	RenderOffset() int

	Render(element Element) error

	CheckHit(x, y int) Element
	AddToHitGrid(element Element, x, y, width, height int)
	PushHitGridScissors(x, y, width, height int)
	PopHitGridScissors()

	UpdateCapabilities(buf []byte) error

	Close() error
}
