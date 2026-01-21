package runtime

import (
	"sync"

	"github.com/AnatoleLucet/loom-term/core/types"
)

type RenderContext struct {
	renderMu sync.Mutex

	runtime *TermRuntime
}

func NewRenderContext(runtime *TermRuntime) *RenderContext {
	return &RenderContext{
		runtime: runtime,
	}
}

func (rc *RenderContext) LockRender() {
	rc.renderMu.Lock()
}

func (rc *RenderContext) TryLockRender() bool {
	return rc.renderMu.TryLock()
}

func (rc *RenderContext) UnlockRender() {
	rc.renderMu.Unlock()
}

func (rc *RenderContext) Render() error {
	return rc.runtime.Render()
}

func (rc *RenderContext) AddToHitGrid(element types.Element, x, y, width, height int) {
	rc.runtime.addToHitGrid(element, x, y, width, height)
}

func (rc *RenderContext) PushHitGridScissors(x, y, width, height int) {
	rc.runtime.pushHitGridScissors(x, y, width, height)
}

func (rc *RenderContext) PopHitGridScissors() {
	rc.runtime.popHitGridScissors()
}
