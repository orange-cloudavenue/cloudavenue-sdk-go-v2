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
	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path vdc.go -output vdc

func init() {
	// ListVDC
	cav.Endpoint{
		DocumentationURL: docURLCerberus,
		Name:             "ListVDC",
		Description:      "List VDCs",
		Method:           cav.MethodGET,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusVDCs,
		ResponseType:     itypes.APIResponseListVDC{},
	}.Register()

	// GetVDC
	cav.Endpoint{
		DocumentationURL: docURLCerberus,
		Name:             "GetVDC",
		Description:      "Get VDC",
		Method:           cav.MethodGET,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusVDCByName,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVdcName,
				Description: "The name of the VDC.",
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseGetVDC{},
	}.Register()

	// GetVDCMetadata is retained for compatibility with clients that inspect the
	// VMware metadata endpoint directly. VDC reads no longer depend on it.
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-VdcMetadata.html",
		Name:             "GetVDCMetadata",
		Description:      "Get VDC Metadata",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVDCMetadata,
		PathParams: []cav.PathParam{{
			Name: pathParamVdcID, Description: descVDCID, Required: true,
			ValidatorFunc: func(value string) error { return validators.New().Var(value, urnVDC) },
			TransformFunc: extractor.ExtractUUID,
		}},
		ResponseType: itypes.APIResponseGetVDCMetadatas{},
	}.Register()

	// CreateVDC
	cav.Endpoint{
		DocumentationURL: docURLCerberus,
		Name:             "CreateVDC",
		Description:      "Create a new Org VDC",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusVDCs,
		BodyRequestType:  itypes.APIRequestCreateVDC{},
		ResponseType:     cav.Job{},
	}.Register()

	// UpdateVDC
	cav.Endpoint{
		DocumentationURL: docURLCerberus,
		Name:             "UpdateVDC",
		Description:      "Update an existing Org VDC",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusVDCByName,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVdcName,
				Description: "The name of the VDC to update.",
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestUpdateVDC{},
		ResponseType:    cav.Job{},
	}.Register()

	// DeleteVDC
	cav.Endpoint{
		DocumentationURL: docURLCerberus,
		Name:             "DeleteVDC",
		Description:      "Delete an existing Org VDC",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusVDCByName,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVdcName,
				Description: "The name of the VDC to delete.",
				Required:    true,
			},
		},
		ResponseType: cav.Job{},
	}.Register()
}
