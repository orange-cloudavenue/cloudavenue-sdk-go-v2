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
	ParamsGetEdgeGatewayFirewall struct {
		// EdgeGatewayID is the ID of the Edge Gateway whose owning VDC Group provides the Distributed Firewall.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owning VDC Group provides the Distributed Firewall.
		EdgeGatewayName string
	}

	ParamsCreateEdgeGatewayFirewall struct {
		// EdgeGatewayID is the ID of the Edge Gateway whose owning VDC Group provides the Distributed Firewall.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owning VDC Group provides the Distributed Firewall.
		EdgeGatewayName string

		// Enabled activates or deactivates the Distributed Firewall.
		Enabled bool

		// Rules is the list of Distributed Firewall rules.
		Rules []ParamsFirewallRule
	}

	ParamsUpdateEdgeGatewayFirewall struct {
		// EdgeGatewayID is the ID of the Edge Gateway whose owning VDC Group provides the Distributed Firewall.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owning VDC Group provides the Distributed Firewall.
		EdgeGatewayName string

		// Enabled activates or deactivates the Distributed Firewall.
		Enabled bool

		// Rules is the full list of Distributed Firewall rules.
		Rules []ParamsFirewallRule
	}

	ParamsDeleteEdgeGatewayFirewall struct {
		// EdgeGatewayID is the ID of the Edge Gateway whose owning VDC Group provides the Distributed Firewall.
		EdgeGatewayID string

		// EdgeGatewayName is the name of the Edge Gateway whose owning VDC Group provides the Distributed Firewall.
		EdgeGatewayName string
	}
)
