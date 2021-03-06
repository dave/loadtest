package loadtest

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
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
	// Cancel stops the load testing run
	Cancel context.CancelFunc
}

func (t *Tester) SetRate(rate int) {
	t.Rate = rate
	fmt.Fprintf(t.Log, "Setting rate to %d\n", rate)
}

// Start starts the load testing run
func (t Tester) Start(ctx context.Context) {

	fmt.Fprintln(t.Log, "Starting...")

	// We can use these channels to exit with an error, or exit gracefully
	done := make(chan struct{})

	// Some counters so we can show a nice status. Remember, increment these with atomic.AddUint64 if concurrent
	var startedCount uint64 = 0
	var finishedCount uint64 = 0
	var errorCount uint64 = 0

	// We have to display the status in a few places, so make it a function
	status := func() {
		fmt.Fprintf(
			t.Log,
			"Started: %d, finished: %d, errors: %d\n",
			atomic.LoadUint64(&startedCount),
			atomic.LoadUint64(&finishedCount),
			atomic.LoadUint64(&errorCount),
		)
	}

	// Let's print the status every second
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			status()
			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	// Here we start a goroutine to send the messages
	go func() {

		defer close(done)

		// The count of sent messages
		var count int

		// Create a waitgroup so we can wait for outstanding messages to finish sending
		wg := &sync.WaitGroup{}

		ticker := time.NewTicker(time.Second / time.Duration(t.Rate))
		defer ticker.Stop()

		for {
			// Loop around until we have sent all messages (or forever if t.Number == 0)
			if t.Number > 0 && count >= t.Number {
				break
			}
			count++

			select {
			case <-ticker.C:

				// The database will signal on this channel when it's done processing the message
				finished := make(chan error)

				// We send the payload to the database
				t.Database.Send(ctx, "foo", finished)

				// We increment the started counter, and add to the waitgroup
				atomic.AddUint64(&startedCount, 1)
				wg.Add(1)

				// We wait for the database to finish, increment the counters and waitgroup
				go func() {
					defer wg.Done()
					select {
					case err := <-finished:
						atomic.AddUint64(&finishedCount, 1)
						if err != nil {
							atomic.AddUint64(&errorCount, 1)
						}
					case <-ctx.Done():
						return
					}
				}()

			case <-ctx.Done():
				return
			}
		}

		wg.Wait()
	}()

	select {
	case <-done:
		status()
		fmt.Fprintln(t.Log, "Finished...")
	case <-ctx.Done():
		fmt.Fprintln(t.Log, "Done...")
	}
}

// Start stops the load testing run
func (t Tester) Stop() {
	if t.Cancel == nil {
		// if not started
		return
	}
	fmt.Fprintln(t.Log, "Stopping...")
	t.Cancel()
}

const DEBUG = false

func init() {
	if DEBUG {
		go func() {
			// debug to see if goroutines aren't being closed...
			ticker := time.NewTicker(time.Millisecond * 200)
			for range ticker.C {
				fmt.Println("runtime.NumGoroutine(): ", runtime.NumGoroutine())
			}
		}()
	}
}
