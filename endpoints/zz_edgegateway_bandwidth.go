package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// UpdateEdgeGatewayBandwidth - Update EdgeGateway Bandwidth
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/put_api_customers_v2_0_edges__edge_id_ 
func UpdateEdgeGatewayBandwidth() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateEdgeGatewayBandwidth")
}
