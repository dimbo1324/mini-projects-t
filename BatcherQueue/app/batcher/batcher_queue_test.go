package batcher

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

//
// === BASIC FLUSH BEHAVIOR ===
//

func TestFlush_WhenBatchIsFull(t *testing.T) {
	t.Parallel()

	var (
		mu    sync.Mutex
		calls [][]int
	)

	handler := func(batch []int) {
		// Копируем содержимое, чтобы сохранить снимок батча
		mu.Lock()
		defer mu.Unlock()

		calls = append(calls, batch)
	}

	b := NewBatcher(3, 10*time.Second, handler)
	t.Cleanup(b.Close)

	// Добавляем 3 элемента → должен быть автоматический flush
	b.Add(1)
	b.Add(2)
	b.Add(3)

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(calls) != 1 {
		t.Fatalf("expected 1 flush, got %d", len(calls))
	}

	if want := []int{1, 2, 3}; fmt.Sprint(calls[0]) != fmt.Sprint(want) {
		t.Fatalf("unexpected batch: got %v, want %v", calls[0], want)
	}
}

func TestFlush_WhenTimeoutExpires(t *testing.T) {
	t.Parallel()

	var (
		mu      sync.Mutex
		calls   [][]string
		flushed = make(chan struct{}, 1)
	)

	handler := func(batch []string) {
		mu.Lock()
		defer mu.Unlock()

		calls = append(calls, batch)
		flushed <- struct{}{}
	}

	b := NewBatcher(5, 100*time.Millisecond, handler)
	t.Cleanup(b.Close)

	b.Add("a")
	b.Add("b")

	select {
	case <-flushed:
		// всё хорошо
	case <-time.After(500 * time.Millisecond):
		t.Fatal("flush by timeout did not happen")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(calls) != 1 {
		t.Fatalf("expected flush by timeout, got %d calls", len(calls))
	}
}

func TestBatcher_NeverExceedsCapacity(t *testing.T) {
	t.Parallel()

	var (
		mu    sync.Mutex
		calls [][]int
	)

	handler := func(batch []int) {
		mu.Lock()
		defer mu.Unlock()
		calls = append(calls, batch)
	}

	batchSize := 3
	b := NewBatcher(batchSize, 500*time.Millisecond, handler)
	t.Cleanup(b.Close)

	b.Add(1, 2, 3, 4, 5, 6, 7) // больше чем один батч
	time.Sleep(600 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(calls) != 3 {
		t.Fatalf("expected 3 flushes, got %d", len(calls))
	}

	// Проверяем, что батчи не превышают заданный размер
	for _, batch := range calls {
		if len(batch) > batchSize {
			t.Fatalf("batch exceeded capacity: %v", batch)
		}
	}
}

//
// === CORNER CASES ===
//

func TestFlush_NothingAdded(t *testing.T) {
	t.Parallel()

	// Проверяем, что обработчик не вызывается, если ничего не добавлялось
	called := false
	b := NewBatcher(3, 50*time.Millisecond, func(_ []string) {
		called = true
	})
	defer b.Close()

	time.Sleep(100 * time.Millisecond)

	if called {
		t.Fatalf("handler was called, but nothing was added")
	}
}

func TestFlush_ContextCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	flushed := false
	handler := func(batch []int) {
		flushed = true
	}
	b := NewBatcher(5, time.Second, handler)

	// Закрытие по таймеру
	go func() {
		<-ctx.Done()
		b.Close()
	}()

	time.Sleep(200 * time.Millisecond)

	if flushed {
		t.Fatalf("should not flush — expected to be canceled before timeout reached")
	}
}

//
// === HIGH LOAD ===
//

func TestFlush_MultipleBatches(t *testing.T) {
	t.Parallel()

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		total    int
		maxBatch int
	)

	const (
		tasks     = 100
		batchSize = 4
	)

	wg.Add(tasks)
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	handler := func(items []int) {
		defer wg.Add(-len(items))

		mu.Lock()
		total += len(items)
		if len(items) > maxBatch {
			maxBatch = len(items)
		}
		mu.Unlock()
	}

	b := NewBatcher(batchSize, 500*time.Millisecond, handler)
	t.Cleanup(b.Close)

	for i := range tasks {
		b.Add(i)
	}

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("flushes did not complete in time")
	}

	mu.Lock()
	defer mu.Unlock()
	if total != tasks {
		t.Fatalf("expected %d items total, got %d", tasks, total)
	}
	if maxBatch > batchSize {
		t.Fatalf("batch exceeded capacity: %d", maxBatch)
	}
}

//
// === GRACEFUL SHUTDOWN ===
//

func TestShutdown_FlushesRemainingItems(t *testing.T) {
	t.Parallel()

	var (
		called bool
		mu     sync.Mutex
		batch  []string
	)

	handler := func(items []string) {
		mu.Lock()
		defer mu.Unlock()
		called = true
		batch = append([]string(nil), items...)
	}

	b := NewBatcher(10, 5*time.Second, handler)

	b.Add("x")
	b.Add("y")

	time.Sleep(20 * time.Millisecond)
	b.Close() // Закрытие должно вызвать flush

	mu.Lock()
	defer mu.Unlock()
	if !called {
		t.Fatalf("handler not called on shutdown with remaining items")
	}
	if len(batch) != 2 {
		t.Fatalf("expected 2 items in final batch, got %v", batch)
	}
}

func TestShutdown_WaitsForHandler(t *testing.T) {
	t.Parallel()

	var (
		called bool
		mu     sync.Mutex
	)

	handler := func(items []int) {
		time.Sleep(100 * time.Millisecond) // Симулируем долгую обработку

		mu.Lock()
		called = true
		mu.Unlock()
	}

	b := NewBatcher(3, time.Second, handler)
	runtime.Gosched()

	b.Add(1)
	b.Add(2)
	b.Add(3)

	time.Sleep(10 * time.Millisecond)

	b.Close() // Должен дождаться окончания handler

	mu.Lock()
	defer mu.Unlock()
	if !called {
		t.Fatalf("handler didn't finish before shutdown")
	}
}
