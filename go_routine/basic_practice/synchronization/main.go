package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1) // Increments the counter for each goroutine
		go task(i, &wg)
	}

	fmt.Println("All tasks are arranged")
	wg.Wait() // Blocks until all goroutines finish
	fmt.Println("All tasks completed")
}

func task(id int, wg *sync.WaitGroup) {
	defer wg.Done() // Decrements the counter when the goroutine completes
	time.Sleep(2 * time.Second)
	fmt.Printf("Task %d is running\n", id)
}
