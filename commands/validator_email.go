/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commands

type validatorEmail struct{}

// EmailValidator valide qu'une chaîne est une adresse email valide.
func ValidatorEmail() Validator {
	return &validatorEmail{}
}

func (v *validatorEmail) GetKey() string {
	return "email"
}

func (v *validatorEmail) GetDescription() string {
	return "Validates that the value is a valid email address."
}

func (v *validatorEmail) GetMarkdownDescription() string {
	return "Validates that the value is a valid email address. (e.g., `foo@bar.com`)"
}
