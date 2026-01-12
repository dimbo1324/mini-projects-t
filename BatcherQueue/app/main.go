package main

import (
	"fmt"
	"time"

	"github.com/moguchev/stepik/4/4.6/HW/batcher"
)

func main() {
	handler := func(batch []string) {
		fmt.Println(">>> Flushed:", batch)
		time.Sleep(500 * time.Millisecond)
	}
	b := batcher.NewBatcher(5, 2*time.Second, handler)
	for i := 1; i <= 12; i++ {
		b.Add(fmt.Sprintf("event-%d", i))
		time.Sleep(300 * time.Millisecond)
	}
	b.Close()
	fmt.Println("Batcher gracefully shut down.")
}
