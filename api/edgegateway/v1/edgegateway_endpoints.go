package edgegateway

import (
	"fmt"
	"regexp"
	"strconv"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/urn"
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

func init() {
	// GET - EdgeGateway
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/getEdgeById",
		Name:             "EdgeGateway",
		Description:      "Get EdgeGateway",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientVmware,
		PathTemplate:     "/cloudapi/1.0.0/edgeGateways/{edgeId}",
		PathParams: []cav.PathParam{
			{
				Name:        "edgeId",
				Description: "The ID of the edge gateway.",
				Required:    true,
			},
		},
		BodyResponseType: apiResponseEdgegateway{},
	}.Register()

	// GET - EdgeGateway Query
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/38.1/doc/types/QueryResultEdgeGatewayRecordType.html",
		Name:             "QueryEdgeGateway",
		Description:      "Query EdgeGateway",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientVmware,
		PathTemplate:     "/api/query",
		QueryParams: []cav.QueryParam{
			{
				Name:        "type",
				Description: "The type of object to query",
				Value:       "edgeGateway",
			},
			{
				Name:        "filter",
				Description: "The filter to apply to the query",
				Required:    false,
				ValidatorFunc: func(value string) error {
					// check if the value is a valid key==value pair
					x := regexp.MustCompile(`^[a-zA-Z0-9_]+==.*`)
					if !x.MatchString(value) {
						return fmt.Errorf("invalid filter format, expected key==value")
					}

					return nil
				},
			},
		},
		PathParams:       nil,
		BodyRequestType:  nil,
		BodyResponseType: apiResponseQueryEdgeGateway{},
		RequestMiddlewares: []resty.RequestMiddleware{
			func(_ *resty.Client, req *resty.Request) error {
				// Set the Accept header to application/*+json;version=38.1
				req.SetHeader("Accept", "application/*+json;version=38.1")
				return nil
			},
		},
		ResponseMiddlewares: []resty.ResponseMiddleware{
			func(_ *resty.Client, resp *resty.Response) error {
				r := resp.Result().(*apiResponseQueryEdgeGateway)

				if len(r.Record) == 0 {
					return fmt.Errorf("no edge gateways found")
				}

				id, err := extractor.ExtractUUID(r.Record[0].HREF)
				if err != nil {
					return err
				}

				r.Record[0].ID = urn.Normalize(urn.EdgeGateway, id).String()
				return nil
			},
		},
	}.Register()

	// GET - EdgeGateway Bandwidth
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/latest/cloudapi/1.0.0/edgeGateways/gatewayId/qos/get/",
		Name:             "Bandwidth",
		Description:      "Get EdgeGateway Bandwidth",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientVmware,
		PathTemplate:     "/cloudapi/1.0.0/edgeGateways/{edgeId}/qos",
		PathParams: []cav.PathParam{
			{
				Name:        "edgeId",
				Description: "The ID of the edge gateway.",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "required,urn=edgeGateway")
				},
			},
		},
		ResponseMiddlewares: []resty.ResponseMiddleware{
			func(_ *resty.Client, resp *resty.Response) error {
				r := resp.Result().(*apiResponseBandwidth)

				// bandwidth has format: qosgw005mbps
				// Extract the numeric part and convert it to int
				if r.EgressProfile.Name != "" {
					re := regexp.MustCompile(`qosgw(\d+)mbps`)
					matches := re.FindStringSubmatch(r.EgressProfile.Name)
					if len(matches) > 1 {
						// Convert the matched string to int
						bandwidth, err := strconv.Atoi(matches[1])
						if err != nil {
							r.Bandwidth = nil
						} else {
							r.Bandwidth = &bandwidth
						}
						return nil
					}

					// If profile name does not match the expected format
					// this is an error case
					return fmt.Errorf("invalid egress profile name format: %s", r.EgressProfile.Name)
				}
				// No egress profile, no bandwidth limit
				return nil
			},
		},
		QueryParams:      nil,
		BodyRequestType:  nil,
		BodyResponseType: apiResponseBandwidth{},
	}.Register()

	// POST - Create EdgeGateway from VDC
	// cav.Endpoint{
	// 	DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/createVdcEdge",
	// 	Name:             "EdgeGatewayVDC",
	// 	Description:      "Create EdgeGateway from VDC",
	// 	Method:           cav.MethodPOST,
	// 	SubClient:        cav.ClientVmware,
	// 	PathTemplate:     "/api/customers/v2.0/vdcs/{vdc-name}/edges",
	// 	PathParams: []cav.PathParam{
	// 		{
	// 			Name:        "vdc-name",
	// 			Description: "The name of the VDC where the edge gateway will be created.",
	// 			Required:    true,
	// 		},
	// 	},
	// 	QueryParams:      nil,
	// 	BodyRequestType:  edgegatewayAPICreateRequest{},
	// 	BodyResponseType: edgegatewayAPICreateResponse{},
	// }.Register()

	// Delete EdgeGateway
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/deleteEdge",
		Name:             "EdgeGateway",
		Description:      "Delete EdgeGateway",
		Method:           cav.MethodDELETE,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/edges/{edgeId}",
		PathParams: []cav.PathParam{
			{
				Name:        "edgeId",
				Description: "The ID of the edge gateway.",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "required,urn=edgeGateway")
				},
			},
		},
		QueryParams:      nil,
		BodyRequestType:  nil,
		BodyResponseType: cav.Job{},
	}.Register()

}
