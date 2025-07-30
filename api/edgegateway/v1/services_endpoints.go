package edgegateway

import (
	"encoding/json"
	"net/http"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

//go:generate endpoint-generator -path services_endpoints.go

func init() {
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/getNetworkHierarchy",
		Name:             "GetEdgeGatewayServices",
		Description:      "Get EdgeGateway Network Services",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/network",
		BodyResponseType: apiResponseNetworkServices{},
		QueryParams: []cav.QueryParam{
			// Query parameters are not used in this endpoint, but can be added
			// for the mock response if needed
			{
				Name:        "edgeId",
				Description: "The ID of the edge gateway to get network services for",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=edgeGateway")
				},
			},
			{
				Name:        "edgeName",
				Description: "The name of the edge gateway to get network services for",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "resource_name=edgegateway")
				},
			},
		},
		RequestMiddlewares: []resty.RequestMiddleware{
			func(_ *resty.Client, req *resty.Request) error {
				edgeID := req.QueryParams.Get("edgeId")
				if edgeID != "" {
					// Remove the URN prefix if it exists
					id, _ := extractor.ExtractUUID(edgeID)
					// Set the edgeId query parameter to the extracted ID
					req.QueryParams.Set("edgeId", id)
				}
				return nil
			},
		},
		MockResponseFunc: func(w http.ResponseWriter, r *http.Request) {
			// One of the two must be filled in. The validator makes sure of this.
			edgeID := r.URL.Query().Get("edgeId")
			edgeName := r.URL.Query().Get("edgeName")

			var data apiResponseNetworkServices

			data = apiResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Name: generator.MustGenerate("{resource_name:t0}"),
					Children: []apiResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: edgeName,
							Properties: struct {
								// EdgeGateway
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"` // The UUID of the edge gateway
							}{
								RateLimit: 5,
								EdgeUUID:  edgeID,
							},
							Children: []apiResponseNetworkServicesSubChildren{
								{
									Type: "load-balancer",
									Name: generator.MustGenerate("{uuid}"),
									Properties: struct {
										// Load Balancer
										ClassOfService     string `json:"classOfService,omitempty"`
										MaxVirtualServices int    `json:"maxVirtualServices,omitempty"`

										// Public IP
										IP        string `json:"ip,omitempty" fake:"{ipv4address}"`
										Announced bool   `json:"announced,omitempty" fake:"true"`

										// Service
										Ranges []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"` // The network in ip/cidr format
									}{
										ClassOfService:     "PREMIUM",
										MaxVirtualServices: 10,
									},
								},
								{
									ServiceID:   generator.MustGenerate("ip-{regex:[1-9]{2}}-{regex:[1-9]{2}}-{regex:[1-9]{2}}-{regex:[1-9]{2}}"),
									Type:        "service",
									Name:        "internet",
									DisplayName: "internet",
									Properties: struct {
										// Load Balancer
										ClassOfService     string `json:"classOfService,omitempty"`
										MaxVirtualServices int    `json:"maxVirtualServices,omitempty"`

										// Public IP
										IP        string `json:"ip,omitempty" fake:"{ipv4address}"`
										Announced bool   `json:"announced,omitempty" fake:"true"`

										// Service
										Ranges []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"` // The network in ip/cidr format
									}{
										IP:        generator.MustGenerate("{ipv4address}"),
										Announced: true,
									},
								},
								{
									ServiceID:   generator.MustGenerate("{resource_name:edgegateway}-cav-services"),
									Type:        "service",
									Name:        "cav-services",
									DisplayName: "Cloud Avenue Services",
									Properties: struct {
										// Load Balancer
										ClassOfService     string `json:"classOfService,omitempty"`
										MaxVirtualServices int    `json:"maxVirtualServices,omitempty"`

										// Public IP
										IP        string `json:"ip,omitempty" fake:"{ipv4address}"`
										Announced bool   `json:"announced,omitempty" fake:"true"`

										// Service
										Ranges []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"` // The network in ip/cidr format
									}{
										Ranges: []string{
											generator.MustGenerate("{ipv4address}/{intrange:24,32}"),
										},
									},
								},
							},
						},
					},
				},
			}

			bodyEncoded, err := json.Marshal(data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Return a mock response
			w.Header().Set("Content-Type", "application/json")
			// ignore write body error for mock response
			w.Write(bodyEncoded) //nolint:errcheck
		},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/addNetworkConnectivity",
		Name:             "EnableCloudavenueServices",
		Description:      "Enable Cloud Avenue Services",
		Method:           cav.MethodPOST,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/services",
		BodyResponseType: cav.Job{},
		BodyRequestType:  apiRequestNetworkServicesCavSvc{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/deleteNetworkService",
		Name:             "DisableCloudavenueServices",
		Description:      "Disable Cloud Avenue Services",
		Method:           cav.MethodDELETE,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/services/{serviceId}",
		PathParams: []cav.PathParam{
			{
				Name:        "serviceId",
				Description: "The ID of the service to delete",
				Required:    true,
			},
		},
		BodyResponseType: cav.Job{},
	}.Register()
}
