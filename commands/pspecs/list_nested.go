package pspecs

import (
	"fmt"
	"reflect"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/validator"
)

var (
	_ ParamSpec       = (*ListNested)(nil)
	_ ParamSpecNested = (*ListNested)(nil)
)

type ListNested struct {
	Name        string
	Description string
	Required    bool
	Example     string
	Validators  []validator.Validator
	ItemsSpec   []ParamSpec

	paramSpecNotation string
}

func (s ListNested) GetName() string {
	return s.Name
}

func (s *ListNested) SetName(name string) {
	s.Name = name
}

func (s *ListNested) GetParamSpecNotation() string {
	if s.paramSpecNotation != "" {
		return s.paramSpecNotation
	}
	return s.Name + ".{index}"
}

func (s *ListNested) SetParamSpecNotation(notation string) {
	s.paramSpecNotation = notation
}

func (s ListNested) GetDescription() string {
	return s.Description
}

func (s ListNested) IsRequired() bool {
	return s.Required
}

func (s ListNested) GetExample() any {
	return s.Example
}

func (s ListNested) GetValidators() []validator.Validator {
	return s.Validators
}

func (s ListNested) GetItemsSpec() []ParamSpec {
	var items []ParamSpec
	for i := range s.ItemsSpec {
		item := s.ItemsSpec[i]
		item.SetParamSpecNotation(fmt.Sprintf("%s.{index}.%s", s.GetName(), item.GetName()))
		items = append(items, item)
	}
	return items
}

func (s ListNested) GetType() reflect.Value {
	return reflect.ValueOf([]any{})
}
