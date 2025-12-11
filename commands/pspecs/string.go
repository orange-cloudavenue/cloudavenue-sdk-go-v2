package pspecs

import (
	"reflect"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/validator"
)

var _ ParamSpec = (*String)(nil)

type String struct {
	Name        string
	Description string
	Required    bool
	Example     any
	Validators  []validator.Validator

	paramSpecNotation string
}

func (s String) GetName() string {
	return s.Name
}

func (s *String) SetName(name string) {
	s.Name = name
}

func (s String) GetParamSpecNotation() string {
	if s.paramSpecNotation != "" {
		return s.paramSpecNotation
	}
	return s.Name
}

func (s *String) SetParamSpecNotation(notation string) {
	s.paramSpecNotation = notation
}

func (s String) GetDescription() string {
	return s.Description
}

func (s String) IsRequired() bool {
	return s.Required
}

func (s String) GetExample() any {
	return s.Example
}

func (s String) GetValidators() []validator.Validator {
	return s.Validators
}

func (s String) GetType() reflect.Value {
	return reflect.ValueOf("")
}
