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
	AppPortProfileProtocolICMPv4 = "ICMPv4"
	AppPortProfileProtocolICMPv6 = "ICMPv6"
	AppPortProfileProtocolTCP    = "TCP"
	AppPortProfileProtocolUDP    = "UDP"

	AppPortProfileScopeTenant   = "TENANT"
	AppPortProfileScopeProvider = "PROVIDER"
	AppPortProfileScopeSystem   = "SYSTEM"
)

// AppPortProfilePortRegex validates a single port (1-65535) or a port range
// ("low-high", each side 1-65535) as used in Application Port Profile
// destination ports.
const AppPortProfilePortRegex = `^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])(-([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5]))?$`

// AppPortProfileProtocols lists all valid protocol values.
var AppPortProfileProtocols = []string{
	AppPortProfileProtocolICMPv4,
	AppPortProfileProtocolICMPv6,
	AppPortProfileProtocolTCP,
	AppPortProfileProtocolUDP,
}

type (
	ParamsListAppPortProfile struct {
		// VDCGroupID is the ID of the Vdc Group owning the Application Port Profiles.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning the Application Port Profiles.
		VDCGroupName string
	}

	ParamsGetAppPortProfile struct {
		// ID is the unique identifier of the Application Port Profile.
		ID string

		// Name is the name of the Application Port Profile.
		Name string

		// VDCGroupID is the ID of the Vdc Group owning this Application Port Profile.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Application Port Profile.
		VDCGroupName string
	}

	ParamsCreateAppPortProfile struct {
		// Name is the name of the Application Port Profile.
		Name string

		// Description is the description of the Application Port Profile.
		Description string

		// VDCGroupID is the ID of the Vdc Group owning this Application Port Profile.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Application Port Profile.
		VDCGroupName string

		// ApplicationPorts is the list of protocol/destination-ports pairs for this Application Port Profile.
		ApplicationPorts []ParamsAppPortProfilePort
	}

	ParamsAppPortProfilePort struct {
		// Protocol is the protocol of the port (ICMPv4, ICMPv6, TCP, UDP).
		Protocol string

		// DestinationPorts is the list of destination ports or port ranges. Must be empty for ICMPv4/ICMPv6, required for TCP/UDP.
		DestinationPorts []string
	}

	ParamsUpdateAppPortProfile struct {
		// ID is the unique identifier of the Application Port Profile to update.
		ID string

		// Name is the name of the Application Port Profile to update.
		Name string

		// Description is the new description of the Application Port Profile.
		Description string

		// ApplicationPorts is the full list of protocol/destination-ports pairs for this Application Port Profile.
		ApplicationPorts []ParamsAppPortProfilePort
	}

	ParamsDeleteAppPortProfile struct {
		// ID is the unique identifier of the Application Port Profile to delete.
		ID string

		// Name is the name of the Application Port Profile to delete.
		Name string

		// VDCGroupID is the ID of the Vdc Group owning this Application Port Profile.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Application Port Profile.
		VDCGroupName string
	}

	// * List
	ModelListAppPortProfile struct {
		AppPortProfiles []ModelGetAppPortProfile `documentation:"List of Application Port Profiles"`
	}

	// * Get
	ModelGetAppPortProfile struct {
		ID          string `documentation:"ID of the Application Port Profile"`
		Name        string `documentation:"Name of the Application Port Profile"`
		Description string `documentation:"Description of the Application Port Profile"`
		Scope       string `documentation:"Scope of the Application Port Profile (TENANT, PROVIDER, SYSTEM)"`
		OrgID       string `documentation:"ID of the Org owning the Application Port Profile"`

		ApplicationPorts []ModelGetAppPortProfilePort `documentation:"List of protocol/destination-ports pairs"`
	}

	ModelGetAppPortProfilePort struct {
		Protocol         string   `documentation:"Protocol of the port (ICMPv4, ICMPv6, TCP, UDP)"`
		DestinationPorts []string `documentation:"List of destination ports or port ranges"`
	}
)
