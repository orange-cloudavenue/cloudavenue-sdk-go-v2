package pspecs

import (
	"reflect"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/validator"
)

type (
	ParamSpec interface {
		GetName() string
		SetName(name string)
		GetParamSpecNotation() string
		SetParamSpecNotation(notation string)
		GetDescription() string
		IsRequired() bool
		GetExample() any
		GetValidators() []validator.Validator

		GetType() reflect.Value
	}

	ParamSpecNested interface {
		ParamSpec
		GetItemsSpec() []ParamSpec
	}

	Params []ParamSpec
)
