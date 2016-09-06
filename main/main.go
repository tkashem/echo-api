package main

import (
	"flag"

	"github.com/rue-tkashem/echo-api/service"
)

func main() {
	flag.Parse()

	service.ListenAndServe()
}
