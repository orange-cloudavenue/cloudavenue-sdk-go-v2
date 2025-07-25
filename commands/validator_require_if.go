package commands

import "strings"

type requireIfParamIsNull struct {
	paramName string
}

func ValidatorRequiredIfParamIsNull(paramName string) Validator {
	return &requireIfParamIsNull{paramName: paramName}
}

func (v *requireIfParamIsNull) GetKey() string {
	return "required_if_null=" + v.paramName
}

func (v *requireIfParamIsNull) GetDescription() string {
	return "The value is required if the parameter '" + v.paramName + "' is null."
}

func (v *requireIfParamIsNull) GetMarkdownDescription() string {
	return "The value is required if the parameter `" + v.paramName + "` is null."
}

type requireIfParamIsOneOf struct {
	paramName string
	values    []string
}

func ValidatorRequiredIfParamIsOneOf(paramName string, values ...string) Validator {
	return &requireIfParamIsOneOf{paramName: paramName, values: values}
}

func (v *requireIfParamIsOneOf) GetKey() string {
	return "required_if_oneof=" + v.paramName + ":" + strings.Join(v.values, ",")
}

func (v *requireIfParamIsOneOf) GetDescription() string {
	return "The value is required if the parameter '" + v.paramName + "' is one of: " + strings.Join(v.values, ", ")
}

func (v *requireIfParamIsOneOf) GetMarkdownDescription() string {
	return "The value is required if the parameter `" + v.paramName + "` is one of: " + strings.Join(wrapBackquoteEach(v.values), ", ")
}
