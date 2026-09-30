/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"

// DhcpConfig represents the DHCP configuration of an Org VDC Network.
type DhcpConfig struct {
	ID                       string                 `json:"id,omitempty"`
	VDCNetworkID             string                 `json:"vdcNetworkId,omitempty"`
	DHCPServerIPAddress      string                 `json:"dhcpServerIpAddress,omitempty"`
	DHCPServerPort           int                    `json:"dhcpServerPort,omitempty"`
	DHCPRelayServerIPAddress string                 `json:"dhcpRelayServerIpAddress,omitempty"`
	DHCPEnabled              bool                   `json:"dhcpEnabled,omitempty"`
	DHCPServerCredentials    *DHCPServerCredentials `json:"dhcpServerCredentials,omitempty"`
	DHCPLeaseTime            int                    `json:"dhcpLeaseTime,omitempty"`
	DHCPIPAddress            string                 `json:"dhcpIpRange,omitempty"`
	DHCPPoolIPAddress        string                 `json:"dhcpPoolIpAddress,omitempty"`
	DHCPSubnetMask           string                 `json:"dhcpSubnetMask,omitempty"`
	DHCPDefaultGateway       string                 `json:"dhcpDefaultGateway,omitempty"`
	DHCPDNS1IPAddress        string                 `json:"dhcpDns1IpAddress,omitempty"`
	DHCPDNS2IPAddress        string                 `json:"dhcpDns2IpAddress,omitempty"`
	DHCPSearchDomain         string                 `json:"dhcpSearchDomain,omitempty"`
	DHCPOptions              []DHCPOption           `json:"dhcpOptions,omitempty"`
}

// DHCPServerCredentials holds the credentials used by the DHCP server.
type DHCPServerCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// DHCPOption represents a single DHCP option (code + value).
type DHCPOption struct {
	Code  string `json:"code,omitempty"`
	Value string `json:"value,omitempty"`
}

// APIRequestDhcpConfig is the request payload shape used by PUT /dhcp/put/.
// It mirrors DhcpConfig and is used as BodyRequestType.
type APIRequestDhcpConfig = DhcpConfig

// ToModel converts the internal DhcpConfig to the public model.
func (r *DhcpConfig) ToModel() types.ModelDhcpConfig {
	m := types.ModelDhcpConfig{
		ID:                       r.ID,
		VDCNetworkID:             r.VDCNetworkID,
		DHCPServerIPAddress:      r.DHCPServerIPAddress,
		DHCPServerPort:           r.DHCPServerPort,
		DHCPRelayServerIPAddress: r.DHCPRelayServerIPAddress,
		DHCPEnabled:              r.DHCPEnabled,
		DHCPLeaseTime:            r.DHCPLeaseTime,
		DHCPIPAddress:            r.DHCPIPAddress,
		DHCPPoolIPAddress:        r.DHCPPoolIPAddress,
		DHCPSubnetMask:           r.DHCPSubnetMask,
		DHCPDefaultGateway:       r.DHCPDefaultGateway,
		DHCPDNS1IPAddress:        r.DHCPDNS1IPAddress,
		DHCPDNS2IPAddress:        r.DHCPDNS2IPAddress,
		DHCPSearchDomain:         r.DHCPSearchDomain,
	}

	if r.DHCPServerCredentials != nil {
		creds := types.ModelDhcpServerCredentials{
			Username: r.DHCPServerCredentials.Username,
			Password: r.DHCPServerCredentials.Password,
		}
		m.DHCPServerCredentials = &creds
	}

	m.DHCPOptions = make([]types.ModelDhcpOption, 0, len(r.DHCPOptions))
	for _, opt := range r.DHCPOptions {
		m.DHCPOptions = append(m.DHCPOptions, types.ModelDhcpOption{
			Code:  opt.Code,
			Value: opt.Value,
		})
	}

	return m
}
