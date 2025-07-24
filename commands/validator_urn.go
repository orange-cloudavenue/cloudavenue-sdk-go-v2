package commands

import "fmt"

type validatorURN struct {
	urnFormat string
}

func ValidatorURN(format string) Validator {
	return &validatorURN{urnFormat: format}
}

func (v *validatorURN) GetKey() string {
	return fmt.Sprintf("urn=%s", v.urnFormat)
}

func (v *validatorURN) GetDescription() string {
	return fmt.Sprintf("Validates that the value is a valid URN (%s).", v.urnFormat)
}

func (v *validatorURN) GetMarkdownDescription() string {
	return fmt.Sprintf("Validates that the value is a valid URN (`%s`).", v.urnFormat)
}
