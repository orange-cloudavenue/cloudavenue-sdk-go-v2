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
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path vdcnetwork.go -output vdcnetwork

// VDCNetwork endpoints back Org VDC Networks (both routed and isolated),
// differentiated only by the "networkType" field on the wire and by the
// optional "connection" field (present only for routed networks). Ownership
// is expressed via a single OwnerRef, which accepts either a VDC URN or a
// VDCGroup URN. Unlike FirewallGroup and AppPortProfile, VCD enforces a
// maximum pageSize of 32 for this resource. Create/Update/Delete are
// synchronous (no async job).
func init() {
	// ListVDCNetwork
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/get/",
		Name:             "ListVDCNetwork",
		Description:      "List Org VDC Networks (routed and isolated)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/orgVdcNetworks",
		QueryParams: []cav.QueryParam{
			{
				Name:        queryParamFilter,
				Description: "Filter to apply to the list of Org VDC Networks. Format: key==value;key==value. Allowed keys: name, ownerRef.id, id.",
			},
			{
				Name:        queryParamPageSize,
				Description: "The number of items per page. VCD enforces a maximum of 32 for this resource.",
				Value:       pageSize32,
			},
		},
		ResponseType: itypes.APIResponseListVDCNetwork{},
	}.Register()

	// GetVDCNetwork
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/vdcNetworkId/get/",
		Name:             "GetVDCNetwork",
		Description:      "Get an Org VDC Network (routed or isolated)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathOrgVDCNetworks,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVDCNetworkID,
				Description: "ID of the Org VDC Network to get",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, urnNetwork)
				},
			},
		},
		ResponseType: itypes.APIResponseVDCNetwork{},
	}.Register()

	// CreateVDCNetwork
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/post/",
		Name:             "CreateVDCNetwork",
		Description:      "Create an Org VDC Network (routed or isolated)",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/orgVdcNetworks",
		BodyRequestType:  itypes.APIRequestVDCNetwork{},
		ResponseType:     itypes.APIResponseVDCNetwork{},
	}.Register()

	// UpdateVDCNetwork
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/vdcNetworkId/put/",
		Name:             "UpdateVDCNetwork",
		Description:      "Update an Org VDC Network (routed or isolated)",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathOrgVDCNetworks,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVDCNetworkID,
				Description: "ID of the Org VDC Network to update",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, urnNetwork)
				},
			},
		},
		BodyRequestType: itypes.APIRequestVDCNetwork{},
		ResponseType:    itypes.APIResponseVDCNetwork{},
	}.Register()

	// DeleteVDCNetwork
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/vdcNetworkId/delete/",
		Name:             "DeleteVDCNetwork",
		Description:      "Delete an Org VDC Network (routed or isolated)",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathOrgVDCNetworks,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVDCNetworkID,
				Description: "ID of the Org VDC Network to delete",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, urnNetwork)
				},
			},
		},
	}.Register()
}
