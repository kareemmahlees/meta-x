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

	singleJsonBodyRes := DecodeBody[map[string]any](testBody)

	val, ok := singleJsonBodyRes["name"]
	assert.True(t, ok)
	assert.Equal(t, val, "foo")

	val, ok = singleJsonBodyRes["age"]
	assert.True(t, ok)
	assert.Equal(t, int(val.(float64)), 123)

	testBody = io.NopCloser(bytes.NewBufferString(`
	[ 
		{
			"name":"foo",
			"age" : 123
		},
		{
			"name":"bar",
			"age" : 123
		}
	 ]`))

	listOfJsonRes := DecodeBody[[]map[string]any](testBody)

	val, ok = listOfJsonRes[1]["name"]
	assert.True(t, ok)
	assert.Equal(t, val, "bar")

	val, ok = singleJsonBodyRes["age"]
	assert.True(t, ok)
	assert.Equal(t, int(val.(float64)), 123)
}
