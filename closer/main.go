package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func doWork(cl *Closer) {
	defer cl.Done()

	i := 0
	for {
		select {
		case <-cl.HasBeenClosed():
			fmt.Println("Exiting doWork. Bye Bye!")
			return
		default:
			// Simulate some work.
			time.Sleep(time.Second)
			fmt.Printf("Processed Items = %+v\n", i)
			i++
		}
	}
}

func main() {
	cl := NewCloser(1)

	// Gracefully exit on CTRL+C.
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT)
	go func() {
		<-sigCh
		// Signal the goroutine to stop.
		cl.Signal()
	}()

	go doWork(cl)

	// Wait for the goroutines to finish.
	cl.Wait()
}
