package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// GetEdgeGatewayServices - Get EdgeGateway Network Services
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/getNetworkHierarchy 
func GetEdgeGatewayServices() *cav.Endpoint {
	return cav.MustGetEndpoint("GetEdgeGatewayServices")
}
// EnableCloudavenueServices - Enable Cloud Avenue Services
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/addNetworkConnectivity 
func EnableCloudavenueServices() *cav.Endpoint {
	return cav.MustGetEndpoint("EnableCloudavenueServices")
}
// DisableCloudavenueServices - Disable Cloud Avenue Services
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/deleteNetworkService 
func DisableCloudavenueServices() *cav.Endpoint {
	return cav.MustGetEndpoint("DisableCloudavenueServices")
}
