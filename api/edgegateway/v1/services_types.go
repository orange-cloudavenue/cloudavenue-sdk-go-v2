/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"net"

	"github.com/orange-cloudavenue/common-go/urn"
)

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

	// * apiResponse
	apiResponseNetworkServices []struct {
		Type     string                               `json:"type" fake:"tier-0-vrf"`
		Name     string                               `json:"name" fake:"{resource_name:t0}"`
		Children []apiResponseNetworkServicesChildren `json:"children,omitempty" fakesize:"1"`
	}

	apiResponseNetworkServicesChildren struct {
		Type        string `json:"type" fake:"edge-gateway"`
		Name        string `json:"name,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		Properties  struct {
			// EdgeGateway
			RateLimit int    `json:"rateLimit,omitempty"`
			EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgeGateway}"` // The UUID of the edge gateway
		} `json:"properties,omitempty"`
		Children  []apiResponseNetworkServicesSubChildren `json:"children,omitempty" fakesize:"6"`
		ServiceID string                                  `json:"serviceId,omitempty"`
	}
	apiResponseNetworkServicesSubChildren struct {
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

	// * apiRequest

	apiRequestNetworkServicesCavSvc struct {
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

func (ap *apiResponseNetworkServices) toModel(params ParamsEdgeGateway) *ModelEdgeGatewayServices {
	if ap == nil || len(*ap) == 0 {
		return nil
	}

	data := &ModelEdgeGatewayServices{
		Service:      nil,
		LoadBalancer: nil,
		PublicIP:     nil,
	}

	// Parse the original response and populate the NetworkServicesModel
	for _, ns := range *ap {
		for _, child := range ns.Children {
			if child.Type == "edge-gateway" && (child.Properties.EdgeUUID == urn.ExtractUUID(params.ID) || child.Name == params.Name) {
				// Found the edge gateway
				// iterate over the children to find the services
				for _, service := range child.Children {
					switch service.Type {
					case "load-balancer":
						// Found load balancer service
						data.LoadBalancer = &ModelEdgeGatewayServicesLoadBalancer{
							ModelEdgeGatewayServicesSvc: ModelEdgeGatewayServicesSvc{
								ID:   service.Name,        // The name is the ID
								Name: service.DisplayName, // The display name is the name
							},
							ClassOfService:     service.Properties.ClassOfService,
							MaxVirtualServices: service.Properties.MaxVirtualServices,
						}
					case "service":
						// service is a generic service
						// the name of the service define the type of service
						switch service.Name {
						case "cav-services", "cav_services": // Match both cav-services and cav_services
							// Found cav-services
							data.Service = &ModelCloudavenueServices{
								ModelEdgeGatewayServicesSvc: ModelEdgeGatewayServicesSvc{
									ID:   service.ServiceID,   // The ServiceID is the ID
									Name: service.DisplayName, // The display name is the name
								},
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
								ServiceDetails: ListOfServices,
							}

						case "internet":
							// Found internet service
							publicIP := &ModelEdgeGatewayServicesPublicIP{
								ModelEdgeGatewayServicesSvc: ModelEdgeGatewayServicesSvc{
									ID:   service.ServiceID,     // The ServiceID is the ID
									Name: service.Properties.IP, // The IP don't have a name use IP instead
								},
								IP:        service.Properties.IP,
								Announced: service.Properties.Announced,
							}

							// Prevent nil pointer dereference
							if data.PublicIP == nil {
								data.PublicIP = make([]*ModelEdgeGatewayServicesPublicIP, 0)
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
