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

//go:generate endpoint-generator -path firewall.go -output firewall

// Firewall (Distributed Firewall / DFW) endpoints are scoped to a single
// VDC Group and cover three distinct sub-resources: the DFW policies
// (enabled state + default policy), the default policy itself, and the
// bulk set of DFW rules. A DFW policy object is a singleton per VDC Group
// (there is no list/collection semantic). The rules endpoint's "policy"
// path segment is hardcoded to "default" since that is the only supported
// policy at the time these endpoints were introduced. All operations are
// synchronous (no async job).
func init() {
	// GetDfwPolicies
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/get/",
		Name:             "GetDfwPolicies",
		Description:      "Get the Distributed Firewall policies (enabled state and default policy) of a VDC Group",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}/dfwPolicies",
		PathParams: []cav.PathParam{
			{
				Name:        "vdcGroupId",
				Description: "ID of the VDC Group",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=vdcGroup")
				},
			},
		},
		ResponseType: itypes.ApiDfwPolicies{},
	}.Register()

	// UpdateDfwPolicies
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/put/",
		Name:             "UpdateDfwPolicies",
		Description:      "Update the Distributed Firewall policies (enabled state) of a VDC Group",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}/dfwPolicies",
		PathParams: []cav.PathParam{
			{
				Name:        "vdcGroupId",
				Description: "ID of the VDC Group",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=vdcGroup")
				},
			},
		},
		BodyRequestType: itypes.ApiDfwPolicies{},
		ResponseType:    itypes.ApiDfwPolicies{},
	}.Register()

	// UpdateDfwDefaultPolicy
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/default/put/",
		Name:             "UpdateDfwDefaultPolicy",
		Description:      "Update the default Distributed Firewall policy of a VDC Group",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}/dfwPolicies/default",
		PathParams: []cav.PathParam{
			{
				Name:        "vdcGroupId",
				Description: "ID of the VDC Group",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=vdcGroup")
				},
			},
		},
		BodyRequestType: itypes.ApiDfwDefaultPolicy{},
		ResponseType:    itypes.ApiDfwDefaultPolicy{},
	}.Register()

	// GetDfwRules
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/default/rules/get/",
		Name:             "GetDfwRules",
		Description:      "Get the Distributed Firewall rules of a VDC Group",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}/dfwPolicies/default/rules",
		PathParams: []cav.PathParam{
			{
				Name:        "vdcGroupId",
				Description: "ID of the VDC Group",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=vdcGroup")
				},
			},
		},
		ResponseType: itypes.ApiDistributedFirewallRules{},
	}.Register()

	// UpdateDfwRules
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/default/rules/put/",
		Name:             "UpdateDfwRules",
		Description:      "Replace (bulk) the Distributed Firewall rules of a VDC Group",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}/dfwPolicies/default/rules",
		PathParams: []cav.PathParam{
			{
				Name:        "vdcGroupId",
				Description: "ID of the VDC Group",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=vdcGroup")
				},
			},
		},
		BodyRequestType: itypes.ApiDistributedFirewallRules{},
		ResponseType:    itypes.ApiDistributedFirewallRules{},
	}.Register()
}
