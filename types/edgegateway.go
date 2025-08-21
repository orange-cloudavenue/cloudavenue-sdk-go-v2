/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

// * Models

type (
	// ModelEdgeGateways represents a list of edge gateways.
	ModelEdgeGateways struct {
		EdgeGateways []ModelEdgeGateway
	}
	// ModelEdgeGateway represents the model of an edge gateway.
	ModelEdgeGateway struct {
		ID string `documentation:"ID of the edge gateway"`

		// Name of edge gateway
		Name string `documentation:"Name of the edge gateway"`

		// Description of edge gateway
		Description string `documentation:"Description of the edge gateway"`

		// OwnerRef defines VDC or VDC Group that this network belongs to.
		OwnerRef *ModelObjectReference `documentation:"VDC or VDC Group that this edge gateway belongs to"`

		// UplinkT0 defines the T0 router name that this edge gateway is connected to.
		UplinkT0 *ModelObjectReference `documentation:"T0 router name that this edge gateway is connected to"`
	}

	// Bandwidth in Mbps.
	ModelBandwidth int
)

// * Functions Parameters

type (
	ParamsEdgeGateway struct {
		ID   string `validate:"required_if_null=Name,omitempty,urn=edgegateway" fake:"{urn:edgegateway}"`
		Name string `validate:"required_if_null=ID,omitempty,resource_name=edgegateway" fake:"{resource_name:edgegateway}"`
	}

	ParamsCreateEdgeGateway struct {
		// OwnerType is the type of the owner of the edge gateway.
		// It can be either "vdc" or "vdcgroup".
		OwnerType string `fake:"{randomstring:[vdc]}"`

		// OwnerName is the VDC or VDC Group that this edge gateway belongs to.
		OwnerName string `fake:"{vdc_name}"`

		// Name is the name of the T0 router that this edge gateway will be connected to.
		// If not provided and only if one T0 router is available,
		// the first T0 router will be used.
		T0Name string `fake:"{resource_name:t0}"`

		// Bandwidth is the bandwidth limit in Mbps.
		// If not provided, default bandwidth will be used (5 Mbps).
		// You can get bandwidth available values for the new edge gateway
		// by calling ListT0().
		// Unlimited bandwidth is allowed if the T0 is DEDICATED.
		Bandwidth int `fake:"5"`
	}

	ParamsUpdateEdgeGateway struct {
		// ID is the ID of the edge gateway to update.
		ID string `fake:"{urn:edgegateway}"`

		// Name is the new name of the edge gateway.
		Name string `fake:"{resource_name:edgegateway}"`

		// Bandwidth is the new bandwidth limit in Mbps.
		Bandwidth int `fake:"5"`
	}
)
