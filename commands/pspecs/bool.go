package pspecs

import (
	"reflect"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/validator"
)

var _ ParamSpec = (*Bool)(nil)

type Bool struct {
	Name        string
	Description string
	Required    bool
	Example     any
	Validators  []validator.Validator

	paramSpecNotation string
}

func (s Bool) GetName() string {
	return s.Name
}

func (s *Bool) SetName(name string) {
	s.Name = name
}

func (s Bool) GetParamSpecNotation() string {
	if s.paramSpecNotation != "" {
		return s.paramSpecNotation
	}
	return s.Name
}

func (s *Bool) SetParamSpecNotation(notation string) {
	s.paramSpecNotation = notation
}

func (s Bool) GetDescription() string {
	return s.Description
}

func (s Bool) IsRequired() bool {
	return s.Required
}

func (s Bool) GetExample() any {
	return s.Example
}

func (s Bool) GetValidators() []validator.Validator {
	return s.Validators
}

func (s Bool) GetType() reflect.Value {
	return reflect.ValueOf(false)
}
