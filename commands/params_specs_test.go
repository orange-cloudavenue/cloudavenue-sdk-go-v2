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
	"errors"
	"testing"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/pspecs"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/validator"
)

type testBuildAndValidateDynamicStruct struct {
	name      string
	paramsDef pspecs.Params
	params    any
	expected  error
}

type (
	testBuildAndValidateDynamicStruct_user struct {
		Name string
		Age  int
	}

	testBuildAndValidateDynamicStruct_users struct {
		Users []testBuildAndValidateDynamicStruct_user
	}
)

func TestBuildAndValidateDynamicStruct(t *testing.T) {
	tests := []testBuildAndValidateDynamicStruct{
		{
			name: "valid params",
			paramsDef: pspecs.Params{
				&pspecs.String{
					Name: "username",
				},
				&pspecs.Int{
					Name: "age",
					Validators: []validator.Validator{
						validator.ValidatorBetween(18, 99),
					},
				},
			},
			params: map[string]any{
				"username": "bob",
				"age":      25,
			},
			expected: nil,
		},
		{
			name: "invalid params",
			paramsDef: pspecs.Params{
				&pspecs.String{
					Name: "username",
				},
				&pspecs.Int{
					Name:       "age",
					Validators: []validator.Validator{validator.ValidatorBetween(18, 99)},
				},
			},
			params: map[string]any{
				"username": "ab",
				"age":      15,
			},
			expected: errors.New("validation failed"),
		},
		{
			name: "valid nested params",
			paramsDef: pspecs.Params{
				&pspecs.ListNested{
					Name: "Users",
					ItemsSpec: []pspecs.ParamSpec{
						&pspecs.String{
							Name:       "Name",
							Validators: []validator.Validator{},
						},
						&pspecs.Int{
							Name: "Age",
							Validators: []validator.Validator{
								validator.ValidatorBetween(18, 99),
							},
						},
					},
					Validators: []validator.Validator{},
				},
			},
			params: testBuildAndValidateDynamicStruct_users{
				Users: []testBuildAndValidateDynamicStruct_user{
					{
						Name: "alice",
						Age:  30,
					},
					{
						Name: "bob",
						Age:  45,
					},
				},
			},
			expected: nil,
		},
		{
			name: "invalid nested params",
			paramsDef: pspecs.Params{
				&pspecs.ListNested{
					Name: "Users",
					ItemsSpec: []pspecs.ParamSpec{
						&pspecs.String{
							Name:       "Name",
							Required:   true,
							Validators: []validator.Validator{},
						},
						&pspecs.Int{
							Name: "Age",
							Validators: []validator.Validator{
								validator.ValidatorBetween(18, 99),
							},
						},
					},
					Validators: []validator.Validator{},
				},
			},
			params: testBuildAndValidateDynamicStruct_users{
				Users: []testBuildAndValidateDynamicStruct_user{
					{
						Name: "alice",
						Age:  30,
					},
					{
						Age: 45,
					},
				},
			},
			expected: errors.New("validation failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := buildAndValidateDynamicStruct(tt.paramsDef, tt.params)
			if tt.expected == nil && err != nil {
				t.Errorf("expected no error, got: %v", err)
			} else if tt.expected != nil && err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}
