package tests

import (
	"testing"

	"github.com/rue-tkashem/echo-api/integration/client"

	"github.com/emicklei/forest"
)

func TestEcho_HappyPath_Success(t *testing.T) {
	c := client.New(endpoint)

	response := c.Echo(t, "Johny")

	forest.ExpectStatus(t, response, 200)
}
