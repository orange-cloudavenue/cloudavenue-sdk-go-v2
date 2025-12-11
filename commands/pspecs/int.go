package pspecs

import (
	"reflect"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/validator"
)

var _ ParamSpec = (*Int)(nil)

type Int struct {
	Name        string
	Description string
	Required    bool
	Example     any
	Validators  []validator.Validator

	paramSpecNotation string
}

func (s Int) GetName() string {
	return s.Name
}

func (s *Int) SetName(name string) {
	s.Name = name
}

func (s Int) GetParamSpecNotation() string {
	if s.paramSpecNotation != "" {
		return s.paramSpecNotation
	}
	return s.Name
}

func (s *Int) SetParamSpecNotation(notation string) {
	s.paramSpecNotation = notation
}

func (s Int) GetDescription() string {
	return s.Description
}

func (s Int) IsRequired() bool {
	return s.Required
}

func (s Int) GetExample() any {
	return s.Example
}

func (s Int) GetValidators() []validator.Validator {
	return s.Validators
}

func (s Int) GetType() reflect.Value {
	return reflect.ValueOf(0)
}
