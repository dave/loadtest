package loadtest

import (
	"context"
)

// Database provides an interface that can be used to mock the InfluxDB service
type Database interface {
	Send(ctx context.Context, datapoint interface{}, finished chan error)
}

// InfluxDbDatabase is a real InfluxDB database
type InfluxDbDatabase struct {
	// TODO
}

// Send sends the datapoint to the mock database, and returns via the finished channel
func (i InfluxDbDatabase) Send(ctx context.Context, datapoint interface{}, finished chan error) {
	finished <- nil // TODO
}
