package batcher

import (
	"sync"
	"time"
)

type Handler[T any] func([]T)
type Batcher[T any] struct {
	capacity int
	interval time.Duration
	handler  Handler[T]
	mu       sync.Mutex
	buf      []T
	wake     chan struct{}
	closeCh  chan struct{}
	wg       sync.WaitGroup
	workerWg sync.WaitGroup
	closed   bool
}

func NewBatcher[T any](capacity int, interval time.Duration, handler Handler[T]) *Batcher[T] {
	if capacity <= 0 {
		panic("capacity must be > 0")
	}
	b := &Batcher[T]{
		capacity: capacity,
		interval: interval,
		handler:  handler,
		buf:      make([]T, 0, capacity),
		wake:     make(chan struct{}, 1),
		closeCh:  make(chan struct{}),
	}
	b.workerWg.Add(1)
	go b.run()
	return b
}
func (b *Batcher[T]) Add(items ...T) {
	if len(items) == 0 {
		return
	}
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return
	}
	b.buf = append(b.buf, items...)
	for len(b.buf) >= b.capacity {
		batch := make([]T, b.capacity)
		copy(batch, b.buf[:b.capacity])
		b.buf = b.buf[b.capacity:]
		b.wg.Add(1)
		go func(batch []T) {
			defer b.wg.Done()
			b.handler(batch)
		}(batch)
	}
	select {
	case b.wake <- struct{}{}:
	default:
	}
	b.mu.Unlock()
}
func (b *Batcher[T]) Close() {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		b.workerWg.Wait()
		b.wg.Wait()
		return
	}
	b.closed = true
	close(b.closeCh)
	select {
	case b.wake <- struct{}{}:
	default:
	}
	b.mu.Unlock()
	b.workerWg.Wait()
	b.wg.Wait()
}
func (b *Batcher[T]) run() {
	defer b.workerWg.Done()
	timer := time.NewTimer(b.interval)
	defer timer.Stop()
	for {
		b.mu.Lock()
		if b.closed && len(b.buf) == 0 {
			b.mu.Unlock()
			return
		}
		b.mu.Unlock()
		select {
		case <-b.wake:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(b.interval)
		case <-timer.C:
			b.mu.Lock()
			for len(b.buf) > 0 {
				n := b.capacity
				if len(b.buf) < n {
					n = len(b.buf)
				}
				batch := make([]T, n)
				copy(batch, b.buf[:n])
				b.buf = b.buf[n:]
				b.wg.Add(1)
				go func(batch []T) {
					defer b.wg.Done()
					b.handler(batch)
				}(batch)
			}
			b.mu.Unlock()
			timer.Reset(b.interval)
		case <-b.closeCh:
			b.mu.Lock()
			for len(b.buf) > 0 {
				n := b.capacity
				if len(b.buf) < n {
					n = len(b.buf)
				}
				batch := make([]T, n)
				copy(batch, b.buf[:n])
				b.buf = b.buf[n:]
				b.wg.Add(1)
				go func(batch []T) {
					defer b.wg.Done()
					b.handler(batch)
				}(batch)
			}
			b.mu.Unlock()
			return
		}
	}
}
