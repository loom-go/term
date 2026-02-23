package stdio

import (
	"io"
	"sync"
)

type Reader struct {
	input io.Reader
	mu    sync.RWMutex
	subs  []chan<- []byte
}

func NewReader(input io.Reader) *Reader {
	r := &Reader{input: input}

	go r.run()

	return r
}

func (r *Reader) Listen(bufSize int) <-chan []byte {
	ch := make(chan []byte, bufSize)
	r.add(ch)
	return ch
}

func (r *Reader) run() {
	buf := make([]byte, 1024)
	for {
		n, err := r.input.Read(buf)

		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			r.broadcast(data)
		}

		if err != nil {
			r.closeAll()
			return
		}
	}
}

func (r *Reader) broadcast(data []byte) {
	r.mu.RLock()
	subs := make([]chan<- []byte, len(r.subs))
	copy(subs, r.subs)
	r.mu.RUnlock()

	for _, ch := range subs {
		select {
		case ch <- data:
		default:
		}
	}
}

func (r *Reader) add(ch chan<- []byte) {
	r.mu.Lock()
	r.subs = append(r.subs, ch)
	r.mu.Unlock()
}

func (r *Reader) closeAll() {
	r.mu.Lock()
	for _, ch := range r.subs {
		close(ch)
	}
	r.subs = nil
	r.mu.Unlock()
}
