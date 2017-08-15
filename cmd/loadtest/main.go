package main

import (
	"os"

	"context"

	"os/signal"

	"github.com/dave/loadtest"
	"github.com/dave/loadtest/mockdb"
)

func main() {

	// Make a ctx that will cancel with ctrl+c
	// Copied from: https://medium.com/@matryer/make-ctrl-c-cancel-the-context-context-bd006a8ad6ff

	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	tester := loadtest.Tester{
		Rate: 1000,
		//Database: loadtest.InfluxDbDatabase{},
		Database: mockdb.MockDatabase{
			MinResponseTime: 250,
			MaxResponseTime: 500,
			ErrorPercentage: 5,
		},
		Log:    os.Stdout,
		Number: 10000,
		Cancel: cancel,
	}
	tester.Start(ctx)

}
