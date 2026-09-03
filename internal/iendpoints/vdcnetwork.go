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

// VdcNetwork endpoints back Org VDC Networks (both routed and isolated),
// differentiated only by the "networkType" field on the wire and by the
// optional "connection" field (present only for routed networks). Ownership
// is expressed via a single OwnerRef, which accepts either a VDC URN or a
// VdcGroup URN. Unlike FirewallGroup and AppPortProfile, VCD enforces a
// maximum pageSize of 32 for this resource. Create/Update/Delete are
// synchronous (no async job).
func init() {
	// ListVdcNetwork
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/get/",
		Name:             "ListVdcNetwork",
		Description:      "List Org VDC Networks (routed and isolated)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/orgVdcNetworks",
		QueryParams: []cav.QueryParam{
			{
				Name:        "filter",
				Description: "Filter to apply to the list of Org VDC Networks. Format: key==value;key==value. Allowed keys: name, ownerRef.id, id.",
			},
			{
				Name:        "pageSize",
				Description: "The number of items per page. VCD enforces a maximum of 32 for this resource.",
				Value:       "32",
			},
		},
		ResponseType: itypes.ApiResponseListVdcNetwork{},
	}.Register()

	// GetVdcNetwork
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/vdcNetworkId/get/",
		Name:             "GetVdcNetwork",
		Description:      "Get an Org VDC Network (routed or isolated)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/orgVdcNetworks/{vdcNetworkId}",
		PathParams: []cav.PathParam{
			{
				Name:        "vdcNetworkId",
				Description: "ID of the Org VDC Network to get",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=network")
				},
			},
		},
		ResponseType: itypes.ApiResponseVdcNetwork{},
	}.Register()

	// CreateVdcNetwork
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/post/",
		Name:             "CreateVdcNetwork",
		Description:      "Create an Org VDC Network (routed or isolated)",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/orgVdcNetworks",
		BodyRequestType:  itypes.ApiRequestVdcNetwork{},
		ResponseType:     itypes.ApiResponseVdcNetwork{},
	}.Register()

	// UpdateVdcNetwork
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/vdcNetworkId/put/",
		Name:             "UpdateVdcNetwork",
		Description:      "Update an Org VDC Network (routed or isolated)",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/orgVdcNetworks/{vdcNetworkId}",
		PathParams: []cav.PathParam{
			{
				Name:        "vdcNetworkId",
				Description: "ID of the Org VDC Network to update",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=network")
				},
			},
		},
		BodyRequestType: itypes.ApiRequestVdcNetwork{},
		ResponseType:    itypes.ApiResponseVdcNetwork{},
	}.Register()

	// DeleteVdcNetwork
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/vdcNetworkId/delete/",
		Name:             "DeleteVdcNetwork",
		Description:      "Delete an Org VDC Network (routed or isolated)",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/orgVdcNetworks/{vdcNetworkId}",
		PathParams: []cav.PathParam{
			{
				Name:        "vdcNetworkId",
				Description: "ID of the Org VDC Network to delete",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=network")
				},
			},
		},
	}.Register()
}
