package loadtest

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

// Setup describes the setup provided by the command line / REST API
type Tester struct {
	// Rate of messages per second
	Rate int
	// Database interface
	Database Database
	// Log writer
	Log io.Writer
	// Number of messages to send. Set to 0 for infinite
	Number int
}

// Start starts the load testing run
func (t Tester) Start(ctx context.Context) {
	fmt.Fprintln(t.Log, "Starting...")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		select {
		case <-time.After(time.Second * 100):
		case <-ctx.Done():
		}
		wg.Done()
	}()
	wg.Wait()
}

// Start stops the load testing run
func (t Tester) Stop() {
	fmt.Fprintln(t.Log, "Stopping...")
}
