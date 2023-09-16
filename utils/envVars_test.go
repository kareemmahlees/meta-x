package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvVar(t *testing.T) {
	os.Setenv("TEST_VAR", "test_val")
	val, err := GetEnvVar("TEST_VAR", false)
	assert.Nil(t, err)
	assert.Equal(t, val, "test_val")

	val, err = GetEnvVar("DOESN'T_EXIST", false)
	assert.NotNil(t, err)
	assert.Equal(t, val, "")

	val ,err = GetEnvVar("SHOULDN'T_ERROR",true)
	assert.Nil(t,err)
	assert.Equal(t,val,"")
}
