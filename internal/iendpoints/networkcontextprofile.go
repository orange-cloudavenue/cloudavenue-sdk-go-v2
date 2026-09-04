/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path networkcontextprofile.go -output networkcontextprofile

// NetworkContextProfile represents an NSX-T Layer-7 application/context
// profile that can be referenced by Distributed Firewall rules. Unlike the
// FirewallGroup family (which uses a single OwnerRef), NetworkContextProfile
// uses a dedicated pathParamVDCGroupID (or "orgVdcId" for plain VDCs, not used by
// this VDCGroup-scoped package) named filter key for scoping. Create and
// Update are asynchronous (VCD returns 202 + task); List/Get/Delete and the
// read-only "/attributes" static catalog sub-resource are synchronous.
func init() {
	// ListNetworkContextProfile
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/networkContextProfiles/get/",
		Name:             "ListNetworkContextProfile",
		Description:      "List Network Context Profiles",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/networkContextProfiles",
		QueryParams: []cav.QueryParam{
			{
				Name:        queryParamFilter,
				Description: "Filter to apply to the list of Network Context Profiles. Format: key==value;key==value. Allowed keys: vdcGroupId, name, scope.",
			},
			{
				Name:        queryParamPageSize,
				Description: descPageSize,
				Value:       pageSize100,
			},
		},
		ResponseType: itypes.APIResponseListNetworkContextProfile{},
	}.Register()

	// GetNetworkContextProfile
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/networkContextProfiles/networkContextProfileId/get/",
		Name:             "GetNetworkContextProfile",
		Description:      "Get a Network Context Profile",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathNetworkContextProfiles,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamNetworkContextProfileID,
				Description: "ID of the Network Context Profile to get",
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseNetworkContextProfile{},
	}.Register()

	// CreateNetworkContextProfile
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/networkContextProfiles/post/",
		Name:             "CreateNetworkContextProfile",
		Description:      "Create a Network Context Profile",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/networkContextProfiles",
		BodyRequestType:  itypes.APIRequestNetworkContextProfile{},
		ResponseType:     cav.Job{},
	}.Register()

	// UpdateNetworkContextProfile
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/networkContextProfiles/networkContextProfileId/put/",
		Name:             "UpdateNetworkContextProfile",
		Description:      "Update a Network Context Profile",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathNetworkContextProfiles,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamNetworkContextProfileID,
				Description: "ID of the Network Context Profile to update",
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestNetworkContextProfile{},
		ResponseType:    cav.Job{},
	}.Register()

	// DeleteNetworkContextProfile
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/networkContextProfiles/networkContextProfileId/delete/",
		Name:             "DeleteNetworkContextProfile",
		Description:      "Delete a Network Context Profile",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathNetworkContextProfiles,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamNetworkContextProfileID,
				Description: "ID of the Network Context Profile to delete",
				Required:    true,
			},
		},
	}.Register()

	// GetNetworkContextProfileAttributes
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/networkContextProfiles/attributes/get/",
		Name:             "GetNetworkContextProfileAttributes",
		Description:      "Get the static reference catalog (App IDs, Domain Names) of attributes usable in Network Context Profiles",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/networkContextProfiles/attributes",
		QueryParams: []cav.QueryParam{
			{
				Name:        queryParamFilter,
				Description: "Filter to apply. Format: key==value. Allowed keys: vdcGroupId.",
			},
		},
		ResponseType: itypes.APINetworkContextProfileAttributesResponse{},
	}.Register()
}
