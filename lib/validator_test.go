package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStruct(t *testing.T) {
	testStructs := []struct {
		Name       any `validate:"required,alpha,min=5"`
		Valid      any `validate:"omitempty,boolean"`
		FailureTag string
	}{
		{Name: "abc", Valid: true, FailureTag: "min"},       // should fail on min
		{Name: 123, Valid: true, FailureTag: "alpha"},       // should fail on alpha
		{Valid: true, FailureTag: "required"},               // should fail on required
		{Name: "kareem", Valid: 123, FailureTag: "boolean"}, // should fail on bool
		{Name: "kareem"}, // shouldn't fail
	}
	for idx, test := range testStructs {
		errs := ValidateStruct(test)
		if idx == 4 {
			assert.Zero(t, len(errs))
			continue
		}
		assert.Greater(t, len(errs), 0)
		var failedTags []string
		for _, err := range errs {
			failedTags = append(failedTags, err.Tag)
		}
		assert.Contains(t, failedTags, test.FailureTag)
	}
}
