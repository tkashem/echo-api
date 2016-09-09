package tests

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/rue-tkashem/echo-api/service"
)

var (
	endpoint string
	inproc   bool
)

func setup() {
	flag.StringVar(&endpoint, "endpoint", "http://localhost:3000", "target endpoint")
	flag.BoolVar(&inproc, "inproc", false, "whether you want to host the service in process")

	fmt.Println("endpoint: ", endpoint)
	fmt.Println("functional coverage: ", inproc)

	if inproc {
		go service.ListenAndServe()
	}
}

func shutdown() {
}

// TestMain has custom setup and shutdown
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()

	os.Exit(code)
}
