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
	// Canonical documentation roots. Per-operation links drift and are not
	// stable. VMware is pinned to 39.1 to match the API version the client sends.
	// OSE and NetBackup have no public documentation and keep their own URLs.
	//
	// These duplicate cav.DocURLVMware and cav.DocURLCerberus on purpose: the
	// endpoint generator renders a cross-package reference as "pkg.Name", not as
	// the value, so aliasing them would blank every generated comment here.
	headerAccept      = "Accept"
	headerContentType = "Content-Type"
	headerXML         = "application/xml"
	descUserURN       = "User URN"
	docURLVMware      = "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/"
	docURLCerberus    = "https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API"

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
