package commands

import (
	"strings"
)

type validatorOneOf struct {
	values []string
}

func ValidatorOneOf(values ...string) Validator {
	return &validatorOneOf{values: values}
}

func (v *validatorOneOf) GetKey() string {
	// format values to oneof=value1 value2 value3
	return "oneof=" + strings.Join(v.values, " ")
}

func (v *validatorOneOf) GetDescription() string {
	return "Validates that the value is one of: " + strings.Join(v.values, ", ")
}

func (v *validatorOneOf) GetMarkdownDescription() string {
	return "Validates that the value is one of: " + strings.Join(wrapBackquoteEach(v.values), ", ")
}
