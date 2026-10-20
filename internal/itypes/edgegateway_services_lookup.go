/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

const (
	// networkTypeTier0VRF is the type of the root node of the network
	// hierarchy.
	networkTypeTier0VRF = "tier-0-vrf"

	// networkTypeEdgeGateway is the type of the second level of the network
	// hierarchy, holding the services attached to an edge gateway.
	networkTypeEdgeGateway = "edge-gateway"

	// networkTypeService is the type of a network service attached to an edge
	// gateway.
	networkTypeService = "service"

	// serviceNameInternet is the name of the network service backing an
	// allocated public IP.
	serviceNameInternet = "internet"
)

// PublicIPServiceID walks the raw org-wide network hierarchy and returns the
// real service identifier of the internet service carrying ip.
//
// The hierarchy is returned in a single call and is not scoped to one edge
// gateway, so every tier-0-vrf and every edge gateway is inspected instead of
// being filtered client-side. The identifier cannot be derived from the IP
// address, so it must always be read from the payload.
//
// The second return value reports whether a matching internet service exists at
// all. That lets callers tell an unknown IP apart from a service whose
// identifier is missing from the payload, and refuse both instead of guessing.
func (ap *APIResponseNetworkServices) PublicIPServiceID(ip string) (serviceID string, found bool) {
	if ap == nil {
		return "", false
	}

	for _, vrf := range *ap {
		if vrf.Type != networkTypeTier0VRF {
			continue
		}

		for _, child := range vrf.Children {
			if child.Type != networkTypeEdgeGateway {
				continue
			}

			for _, service := range child.Children {
				if service.Type != networkTypeService || service.Name != serviceNameInternet {
					continue
				}

				if service.Properties.IP == ip {
					return service.ServiceID, true
				}
			}
		}
	}

	return "", false
}
