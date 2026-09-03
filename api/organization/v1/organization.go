/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package organization

import (
	"context"
	"fmt"
	"net/mail"

	"golang.org/x/sync/errgroup"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opGetOrganization    = "Organization.Get"
	opUpdateOrganization = "Organization.Update"
)

// GetOrganization gets detailed information about organization.
func (c *Client) GetOrganization(ctx context.Context) (*types.ModelGetOrganization, error) {
	logger := c.logger.WithGroup("GetOrganization")

	var (
		errG       errgroup.Group
		org        *types.ModelGetOrganization
		orgDetails *types.ModelGetOrganization
	)

	errG.Go(func() error {
		resp, err := c.c.Do(ctx, endpoints.GetOrganization())
		if err != nil {
			return fmt.Errorf("%s: get organization info: %w", opGetOrganization, err)
		}
		org = resp.Result().(*itypes.ApiResponseGetOrg).ToModel()
		return nil
	})

	errG.Go(func() error {
		resp, err := c.c.Do(ctx, endpoints.GetOrganizationDetails())
		if err != nil {
			return fmt.Errorf("%s: get organization details: %w", opGetOrganization, err)
		}
		orgDetails = resp.Result().(*itypes.ApiResponseGetOrgs).ToModel()
		if orgDetails == nil {
			return fmt.Errorf("%s: organization not found", opGetOrganization)
		}
		return nil
	})

	if err := errG.Wait(); err != nil {
		return nil, err
	}

	logger.DebugContext(ctx, "Successfully retrieved organization information")

	return &types.ModelGetOrganization{
		ID:                  orgDetails.ID,
		Name:                org.Name,
		FullName:            org.FullName,
		Description:         org.Description,
		Enabled:             org.Enabled,
		Email:               org.Email,
		InternetBillingMode: org.InternetBillingMode,
		Resources: types.ModelGetOrganizationResources{
			Vdc:       orgDetails.Resources.Vdc,
			Catalog:   orgDetails.Resources.Catalog,
			Vapp:      orgDetails.Resources.Vapp,
			VMRunning: orgDetails.Resources.VMRunning,
			User:      orgDetails.Resources.User,
			Disk:      orgDetails.Resources.Disk,
		},
	}, nil
}

// UpdateOrganization updates existing organization details.
func (c *Client) UpdateOrganization(ctx context.Context, p types.ParamsUpdateOrganization) (*types.ModelGetOrganization, error) {
	if err := validateUpdateOrganizationParams(p); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateOrganization, err)
	}

	if p.FullName == "" && p.Email == "" && p.InternetBillingMode == "" && p.Description == nil {
		return nil, fmt.Errorf("%s: no parameters provided for organization update", opUpdateOrganization)
	}

	logger := c.logger.WithGroup("UpdateOrganization")

	orgDefault, err := c.c.Do(ctx, endpoints.GetOrganization())
	if err != nil {
		return nil, fmt.Errorf("%s: get current organization: %w", opUpdateOrganization, err)
	}
	data := orgDefault.Result().(*itypes.ApiResponseGetOrg).ToModel()

	reqBody := &itypes.ApiRequestUpdateOrg{
		FullName:            data.FullName,
		Description:         data.Description,
		CustomerMail:        data.Email,
		InternetBillingMode: data.InternetBillingMode,
	}
	if p.FullName != "" {
		reqBody.FullName = p.FullName
	}
	if p.Description != nil {
		reqBody.Description = *p.Description
	}
	if p.Email != "" {
		reqBody.CustomerMail = p.Email
	}
	if p.InternetBillingMode != "" {
		reqBody.InternetBillingMode = p.InternetBillingMode
	}

	if _, err = c.c.Do(ctx, endpoints.UpdateOrganization(), cav.SetBody(reqBody)); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateOrganization, err)
	}

	logger.DebugContext(ctx, "Successfully initiated organization update")

	return c.GetOrganization(ctx)
}

func validateUpdateOrganizationParams(p types.ParamsUpdateOrganization) error {
	if len(p.FullName) > 129 {
		return fmt.Errorf("full_name exceeds 129 characters")
	}

	if p.Description != nil && len(*p.Description) > 257 {
		return fmt.Errorf("description exceeds 257 characters")
	}

	if p.Email != "" {
		if _, err := mail.ParseAddress(p.Email); err != nil {
			return fmt.Errorf("email is invalid: %w", err)
		}
	}

	if p.InternetBillingMode != "" && p.InternetBillingMode != "PAYG" && p.InternetBillingMode != "TRAFFIC_VOLUME" {
		return fmt.Errorf("internet_billing_mode must be one of PAYG, TRAFFIC_VOLUME")
	}

	return nil
}
