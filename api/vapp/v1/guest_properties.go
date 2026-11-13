/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	"context"
	"fmt"
	"sort"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opGetVAppGuestProperties    = "VApp.GuestProperties.Get"
	opUpdateVAppGuestProperties = "VApp.GuestProperties.Update"
)

var (
	getVAppProductSectionsOp = cav.Operation[getVAppByIDParams, *itypes.APIResponseGuestProperties]{
		Name:     opGetVAppGuestProperties,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.GetVAppProductSections(),
		Validate: func(p getVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			return nil
		},
		RequestOptions: func(p getVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.GetVAppProductSections()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Extract: func(resp *cav.Response, _ getVAppByIDParams) (*itypes.APIResponseGuestProperties, error) {
			section, ok := resp.Result().(*itypes.APIResponseGuestProperties)
			if !ok || section == nil {
				return nil, fmt.Errorf("unexpected get response type %T", resp.Result())
			}
			return section, nil
		},
	}
	updateVAppGuestPropertiesOp = cav.Operation[updateVAppGuestPropertiesParams, cav.Job]{
		Name:     opUpdateVAppGuestProperties,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.UpdateVAppProductSections(),
		Validate: func(p updateVAppGuestPropertiesParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			return nil
		},
		RequestOptions: func(p updateVAppGuestPropertiesParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.UpdateVAppProductSections()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Transform: func(p updateVAppGuestPropertiesParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ updateVAppGuestPropertiesParams) (cav.Job, error) {
			return extractVMwareJob(resp, "replace guest properties")
		},
	}
)

// GetVAppGuestProperties returns the ProductSection guest properties for a VApp.
func (c *Client) GetVAppGuestProperties(ctx context.Context, params types.ParamsGetVAppGuestProperties) (*types.ModelVAppGuestProperties, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetVAppGuestProperties, err)
	}

	vappID, err := c.resolveVAppID(ctx, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVAppGuestProperties, err)
	}

	resp, err := cav.Execute(ctx, c.c, getVAppProductSectionsOp, getVAppByIDParams{ID: vappID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVAppGuestProperties, err)
	}

	model := resp.ToModel()
	return &model, nil
}

// UpdateVAppGuestProperties replaces the ProductSection guest properties for a VApp.
func (c *Client) UpdateVAppGuestProperties(ctx context.Context, params types.ParamsUpdateVAppGuestProperties) (*types.ModelVAppGuestProperties, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateVAppGuestProperties, err)
	}

	vappID, err := c.resolveVAppID(ctx, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateVAppGuestProperties, err)
	}

	body := toAPIResponseGuestProperties(params.GuestProperties)
	job, err := cav.Execute(ctx, c.c, updateVAppGuestPropertiesOp, updateVAppGuestPropertiesParams{ID: vappID, Body: body})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateVAppGuestProperties, err)
	}
	if err := awaitVAppJob(ctx, c.c, opUpdateVAppGuestProperties, job.ID); err != nil {
		return nil, err
	}
	return c.GetVAppGuestProperties(ctx, types.ParamsGetVAppGuestProperties{ID: vappID})
}

func toAPIResponseGuestProperties(model types.ModelVAppGuestProperties) itypes.APIResponseGuestProperties {
	result := itypes.APIResponseGuestProperties{
		XMLNS:                        vcloudXMLNS,
		XMLNSXsi:                     xmlSchemaInstanceNS,
		XsiNoNamespaceSchemaLocation: "http://www.vmware.com/vcloud/v1.5/vmguest.xsd",
		ProductSectionList:           make([]itypes.APIResponseProductSection, 0, len(model.Sections)),
	}
	sections := append([]types.ModelVAppGuestPropertySection(nil), model.Sections...)
	sort.SliceStable(sections, func(i, j int) bool {
		if sections[i].Class != sections[j].Class {
			return sections[i].Class < sections[j].Class
		}
		return sections[i].Instance < sections[j].Instance
	})

	for _, section := range sections {
		keys := make([]string, 0, len(section.Properties))
		for key := range section.Properties {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		properties := make([]itypes.APIResponseProductProperty, 0, len(section.Properties))
		for _, key := range keys {
			properties = append(properties, itypes.APIResponseProductProperty{Key: key, Value: section.Properties[key]})
		}
		result.ProductSectionList = append(result.ProductSectionList, itypes.APIResponseProductSection{
			Class:    section.Class,
			Instance: section.Instance,
			Property: properties,
		})
	}
	return result
}
