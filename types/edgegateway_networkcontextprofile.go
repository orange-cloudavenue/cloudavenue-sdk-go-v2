/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

type (
	ParamsListEdgeGatewayNetworkContextProfile struct {
		// EdgeGatewayID is the ID of the Edge Gateway whose owner scopes the Network Context Profiles.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owner scopes the Network Context Profiles.
		EdgeGatewayName string
	}

	ParamsGetEdgeGatewayNetworkContextProfile struct {
		// ID is the unique identifier of the Network Context Profile.
		ID string

		// Name is the name of the Network Context Profile.
		Name string

		// EdgeGatewayID is the ID of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayName string
	}

	ParamsCreateEdgeGatewayNetworkContextProfile struct {
		// Name is the name of the Network Context Profile.
		Name string

		// Description is the description of the Network Context Profile.
		Description string

		// EdgeGatewayID is the ID of the Edge Gateway whose owner will scope the Network Context Profile.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owner will scope the Network Context Profile.
		EdgeGatewayName string

		// Attributes is the list of Layer 7 attributes for this profile.
		Attributes []ParamsNetworkContextProfileAttribute
	}

	ParamsUpdateEdgeGatewayNetworkContextProfile struct {
		// ID is the unique identifier of the Network Context Profile to update.
		ID string

		// Name is the name of the Network Context Profile to update when ID is not provided.
		Name string

		// EdgeGatewayID is the ID of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayName string

		// Description is the new description of the Network Context Profile.
		Description string

		// Attributes is the full list of Layer 7 attributes for this profile.
		Attributes []ParamsNetworkContextProfileAttribute
	}

	ParamsDeleteEdgeGatewayNetworkContextProfile struct {
		// ID is the unique identifier of the Network Context Profile to delete.
		ID string

		// Name is the name of the Network Context Profile to delete when ID is not provided.
		Name string

		// EdgeGatewayID is the ID of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayName string
	}

	ParamsGetEdgeGatewayNetworkContextProfileAttributes struct {
		// EdgeGatewayID is the ID of the Edge Gateway whose owner scopes the attribute catalog.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owner scopes the attribute catalog.
		EdgeGatewayName string
	}
)
