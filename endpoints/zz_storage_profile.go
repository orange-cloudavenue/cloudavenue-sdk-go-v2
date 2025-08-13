package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// ListStorageProfiles - List VDC Storage Profiles
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/queries/orgVdcStorageProfile.html 
func ListStorageProfiles() *cav.Endpoint {
	return cav.MustGetEndpoint("ListStorageProfiles")
}
