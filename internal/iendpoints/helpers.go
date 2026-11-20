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
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

const (
	headerAccept      = "Accept"
	headerContentType = "Content-Type"
	// VMware VCD negotiates both JSON and XML responses by API version. Keep
	// endpoint-specific media types versioned rather than relying on a
	// subclient default that may be overridden by an endpoint.
	headerXML  = "application/xml;version=39.1"
	headerJSON = "application/json;version=39.1"
	// Vendor media types required by the legacy /api entity endpoints: the
	// generic application/xml and application/json accepts are rejected with
	// HTTP 406 by Cloud Director.
	headerVAppJSON              = "application/vnd.vmware.vcloud.vApp+json;version=39.1"
	headerVAppNetworkSectionXML = "application/vnd.vmware.vcloud.networkConfigSection+xml;version=39.1"
	// headerVAppControlAccessXML is the vApp ACL media type: the generic XML
	// accept is rejected with HTTP 406 by Cloud Director.
	headerVAppControlAccessXML = "application/vnd.vmware.vcloud.controlAccess+xml;version=39.1"
	// headerXMLAdminUser is the admin user media type: the admin user record
	// rejects the generic XML accept with HTTP 406.
	headerXMLAdminUser = "application/vnd.vmware.admin.user+xml;version=39.1"
	// headerQueryRecordsJSON is the query API media type: the query endpoint
	// rejects the generic JSON accept with HTTP 406.
	headerQueryRecordsJSON = "application/vnd.vmware.vcloud.query.records+json;version=39.1"
	// headerXMLAny is a permissive XML accept for operations whose response
	// media type varies (e.g. instantiate returns a task or a vApp reference).
	headerXMLAny = "application/*+xml;version=39.1"
	// headerInstantiateParamsXML is the mandatory request media type of the
	// instantiateVAppTemplate action; the generic XML content type yields 415.
	headerInstantiateParamsXML = "application/vnd.vmware.vcloud.instantiateVAppTemplateParams+xml;version=39.1"
	// headerRecomposeVAppParamsXML is the mandatory request media type of the
	// recomposeVApp action (vApp name and description update); the generic XML
	// content type yields 415.
	headerRecomposeVAppParamsXML = "application/vnd.vmware.vcloud.recomposeVAppParams+xml;version=39.1"
	// headerVAppLeaseSettingsSectionXML is the mandatory request media type of the
	// lease settings section update; the generic XML content type yields 415.
	headerVAppLeaseSettingsSectionXML = "application/vnd.vmware.vcloud.leaseSettingsSection+xml;version=39.1"
	// headerComposeVAppParamsXML is the mandatory request media type of the
	// composeVApp action; the generic XML content type yields 415.
	headerComposeVAppParamsXML = "application/vnd.vmware.vcloud.composeVAppParams+xml;version=39.1"
	// headerUndeployVAppParamsXML is the mandatory request media type of the
	// undeploy action; the generic XML content type yields 415.
	headerUndeployVAppParamsXML = "application/vnd.vmware.vcloud.undeployVAppParams+xml;version=39.1"
	// headerProductSectionsXML is the productSections XML media type, the
	// reference representation of the vApp guest properties.
	headerProductSectionsXML = "application/vnd.vmware.vcloud.productSections+xml;version=39.1"
	descUserURN              = "User URN"

	queryParamFilter   = "filter"
	queryParamQ        = "q"
	queryParamFormat   = "format"
	queryParamPage     = "page"
	queryParamPageSize = "pageSize"
	queryParamSortAsc  = "sortAsc"
	queryParamType     = "type"

	pageSize32  = "32"
	pageSize30  = "30"
	pageSize100 = "100"
	pageSize128 = "128"

	formatRecords           = "records"
	sortAscName             = "name"
	typeEdgeGateway         = "edgeGateway"
	typeOrgVDC              = "orgVdc"
	typeOrgVDCStoragePolicy = "orgVdcStoragePolicy"
	typeVApp                = "vApp"

	descPageSize                 = "The number of items per page."
	descPage                     = "Page to fetch, zero offset."
	descIAMPage                  = "Page to fetch"
	descIAMPageSize              = "Results per page to fetch"
	descFormatResponse           = "The format of the response."
	descTypeOfObjectQuery        = "The type of object to query"
	descEdgeGatewayID            = "The ID of the edge gateway."
	descCertificateLibraryItemID = "ID of the certificate library item"
	descFilterNameOrID           = "Filter to apply to the list of VDCs. Format: key==value. Allowed keys: name, id."
	descNetworkContextProfileID  = "ID of the Network Context Profile"
	descTrustedCertificateID     = "ID of the trusted certificate"
	descVDCGroupID               = "ID of the VDC Group"
	descVDCID                    = "The ID of the VDC."
	descVAppID                   = "The ID of the VApp."

	errFilterFormatSingle   = "filter must be in the format 'key==value'"
	errFilterFormatMultiple = "filter must be in the format 'key==value' or 'key1==value1;key2==value2'"
	errFilterKeyNotAllowed  = "filter key '%s' is not allowed"

	urnApplicationPortProfile   = "urn=applicationPortProfile"
	urnCertificateLibraryItem   = "urn=certificateLibraryItem"
	urnEdgeGateway              = "urn=edgegateway"
	urnFirewallGroup            = "urn=firewallGroup"
	urnNetwork                  = "urn=network"
	urnVDC                      = "urn=vdc"
	urnVDCGroup                 = "urn=vdcGroup"
	urnVDCNetwork               = "urn=vdcNetwork"
	urnVDCStoragePolicy         = "urn=vdcstorageProfile"
	urnVApp                     = "urn=vapp"
	urnOrg                      = "urn=org"
	urnEdgeGatewayID            = "urn=edgegateway"
	urnCatalog                  = "urn=catalog"
	ruleRequiredURNEdgeGateway  = "required," + urnEdgeGateway
	ruleResourceNameEdgeGateway = "resource_name=edgegateway"

	descCatalogURN   = "URN of the catalog"
	descVDCNetworkID = "ID of the Org VDC Network"
	descOrgID        = "Organization ID"
	descUserID       = "User ID or name"
	descTokenID      = "Token ID"
	descRoleID       = "Role ID"
	descGlobalRoleID = "Global Role ID"
	queryParamVDC    = "vdc"

	descOrgIDAdmin         = "Organization ID"
	descEdgeGatewayIDAdmin = "Edge Gateway ID"
)

var filterKeysNameOrID = []string{sortAscName, "id"}

func pageSizeQueryParam(value string) cav.QueryParam {
	return cav.QueryParam{
		Name:        queryParamPageSize,
		Description: descPageSize,
		Value:       value,
	}
}

func formatRecordsQueryParam() cav.QueryParam {
	return cav.QueryParam{
		Name:        queryParamFormat,
		Description: descFormatResponse,
		Value:       formatRecords,
	}
}

func typeQueryParam(value string) cav.QueryParam {
	return cav.QueryParam{
		Name:        queryParamType,
		Description: descTypeOfObjectQuery,
		Value:       value,
	}
}

func validateRule(rule string) func(string) error {
	return func(value string) error {
		return validators.New().Var(value, rule)
	}
}

func validateSingleFilterAllowedKeys(value string, allowedKeys []string) error {
	valueSplit := strings.Split(value, "==")
	if len(valueSplit) != 2 {
		return errors.New(errFilterFormatSingle)
	}

	if !slices.Contains(allowedKeys, valueSplit[0]) {
		return fmt.Errorf(errFilterKeyNotAllowed, valueSplit[0])
	}

	return nil
}

func wrapFilterInParentheses(value string) (string, error) {
	return fmt.Sprintf("(%s)", value), nil
}
