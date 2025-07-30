package edgegateway

import (
	"encoding/json"
	"net/http"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/orange-cloudavenue/common-go/validators"
)

//go:generate endpoint-generator -path t0_endpoints.go

func init() {
	// GET - List all T0
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/getNetworkHierarchy",
		Name:             "ListT0",
		Description:      "List T0",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/network",
		BodyResponseType: apiResponseT0s{},

		QueryParams: []cav.QueryParam{
			// Query parameters are not used in this endpoint, but can be added
			// for the mock response if needed
			{
				Name:        "t0Name",
				Description: "The name of the T0",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "resource_name=t0")
				},
			},
			{
				Name:        "edgeGatewayName",
				Description: "The name of the Edge Gateway",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "resource_name=edgegateway")
				},
			},
			{
				Name:        "edgeGatewayID",
				Description: "The ID of the Edge Gateway",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=edgegateway")
				},
			},
		},
		MockResponseFunc: func(w http.ResponseWriter, r *http.Request) {
			t0Name := r.URL.Query().Get("t0Name")
			edgeGatewayName := r.URL.Query().Get("edgeGatewayName")
			edgeGatewayID := r.URL.Query().Get("edgeGatewayID")

			var data apiResponseT0s

			if t0Name == "" {
				t0 := apiResponseT0{}
				_ = generator.Struct(&t0)

				data = apiResponseT0s{
					t0,
				}
			} else {
				data = apiResponseT0s{
					{
						Type: "tier-0-vrf",
						Name: t0Name,
						Children: []apiResponseT0Children{
							{
								Type: "edge-gateway",
								Name: generator.MustGenerate("{resource_name:edgegateway}"),
								Properties: struct {
									RateLimit int    `json:"rateLimit,omitempty" fake:"5"`
									EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"` // The UUID of the edge gateway
								}{
									RateLimit: 5,
									EdgeUUID:  generator.MustGenerate("{urn:edgegateway}"),
								},
							},
						},
					},
				}
			}

			if edgeGatewayName != "" || edgeGatewayID != "" {
				// Append the edge gateway to the T0
				data[0].Children = append(data[0].Children, apiResponseT0Children{
					Type: "edge-gateway",
					Name: func() string {
						if edgeGatewayName != "" {
							return edgeGatewayName
						}
						return generator.MustGenerate("{resource_name:edgegateway}")
					}(),
					Properties: struct {
						RateLimit int    `json:"rateLimit,omitempty" fake:"5"`
						EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"` // The UUID of the edge gateway
					}{
						RateLimit: 5,
						EdgeUUID: func() string {
							if edgeGatewayID != "" {
								return edgeGatewayID
							}
							return generator.MustGenerate("{urn:edgegateway}")
						}(),
					},
				})
			}

			bodyEncoded, _ := json.Marshal(data)

			// Return a mock response
			w.Header().Set("Content-Type", "application/json")
			// ignore write body error for mock response
			w.Write(bodyEncoded) //nolint:errcheck
		},
	}.Register()
}
