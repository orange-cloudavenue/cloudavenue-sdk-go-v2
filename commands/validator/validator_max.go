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

type validatorMax struct {
	max int
}

func ValidatorMax(maxValue int) Validator {
	return &validatorMax{max: maxValue}
}

func (v *validatorMax) GetKey() string {
	return fmt.Sprintf("max=%d", v.max)
}

func (v *validatorMax) GetDescription() string {
	return fmt.Sprintf("Ensures that the input value does not exceed a maximum length of %d characters", v.max)
}

func (v *validatorMax) GetMarkdownDescription() string {
	return fmt.Sprintf("Ensures that the input value does not exceed a maximum length of `%d` characters", v.max)
}
