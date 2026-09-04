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
	queryParamFilter   = "filter"
	queryParamFormat   = "format"
	queryParamPageSize = "pageSize"
	queryParamSortAsc  = "sortAsc"
	queryParamType     = "type"

	pathParamEdgeID                  = "edgeId"
	pathParamNetworkContextProfileID = "networkContextProfileId"
	pathParamServiceID               = "serviceId"
	pathParamTrustedCertificate      = "trustedCertificate"
	pathParamVDCGroupID              = "vdcGroupId"
	pathParamVDCID                   = "vdc-id"
	pathParamVDCName                 = "vdc-name"
	pathParamVAppID                  = "vapp-id"

	pageSize32  = "32"
	pageSize30  = "30"
	pageSize100 = "100"
	pageSize128 = "128"

	formatRecords            = "records"
	sortAscName              = "name"
	typeEdgeGateway          = "edgeGateway"
	typeOrgVDC               = "orgVdc"
	typeOrgVDCStorageProfile = "orgVdcStorageProfile"
	typeVApp                 = "vApp"

	descPageSize                 = "The number of items per page."
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
	urnVDCStorageProfile        = "urn=vdcstorageProfile"
	urnVApp                     = "urn=vapp"
	ruleRequiredURNEdgeGateway  = "required," + urnEdgeGateway
	ruleResourceNameEdgeGateway = "resource_name=edgegateway"

	pathQueryAPI                    = "/api/query"
	pathApplicationPortProfiles     = "/cloudapi/1.0.0/applicationPortProfiles/{appPortProfileId}"
	pathParamAppPortProfileID       = "appPortProfileId"
	pathCertificateLibrary          = "/cloudapi/1.0.0/ssl/certificateLibrary/{id}"
	pathParamCertLibraryItemID      = "certLibraryItemId"
	pathCertificateLibraryConsumers = "/cloudapi/1.0.0/ssl/certificateLibrary/{certLibraryItemId}/consumers"
	pathFirewallGroups              = "/cloudapi/1.0.0/firewallGroups/{firewallGroupId}"
	pathParamFirewallGroupID        = "firewallGroupId"
	pathNetworkContextProfiles      = "/cloudapi/1.0.0/networkContextProfiles/{networkContextProfileId}"
	pathTrustedCertificates         = "/cloudapi/1.0.0/ssl/trustedCertificates/{trustedCertificate}"
	pathOrgVDCNetworks              = "/cloudapi/1.0.0/orgVdcNetworks/{vdcNetworkId}"
	pathParamVDCNetworkID           = "vdcNetworkId"
	pathParamOrgID                  = "orgId"
	descOrgID                       = "Organization ID"
	pathParamUserID                 = "userId"
	descUserID                      = "User ID or name"
	queryParamVDC                   = "vdc"
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
