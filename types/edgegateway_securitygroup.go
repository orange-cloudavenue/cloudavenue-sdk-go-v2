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
	ParamsListEdgeGatewaySecurityGroup struct {
		// EdgeGatewayID is the ID of the Edge Gateway whose owning VDC Group provides the Security Groups.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owning VDC Group provides the Security Groups.
		EdgeGatewayName string
	}

	ParamsGetEdgeGatewaySecurityGroup struct {
		// ID is the unique identifier of the Security Group.
		ID string

		// Name is the name of the Security Group.
		Name string

		// EdgeGatewayID is the ID of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayName string
	}

	ParamsCreateEdgeGatewaySecurityGroup struct {
		// Name is the name of the Security Group.
		Name string

		// Description is the description of the Security Group.
		Description string

		// EdgeGatewayID is the ID of the Edge Gateway whose owning VDC Group will own the Security Group.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owning VDC Group will own the Security Group.
		EdgeGatewayName string

		// Members is the list of Org VDC network references to attach to this Security Group.
		Members []ParamsFirewallGroupMember
	}

	ParamsUpdateEdgeGatewaySecurityGroup struct {
		// ID is the unique identifier of the Security Group.
		ID string

		// Name is the name of the Security Group.
		Name string

		// EdgeGatewayID is the ID of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayName string

		// Description is the description of the Security Group.
		Description string

		// Members is the list of Org VDC network references to attach to this Security Group.
		Members []ParamsFirewallGroupMember
	}

	ParamsDeleteEdgeGatewaySecurityGroup struct {
		// ID is the unique identifier of the Security Group.
		ID string

		// Name is the name of the Security Group.
		Name string

		// EdgeGatewayID is the ID of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway used to scope name-based lookup.
		EdgeGatewayName string
	}
)
