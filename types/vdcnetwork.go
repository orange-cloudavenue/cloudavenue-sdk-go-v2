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
	VdcNetworkTypeRouted   = "NAT_ROUTED"
	VdcNetworkTypeIsolated = "ISOLATED"
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
		IPRanges []ParamsVdcNetworkIPRange
	}

	ParamsVdcNetworkIPRange struct {
		// StartAddress is the first IP address of the range.
		StartAddress string

		// EndAddress is the last IP address of the range.
		EndAddress string
	}

	// * List
	ParamsListVdcNetwork struct {
		// VdcGroupID is the ID of the Vdc Group owning the Org VDC Networks.
		VdcGroupID string

		// VdcGroupName is the name of the Vdc Group owning the Org VDC Networks.
		VdcGroupName string
	}

	// * Isolated networks
	ParamsGetVdcNetworkIsolated struct {
		// ID is the unique identifier of the Org VDC Network.
		ID string

		// Name is the name of the Org VDC Network.
		Name string

		// VdcGroupID is the ID of the Vdc Group owning this Org VDC Network.
		VdcGroupID string

		// VdcGroupName is the name of the Vdc Group owning this Org VDC Network.
		VdcGroupName string
	}

	ParamsCreateVdcNetworkIsolated struct {
		// Name is the name of the Org VDC Network.
		Name string

		// Description is the description of the Org VDC Network.
		Description string

		// VdcGroupID is the ID of the Vdc Group owning this Org VDC Network.
		VdcGroupID string

		// VdcGroupName is the name of the Vdc Group owning this Org VDC Network.
		VdcGroupName string

		// Subnet is the subnet configuration of the Org VDC Network.
		Subnet ParamsSubnet

		// GuestVlanTaggingAllowed defines if guest VLAN tagging is allowed on this network.
		GuestVlanTaggingAllowed *bool
	}

	ParamsUpdateVdcNetworkIsolated struct {
		// ID is the unique identifier of the Org VDC Network to update.
		ID string

		// Name is the name of the Org VDC Network to update.
		Name string

		// Description is the new description of the Org VDC Network.
		Description string

		// Subnet is the new subnet configuration of the Org VDC Network.
		Subnet *ParamsSubnet

		// GuestVlanTaggingAllowed defines if guest VLAN tagging is allowed on this network.
		GuestVlanTaggingAllowed *bool
	}

	ParamsDeleteVdcNetworkIsolated struct {
		// ID is the unique identifier of the Org VDC Network to delete.
		ID string

		// Name is the name of the Org VDC Network to delete.
		Name string

		// VdcGroupID is the ID of the Vdc Group owning this Org VDC Network.
		VdcGroupID string

		// VdcGroupName is the name of the Vdc Group owning this Org VDC Network.
		VdcGroupName string
	}

	// * Routed networks
	ParamsGetVdcNetworkRouted struct {
		// ID is the unique identifier of the Org VDC Network.
		ID string

		// Name is the name of the Org VDC Network.
		Name string

		// VdcGroupID is the ID of the Vdc Group owning this Org VDC Network.
		VdcGroupID string

		// VdcGroupName is the name of the Vdc Group owning this Org VDC Network.
		VdcGroupName string
	}

	ParamsCreateVdcNetworkRouted struct {
		// Name is the name of the Org VDC Network.
		Name string

		// Description is the description of the Org VDC Network.
		Description string

		// VdcGroupID is the ID of the Vdc Group owning this Org VDC Network.
		VdcGroupID string

		// VdcGroupName is the name of the Vdc Group owning this Org VDC Network.
		VdcGroupName string

		// Subnet is the subnet configuration of the Org VDC Network.
		Subnet ParamsSubnet

		// GuestVlanTaggingAllowed defines if guest VLAN tagging is allowed on this network.
		GuestVlanTaggingAllowed *bool

		// EdgeGatewayID is the ID of the Edge Gateway this network is connected to.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway this network is connected to.
		EdgeGatewayName string
	}

	ParamsUpdateVdcNetworkRouted struct {
		// ID is the unique identifier of the Org VDC Network to update.
		ID string

		// Name is the name of the Org VDC Network to update.
		Name string

		// Description is the new description of the Org VDC Network.
		Description string

		// Subnet is the new subnet configuration of the Org VDC Network.
		Subnet *ParamsSubnet

		// GuestVlanTaggingAllowed defines if guest VLAN tagging is allowed on this network.
		GuestVlanTaggingAllowed *bool

		// EdgeGatewayID is the new ID of the Edge Gateway this network is connected to.
		EdgeGatewayID string

		// EdgeGatewayName is the new name of the Edge Gateway this network is connected to.
		EdgeGatewayName string
	}

	ParamsDeleteVdcNetworkRouted struct {
		// ID is the unique identifier of the Org VDC Network to delete.
		ID string

		// Name is the name of the Org VDC Network to delete.
		Name string

		// VdcGroupID is the ID of the Vdc Group owning this Org VDC Network.
		VdcGroupID string

		// VdcGroupName is the name of the Vdc Group owning this Org VDC Network.
		VdcGroupName string
	}

	// * Model (shared shape for Isolated and Routed networks)
	ModelListVdcNetwork struct {
		VdcNetworks []ModelGetVdcNetwork `documentation:"List of Org VDC Networks"`
	}

	ModelGetVdcNetwork struct {
		ID          string `documentation:"ID of the Org VDC Network"`
		Name        string `documentation:"Name of the Org VDC Network"`
		Description string `documentation:"Description of the Org VDC Network"`
		Status      string `documentation:"Status of the Org VDC Network"`
		NetworkType string `documentation:"Type of the Org VDC Network (NAT_ROUTED or ISOLATED)"`

		OwnerID   string `documentation:"ID of the owner (VDC or VdcGroup) of the Org VDC Network"`
		OwnerName string `documentation:"Name of the owner (VDC or VdcGroup) of the Org VDC Network"`

		Subnet ModelVdcNetworkSubnet `documentation:"Subnet configuration of the Org VDC Network"`

		GuestVlanTaggingAllowed *bool `documentation:"Defines if guest VLAN tagging is allowed on this network"`
		Shared                  *bool `documentation:"Defines if this network is shared across the VdcGroup"`

		EdgeGatewayID   string `documentation:"ID of the Edge Gateway this network is connected to (only for routed networks)"`
		EdgeGatewayName string `documentation:"Name of the Edge Gateway this network is connected to (only for routed networks)"`
	}

	ModelVdcNetworkSubnet struct {
		Gateway      string `documentation:"Gateway IP address of the subnet"`
		PrefixLength int    `documentation:"Prefix length of the subnet"`
		DNSServer1   string `documentation:"Primary DNS server for the subnet"`
		DNSServer2   string `documentation:"Secondary DNS server for the subnet"`
		DNSSuffix    string `documentation:"DNS suffix for the subnet"`

		IPRanges []ModelVdcNetworkIPRange `documentation:"List of static IP pools available in this subnet"`
	}

	ModelVdcNetworkIPRange struct {
		StartAddress string `documentation:"First IP address of the range"`
		EndAddress   string `documentation:"Last IP address of the range"`
	}
)
