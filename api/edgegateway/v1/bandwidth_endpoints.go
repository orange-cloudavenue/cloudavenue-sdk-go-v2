package edgegateway

import (
	"time"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/validators"
)

//go:generate endpoint-generator -path bandwidth_endpoints.go

func init() {
	// // * GetEdgeGatewayBandwidth
	// cav.Endpoint{
	// 	DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/latest/cloudapi/1.0.0/edgeGateways/gatewayId/qos/get/",
	// 	Name:             "GetEdgeGatewayBandwidth",
	// 	Description:      "Get EdgeGateway Bandwidth",
	// 	Method:           cav.MethodGET,
	// 	SubClient:        cav.ClientVmware,
	// 	PathTemplate:     "/cloudapi/1.0.0/edgeGateways/{edgeId}/qos",
	// 	PathParams: []cav.PathParam{
	// 		{
	// 			Name:        "edgeId",
	// 			Description: "The ID of the edge gateway.",
	// 			Required:    true,
	// 			ValidatorFunc: func(value string) error {
	// 				return validators.New().Var(value, "required,urn=edgeGateway")
	// 			},
	// 		},
	// 	},
	// 	ResponseMiddlewares: []resty.ResponseMiddleware{
	// 		func(_ *resty.Client, resp *resty.Response) error {
	// 			r := resp.Result().(*apiResponseBandwidth)

	// 			// bandwidth has format: qosgw005mbps
	// 			// Extract the numeric part and convert it to int
	// 			if r.EgressProfile.Name != "" {
	// 				re := regexp.MustCompile(`qosgw(\d+)mbps`)
	// 				matches := re.FindStringSubmatch(r.EgressProfile.Name)
	// 				if len(matches) > 1 {
	// 					// Convert the matched string to int
	// 					bandwidth, err := strconv.Atoi(matches[1])
	// 					if err != nil {
	// 						r.Bandwidth = nil
	// 					} else {
	// 						r.Bandwidth = &bandwidth
	// 					}
	// 					return nil
	// 				}

	// 				// If profile name does not match the expected format
	// 				// this is an error case
	// 				return fmt.Errorf("invalid egress profile name format: %s", r.EgressProfile.Name)
	// 			}
	// 			// No egress profile, no bandwidth limit
	// 			return nil
	// 		},
	// 	},
	// 	QueryParams:      nil,
	// 	BodyRequestType:  nil,
	// 	BodyResponseType: apiResponseBandwidth{},
	// }.Register()

	// * UpdateEdgeGatewayBandwidth
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/put_api_customers_v2_0_edges__edge_id_",
		Name:             "UpdateEdgeGatewayBandwidth",
		Description:      "Update EdgeGateway Bandwidth",
		Method:           cav.MethodPUT,
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
				TransformFunc: func(value string) (string, error) {
					// Transform the value to a uuidv4 format
					return extractor.ExtractUUID(value)
				},
			},
		},
		QueryParams:      nil,
		BodyRequestType:  apiRequestBandwidth{},
		BodyResponseType: cav.Job{},
		JobOptions: &cav.JobOptions{
			PollInterval: time.Second * 1,
		},
	}.Register()
}
