/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commands

type validatorOmitempty struct{}

func ValidatorOmitempty() Validator {
	return &validatorOmitempty{}
}

func (v *validatorOmitempty) GetKey() string {
	return "omitempty"
}

func (v *validatorOmitempty) GetDescription() string {
	return ""
}

func (v *validatorOmitempty) GetMarkdownDescription() string {
	return ""
}
