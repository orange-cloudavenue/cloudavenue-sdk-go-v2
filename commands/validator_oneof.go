/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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
