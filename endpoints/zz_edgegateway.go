package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

func GetEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("GetEdgeGateway")
}
func QueryEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("QueryEdgeGateway")
}
func GetEdgeGatewayBandwidth() *cav.Endpoint {
	return cav.MustGetEndpoint("GetEdgeGatewayBandwidth")
}
func CreateEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateEdgeGateway")
}
func DeleteEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteEdgeGateway")
}
func ListEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("ListEdgeGateway")
}
