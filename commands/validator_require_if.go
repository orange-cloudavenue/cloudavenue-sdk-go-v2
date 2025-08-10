package commands

import "strings"

type requireIfParamIsNull struct {
	paramName []string
}

func ValidatorRequiredIfParamIsNull(paramName ...string) Validator {
	return &requireIfParamIsNull{paramName: paramName}
}

func (v *requireIfParamIsNull) GetKey() string {
	return "required_if_null=" + v.formatParamName(" ", "")
}

func (v *requireIfParamIsNull) GetDescription() string {
	if len(v.paramName) > 1 {
		return "The value is required if one of the parameters " + v.formatParamName(", ", "") + " is null."
	}

	return "The value is required if the parameter " + v.formatParamName(", ", "") + " is null."
}

func (v *requireIfParamIsNull) GetMarkdownDescription() string {
	if len(v.paramName) > 1 {
		return "The value is required if one of the parameters " + v.formatParamName(", ", "`") + " is null."
	}

	return "The value is required if the parameter " + v.formatParamName(", ", "`") + " is null."
}

func (v *requireIfParamIsNull) formatParamName(delimiter, charAroundParamName string) string {
	if len(v.paramName) == 0 {
		return ""
	}
	if len(v.paramName) == 1 {
		return charAroundParamName + v.paramName[0] + charAroundParamName
	}
	return charAroundParamName + strings.Join(v.paramName, charAroundParamName+delimiter+charAroundParamName) + charAroundParamName
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
