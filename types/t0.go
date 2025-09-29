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
	// * Params

	ParamsGetT0 struct {
		Name            string
		EdgeGatewayID   string
		EdgeGatewayName string
	}

	//* Model

	// ModelT0s represents the model for T0 routers.
	ModelT0s struct {
		Count int
		T0s   []ModelT0
	}

	// ModelT0 represents a T0 router model.
	// It contains the name, class of service, bandwidth, summary, and associated edge gateways.
	ModelT0 struct {
		// Name defines the name of the T0 router.
		Name string `documentation:"Name of the T0 router"`

		// ClassOfService defines the class of service for the T0 router.
		ClassOfService string `documentation:"Class of service for the T0 router."`

		// Bandwidth defines the bandwidth for the T0 router.
		Bandwidth ModelT0Bandwidth `documentation:"Bandwidth for the T0 router in Mbps."`

		// Summary provides a summary of the T0 router.
		Summary ModelT0Summary `documentation:"Summary of the T0 router."`

		// EdgeGateways contains the list of edge gateways associated with the T0 router.
		EdgeGateways []ModelT0EdgeGateway `documentation:"List of edge gateways associated with the T0 router."`
	}

	// ModelT0Summary represents a summary of a T0 router.
	ModelT0Summary struct {
		EdgeGateways    int `documentation:"Number of edge gateways associated with the T0 router."`
		MaxEdgeGateways int `documentation:"Maximum number of edge gateways allowed for the T0 router."`
	}

	// ModelT0Bandwidth represents the bandwidth model for a T0 router.
	ModelT0Bandwidth struct {
		// Capacity defines the total bandwidth capacity in Mbps.
		// This is the maximum bandwidth that can be allocated to the T0 router.
		// It is a limit imposed by the Class of Service.
		Capacity int `documentation:"Total bandwidth capacity for the T0 router in Mbps. This is the maximum bandwidth that can be allocated to the T0 router. It is a limit imposed by the Class of Service."`

		// Provisioned defines the amount of bandwidth that has been provisioned in the each edge gateway.
		// This is the total bandwidth that has been allocated to the T0 router across all edge gateways.
		// It is the sum of the bandwidth allocated to each edge gateway.
		Provisioned int `documentation:"Total bandwidth provisioned for the T0 router across all edge gateways in Mbps."`

		// Remaining defines the remaining bandwidth that can be allocated to the edge gateways.
		Remaining int `documentation:"Remaining bandwidth that can be allocated to the edge gateways in Mbps. This is calculated as Capacity - Provisioned."`

		// AllowedBandwidthValues returns the allowed bandwidth values for the T0 router.
		// It's used to determine the available bandwidth options for the new edge gateway.
		// It returns a slice of integers representing the allowed bandwidth values in Mbps.
		// If values are empty, no bandwidth remaining to allocate.
		AllowedBandwidthValues []int `documentation:"Allowed bandwidth values for the T0 router in Mbps. This is used to determine the available bandwidth options for the new edge gateway. If empty, no bandwidth remaining to allocate."`

		// AllowUnlimited indicates if unlimited bandwidth is allowed for the T0 router.
		// This is true if the T0 router is DEDICATED.
		AllowUnlimited bool `documentation:"Indicates if unlimited bandwidth is allowed for the T0 router. This is true if the T0 router is DEDICATED."`
	}

	// ModelT0EdgeGateway represents an edge gateway associated with a T0 router.
	ModelT0EdgeGateway struct {
		// ID defines the unique identifier of the edge gateway.
		// It is a URN that uniquely identifies the edge gateway.
		ID string `documentation:"Unique identifier of the edge gateway. It is a URN that uniquely identifies the edge gateway."`

		// Name defines the name of the edge gateway.
		Name string `documentation:"Name of the edge gateway."`

		// Bandwidth defines the bandwidth allocated to the edge gateway.
		// The value is in Mbps.
		Bandwidth int `documentation:"Bandwidth allocated to the edge gateway in Mbps."`

		// AllowedBandwidthValues returns the allowed bandwidth values for the edge gateway.
		AllowedBandwidthValues []int `documentation:"Allowed bandwidth values for the edge gateway in Mbps."`
	}
)
