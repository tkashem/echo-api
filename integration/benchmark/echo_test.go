package benchmark

import (
	"testing"

	"github.com/rue-tkashem/echo-api/integration/client"
)

func BenchmarkEcho(b *testing.B) {

	t := &testing.T{}

	// setup
	c := client.New(endpoint)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Echo(t, "ping")
		}
	})
}
