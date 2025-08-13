package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// ListT0 - List T0
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/getNetworkHierarchy 
func ListT0() *cav.Endpoint {
	return cav.MustGetEndpoint("ListT0")
}
