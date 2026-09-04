/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

const (
	VDCNetworkTypeRouted   = "NAT_ROUTED"
	VDCNetworkTypeIsolated = "ISOLATED"
)

type (
	// * Subnet (shared shape used by both Isolated and Routed networks)
	ParamsSubnet struct {
		// Gateway is the gateway IP address of the subnet.
		Gateway string

		// PrefixLength is the prefix length of the subnet.
		PrefixLength int

		// DNSServer1 is the primary DNS server for the subnet.
		DNSServer1 string

		// DNSServer2 is the secondary DNS server for the subnet.
		DNSServer2 string

		// DNSSuffix is the DNS suffix for the subnet.
		DNSSuffix string

		// IPRanges is the list of static IP pools available in this subnet.
		IPRanges []ParamsVDCNetworkIPRange
	}

	ParamsVDCNetworkIPRange struct {
		// StartAddress is the first IP address of the range.
		StartAddress string

		// EndAddress is the last IP address of the range.
		EndAddress string
	}

	// * List
	ParamsListVDCNetwork struct {
		// VDCGroupID is the ID of the Vdc Group owning the Org VDC Networks.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning the Org VDC Networks.
		VDCGroupName string
	}

	// * Isolated networks
	ParamsGetVDCNetworkIsolated struct {
		// ID is the unique identifier of the Org VDC Network.
		ID string

		// Name is the name of the Org VDC Network.
		Name string

		// VDCGroupID is the ID of the Vdc Group owning this Org VDC Network.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Org VDC Network.
		VDCGroupName string
	}

	ParamsCreateVDCNetworkIsolated struct {
		// Name is the name of the Org VDC Network.
		Name string

		// Description is the description of the Org VDC Network.
		Description string

		// VDCGroupID is the ID of the Vdc Group owning this Org VDC Network.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Org VDC Network.
		VDCGroupName string

		// Subnet is the subnet configuration of the Org VDC Network.
		Subnet ParamsSubnet

		// GuestVLANTaggingAllowed defines if guest VLAN tagging is allowed on this network.
		GuestVLANTaggingAllowed *bool
	}

	ParamsUpdateVDCNetworkIsolated struct {
		// ID is the unique identifier of the Org VDC Network to update.
		ID string

		// Name is the name of the Org VDC Network to update.
		Name string

		// Description is the new description of the Org VDC Network.
		Description string

		// Subnet is the new subnet configuration of the Org VDC Network.
		Subnet *ParamsSubnet

		// GuestVLANTaggingAllowed defines if guest VLAN tagging is allowed on this network.
		GuestVLANTaggingAllowed *bool
	}

	ParamsDeleteVDCNetworkIsolated struct {
		// ID is the unique identifier of the Org VDC Network to delete.
		ID string

		// Name is the name of the Org VDC Network to delete.
		Name string

		// VDCGroupID is the ID of the Vdc Group owning this Org VDC Network.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Org VDC Network.
		VDCGroupName string
	}

	// * Routed networks
	ParamsGetVDCNetworkRouted struct {
		// ID is the unique identifier of the Org VDC Network.
		ID string

		// Name is the name of the Org VDC Network.
		Name string

		// VDCGroupID is the ID of the Vdc Group owning this Org VDC Network.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Org VDC Network.
		VDCGroupName string
	}

	ParamsCreateVDCNetworkRouted struct {
		// Name is the name of the Org VDC Network.
		Name string

		// Description is the description of the Org VDC Network.
		Description string

		// VDCGroupID is the ID of the Vdc Group owning this Org VDC Network.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Org VDC Network.
		VDCGroupName string

		// Subnet is the subnet configuration of the Org VDC Network.
		Subnet ParamsSubnet

		// GuestVLANTaggingAllowed defines if guest VLAN tagging is allowed on this network.
		GuestVLANTaggingAllowed *bool

		// EdgeGatewayID is the ID of the Edge Gateway this network is connected to.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway this network is connected to.
		EdgeGatewayName string
	}

	ParamsUpdateVDCNetworkRouted struct {
		// ID is the unique identifier of the Org VDC Network to update.
		ID string

		// Name is the name of the Org VDC Network to update.
		Name string

		// Description is the new description of the Org VDC Network.
		Description string

		// Subnet is the new subnet configuration of the Org VDC Network.
		Subnet *ParamsSubnet

		// GuestVLANTaggingAllowed defines if guest VLAN tagging is allowed on this network.
		GuestVLANTaggingAllowed *bool

		// EdgeGatewayID is the new ID of the Edge Gateway this network is connected to.
		EdgeGatewayID string

		// EdgeGatewayName is the new name of the Edge Gateway this network is connected to.
		EdgeGatewayName string
	}

	ParamsDeleteVDCNetworkRouted struct {
		// ID is the unique identifier of the Org VDC Network to delete.
		ID string

		// Name is the name of the Org VDC Network to delete.
		Name string

		// VDCGroupID is the ID of the Vdc Group owning this Org VDC Network.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Org VDC Network.
		VDCGroupName string
	}

	// * Model (shared shape for Isolated and Routed networks)
	ModelListVDCNetwork struct {
		VDCNetworks []ModelGetVDCNetwork `documentation:"List of Org VDC Networks"`
	}

	ModelGetVDCNetwork struct {
		ID          string `documentation:"ID of the Org VDC Network"`
		Name        string `documentation:"Name of the Org VDC Network"`
		Description string `documentation:"Description of the Org VDC Network"`
		Status      string `documentation:"Status of the Org VDC Network"`
		NetworkType string `documentation:"Type of the Org VDC Network (NAT_ROUTED or ISOLATED)"`

		OwnerID   string `documentation:"ID of the owner (VDC or VDCGroup) of the Org VDC Network"`
		OwnerName string `documentation:"Name of the owner (VDC or VDCGroup) of the Org VDC Network"`

		Subnet ModelVDCNetworkSubnet `documentation:"Subnet configuration of the Org VDC Network"`

		GuestVLANTaggingAllowed *bool `documentation:"Defines if guest VLAN tagging is allowed on this network"`
		Shared                  *bool `documentation:"Defines if this network is shared across the VDCGroup"`

		EdgeGatewayID   string `documentation:"ID of the Edge Gateway this network is connected to (only for routed networks)"`
		EdgeGatewayName string `documentation:"Name of the Edge Gateway this network is connected to (only for routed networks)"`
	}

	ModelVDCNetworkSubnet struct {
		Gateway      string `documentation:"Gateway IP address of the subnet"`
		PrefixLength int    `documentation:"Prefix length of the subnet"`
		DNSServer1   string `documentation:"Primary DNS server for the subnet"`
		DNSServer2   string `documentation:"Secondary DNS server for the subnet"`
		DNSSuffix    string `documentation:"DNS suffix for the subnet"`

		IPRanges []ModelVDCNetworkIPRange `documentation:"List of static IP pools available in this subnet"`
	}

	ModelVDCNetworkIPRange struct {
		StartAddress string `documentation:"First IP address of the range"`
		EndAddress   string `documentation:"Last IP address of the range"`
	}
)
