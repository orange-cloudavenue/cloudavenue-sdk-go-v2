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
	// * Model

	ModelEdgeGatewayServices struct {
		LoadBalancer *ModelEdgeGatewayServicesLoadBalancer
		PublicIP     []*ModelEdgeGatewayServicesPublicIP
		Service      *ModelCloudavenueServices
	}

	ModelEdgeGatewayServicesSvc struct {
		// ID is the identifier of the network service
		ID string `documentation:"Identifier of the network service"`
		// Name is the name of the network service
		Name string `documentation:"Name of the network service"`
	}

	ModelEdgeGatewayServicesLoadBalancer struct {
		ModelEdgeGatewayServicesSvc

		ClassOfService     string `documentation:"Class of service for the load balancer"`
		MaxVirtualServices int    `documentation:"Maximum number of virtual services"`
	}

	ModelEdgeGatewayServicesPublicIP struct {
		ModelEdgeGatewayServicesSvc

		// IP is the public IP address
		IP string `documentation:"Public IP address"`
		// Announced represents if the public IP address is announced
		Announced bool `documentation:"Indicates if the public IP address is announced"`
	}

	ModelCloudavenueServices struct {
		ModelEdgeGatewayServicesSvc

		// Network is the network of the service ip/cidr
		Network string `documentation:"Network of the service in IP/CIDR format"`
		// DedicatedIPForService is the dedicated IP for the service
		// Used for the NAT to connect to the service
		DedicatedIPForService string `documentation:"Dedicated IP for the service in IP format (Used for the NAT to connect to the service)"`
		// Services is the list of services
		ServiceDetails []ModelCloudavenueServiceDetails `documentation:"List of services details"`
	}

	ModelNetworkServicesSvcServiceDetailsPorts struct {
		// Port is the port of the service
		Port int `documentation:"Port of the service"`
		// Protocol is the protocol of the service
		Protocol string `documentation:"Protocol of the service"`
	}

	ModelCloudavenueServiceDetails struct {
		// Category is the category of the service
		Category string
		// Network is the network of the service
		Network string
		// Services is the list of services
		Services []ModelCloudavenueServiceDetail
	}

	ModelCloudavenueServiceDetail struct {
		// Name is the name of the service
		Name string
		// Description
		Description string
		// DocumentationURL is the URL of the documentation
		DocumentationURL string
		// IP is the IP address of the service
		IP []string
		// FQDN is the FQDN of the service
		FQDN  []string
		Ports []ModelNetworkServicesSvcServiceDetailsPorts
	}
)
