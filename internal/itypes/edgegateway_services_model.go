/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import (
	"net"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func populateEdgeGatewayServices(data *types.ModelEdgeGatewayServices, services []ApiResponseNetworkServicesSubChildren) {
	for _, service := range services {
		switch service.Type {
		case "load-balancer":
			data.LoadBalancer = &types.ModelEdgeGatewayServicesLoadBalancer{
				ID:                 service.Name,
				Name:               service.DisplayName,
				ClassOfService:     service.Properties.ClassOfService,
				MaxVirtualServices: service.Properties.MaxVirtualServices,
			}
		case "service":
			populateGenericEdgeGatewayService(data, service)
		}
	}
}

func populateGenericEdgeGatewayService(data *types.ModelEdgeGatewayServices, service ApiResponseNetworkServicesSubChildren) {
	switch service.Name {
	case "cav-services", "cav_services":
		data.Services = &types.ModelCloudavenueServices{
			ID:        service.ServiceID,
			Name:      service.DisplayName,
			Network:   firstServiceRange(service.Properties.Ranges),
			IPAddress: firstServiceIPAddress(service.Properties.Ranges),
			Services:  ListOfServices,
		}
	case "internet":
		data.PublicIP = append(data.PublicIP, &types.ModelEdgeGatewayServicesPublicIP{
			ID:        service.ServiceID,
			Name:      service.Properties.IP,
			IP:        service.Properties.IP,
			Announced: service.Properties.Announced,
		})
	}
}

func firstServiceRange(ranges []string) string {
	if len(ranges) == 0 {
		return ""
	}

	return ranges[0]
}

func firstServiceIPAddress(ranges []string) string {
	if len(ranges) == 0 {
		return ""
	}

	ip, _, err := net.ParseCIDR(ranges[0])
	if err != nil {
		return ""
	}

	return ip.String()
}
