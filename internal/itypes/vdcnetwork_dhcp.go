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

// DhcpConfig is the VMware CloudAPI VdcNetworkDhcpConfig wire model.
type DhcpConfig struct {
	Enabled    bool       `json:"enabled"`
	LeaseTime  int64      `json:"leaseTime,omitempty"`
	DhcpPools  []DhcpPool `json:"dhcpPools,omitempty"`
	Mode       string     `json:"mode,omitempty"`
	IPAddress  string     `json:"ipAddress,omitempty"`
	DNSServers []string   `json:"dnsServers,omitempty"`
}

// DhcpPool is a DHCP address pool in the CloudAPI response.
type DhcpPool struct {
	Enabled          bool    `json:"enabled"`
	IPRange          IPRange `json:"ipRange"`
	MaxLeaseTime     int64   `json:"maxLeaseTime,omitempty"`
	DefaultLeaseTime int64   `json:"defaultLeaseTime,omitempty"`
}

// IPRange is an inclusive IP address range.
type IPRange struct {
	StartAddress string `json:"startAddress,omitempty"`
	EndAddress   string `json:"endAddress,omitempty"`
}

// APIRequestDhcpConfig is the request payload for PUT /dhcp.
type APIRequestDhcpConfig = DhcpConfig

// ToModel converts the internal DHCP configuration to the public model.
func (r *DhcpConfig) ToModel() types.ModelDhcpConfig {
	m := types.ModelDhcpConfig{
		Enabled:           r.Enabled,
		LeaseTime:         r.LeaseTime,
		Mode:              r.Mode,
		ListenerIPAddress: r.IPAddress,
		DNSServers:        append([]string(nil), r.DNSServers...),
		Pools:             make([]types.ModelDhcpPool, 0, len(r.DhcpPools)),
	}

	for _, pool := range r.DhcpPools {
		m.Pools = append(m.Pools, types.ModelDhcpPool{
			StartAddress: pool.IPRange.StartAddress,
			EndAddress:   pool.IPRange.EndAddress,
		})
	}

	return m
}
