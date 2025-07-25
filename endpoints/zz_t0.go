package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

func ListT0() *cav.Endpoint {
	return cav.MustGetEndpoint("ListT0")
}
