package commands

import (
	"strconv"
)

type validatorBetween struct {
	min int
	max int
}

func ValidatorBetween(minValue, maxValue int) Validator {
	return &validatorBetween{min: minValue, max: maxValue}
}

func (v *validatorBetween) GetKey() string {
	// format values to between=value1 value2 value3
	return "min=" + strconv.Itoa(v.min) + " max=" + strconv.Itoa(v.max)
}

func (v *validatorBetween) GetDescription() string {
	return "Validates that the value is between " + strconv.Itoa(v.min) + " and " + strconv.Itoa(v.max)
}

func (v *validatorBetween) GetMarkdownDescription() string {
	return "Validates that the value is between " + strconv.Itoa(v.min) + " and " + strconv.Itoa(v.max)
}
