package main

import (
	"fmt"
	"time"
)

// use defer to close channel properly, prevent channel keeps running and causes memory leak. Ensures resources are cleaned up.
func main() {
	doneCh := make(chan string)
	go worker(doneCh)

	// Wait for the goroutine to finish
	result := <-doneCh
	fmt.Println("Result from worker:", result)
}

func worker(doneCh chan string) {
	defer close(doneCh)
	fmt.Println("Starting work...")
	time.Sleep(2 * time.Second)
	doneCh <- "Work done!"
}
