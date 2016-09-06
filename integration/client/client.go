package client

import (
	"net/http"
	"testing"

	"github.com/emicklei/forest"
)

// New returns a new instance of the client ready to be used
func New(endpoint string) Invoker {
	return Invoker{url: endpoint}
}

// Invoker invokes account service operations
type Invoker struct {
	url string
}

// Echo creates an account
func (c Invoker) Echo(t *testing.T, message string) *http.Response {
	server := forest.NewClient(c.url, new(http.Client))
	resource := forest.Path("/echo/" + message)

	return server.GET(t, resource)
}
