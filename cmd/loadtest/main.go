package main

import (
	"os"

	"github.com/dave/loadtest"
)

func main() {
	tester := loadtest.Tester{
		Rate:     10,
		Database: loadtest.InfluxDbDatabase{},
		Log:      os.Stdout,
	}
	tester.Start()
}
