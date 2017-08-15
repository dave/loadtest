package loadtest_test

import (
	"bytes"
	"testing"

	"context"

	"sync"

	"github.com/dave/loadtest"
)

func TestStart(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	tester := loadtest.Tester{
		Rate: 10,
		Database: mockDatabase{
			minResponseTime: 10,
			maxResponseTime: 20,
		},
		Log: b,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		tester.Start(ctx)
		wg.Done()
	}()
	cancel()
	wg.Wait()
	if b.String() != "Starting...\nCancelled...\n" {
		t.Fatalf("Failed starting, got %#v", b.String())
	}
}

func TestStop(t *testing.T) {
	b := &bytes.Buffer{}
	tester := loadtest.Tester{
		Rate: 10,
		Database: mockDatabase{
			minResponseTime: 10,
			maxResponseTime: 20,
		},
		Log: b,
	}
	tester.Stop()
	if b.String() != "Stopping...\n" {
		t.Fatal("Failed stopping")
	}
}
