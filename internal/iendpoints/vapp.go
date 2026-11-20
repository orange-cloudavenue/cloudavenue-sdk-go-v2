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
	"fmt"
	"strings"

	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/urn"

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
		Headers:          map[string]string{headerAccept: headerQueryRecordsJSON},
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

	// DiscardSuspendedStateVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-DiscardVAppState.html",
		Name:             "DiscardSuspendedStateVApp",
		Description:      "Discard suspended state of a VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathVAppDiscardSuspendedState,
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
		Headers:          map[string]string{headerAccept: headerVAppJSON},
		PathTemplate:     pathVAppByID,
		PathParams: []cav.PathParam{
			{
				Name:          pathParamVAppID,
				Description:   descVAppID,
				Required:      true,
				ValidatorFunc: validateRule(urnVApp),
				TransformFunc: toLegacyVAppID,
			},
		},
		ResponseType: itypes.APIResponseGetVApp{},
	}.Register()

	// GetVAppACL
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-ControlAccessParams-vApp.html",
		Name:             "GetVAppACL",
		Description:      "Get vApp ACL",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerVAppControlAccessXML, headerContentType: headerVAppControlAccessXML},
		PathTemplate:     pathVAppControlAccess,
		PathParams:       vAppActionPathParams(),
		ResponseType:     itypes.ControlAccessParams{},
	}.Register()

	// SetVAppACL
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-ControlAccess.html",
		Name:             "SetVAppACL",
		Description:      "Set vApp ACL",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerVAppControlAccessXML, headerContentType: headerVAppControlAccessXML},
		PathTemplate:     pathVAppControlAccessAction,
		PathParams:       vAppActionPathParams(),
		BodyRequestType:  itypes.ControlAccessParams{},
		ResponseType:     itypes.ControlAccessParams{},
	}.Register()

	// GetVAppNetworkConfigSection
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-NetworkConfigSection-vApp.html",
		Name:             "GetVAppNetworkConfigSection",
		Description:      "Get the runtime network configuration section of a VApp",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerVAppNetworkSectionXML},
		PathTemplate:     pathVAppNetworkConfigSection,
		PathParams:       vAppActionPathParams(),
		ResponseType:     itypes.APIResponseVAppNetworkConfigSection{},
	}.Register()

	// UpdateVAppNetworkConfigSection
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-NetworkConfigSection-vApp.html",
		Name:             "UpdateVAppNetworkConfigSection",
		Description:      "Replace the runtime network configuration section of a VApp",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXMLAny, headerContentType: headerVAppNetworkSectionXML},
		PathTemplate:     pathVAppNetworkConfigSection,
		PathParams:       vAppActionPathParams(),
		BodyRequestType:  itypes.APIResponseVAppNetworkConfigSection{},
		ResponseType:     cav.Job{},
	}.Register()

	// ComposeVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-ComposeVApp.html",
		Name:             "ComposeVApp",
		Description:      "Compose a new VApp",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXMLAny, headerContentType: headerComposeVAppParamsXML},
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
		Headers:          map[string]string{headerAccept: headerXMLAny, headerContentType: headerInstantiateParamsXML},
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

	// RecomposeVApp updates the description of a vApp. A PUT on /api/vApp/{id}
	// answers 415 on the real API: the name and description are changed with the
	// recomposeVApp action (same contract as go-vcloud-director UpdateNameDescription).
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-RecomposeVApp.html",
		Name:             "RecomposeVApp",
		Description:      "Update the name and the description of a VApp (recompose action)",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXMLAny, headerContentType: headerRecomposeVAppParamsXML},
		PathTemplate:     pathVAppRecompose,
		PathParams:       vAppActionPathParams(),
		BodyRequestType:  itypes.APIRequestRecomposeVApp{},
		ResponseType:     cav.Job{},
	}.Register()

	// UpdateVAppLeaseSettingsSection replaces the lease settings of a vApp.
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-LeaseSettingsSection.html",
		Name:             "UpdateVAppLeaseSettingsSection",
		Description:      "Replace the lease settings section of a VApp",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXMLAny, headerContentType: headerVAppLeaseSettingsSectionXML},
		PathTemplate:     pathVAppLeaseSettingsSection,
		PathParams:       vAppActionPathParams(),
		BodyRequestType:  itypes.APIRequestLeaseSettingsSection{},
		ResponseType:     cav.Job{},
	}.Register()

	// DeleteVApp
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/DELETE-VApp.html",
		Name:             "DeleteVApp",
		Description:      "Delete an existing VApp",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXMLAny},
		PathTemplate:     pathVAppByID,
		PathParams: []cav.PathParam{
			{
				Name:          pathParamVAppID,
				Description:   descVAppID,
				Required:      true,
				ValidatorFunc: validateRule(urnVApp),
				TransformFunc: toLegacyVAppID,
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
				TransformFunc: toLegacyVAppID,
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
				TransformFunc: toLegacyVAppID,
			},
		},
		Headers:         map[string]string{headerAccept: headerXMLAny, headerContentType: headerUndeployVAppParamsXML},
		BodyRequestType: itypes.APIRequestUndeployVApp{},
		ResponseType:    cav.Job{},
	}.Register()

	// GetVAppProductSections
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-ProductSections-vApp.html",
		Name:             "GetVAppProductSections",
		Description:      "Get the ProductSection guest properties of a VApp",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerProductSectionsJSON},
		PathTemplate:     pathVAppProductSections,
		PathParams:       vAppActionPathParams(),
		ResponseType:     itypes.APIResponseGuestProperties{},
	}.Register()

	// UpdateVAppProductSections
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-ProductSections-vApp.html",
		Name:             "UpdateVAppProductSections",
		Description:      "Replace the ProductSection guest properties of a VApp",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerProductSectionsJSON, headerContentType: headerProductSectionsJSON},
		PathTemplate:     pathVAppProductSections,
		PathParams:       vAppActionPathParams(),
		BodyRequestType:  itypes.APIResponseGuestProperties{},
		ResponseType:     cav.Job{},
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
			TransformFunc: toLegacyVAppID,
		},
	}
}

// toLegacyVAppID converts a VApp URN or bare UUID to the legacy
// "/api/vApp/vapp-<uuid>" path identifier. Cloud Director answers the bare
// UUID form with a misleading 403 "Access is forbidden".
func toLegacyVAppID(value string) (string, error) {
	if strings.HasPrefix(value, "vapp-") && urn.IsUUIDV4(strings.TrimPrefix(value, "vapp-")) {
		return value, nil
	}
	uuid := value
	switch {
	case urn.IsVAPP(value):
		uuid = urn.ExtractUUID(value)
	case urn.IsUUIDV4(value):
	default:
		return "", fmt.Errorf("expected a vApp URN or bare UUID, got %q", value)
	}
	return "vapp-" + uuid, nil
}
