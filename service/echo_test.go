package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcho_HappyPath_Success(t *testing.T) {
	output := echo("Johny")

	assert.Equal(t, "Hello Johny!", output)
}
