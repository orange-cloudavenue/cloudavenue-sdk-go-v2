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
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListTrustedCertificate   = "TrustedCertificate.List"
	opGetTrustedCertificate    = "TrustedCertificate.Get"
	opCreateTrustedCertificate = "TrustedCertificate.Create"
	opUpdateTrustedCertificate = "TrustedCertificate.Update"
	opDeleteTrustedCertificate = "TrustedCertificate.Delete"
)

// ListTrustedCertificate lists trusted certificates for organization.
func (c *Client) ListTrustedCertificate(ctx context.Context, params types.ParamsListTrustedCertificate) (*types.ModelListTrustedCertificate, error) {
	ep := endpoints.ListTrustedCertificate()
	opts := make([]cav.EndpointRequestOption, 0, 1)
	if params.Name != "" {
		opts = append(opts, cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("alias==%s", params.Name)))
	}

	resp, err := c.c.Do(ctx, ep, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListTrustedCertificate, err)
	}

	return resp.Result().(*itypes.ApiResponseListTrustedCertificate).ToModel(), nil
}

// GetTrustedCertificate gets trusted certificate by URN or alias.
func (c *Client) GetTrustedCertificate(ctx context.Context, params types.ParamsGetTrustedCertificate) (*types.ModelGetTrustedCertificate, error) {
	certificate, err := findTrustedCertificate(ctx, c, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetTrustedCertificate, err)
	}

	model := certificate.ToModel()
	return &model, nil
}

// CreateTrustedCertificate creates trusted certificate for organization.
func (c *Client) CreateTrustedCertificate(ctx context.Context, params types.ParamsCreateTrustedCertificate) (*types.ModelGetTrustedCertificate, error) {
	if err := validateCreateTrustedCertificate(ctx, c, params); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateTrustedCertificate, err)
	}

	return cav.Execute(ctx, c.c, cav.Operation[types.ParamsCreateTrustedCertificate, *types.ModelGetTrustedCertificate]{
		Name:     opCreateTrustedCertificate,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.CreateTrustedCertificate(),
		Validate: validateCreateTrustedCertificateParams,
		Transform: func(p types.ParamsCreateTrustedCertificate) (any, error) {
			return itypes.ApiRequestTrustedCertificate{Alias: p.Name, Certificate: p.Certificate}, nil
		},
		Extract: func(resp *cav.Response, _ types.ParamsCreateTrustedCertificate) (*types.ModelGetTrustedCertificate, error) {
			certificate, ok := resp.Result().(*itypes.ApiResponseTrustedCertificate)
			if !ok || certificate == nil {
				return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateTrustedCertificate, resp.Result())
			}
			model := certificate.ToModel()
			return &model, nil
		},
	}, params)
}

// UpdateTrustedCertificate updates trusted certificate alias or certificate body.
func (c *Client) UpdateTrustedCertificate(ctx context.Context, params types.ParamsUpdateTrustedCertificate) (*types.ModelGetTrustedCertificate, error) {
	if params.Name == "" && params.Certificate == "" {
		return nil, fmt.Errorf("%s: no parameters provided for trusted certificate update", opUpdateTrustedCertificate)
	}

	current, err := findTrustedCertificate(ctx, c, params.ID, params.LookupName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve: %w", opUpdateTrustedCertificate, err)
	}

	body := itypes.ApiRequestTrustedCertificate{ID: current.ID, Alias: current.Alias, Certificate: current.Certificate}
	if params.Name != "" {
		body.Alias = params.Name
	}
	if params.Certificate != "" {
		body.Certificate = params.Certificate
	}

	ep := endpoints.UpdateTrustedCertificate()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID), cav.SetBody(body)); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateTrustedCertificate, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteTrustedCertificate deletes trusted certificate by URN or alias.
func (c *Client) DeleteTrustedCertificate(ctx context.Context, params types.ParamsDeleteTrustedCertificate) error {
	current, err := findTrustedCertificate(ctx, c, params.ID, params.Name)
	if err != nil {
		return fmt.Errorf("%s: resolve: %w", opDeleteTrustedCertificate, err)
	}

	ep := endpoints.DeleteTrustedCertificate()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID)); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteTrustedCertificate, err)
	}

	return nil
}

func findTrustedCertificate(ctx context.Context, c *Client, id, name string) (*itypes.ApiResponseTrustedCertificate, error) {
	if id != "" {
		ep := endpoints.GetTrustedCertificate()
		resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], id))
		if err != nil {
			return nil, err
		}
		return resp.Result().(*itypes.ApiResponseTrustedCertificate), nil
	}

	if name == "" {
		return nil, &pkgerrors.APIError{Operation: "GetTrustedCertificate", StatusCode: 400, Message: "id or name is required"}
	}

	ep := endpoints.ListTrustedCertificate()
	resp, err := c.c.Do(ctx, ep, cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("alias==%s", name)))
	if err != nil {
		return nil, err
	}

	list := resp.Result().(*itypes.ApiResponseListTrustedCertificate)
	if len(list.Values) == 0 {
		return nil, &pkgerrors.APIError{Operation: "GetTrustedCertificate", StatusCode: 404, Message: fmt.Sprintf("trusted certificate %q not found", name)}
	}
	if len(list.Values) > 1 {
		return nil, pkgerrors.Newf("multiple trusted certificates found for %q", name)
	}

	return &list.Values[0], nil
}

func validateCreateTrustedCertificate(ctx context.Context, c *Client, p types.ParamsCreateTrustedCertificate) error {
	if err := validateCreateTrustedCertificateParams(p); err != nil {
		return err
	}

	list, err := c.ListTrustedCertificate(ctx, types.ParamsListTrustedCertificate{Name: p.Name})
	if err != nil {
		return err
	}

	if list != nil {
		for _, certificate := range list.Certificates {
			if strings.EqualFold(certificate.Name, p.Name) {
				return fmt.Errorf("trusted certificate %q already exists", p.Name)
			}
		}
	}

	return nil
}

func validateCreateTrustedCertificateParams(p types.ParamsCreateTrustedCertificate) error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Certificate == "" {
		return fmt.Errorf("certificate is required")
	}
	return nil
}
