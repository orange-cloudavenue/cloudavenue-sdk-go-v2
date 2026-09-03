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
	ParamsListEdgeGatewayAppPortProfile struct {
		// EdgeGatewayID is the ID of the Edge Gateway whose owner scopes the Application Port Profiles.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owner scopes the Application Port Profiles.
		EdgeGatewayName string
	}

	ParamsGetEdgeGatewayAppPortProfile struct {
		// ID is the unique identifier of the Application Port Profile.
		ID string

		// Name is the name of the Application Port Profile.
		Name string

		// EdgeGatewayID is the ID of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayName string
	}

	ParamsCreateEdgeGatewayAppPortProfile struct {
		// Name is the name of the Application Port Profile.
		Name string

		// Description is the description of the Application Port Profile.
		Description string

		// EdgeGatewayID is the ID of the Edge Gateway whose owner will scope the Application Port Profile.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owner will scope the Application Port Profile.
		EdgeGatewayName string

		// ApplicationPorts is the list of protocol/destination-ports pairs for this Application Port Profile.
		ApplicationPorts []ParamsAppPortProfilePort
	}

	ParamsUpdateEdgeGatewayAppPortProfile struct {
		// ID is the unique identifier of the Application Port Profile to update.
		ID string

		// Name is the name of the Application Port Profile to update when ID is not provided.
		Name string

		// EdgeGatewayID is the ID of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayName string

		// Description is the new description of the Application Port Profile.
		Description string

		// ApplicationPorts is the full list of protocol/destination-ports pairs for this Application Port Profile.
		ApplicationPorts []ParamsAppPortProfilePort
	}

	ParamsDeleteEdgeGatewayAppPortProfile struct {
		// ID is the unique identifier of the Application Port Profile to delete.
		ID string

		// Name is the name of the Application Port Profile to delete when ID is not provided.
		Name string

		// EdgeGatewayID is the ID of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayName string
	}
)
