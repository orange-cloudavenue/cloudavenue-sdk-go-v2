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

//go:generate endpoint-generator -path firewallgroup.go -output firewallgroup

// FirewallGroup endpoints back three distinct resources (Security Group,
// IP Set, and Dynamic Security Group) that all share the same NSX-T Firewall
// Group REST resource, differentiated only by the "typeValue" field on the
// wire and by the OwnerRef (which accepts either an EdgeGateway URN or a
// VdcGroup URN). Unlike most cloudapi resources, VCD does not follow the
// regular REST scheme for listing Firewall Groups: the list endpoint
// requires a "/summaries" suffix. Create/Update/Delete are synchronous
// (no async job).
func init() {
	// ListFirewallGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/firewallGroups/summaries/get/",
		Name:             "ListFirewallGroup",
		Description:      "List Firewall Groups (Security Groups, IP Sets, Dynamic Security Groups)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/firewallGroups/summaries",
		QueryParams: []cav.QueryParam{
			{
				Name:        "filter",
				Description: "Filter to apply to the list of Firewall Groups. Format: key==value;key==value. Allowed keys: typeValue, ownerRef.id, _context, name, id.",
			},
			{
				Name:        "pageSize",
				Description: "The number of items per page.",
				Value:       "100",
			},
		},
		ResponseType: itypes.ApiResponseListFirewallGroup{},
	}.Register()

	// GetFirewallGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/firewallGroups/firewallGroupId/get/",
		Name:             "GetFirewallGroup",
		Description:      "Get a Firewall Group (Security Group, IP Set, or Dynamic Security Group)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/firewallGroups/{firewallGroupId}",
		PathParams: []cav.PathParam{
			{
				Name:        "firewallGroupId",
				Description: "ID of the Firewall Group to get",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=firewallGroup")
				},
			},
		},
		ResponseType: itypes.ApiResponseFirewallGroup{},
	}.Register()

	// CreateFirewallGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/firewallGroups/post/",
		Name:             "CreateFirewallGroup",
		Description:      "Create a Firewall Group (Security Group, IP Set, or Dynamic Security Group)",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/firewallGroups",
		BodyRequestType:  itypes.ApiRequestFirewallGroup{},
		ResponseType:     itypes.ApiResponseFirewallGroup{},
	}.Register()

	// UpdateFirewallGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/firewallGroups/firewallGroupId/put/",
		Name:             "UpdateFirewallGroup",
		Description:      "Update a Firewall Group (Security Group, IP Set, or Dynamic Security Group)",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/firewallGroups/{firewallGroupId}",
		PathParams: []cav.PathParam{
			{
				Name:        "firewallGroupId",
				Description: "ID of the Firewall Group to update",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=firewallGroup")
				},
			},
		},
		BodyRequestType: itypes.ApiRequestFirewallGroup{},
		ResponseType:    itypes.ApiResponseFirewallGroup{},
	}.Register()

	// DeleteFirewallGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/firewallGroups/firewallGroupId/delete/",
		Name:             "DeleteFirewallGroup",
		Description:      "Delete a Firewall Group (Security Group, IP Set, or Dynamic Security Group)",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/firewallGroups/{firewallGroupId}",
		PathParams: []cav.PathParam{
			{
				Name:        "firewallGroupId",
				Description: "ID of the Firewall Group to delete",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=firewallGroup")
				},
			},
		},
	}.Register()
}
