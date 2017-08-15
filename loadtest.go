package loadtest

import (
	"fmt"
	"io"
)

// Setup describes the setup provided by the command line / REST API
type Tester struct {
	// Rate of messages per second
	Rate int
	// Database interface
	Database Database
	// Log writer
	Log io.Writer
}

// Start starts the load testing run
func (t Tester) Start() {
	fmt.Fprintln(t.Log, "Starting...")
}

// Start stops the load testing run
func (t Tester) Stop() {
	fmt.Fprintln(t.Log, "Stopping...")
}
