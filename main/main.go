package main

import (
	"flag"

	"github.com/tkashem/echo-api/service"
)

func main() {
	flag.Parse()

	service.ListenAndServe()
}
