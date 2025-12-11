/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package validator

import (
	"fmt"
)

type validatorBetween struct {
	min int
	max int
}

func ValidatorBetween(minValue, maxValue int) Validator {
	return &validatorBetween{min: minValue, max: maxValue}
}

func (v *validatorBetween) GetKey() string {
	return fmt.Sprintf("min=%d,max=%d", v.min, v.max)
}

func (v *validatorBetween) GetDescription() string {
	return fmt.Sprintf("Validates that the value is between %d and %d", v.min, v.max)
}

func (v *validatorBetween) GetMarkdownDescription() string {
	return fmt.Sprintf("Validates that the value is between `%d` and `%d`", v.min, v.max)
}
