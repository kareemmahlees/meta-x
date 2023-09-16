package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStyle(t *testing.T) {
	style := NewStyle("Hello World", "#ffffff")
	assert.Equal(t, "Hello World", style.Value())
}
