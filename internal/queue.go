package internal

type renderRequest struct {
	value *Element
	done  chan error
}

type RenderQueue struct {
	ch chan *renderRequest
}

func NewRenderQueue(handler func(elem *Element) error) *RenderQueue {
	rq := &RenderQueue{
		ch: make(chan *renderRequest),
	}

	go func() {
		for req := range rq.ch {
			err := handler(req.value)
			req.done <- err
			close(req.done)
		}
	}()

	return rq
}

func (rq *RenderQueue) Enqueue(value *Element) error {
	done := make(chan error, 1)

	rq.ch <- &renderRequest{
		value: value,
		done:  done,
	}

	return <-done
}
