package loadtest_test

import (
	"bytes"
	"testing"

	"github.com/dave/loadtest"
)

func TestFoo(t *testing.T) {
	b := &bytes.Buffer{}
	loadtest.Start(b)
	if b.String() != "Starting...\n" {
		t.Fail()
	}
}
