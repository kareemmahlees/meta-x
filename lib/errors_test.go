package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseError400(t *testing.T) {
	respErr := ResponseError400("error 400")

	val, ok := respErr["status"]
	assert.True(t, ok)
	assert.Equal(t, val, 400)

	assert.Equal(t, "error 400", respErr["error"])
}

func TestResponseError500(t *testing.T) {
	respErr := ResponseError500("error 500")

	val, ok := respErr["status"]
	assert.True(t, ok)
	assert.Equal(t, val, 500)

	assert.Equal(t, "error 500", respErr["error"])
}
