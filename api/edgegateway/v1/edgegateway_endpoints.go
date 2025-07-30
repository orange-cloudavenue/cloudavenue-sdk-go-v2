package edgegateway

import (
	"fmt"
	"regexp"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/urn"
	"github.com/orange-cloudavenue/common-go/validators"
)

//go:generate endpoint-generator -path edgegateway_endpoints.go

func init() {
	// GetEdgeGateway
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/latest/cloudapi/1.0.0/edgeGateways/gatewayId/get/",
		Name:             "GetEdgeGateway",
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

	// QueryEdgeGateway
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

	// CreateEdgeGateway
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/createVdcEdge",
		Name:             "CreateEdgeGateway",
		Description:      "Create EdgeGateway",
		Method:           cav.MethodPOST,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/{vdc-type}/{vdc-name}/edges",
		PathParams: []cav.PathParam{
			{
				Name:        "vdc-type",
				Description: "The type of the VDC where the edge gateway will be created.",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "oneof=vdc vdcgroup")
				},
				TransformFunc: func(value string) (string, error) {
					switch value {
					case "vdc":
						return "vdcs", nil
					case "vdcgroup":
						return "vdcgroups", nil
					}
					return "", fmt.Errorf("invalid vdc-type: %s", value)
				},
			},
			{
				Name:        "vdc-name",
				Description: "The name of the VDC where the edge gateway will be created.",
				Required:    true,
			},
		},
		QueryParams:      nil,
		BodyRequestType:  apiRequestEdgeGateway{},
		BodyResponseType: cav.Job{},
	}.Register()

	// DeleteEdgeGateway
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/deleteEdge",
		Name:             "DeleteEdgeGateway",
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
					return validators.New().Var(value, "required,urn=edgegateway")
				},
				TransformFunc: func(value string) (string, error) {
					// Transform the value to a uuidv4 format
					return extractor.ExtractUUID(value)
				},
			},
		},
		QueryParams:      nil,
		BodyRequestType:  nil,
		BodyResponseType: cav.Job{},
	}.Register()

	// ListEdgeGateway
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/latest/cloudapi/1.0.0/edgeGateways/get/",
		Name:             "ListEdgeGateway",
		Description:      "List EdgeGateways",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientVmware,
		PathTemplate:     "/cloudapi/1.0.0/edgeGateways",
		PathParams:       nil,
		QueryParams: []cav.QueryParam{
			{
				Name:        "pageSize",
				Description: "The number of items to return per page.",
				Value:       "128",
			},
		},
		BodyResponseType: apiResponseEdgegateways{},
	}.Register()
}
