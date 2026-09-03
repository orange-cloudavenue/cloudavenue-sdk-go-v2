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

	"github.com/orange-cloudavenue/common-go/extractor"
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

	pathCustomersEdgesByID = "/api/customers/v2.0/edges/{edgeId}"
	pathCustomersVDCByName = "/api/customers/v2.0/vdcs/{vdc-name}"

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

	ruleResourceNameEdgeGateway = "resource_name=edgegateway"

	errFilterFormatSingle   = "filter must be in the format 'key==value'"
	errFilterFormatMultiple = "filter must be in the format 'key==value' or 'key1==value1;key2==value2'"
	errFilterKeyNotAllowed  = "filter key '%s' is not allowed"

	urnApplicationPortProfile  = "urn=applicationPortProfile"
	urnCertificateLibraryItem  = "urn=certificateLibraryItem"
	urnEdgeGateway             = "urn=edgegateway"
	urnFirewallGroup           = "urn=firewallGroup"
	urnNetwork                 = "urn=network"
	urnVDC                     = "urn=vdc"
	urnVDCGroup                = "urn=vdcGroup"
	urnVDCStorageProfile       = "urn=vdcstorageProfile"
	urnVApp                    = "urn=vapp"
	ruleRequiredURNEdgeGateway = "required," + urnEdgeGateway
)

var filterKeysNameOrID = []string{"name", "id"}

func requiredPathParam(name, description string) cav.PathParam {
	return cav.PathParam{
		Name:        name,
		Description: description,
		Required:    true,
	}
}

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

func filterQueryParam(description string) cav.QueryParam {
	return cav.QueryParam{
		Name:        queryParamFilter,
		Description: description,
	}
}

func requiredURNPathParam(name, description, rule string) cav.PathParam {
	param := requiredPathParam(name, description)
	param.ValidatorFunc = validateRule(rule)

	return param
}

func validateRule(rule string) func(string) error {
	return func(value string) error {
		return validators.New().Var(value, rule)
	}
}

func extractUUID(value string) (string, error) {
	return extractor.ExtractUUID(value)
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
