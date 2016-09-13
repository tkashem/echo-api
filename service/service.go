package service

import (
	"flag"
	"net/http"

	"goji.io"
	"goji.io/pat"
)

var (
	bindTo string
)

func init() {
	flag.StringVar(&bindTo, "listen", ":3000", "host:port to bind to")
}

// ListenAndServe initializes and starts the service
func ListenAndServe() {
	handler := httpHandlerImpl{}
	mux := goji.NewMux()

	mux.HandleFuncC(pat.Get("/echo/:name"), handler.echo)

	http.ListenAndServe(bindTo, mux)
}
