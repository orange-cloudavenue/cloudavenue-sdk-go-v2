/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type executeParams struct {
	Name string `validate:"required"`
}

type executeResult struct {
	Value string
}

func init() {
	Endpoint{
		Name:             "ExecuteOperationTest",
		Description:      "Test endpoint for typed operation execution",
		Method:           MethodGET,
		Backend:          BackendVMware,
		PathTemplate:     "/test/execute",
		PathParams:       []PathParam{},
		QueryParams:      []QueryParam{},
		DocumentationURL: "https://example.invalid/execute",
	}.Register()
}

func TestExecute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, err := newMockClient()
		require.NoError(t, err)
		defer func() { require.NoError(t, client.Close()) }()

		op := Operation[executeParams, executeResult]{
			Name:     "Organization.Get",
			Backend:  BackendVMware,
			Endpoint: MustGetEndpoint("ExecuteOperationTest"),
			Validate: func(p executeParams) error {
				if p.Name == "" {
					return fmt.Errorf("name required")
				}
				return nil
			},
			Transform: func(p executeParams) (any, error) {
				return nil, nil
			},
			Extract: func(resp *Response, p executeParams) (executeResult, error) {
				if resp == nil || resp.Raw == nil {
					return executeResult{}, fmt.Errorf("nil response")
				}
				return executeResult{Value: p.Name}, nil
			},
		}

		result, err := Execute(t.Context(), client, op, executeParams{Name: "ok"})
		require.NoError(t, err)
		require.Equal(t, "ok", result.Value)
	})

	t.Run("validate error wrapped", func(t *testing.T) {
		client, err := newMockClient()
		require.NoError(t, err)
		defer func() { require.NoError(t, client.Close()) }()

		op := Operation[executeParams, executeResult]{
			Name:     "Organization.Get",
			Backend:  BackendVMware,
			Endpoint: MustGetEndpoint("ExecuteOperationTest"),
			Validate: func(p executeParams) error {
				if p.Name == "" {
					return fmt.Errorf("boom")
				}
				return nil
			},
			Transform: func(p executeParams) (any, error) { return nil, nil },
			Extract: func(resp *Response, p executeParams) (executeResult, error) {
				return executeResult{}, nil
			},
		}

		_, err = Execute(t.Context(), client, op, executeParams{})
		require.EqualError(t, err, "Organization.Get: validate: boom")
	})
}

func TestPartialSuccessError(t *testing.T) {
	wrapped := &PartialSuccessError[string]{Result: "id", Err: fmt.Errorf("patch failed")}
	require.EqualError(t, wrapped, "patch failed")
	require.EqualError(t, wrapped.Unwrap(), "patch failed")
}
