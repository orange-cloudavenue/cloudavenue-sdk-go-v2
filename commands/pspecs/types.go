/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package pspecs

import (
	"reflect"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/validator"
)

type (
	ParamSpec interface {
		GetName() string
		SetName(name string)
		GetParamSpecNotation() string
		SetParamSpecNotation(notation string)
		GetDescription() string
		IsRequired() bool
		GetExample() any
		GetValidators() []validator.Validator

		GetType() reflect.Value
	}

	ParamSpecNested interface {
		ParamSpec
		GetItemsSpec() []ParamSpec
	}

	Params []ParamSpec
)
