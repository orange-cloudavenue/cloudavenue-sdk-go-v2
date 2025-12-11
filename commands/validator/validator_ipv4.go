/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package validator

type validatorIPV4 struct{}

func ValidatorIPV4() Validator {
	return &validatorIPV4{}
}

func (v *validatorIPV4) GetKey() string {
	return "ipv4"
}

func (v *validatorIPV4) GetDescription() string {
	return "Validates that the value is a valid IPv4 address."
}

func (v *validatorIPV4) GetMarkdownDescription() string {
	return "Validates that the value is a valid IPv4 address. (E.g. 192.168.1.1 or 203.0.113.0)"
}
