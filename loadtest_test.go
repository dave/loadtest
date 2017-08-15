package loadtest_test

import (
	"bytes"
	"testing"

	"github.com/dave/loadtest"
)

func TestStart(t *testing.T) {
	b := &bytes.Buffer{}
	tester := loadtest.Tester{
		Rate: 10,
		Database: mockDatabase{
			minResponseTime: 10,
			maxResponseTime: 20,
		},
		Log: b,
	}
	tester.Start()
	if b.String() != "Starting...\n" {
		t.Fatal("Failed starting")
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
