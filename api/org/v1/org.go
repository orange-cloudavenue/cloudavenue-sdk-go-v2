/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package org

import (
	"context"

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

func init() {
	cav.Endpoint{
		Name:         "GetOrganization",
		Method:       cav.MethodGET,
		SubClient:    cav.ClientVmware,
		PathTemplate: "/cloudapi/1.0.0/orgs/{orgUrn}",
		PathParams: []cav.PathParam{
			{
				Name:        "orgUrn",
				Description: "The organization URN (Uniform Resource Name) to identify the organization.",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn_rfc2141")
				},
			},
		},
		QueryParams:      []cav.QueryParam{},
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/orgs/orgUrn/get/",
		BodyResponseType: OrgResponse{},
	}.Register()
}

type OrgResponse struct { //nolint:revive
	CanManageOrgs bool   `json:"canManageOrgs" `
	CanPublish    bool   `json:"canPublish"`
	CatalogCount  int64  `json:"catalogCount"`
	Description   string `json:"description"`
	DiskCount     int64  `json:"diskCount"`
	DisplayName   string `json:"displayName"`
	ID            string `json:"id"`
	IsEnabled     bool   `json:"isEnabled"`
	ManagedBy     struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"managedBy"`
	Name           string `json:"name"`
	OrgVdcCount    int64  `json:"orgVdcCount"`
	RunningVMCount int64  `json:"runningVMCount"`
	UserCount      int64  `json:"userCount"`
	VappCount      int64  `json:"vappCount"`
}

// DemoRequest represents a request to the demo cav.
func (o *Org) DemoRequest(ctx context.Context, orgID string) (*OrgResponse, error) {
	demoEndpoint, err := cav.GetEndpoint("GetOrganization", cav.MethodGET)
	if err != nil {
		return nil, err
	}

	resp, err := demoEndpoint.RequestFunc(
		ctx,
		o.c,
		demoEndpoint,
		cav.WithPathParam(demoEndpoint.PathParams[0], orgID),
	)
	if err != nil {
		return nil, err
	}

	if err := o.c.ParseAPIError("Get organization detail", resp); err != nil {
		return nil, err
	}

	return resp.Result().(*OrgResponse), nil
}
