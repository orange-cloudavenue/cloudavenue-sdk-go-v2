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
	opListCertificate          = "Certificate.List"
	opGetCertificate           = "Certificate.Get"
	opCreateCertificate        = "Certificate.Create"
	opUpdateCertificate        = "Certificate.Update"
	opDeleteCertificate        = "Certificate.Delete"
	opListCertificateConsumers = "Certificate.Consumers.List"
	opAddCertificateConsumer   = "Certificate.Consumers.Add"
	opSetCertificateConsumers  = "Certificate.Consumers.Set"
)

// ListCertificate lists certificate library items for organization.
func (c *Client) ListCertificate(ctx context.Context, params types.ParamsListCertificate) (*types.ModelListCertificate, error) {
	ep := endpoints.ListCertificate()
	opts := make([]cav.EndpointRequestOption, 0, 1)
	if params.Name != "" {
		opts = append(opts, cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("alias==%s", params.Name)))
	}

	resp, err := c.c.Do(ctx, ep, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListCertificate, err)
	}

	return resp.Result().(*itypes.ApiResponseListCertificate).ToModel(), nil
}

// GetCertificate gets certificate library item by URN or alias.
func (c *Client) GetCertificate(ctx context.Context, params types.ParamsGetCertificate) (*types.ModelGetCertificate, error) {
	certificate, err := findCertificate(ctx, c, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetCertificate, err)
	}

	model := certificate.ToModel()
	return &model, nil
}

// CreateCertificate creates certificate library item for organization.
func (c *Client) CreateCertificate(ctx context.Context, params types.ParamsCreateCertificate) (*types.ModelGetCertificate, error) {
	if err := validateCreateCertificate(ctx, c, params); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateCertificate, err)
	}

	return cav.Execute(ctx, c.c, cav.Operation[types.ParamsCreateCertificate, *types.ModelGetCertificate]{
		Name:     opCreateCertificate,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.CreateCertificate(),
		Validate: validateCreateCertificateParams,
		Transform: func(p types.ParamsCreateCertificate) (any, error) {
			return itypes.ApiRequestCertificate{
				Alias:                p.Name,
				Description:          p.Description,
				Certificate:          p.Certificate,
				PrivateKey:           p.PrivateKey,
				PrivateKeyPassphrase: p.PrivateKeyPassphrase,
			}, nil
		},
		Extract: func(resp *cav.Response, _ types.ParamsCreateCertificate) (*types.ModelGetCertificate, error) {
			certificate, ok := resp.Result().(*itypes.ApiResponseCertificate)
			if !ok || certificate == nil {
				return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateCertificate, resp.Result())
			}
			model := certificate.ToModel()
			return &model, nil
		},
	}, params)
}

// UpdateCertificate updates certificate alias or description.
// Certificate content and private key remain immutable after creation.
func (c *Client) UpdateCertificate(ctx context.Context, params types.ParamsUpdateCertificate) (*types.ModelGetCertificate, error) {
	if params.Name == "" && params.Description == "" {
		return nil, fmt.Errorf("%s: no parameters provided for certificate update", opUpdateCertificate)
	}

	current, err := findCertificate(ctx, c, params.ID, params.LookupName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve: %w", opUpdateCertificate, err)
	}

	body := itypes.ApiRequestCertificate{
		ID:          current.ID,
		Alias:       current.Alias,
		Certificate: current.Certificate,
		Description: current.Description,
	}
	if params.Name != "" {
		body.Alias = params.Name
	}
	if params.Description != "" {
		body.Description = params.Description
	}

	ep := endpoints.UpdateCertificate()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID), cav.SetBody(body)); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateCertificate, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteCertificate deletes certificate library item by URN or alias.
func (c *Client) DeleteCertificate(ctx context.Context, params types.ParamsDeleteCertificate) error {
	current, err := findCertificate(ctx, c, params.ID, params.Name)
	if err != nil {
		return fmt.Errorf("%s: resolve: %w", opDeleteCertificate, err)
	}

	ep := endpoints.DeleteCertificate()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID)); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteCertificate, err)
	}

	return nil
}

func findCertificate(ctx context.Context, c *Client, id, name string) (*itypes.ApiResponseCertificate, error) {
	if id != "" {
		ep := endpoints.GetCertificate()
		resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], id))
		if err != nil {
			return nil, err
		}
		return resp.Result().(*itypes.ApiResponseCertificate), nil
	}

	if name == "" {
		return nil, &pkgerrors.APIError{Operation: "GetCertificate", StatusCode: 400, Message: "id or name is required"}
	}

	ep := endpoints.ListCertificate()
	resp, err := c.c.Do(ctx, ep, cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("alias==%s", name)))
	if err != nil {
		return nil, err
	}

	list := resp.Result().(*itypes.ApiResponseListCertificate)
	if len(list.Values) == 0 {
		return nil, &pkgerrors.APIError{Operation: "GetCertificate", StatusCode: 404, Message: fmt.Sprintf("certificate %q not found", name)}
	}
	if len(list.Values) > 1 {
		return nil, pkgerrors.Newf("multiple certificates found for %q", name)
	}

	return &list.Values[0], nil
}

func validateCreateCertificate(ctx context.Context, c *Client, p types.ParamsCreateCertificate) error {
	if err := validateCreateCertificateParams(p); err != nil {
		return err
	}

	list, err := c.ListCertificate(ctx, types.ParamsListCertificate{Name: p.Name})
	if err != nil {
		return err
	}

	if list != nil {
		for _, certificate := range list.Certificates {
			if strings.EqualFold(certificate.Name, p.Name) {
				return fmt.Errorf("certificate %q already exists", p.Name)
			}
		}
	}

	return nil
}

func validateCreateCertificateParams(p types.ParamsCreateCertificate) error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Certificate == "" {
		return fmt.Errorf("certificate is required")
	}
	if p.PrivateKeyPassphrase != "" && p.PrivateKey == "" {
		return fmt.Errorf("private_key is required when private_key_passphrase is provided")
	}
	return nil
}
