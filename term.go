package main

import (
	"errors"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/signals"
)

type App struct {
	running bool

	root  *loom.Slot
	owner *signals.Owner
}

func NewApp() *App {
	return &App{
		running: false,
		owner:   signals.NewOwner(),
	}
}

func (a *App) Run(fn func() loom.Node) <-chan any {
	errc := make(chan any, 1)
	if a.running {
		errc <- errors.New("app is already running")
		return errc
	}

	a.owner.OnError(func(err any) {
		errc <- err
	})

	err := a.Render(fn)
	if err != nil {
		errc <- err
		return errc
	}

	a.running = true
	return errc
}

func (a *App) Render(fn func() loom.Node) error {
	err := a.owner.Run(func() error {
		slot, err := loom.Render(nil, fn())
		a.root = slot
		return err
	})

	if err != nil {
		a.owner.Dispose()
		return err
	}

	return nil
}

func (a *App) Close() {
	if a.root != nil {
		a.root.Unmount()
	}

	a.owner.Dispose()
}
