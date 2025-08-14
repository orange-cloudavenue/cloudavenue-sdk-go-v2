package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// ListStorageProfile - List VDC Storage Profiles
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/queries/orgVdcStorageProfile.html 
func ListStorageProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("ListStorageProfile")
}
