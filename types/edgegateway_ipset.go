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
	ParamsListEdgeGatewayIPSet struct {
		// EdgeGatewayID is the ID of the Edge Gateway whose owning VDC Group provides the IP Sets.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owning VDC Group provides the IP Sets.
		EdgeGatewayName string
	}

	ParamsGetEdgeGatewayIPSet struct {
		// ID is the unique identifier of the IP Set.
		ID string

		// Name is the name of the IP Set.
		Name string

		// EdgeGatewayID is the ID of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayName string
	}

	ParamsCreateEdgeGatewayIPSet struct {
		// Name is the name of the IP Set.
		Name string

		// Description is the description of the IP Set.
		Description string

		// EdgeGatewayID is the ID of the Edge Gateway whose owning VDC Group will own the IP Set.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owning VDC Group will own the IP Set.
		EdgeGatewayName string

		// IPAddresses is the list of IP addresses, ranges, or CIDRs in this IP Set.
		IPAddresses []string
	}

	ParamsUpdateEdgeGatewayIPSet struct {
		// ID is the unique identifier of the IP Set.
		ID string

		// Name is the name of the IP Set.
		Name string

		// EdgeGatewayID is the ID of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayName string

		// Description is the description of the IP Set.
		Description string

		// IPAddresses is the list of IP addresses, ranges, or CIDRs in this IP Set.
		IPAddresses []string
	}

	ParamsDeleteEdgeGatewayIPSet struct {
		// ID is the unique identifier of the IP Set.
		ID string

		// Name is the name of the IP Set.
		Name string

		// EdgeGatewayID is the ID of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayName string
	}
)
