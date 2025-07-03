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
	endpointRegistryOptions struct {
		// API is the API name.
		api API
		// Version is the API version.
		version Version
	}

	EndpointRegistryOptions func(*endpointRegistryOptions)
)

// WithExtraProperties allows adding extra properties to the endpoint registry options.
func WithExtraProperties(api API, version Version) EndpointRegistryOptions {
	return func(opts *endpointRegistryOptions) {
		opts.api = api
		opts.version = version
	}
}
