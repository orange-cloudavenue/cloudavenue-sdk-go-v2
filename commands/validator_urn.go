/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commands

import "fmt"

type validatorURN struct {
	urnFormat string
}

func ValidatorURN(format string) Validator {
	return &validatorURN{urnFormat: format}
}

func (v *validatorURN) GetKey() string {
	return fmt.Sprintf("urn=%s", v.urnFormat)
}

func (v *validatorURN) GetDescription() string {
	return fmt.Sprintf("Validates that the value is a valid URN (%s).", v.urnFormat)
}

func (v *validatorURN) GetMarkdownDescription() string {
	return fmt.Sprintf("Validates that the value is a valid URN (`%s`).", v.urnFormat)
}
