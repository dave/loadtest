package main

import (
	"os"

	"github.com/dave/loadtest"
)

func main() {
	loadtest.Start(os.Stdout)
}
