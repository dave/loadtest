package main

import (
	"os"

	"context"

	"os/signal"

	"flag"

	"log"

	"github.com/dave/loadtest"
	"github.com/dave/loadtest/mockdb"
)

func main() {

	number := flag.Int("number", 0, "number of messages to send. Set to 0 for infinite.")
	rate := flag.Int("rate", 1000, "rate in message per second.")
	min := flag.Int("min", 100, "minimum latency of mock database in milliseconds")
	max := flag.Int("max", 200, "maximum latency of mock database in milliseconds")
	errpercent := flag.Int("err", 0, "error percentage of mock database")

	flag.Parse()

	if *min > *max {
		log.Fatal("Minimum latency must be <= maximum")
	}
	if *min < 0 {
		log.Fatal("Minimum latency must be >= 0")
	}
	if *errpercent > 100 || *errpercent < 0 {
		log.Fatal("Error percentage must be between 0 and 100")
	}

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
		Rate: *rate,
		//Database: loadtest.InfluxDbDatabase{},
		Database: mockdb.MockDatabase{
			MinResponseTime: *min,
			MaxResponseTime: *max,
			ErrorPercentage: *errpercent,
		},
		Log:    os.Stdout,
		Number: *number,
		Cancel: cancel,
	}
	tester.Start(ctx)

}
