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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path vapp.go -output vapp

func init() {
	// ListVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-ExecuteQuery.html",
		Name:             "ListVApp",
		Description:      "List VApps",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathQueryAPI,
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

	// DeployVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-DeployVApp.html",
		Name:             "DeployVApp",
		Description:      "Deploy a VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppDeploy,
		PathParams:       vAppActionPathParams(),
		ResponseType:     cav.Job{},
	}.Register()

	// PowerOffVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-PowerOffVApp.html",
		Name:             "PowerOffVApp",
		Description:      "Power off a VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppPowerOff,
		PathParams:       vAppActionPathParams(),
		ResponseType:     cav.Job{},
	}.Register()

	// PowerOnVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-PowerOnVApp.html",
		Name:             "PowerOnVApp",
		Description:      "Power on a VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppPowerOn,
		PathParams:       vAppActionPathParams(),
		ResponseType:     cav.Job{},
	}.Register()

	// RebootVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-RebootVApp.html",
		Name:             "RebootVApp",
		Description:      "Reboot a VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppReboot,
		PathParams:       vAppActionPathParams(),
		ResponseType:     cav.Job{},
	}.Register()

	// ResetVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-ResetVApp.html",
		Name:             "ResetVApp",
		Description:      "Reset a VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppReset,
		PathParams:       vAppActionPathParams(),
		ResponseType:     cav.Job{},
	}.Register()

	// ShutdownVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-ShutdownVApp.html",
		Name:             "ShutdownVApp",
		Description:      "Shut down a VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppShutdown,
		PathParams:       vAppActionPathParams(),
		ResponseType:     cav.Job{},
	}.Register()

	// SuspendVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-SuspendVApp.html",
		Name:             "SuspendVApp",
		Description:      "Suspend a VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppSuspend,
		PathParams:       vAppActionPathParams(),
		ResponseType:     cav.Job{},
	}.Register()

	// GetVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-VApp.html",
		Name:             "GetVApp",
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
				TransformFunc: extractor.ExtractUUID,
			},
		},
		ResponseType: itypes.APIResponseGetVApp{},
	}.Register()

	// ComposeVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-ComposeVApp.html",
		Name:             "ComposeVApp",
		Description:      "Compose a new VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXML, headerContentType: headerXML},
		PathTemplate:     pathVAppCompose,
		PathParams: []cav.PathParam{
			{
				Name:          pathParamVDCID,
				Description:   descVDCID,
				Required:      true,
				ValidatorFunc: validateRule(urnVDC),
				TransformFunc: extractor.ExtractUUID,
			},
		},
		BodyRequestType: itypes.APIRequestComposeVApp{},
		ResponseType:    cav.Job{},
	}.Register()

	// InstantiateVAppTemplate
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-InstantiateVAppTemplate.html",
		Name:             "InstantiateVAppTemplate",
		Description:      "Instantiate a VApp from a template",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXML, headerContentType: headerXML},
		PathTemplate:     pathVAppInstantiateTemplate,
		PathParams: []cav.PathParam{
			{
				Name:          pathParamVDCID,
				Description:   descVDCID,
				Required:      true,
				ValidatorFunc: validateRule(urnVDC),
				TransformFunc: extractor.ExtractUUID,
			},
		},
		BodyRequestType: itypes.APIRequestInstantiateVAppTemplate{},
		ResponseType:    cav.Job{},
	}.Register()

	// UpdateVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-Vm.html",
		Name:             "UpdateVApp",
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
				TransformFunc: extractor.ExtractUUID,
			},
		},
		BodyRequestType: itypes.APIRequestUpdateVApp{},
		ResponseType:    cav.Job{},
	}.Register()

	// DeleteVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/DELETE-VApp.html",
		Name:             "DeleteVApp",
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
				TransformFunc: extractor.ExtractUUID,
			},
		},
		ResponseType: cav.Job{},
	}.Register()

	// RemoveAllNetworks
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-NetworkConfigSection-vApp.html",
		Name:             "RemoveAllNetworks",
		Description:      "Remove all networks from a VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppRemoveAllNetworks,
		PathParams: []cav.PathParam{
			{
				Name:          pathParamVAppID,
				Description:   descVAppID,
				Required:      true,
				ValidatorFunc: validateRule(urnVApp),
				TransformFunc: extractor.ExtractUUID,
			},
		},
		ResponseType: cav.Job{},
	}.Register()

	// UndeployVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-UndeployVApp.html",
		Name:             "UndeployVApp",
		Description:      "Undeploy a VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppUndeploy,
		PathParams: []cav.PathParam{
			{
				Name:          pathParamVAppID,
				Description:   descVAppID,
				Required:      true,
				ValidatorFunc: validateRule(urnVApp),
				TransformFunc: extractor.ExtractUUID,
			},
		},
		BodyRequestType: itypes.APIRequestUndeployVApp{},
		ResponseType:    cav.Job{},
	}.Register()

	// GetVAppLeaseSettings
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-VAppLeaseSettings.html",
		Name:             "GetVAppLeaseSettings",
		Description:      "Get org/VDC-level lease settings",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppLeaseSettings,
		PathParams: []cav.PathParam{
			{
				Name:          pathParamOrgID,
				Description:   descOrgID,
				Required:      true,
				ValidatorFunc: validateRule(urnOrg),
			},
		},
		ResponseType: itypes.APIResponseOrgLeaseSettings{},
	}.Register()

	// UpdateVAppLeaseSettings
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-VAppLeaseSettings.html",
		Name:             "UpdateVAppLeaseSettings",
		Description:      "Update org/VDC-level lease settings",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppLeaseSettings,
		PathParams: []cav.PathParam{
			{
				Name:          pathParamOrgID,
				Description:   descOrgID,
				Required:      true,
				ValidatorFunc: validateRule(urnOrg),
			},
		},
		BodyRequestType: itypes.APIRequestOrgLeaseSettings{},
		ResponseType:    itypes.APIResponseOrgLeaseSettings{},
	}.Register()
}

func vAppActionPathParams() []cav.PathParam {
	return []cav.PathParam{
		{
			Name:          pathParamVAppID,
			Description:   descVAppID,
			Required:      true,
			ValidatorFunc: validateRule(urnVApp),
			TransformFunc: extractor.ExtractUUID,
		},
	}
}
