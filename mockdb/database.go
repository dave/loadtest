package mockdb

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

// MockDatabase
type MockDatabase struct {
	MinResponseTime int
	MaxResponseTime int
	ErrorPercentage int
}

// Send sends the datapoint to the mock database, and returns via the finished channel
func (m MockDatabase) Send(ctx context.Context, datapoint interface{}, finished chan error) {
	go func() {
		// Note use of rand.Intn without rand.Seed() in order to make
		ms := m.MinResponseTime + rand.Intn(m.MaxResponseTime)
		<-time.After(time.Millisecond * time.Duration(ms))
		if rand.Intn(100) < m.ErrorPercentage {
			finished <- errors.New("mock error")
			return
		}
		finished <- nil
	}()
}
