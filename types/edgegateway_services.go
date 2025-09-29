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
		// ID is the identifier of the edge gateway
		ID string `documentation:"Unique identifier of the edge gateway"`
		// Name is the name of the edge gateway
		Name string `documentation:"Edge gateway name"`

		LoadBalancer *ModelEdgeGatewayServicesLoadBalancer `documentation:"Load balancer service details (nil if not enabled)"`
		PublicIP     []*ModelEdgeGatewayServicesPublicIP   `documentation:"List of allocated public IP addresses"`
		Services     *ModelCloudavenueServices             `documentation:"Cloud Avenue service catalog attached to this edge gateway"`
	}

	ModelEdgeGatewayServicesLoadBalancer struct {
		ID                 string `documentation:"Unique identifier of the load balancer"`
		Name               string `documentation:"Load balancer name"`
		ClassOfService     string `documentation:"Service class (tier) applied to the load balancer"`
		MaxVirtualServices int    `documentation:"Maximum number of virtual services supported"`
	}

	ModelEdgeGatewayServicesPublicIP struct {
		ID        string `documentation:"Unique identifier of the public IP"`
		Name      string `documentation:"Public IP logical name"`
		IP        string `documentation:"Public IPv4 address"`
		Announced bool   `documentation:"True if the public IP is advertised via BGP"`
	}

	ModelCloudavenueServices struct {
		ID   string `documentation:"Unique identifier of the Cloud Avenue services configuration"`
		Name string `documentation:"Services configuration name"`

		// Network is the network of the service ip/cidr
		Network string `documentation:"Service network in CIDR notation"`
		// DedicatedIPForService is the dedicated IP for the service
		// Used for the NAT to connect to the service
		IPAddress string `documentation:"Dedicated IPv4 used as NAT entrypoint for the services"`
		// ServicesCatalog is the list of categorized service groups
		Services []ModelCloudavenueServicesCatalog `documentation:"Catalog of Cloud Avenue Service categories and their services"`
	}

	ModelCloudavenueServicesCatalog struct {
		// Category is the category of the service
		Category string `documentation:"Service category grouping related endpoints"`
		// Network is the network of the service
		Network string `documentation:"Category network in CIDR notation"`
		// Services is the list of services
		Services []ModelCloudavenueServicesCatalogService `documentation:"Service instances within the category"`
	}

	ModelCloudavenueServicesCatalogService struct {
		// Name is the name of the service
		Name string `documentation:"Service name"`
		// Description
		Description string `documentation:"Human-readable description"`
		// DocumentationURL is the URL of the documentation
		DocumentationURL string `documentation:"Reference documentation URL"`
		// IP is the IP address of the service
		IPs []string `documentation:"One or more IPv4 addresses used by the service"`
		// FQDN is the FQDN of the service
		FQDNs []string                                     `documentation:"Fully qualified domain names for the service"`
		Ports []ModelCloudavenueServicesCatalogServicePort `documentation:"Exposed service ports"`
	}

	ModelCloudavenueServicesCatalogServicePort struct {
		// Port is the port of the service
		Port int `documentation:"TCP/UDP port number"`
		// Protocol is the protocol of the service
		Protocol string `documentation:"Layer 4 protocol (e.g. TCP, UDP)"`
	}
)
