package service

import (
	"fmt"
)

func echo(in string) string {
	if in == "" {
		in = "World"
	}

	return fmt.Sprintf("Hello %s!", in)
}
