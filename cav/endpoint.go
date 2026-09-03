/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

type (
	// API identifies API family.
	API string
	// Version identifies API version.
	Version string
	// Method identifies HTTP method.
	Method string

	// Endpoint describes backend endpoint metadata used by request execution.
	Endpoint struct {
		// ID uniquely identifies endpoint definition.
		ID string

		// Name is stable registry name.
		Name string `validate:"required,disallow_space,case=PascalCase"` // e.g., "Firewall", "LoadBalancer"

		// Description is human-readable operation summary.
		Description string `validate:"required"`

		// Backend identifies target backend.
		Backend BackendTarget `validate:"required"`

		// Method is HTTP verb used for requests.
		Method Method `validate:"required,oneof=GET POST PUT DELETE PATCH"`

		// PathTemplate is URL path template.
		PathTemplate string `validate:"required"` // e.g., "/v1/edgeGateways/{gatewayId}/firewall/rules"

		// PathParams describes supported path placeholders.
		PathParams []PathParam `validate:"dive"`

		// QueryParams describes supported query parameters.
		QueryParams []QueryParam `validate:"dive"`

		// DocumentationURL points to upstream endpoint documentation.
		DocumentationURL string `validate:"required,url"` // e.g., "https://docs.xx.com/api/v1/xx"

		// BodyRequestType defines expected request body type when applicable.
		BodyRequestType any `validate:"-"`

		// ResponseType defines decoded response body type when applicable.
		ResponseType any `validate:"-"`
	}

	// QueryParam describes one supported query parameter.
	QueryParam struct {
		Name        string `validate:"required,disallow_space"`
		Description string `validate:"required"`

		// Required reports whether parameter must be provided.
		Required bool

		// ValidatorFunc validates parameter value before request execution.
		ValidatorFunc func(value string) error

		// TransformFunc rewrites parameter value before request execution.
		TransformFunc func(value string) (string, error)

		// Value provides fixed parameter value set at registration time.
		Value string
	}

	// PathParam describes one supported path placeholder.
	PathParam struct {
		Name        string `validate:"required,disallow_space"`
		Description string `validate:"required"`

		// Required reports whether parameter must be provided.
		Required bool

		// ValidatorFunc validates parameter value before request execution.
		ValidatorFunc func(value string) error

		// TransformFunc rewrites parameter value before request execution.
		TransformFunc func(value string) (string, error)

		// Value provides fixed parameter value set at registration time.
		Value string
	}
)
