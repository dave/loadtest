package loadtest_test

import (
	"bytes"
	"testing"

	"context"

	"sync"

	"strings"

	"github.com/dave/loadtest"
	"github.com/dave/loadtest/mockdb"
)

func TestStart(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	ctx, _ = context.WithCancel(ctx)
	tester := loadtest.Tester{
		Rate: 1000,
		Database: mockdb.MockDatabase{
			MinResponseTime: 10,
			MaxResponseTime: 20,
		},
		Log:    b,
		Number: 10,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		tester.Start(ctx)
		wg.Done()
	}()
	wg.Wait()
	if !strings.Contains(b.String(), "Started: 10, finished: 10, errors: 0\nFinished...") {
		t.Fatalf("Failed, got %#v", b.String())
	}
}

func TestStart_cancel(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	tester := loadtest.Tester{
		Rate: 1000,
		Database: mockdb.MockDatabase{
			MinResponseTime: 10,
			MaxResponseTime: 20,
		},
		Log: b,
		Cancel: cancel,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		tester.Start(ctx)
		wg.Done()
	}()
	tester.Stop()
	wg.Wait()
	if !strings.Contains(b.String(), "Done...") {
		t.Fatalf("Failed, got %#v", b.String())
	}
}
