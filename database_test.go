package loadtest_test

import (
	"context"
	"math/rand"
	"time"
)

// mockDatabase
type mockDatabase struct {
	minResponseTime int
	maxResponseTime int
	errorPercentage int // TODO
}

// Send sends the datapoint to the mock database, and returns via the finished channel
func (m mockDatabase) Send(ctx context.Context, datapoint interface{}, finished chan error) {
	go func() {
		// Note use of rand.Intn without rand.Seed() in order to make
		ms := m.minResponseTime + rand.Intn(m.maxResponseTime)
		<-time.After(time.Millisecond * time.Duration(ms))
		finished <- nil
	}()
}
