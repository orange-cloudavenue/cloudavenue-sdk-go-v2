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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// ListCertificateConsumers lists entity references consuming certificate library item.
func (c *Client) ListCertificateConsumers(ctx context.Context, params types.ParamsListCertificateConsumers) (*types.ModelListEntityReference, error) {
	certificate, err := findCertificate(ctx, c, params.CertificateID, params.CertificateName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve certificate: %w", opListCertificateConsumers, err)
	}

	ep := endpoints.ListCertificateConsumers()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], certificate.ID))
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListCertificateConsumers, err)
	}

	return resp.Result().(*itypes.APIEntityReferences).ToModel(), nil
}

// AddCertificateConsumer adds consumer reference to certificate library item.
func (c *Client) AddCertificateConsumer(ctx context.Context, params types.ParamsAddCertificateConsumer) (*types.ModelEntityReference, error) {
	certificate, err := findCertificate(ctx, c, params.CertificateID, params.CertificateName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve certificate: %w", opAddCertificateConsumer, err)
	}
	if params.ConsumerID == "" && params.ConsumerName == "" {
		return nil, fmt.Errorf("%s: consumer id or name is required", opAddCertificateConsumer)
	}

	body := itypes.APIEntityReference{ID: params.ConsumerID, Name: params.ConsumerName}
	ep := endpoints.AddCertificateConsumer()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], certificate.ID), cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: add: %w", opAddCertificateConsumer, err)
	}

	entity, ok := resp.Result().(*itypes.APIEntityReference)
	if !ok || entity == nil {
		return nil, fmt.Errorf("%s: unexpected add response type %T", opAddCertificateConsumer, resp.Result())
	}

	model := entity.ToModel()
	return &model, nil
}

// SetCertificateConsumers replaces consumer references for certificate library item.
func (c *Client) SetCertificateConsumers(ctx context.Context, params types.ParamsSetCertificateConsumers) (*types.ModelListEntityReference, error) {
	certificate, err := findCertificate(ctx, c, params.CertificateID, params.CertificateName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve certificate: %w", opSetCertificateConsumers, err)
	}

	body := itypes.APIEntityReferences{Values: make([]itypes.APIEntityReference, 0, len(params.Consumers))}
	for _, consumer := range params.Consumers {
		body.Values = append(body.Values, itypes.APIEntityReference{ID: consumer.ID, Name: consumer.Name})
	}

	ep := endpoints.SetCertificateConsumers()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], certificate.ID), cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: set: %w", opSetCertificateConsumers, err)
	}

	refs, ok := resp.Result().(*itypes.APIEntityReferences)
	if !ok || refs == nil {
		return nil, fmt.Errorf("%s: unexpected set response type %T", opSetCertificateConsumers, resp.Result())
	}

	return refs.ToModel(), nil
}
