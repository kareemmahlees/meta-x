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

func TestValidateVar(t *testing.T) {
	const validationTags = "required,alpha,min=3"

	err := ValidateVar("kareem", validationTags)
	assert.Nil(t, err)

	err = ValidateVar(123, validationTags)
	assert.Error(t, err)
}

func TestNotEmpty(t *testing.T) {
	testStructs := []struct {
		Names []string `validate:"notEmpty"`
	}{
		{Names: []string{"foo", "bar", "baz"}},
		{Names: []string{}},
	}
	for idx, test := range testStructs {
		errs := ValidateStruct(test)
		if idx == 0 {
			assert.Zero(t, len(errs))
		} else if idx == 1 {
			assert.Greater(t, len(errs), 0)
			assert.Equal(t, errs[0].Tag, "notEmpty")
		}
	}
}

func TestUpdateTableDataValidator(t *testing.T) {
	baseStruct := UpdateTableProps{}
	testData := []struct {
		Type     string
		Data     any
		willFail bool
	}{
		{Type: "add", Data: map[string]string{"foo": "bar"}, willFail: false},
		{Type: "add", Data: []string{"foo", "bar"}, willFail: true},
		{Type: "modify", Data: map[string]string{"foo": "bar"}, willFail: false},
		{Type: "modify", Data: []string{"foo", "bar"}, willFail: true},
		{Type: "delete", Data: []string{"foo", "bar"}, willFail: false},
		{Type: "delete", Data: map[string]string{"foo": "bar"}, willFail: true},
	}

	for _, test := range testData {
		baseStruct.Operation.Type = test.Type
		baseStruct.Operation.Data = test.Data

		errs := ValidateStruct(baseStruct)

		if test.willFail {
			assert.Greater(t, len(errs), 0)
		} else {
			assert.Zero(t, len(errs))
		}
	}
}
