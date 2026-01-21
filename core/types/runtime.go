package types

type RenderContext interface {
	LockRender()
	TryLockRender() bool
	UnlockRender()

	Render() error

	AddToHitGrid(element Element, x, y, width, height int)
	PushHitGridScissors(x, y, width, height int)
	PopHitGridScissors()
}

type Runtime interface {
	Root() Element

	RenderContext() RenderContext
	Render() error
	RenderUnsafe() error

	Errors() (errors <-chan error)
	Events() (events <-chan Event, stop func())

	Close() error
}
