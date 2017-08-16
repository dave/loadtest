package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"sync"

	"strconv"

	"github.com/dave/loadtest"
	"github.com/dave/loadtest/mockdb"
)

func main() {

	// Ensure we stop before we start again
	var starter sync.Once

	// Default options
	t := loadtest.Tester{
		Rate: 1000,
		//Database: loadtest.InfluxDbDatabase{},
		Database: mockdb.MockDatabase{
			MinResponseTime: 100,
			MaxResponseTime: 200,
			ErrorPercentage: 5,
		},
		Log:    os.Stdout,
		Number: 0,
	}

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		var success bool
		starter.Do(func() {
			success = true
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			t.Cancel = cancel
			go func() {
				t.Start(ctx)
			}()
		})
		if !success {
			http.Error(w, "Already started. Call /stop to reset...", 500)
		}
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		t.Stop()
		starter = sync.Once{}
	})

	http.HandleFunc("/rate", func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query()["rate"][0]
		rate, err := strconv.Atoi(value)
		if err != nil {
			http.Error(w, "Error converting rate to int", 500)
		}
		t.Stop()
		t.SetRate(rate)
		starter = sync.Once{}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
