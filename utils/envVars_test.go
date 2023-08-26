package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvVar(t *testing.T) {
	os.Setenv("TEST_VAR", "test_val")
	val,err := GetEnvVar("TEST_VAR")
	assert.Nil(t,err)
	assert.Equal(t,val,"test_val")

	val,err = GetEnvVar("DOESN'T_EXIST")
	assert.NotNil(t,err)
	assert.Equal(t,val,"")
}