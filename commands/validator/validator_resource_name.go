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
	"strings"

	"github.com/orange-cloudavenue/common-go/validators"
)

type validatorResourceName struct {
	rn string
}

// ValidatorResourceName creates a validator for resource names based on a specific format.
// List of supported formats can be found in the regex package. (validators.ListCavResourceNames)
func ValidatorResourceName(resourceName string) Validator {
	return &validatorResourceName{rn: resourceName}
}

func (v *validatorResourceName) GetKey() string {
	return fmt.Sprintf("resource_name=%s", v.rn)
}

func (v *validatorResourceName) GetDescription() string {
	for _, resource := range validators.ListCavResourceNames {
		if strings.EqualFold(resource.Key, v.rn) {
			return fmt.Sprintf("Validates that the value is a valid %s.", resource.Description)
		}
	}

	return fmt.Sprintf("Validates that the value is a valid resource name (%s).", v.rn)
}

func (v *validatorResourceName) GetMarkdownDescription() string {
	for _, resource := range validators.ListCavResourceNames {
		if strings.EqualFold(resource.Key, v.rn) {
			return fmt.Sprintf("Validates that the value is a valid %s.", resource.Description)
		}
	}

	return fmt.Sprintf("Validates that the value is a valid resource name (%s).", v.rn)
}
