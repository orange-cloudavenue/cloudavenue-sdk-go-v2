package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

func GetNetworkServices() *cav.Endpoint {
	return cav.MustGetEndpoint("GetNetworkServices")
}
func EnableCloudAvenueServices() *cav.Endpoint {
	return cav.MustGetEndpoint("EnableCloudAvenueServices")
}
func DisableCloudAvenueServices() *cav.Endpoint {
	return cav.MustGetEndpoint("DisableCloudAvenueServices")
}
