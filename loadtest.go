package loadtest

import (
	"fmt"
	"io"
)

func Start(log io.Writer) {
	fmt.Fprintln(log, "Starting...")
}
