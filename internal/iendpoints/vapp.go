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

//go:generate endpoint-generator -path vapp.go -output vapp

func init() {
	const pathVAppByID = "/api/vapp/{vapp-id}"

	// ListVapp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-VApp.html",
		Name:             "ListVapp",
		Description:      "List VApps",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/query",
		QueryParams: []cav.QueryParam{
			{
				Name:          queryParamFilter,
				Description:   descFilterNameOrID,
				ValidatorFunc: func(value string) error { return validateSingleFilterAllowedKeys(value, filterKeysNameOrID) },
				TransformFunc: wrapFilterInParentheses,
			},
			pageSizeQueryParam(pageSize100),
			formatRecordsQueryParam(),
			typeQueryParam(typeVApp),
		},
		ResponseType: itypes.APIResponseListVApp{},
	}.Register()

	// GetVapp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-VApp.html",
		Name:             "GetVapp",
		Description:      "Get VApp",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppByID,
		PathParams: []cav.PathParam{
			{
				Name:          pathParamVAppID,
				Description:   descVAppID,
				Required:      true,
				ValidatorFunc: validateRule(urnVApp),
				TransformFunc: extractUUID,
			},
		},
		ResponseType: itypes.APIResponseGetVApp{},
	}.Register()

	// CreateVapp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-VApp.html",
		Name:             "CreateVapp",
		Description:      "Create a new VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/vdc/{vdc-id}/action/createVApp",
		PathParams: []cav.PathParam{
			{
				Name:          pathParamVDCID,
				Description:   descVDCID,
				Required:      true,
				ValidatorFunc: validateRule(urnVDC),
				TransformFunc: extractUUID,
			},
		},
		BodyRequestType: itypes.APIRequestCreateVApp{},
		ResponseType:    cav.Job{},
	}.Register()

	// UpdateVapp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-VApp.html",
		Name:             "UpdateVapp",
		Description:      "Update an existing VApp",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppByID,
		PathParams: []cav.PathParam{
			{
				Name:          pathParamVAppID,
				Description:   descVAppID,
				Required:      true,
				ValidatorFunc: validateRule(urnVApp),
				TransformFunc: extractUUID,
			},
		},
		BodyRequestType: itypes.APIRequestUpdateVApp{},
		ResponseType:    cav.Job{},
	}.Register()

	// DeleteVapp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/DELETE-VApp.html",
		Name:             "DeleteVapp",
		Description:      "Delete an existing VApp",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppByID,
		PathParams: []cav.PathParam{
			{
				Name:          pathParamVAppID,
				Description:   descVAppID,
				Required:      true,
				ValidatorFunc: validateRule(urnVApp),
				TransformFunc: extractUUID,
			},
		},
		ResponseType: cav.Job{},
	}.Register()

	// RemoveAllNetworks
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-VApp-RemoveAllNetworks.html",
		Name:             "RemoveAllNetworks",
		Description:      "Remove all networks from a VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppByID + "/action/removeAllNetworks",
		PathParams: []cav.PathParam{
			{
				Name:          pathParamVAppID,
				Description:   descVAppID,
				Required:      true,
				ValidatorFunc: validateRule(urnVApp),
				TransformFunc: extractUUID,
			},
		},
		ResponseType: cav.Job{},
	}.Register()

	// UndeployVapp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-VApp-Undeploy.html",
		Name:             "UndeployVapp",
		Description:      "Undeploy a VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppByID + "/action/undeploy",
		PathParams: []cav.PathParam{
			{
				Name:          pathParamVAppID,
				Description:   descVAppID,
				Required:      true,
				ValidatorFunc: validateRule(urnVApp),
				TransformFunc: extractUUID,
			},
		},
		BodyRequestType: itypes.APIRequestUndeployVApp{},
		ResponseType:    cav.Job{},
	}.Register()
}
