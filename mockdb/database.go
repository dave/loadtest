package mockdb

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

// MockDatabase is a mock database to stand in for InfluxDB during development
type MockDatabase struct {
	MinResponseTime int
	MaxResponseTime int
	ErrorPercentage int
}

// Send sends the datapoint to the mock database, and returns via the finished channel
func (m MockDatabase) Send(ctx context.Context, datapoint interface{}, finished chan error) {
	go func() {
		// Note use of rand.Intn without rand.Seed() ... better for testing?
		ms := m.MinResponseTime + rand.Intn(m.MaxResponseTime)
		<-time.After(time.Millisecond * time.Duration(ms))
		if rand.Intn(100) < m.ErrorPercentage {
			finished <- errors.New("mock error")
			return
		}
		finished <- nil
	}()
}
