package utils

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadBody(t *testing.T) {
	testBody := io.NopCloser(bytes.NewBufferString(`
	{
		"name":"foo",
		"age": 123
	}`))

	result := ReadBody[map[string]any](testBody)

	val, ok := result["name"]
	assert.True(t, ok)
	assert.Equal(t, val, "foo")

	val, ok = result["age"]
	assert.True(t, ok)
	assert.Equal(t, int(val.(float64)), 123)
}
