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
	"github.com/orange-cloudavenue/common-go/urn"
)

type (

	// * ApiResponse
	ApiResponseNetworkServices []struct {
		Type     string                               `json:"type" fake:"tier-0-vrf"`
		Name     string                               `json:"name" fake:"{resource_name:t0}"`
		Children []ApiResponseNetworkServicesChildren `json:"children,omitempty" fakesize:"1"`
	}

	ApiResponseNetworkServicesChildren struct {
		Type        string `json:"type" fake:"edge-gateway"`
		Name        string `json:"name,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		Properties  struct {
			// EdgeGateway
			RateLimit int    `json:"rateLimit,omitempty"`
			EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"` // The UUID of the edge gateway
		} `json:"properties,omitempty"`
		Children  []ApiResponseNetworkServicesSubChildren `json:"children,omitempty" fakesize:"6"`
		ServiceID string                                  `json:"serviceId,omitempty"`
	}
	ApiResponseNetworkServicesSubChildren struct {
		Type        string `json:"type" fake:"{randomstring:[load-balancer,service]}"`
		Name        string `json:"name,omitempty" fake:"{randomstring:[cav-services,internet]}"`
		DisplayName string `json:"displayName,omitempty" fake:"{word}"`
		Properties  struct {
			// Load Balancer
			ClassOfService     string `json:"classOfService,omitempty"`
			MaxVirtualServices int    `json:"maxVirtualServices,omitempty"`

			// Public IP
			IP        string `json:"ip,omitempty" fake:"{ipv4address}"`
			Announced bool   `json:"announced,omitempty" fake:"true"`

			// Service
			Ranges []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"` // The network in ip/cidr format
		} `json:"properties,omitempty"`
		ServiceID string `json:"serviceId,omitempty"`
	}

	// * ApiRequest

	ApiRequestNetworkServicesCavSvc struct {
		// NetworkType
		NetworkType string `json:"networkType" default:"cav-services" validate:"required"` // The type of network service to create (load-balancer, service, internet)

		// EdgeGatewayID - The ID of the edge gateway is a UUID and not a URN.
		EdgeGatewayID string `json:"edgeGateway" validate:"required,uuid"`

		// Properties
		Properties struct {
			PrefixLength int `json:"prefixLength,omitempty" validate:"omitempty,min=25,max=28" default:"27"` // The prefix length of the network in CIDR notation
		}
	}
)

func (ap *ApiResponseNetworkServices) ToModel(params types.ParamsEdgeGateway) *types.ModelEdgeGatewayServices {
	if ap == nil || len(*ap) == 0 {
		return nil
	}

	data := &types.ModelEdgeGatewayServices{
		Services:     nil,
		LoadBalancer: nil,
		PublicIP:     nil,
	}

	// Parse the original response and populate the NetworkServicesModel
	for _, ns := range *ap {
		for _, child := range ns.Children {
			if child.Type == "edge-gateway" && (child.Properties.EdgeUUID == urn.ExtractUUID(params.ID) || child.Name == params.Name) {
				// Found the edge gateway
				data.ID = urn.Normalize(urn.EdgeGateway, child.Properties.EdgeUUID).String()
				data.Name = child.Name

				// iterate over the children to find the services
				for _, service := range child.Children {
					switch service.Type {
					case "load-balancer":
						// Found load balancer service
						data.LoadBalancer = &types.ModelEdgeGatewayServicesLoadBalancer{
							ID:                 service.Name,        // The name is the ID
							Name:               service.DisplayName, // The display name is the name
							ClassOfService:     service.Properties.ClassOfService,
							MaxVirtualServices: service.Properties.MaxVirtualServices,
						}
					case "service":
						// service is a generic service
						// the name of the service define the type of service
						switch service.Name {
						case "cav-services", "cav_services": // Match both cav-services and cav_services
							// Found cav-services
							data.Services = &types.ModelCloudavenueServices{
								ID:   service.ServiceID,   // The ServiceID is the ID
								Name: service.DisplayName, // The display name is the name
								Network: func() string {
									if len(service.Properties.Ranges) == 0 {
										return ""
									}

									return service.Properties.Ranges[0] // The first range is the network
								}(),
								DedicatedIPForService: func() string {
									if len(service.Properties.Ranges) == 0 {
										return ""
									}

									// Parse Network (ip/cidr) to get the first IP of the network
									// and use it as the dedicated IP for the service

									ip, _, err := net.ParseCIDR(service.Properties.Ranges[0])
									if err != nil {
										return ""
									}
									return ip.String()
								}(),
								ServicesDetails: ListOfServices,
							}

						case "internet":
							// Found internet service
							publicIP := &types.ModelEdgeGatewayServicesPublicIP{
								ID:        service.ServiceID,     // The ServiceID is the ID
								Name:      service.Properties.IP, // The IP don't have a name use IP instead
								IP:        service.Properties.IP,
								Announced: service.Properties.Announced,
							}

							// Prevent nil pointer dereference
							if data.PublicIP == nil {
								data.PublicIP = make([]*types.ModelEdgeGatewayServicesPublicIP, 0)
							}

							// Append the public IP to the list
							data.PublicIP = append(data.PublicIP, publicIP)
						}
					}
				}
			}
		}
	}

	return data
}

var ListOfServices = []types.ModelCloudavenueServiceDetails{
	{
		Category: "administration",
		Network:  "57.199.209.192/27",
		Services: []types.ModelCloudavenueServiceDetailService{
			{
				Name:        "linux-repository",
				Description: "Linux (Debian, Ubuntu, CentOS) package repository",
				IP:          []string{"57.199.209.214"},
				FQDN:        []string{"repo.service.cav"},
				Ports: []types.ModelCloudavenueServiceDetailServicePort{
					{
						Port:     3142,
						Protocol: "tcp",
					},
				},
			},
			{
				Name:        "rhui-repository",
				Description: "Red Hat (RHUI) package repository",
				IP:          []string{"57.199.209.197"},
				FQDN:        []string{"rhui.service.cav"},
				Ports: []types.ModelCloudavenueServiceDetailServicePort{
					{
						Port:     8080,
						Protocol: "tcp",
					},
				},
			},
			{
				Name:        "windows-repository",
				Description: "Windows (WSUS) package repository",
				IP:          []string{"57.199.209.212"},
				FQDN:        []string{"wsus.service.cav"},
				Ports: []types.ModelCloudavenueServiceDetailServicePort{
					{
						Port:     8530,
						Protocol: "tcp",
					},
					{
						Port:     8531,
						Protocol: "tcp",
					},
				},
			},
			{
				Name:        "windows-kms",
				Description: "Windows (KMS) license server",
				IP: []string{
					"57.199.209.210",
				},
				FQDN: []string{"kms.service.cav"},
				Ports: []types.ModelCloudavenueServiceDetailServicePort{
					{
						Port:     1688,
						Protocol: "tcp",
					},
				},
			},
			{
				Name:        "ntp",
				Description: "Network Time Protocol (NTP) server",
				IP: []string{
					"57.199.209.217",
					"57.199.209.218",
				},
				FQDN: []string{
					"ntp1.service.cav",
					"ntp2.service.cav",
				},
				Ports: []types.ModelCloudavenueServiceDetailServicePort{
					{
						Port:     123,
						Protocol: "udp",
					},
				},
			},
			{
				Name:             "dns-authoritative",
				Description:      "DNS authoritative server. Use for resolving cloudavenue services names",
				DocumentationURL: "https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/practical-sheets/services-area/services-en/service-zone-dns/",
				IP: []string{
					"57.199.209.207",
					"57.199.209.208",
				},
				FQDN: nil,
				Ports: []types.ModelCloudavenueServiceDetailServicePort{
					{
						Port:     53,
						Protocol: "tcp",
					},
					{
						Port:     53,
						Protocol: "udp",
					},
				},
			},
			{
				Name:             "dns-resolver",
				Description:      "DNS resolver. Use for resolving cloudavenue services names and public names",
				DocumentationURL: "https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/practical-sheets/services-area/services-en/service-zone-dns/",
				IP: []string{
					"57.199.209.220",
					"57.199.209.221",
				},
				FQDN: nil,
				Ports: []types.ModelCloudavenueServiceDetailServicePort{
					{
						Port:     53,
						Protocol: "tcp",
					},
					{
						Port:     53,
						Protocol: "udp",
					},
				},
			},
			{
				Name:             "smtp",
				Description:      "SMTP relay. Use for sending emails",
				DocumentationURL: "https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/practical-sheets/services-area/services-en/smtp-service-2/",
				IP: []string{
					"57.199.209.206",
				},
				FQDN: []string{"smtp.service.cav"},
				Ports: []types.ModelCloudavenueServiceDetailServicePort{
					{
						Port:     25,
						Protocol: "tcp",
					},
				},
			},
		},
	},
	{
		Category: "s3",
		Network:  "194.206.55.5/32",
		Services: []types.ModelCloudavenueServiceDetailService{
			{
				Name:             "s3-internal",
				Description:      "S3 internal service. Use for accessing S3 directly from the organization",
				DocumentationURL: "https://cloud.orange-business.com/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/fiches-pratiques/stockage/stockage-objet-s3/guide-de-demarrage/premiere-utilisation-stockage-objet/",
				IP: []string{
					"194.206.55.5",
				},
				FQDN: []string{"s3-region01-priv.cloudavenue.orange-business.com"},
				Ports: []types.ModelCloudavenueServiceDetailServicePort{
					{
						Port:     443,
						Protocol: "tcp",
					},
				},
			},
		},
	},
}
