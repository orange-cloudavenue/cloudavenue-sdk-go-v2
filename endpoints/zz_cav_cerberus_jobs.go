package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// GetJobCerberus - Get Cerberus Job
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Jobs/getJobById 
func GetJobCerberus() *cav.Endpoint {
	return cav.MustGetEndpoint("GetJobCerberus")
}
