/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"context"
	"fmt"
	"slices"
	"time"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/common-go/validators"
	"golang.org/x/sync/errgroup"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opCreateEdgeGatewaySubmit          = "EdgeGateway.Create.Submit"
	opUpdateEdgeGatewayBandwidthCreate = "EdgeGateway.Create.UpdateBandwidth"
	opListVDCByOwnerName               = "EdgeGateway.Create.ListVDCByOwnerName"
	opListVDCGroupByOwnerName          = "EdgeGateway.Create.ListVDCGroupByOwnerName"
)

type createEdgeGatewaySubmitParams struct {
	OwnerType string
	OwnerName string
	T0Name    string
}

type updateCreatedEdgeGatewayBandwidthParams struct {
	EdgeGatewayID string
	Bandwidth     int
}

type listVDCByOwnerNameParams struct {
	OwnerName string
}

type listVDCGroupByOwnerNameParams struct {
	OwnerName string
}

var (
	createEdgeGatewaySubmitOp = cav.Operation[createEdgeGatewaySubmitParams, string]{
		Name:     opCreateEdgeGatewaySubmit,
		Backend:  cav.BackendInfrapi,
		Endpoint: endpoints.CreateEdgeGateway(),
		Validate: func(p createEdgeGatewaySubmitParams) error {
			if p.OwnerType == "" {
				return fmt.Errorf("owner type is required")
			}
			if p.OwnerName == "" {
				return fmt.Errorf("owner name is required")
			}
			if p.T0Name == "" {
				return fmt.Errorf("t0 name is required")
			}

			return nil
		},
		RequestOptions: func(p createEdgeGatewaySubmitParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.CreateEdgeGateway()
			return []cav.EndpointRequestOption{
				cav.WithPathParam(ep.PathParams[0], p.OwnerType),
				cav.WithPathParam(ep.PathParams[1], p.OwnerName),
			}, nil
		},
		Transform: func(p createEdgeGatewaySubmitParams) (any, error) {
			return itypes.ApiRequestEdgeGateway{T0Name: p.T0Name}, nil
		},
		Extract: func(resp *cav.Response, _ createEdgeGatewaySubmitParams) (string, error) {
			created, ok := resp.Result().(*cav.CerberusJobCreatedAPIResponse)
			if !ok || created == nil {
				return "", fmt.Errorf("unexpected create response type %T", resp.Result())
			}

			if created.ID == "" {
				return "", errors.New("job id not found in create response")
			}

			return created.ID, nil
		},
	}
	updateCreatedEdgeGatewayBandwidthOp = cav.Operation[updateCreatedEdgeGatewayBandwidthParams, struct{}]{
		Name:     opUpdateEdgeGatewayBandwidthCreate,
		Backend:  cav.BackendInfrapi,
		Endpoint: endpoints.UpdateEdgeGatewayBandwidth(),
		Validate: func(p updateCreatedEdgeGatewayBandwidthParams) error {
			if p.EdgeGatewayID == "" {
				return fmt.Errorf("edge gateway id is required")
			}
			if p.Bandwidth <= 0 {
				return fmt.Errorf("bandwidth must be greater than 0")
			}

			return nil
		},
		RequestOptions: func(p updateCreatedEdgeGatewayBandwidthParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.UpdateEdgeGatewayBandwidth()
			return []cav.EndpointRequestOption{
				cav.WithPathParam(ep.PathParams[0], p.EdgeGatewayID),
			}, nil
		},
		Transform: func(p updateCreatedEdgeGatewayBandwidthParams) (any, error) {
			return itypes.ApiRequestBandwidth{Bandwidth: p.Bandwidth}, nil
		},
		Extract: func(_ *cav.Response, _ updateCreatedEdgeGatewayBandwidthParams) (struct{}, error) {
			return struct{}{}, nil
		},
	}
	listVDCByOwnerNameOp = cav.Operation[listVDCByOwnerNameParams, *itypes.ApiResponseListVDC]{
		Name:     opListVDCByOwnerName,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.ListVdc(),
		Validate: func(p listVDCByOwnerNameParams) error {
			if p.OwnerName == "" {
				return fmt.Errorf("owner name is required")
			}

			return nil
		},
		RequestOptions: func(p listVDCByOwnerNameParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.ListVdc()
			return []cav.EndpointRequestOption{
				cav.WithQueryParam(ep.QueryParams[0], "name=="+p.OwnerName),
			}, nil
		},
		Extract: func(resp *cav.Response, _ listVDCByOwnerNameParams) (*itypes.ApiResponseListVDC, error) {
			list, ok := resp.Result().(*itypes.ApiResponseListVDC)
			if !ok || list == nil {
				return nil, fmt.Errorf("unexpected list VDC response type %T", resp.Result())
			}

			return list, nil
		},
	}
	listVDCGroupByOwnerNameOp = cav.Operation[listVDCGroupByOwnerNameParams, *itypes.ApiResponseListVdcGroup]{
		Name:     opListVDCGroupByOwnerName,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.ListVdcGroup(),
		Validate: func(p listVDCGroupByOwnerNameParams) error {
			if p.OwnerName == "" {
				return fmt.Errorf("owner name is required")
			}

			return nil
		},
		RequestOptions: func(p listVDCGroupByOwnerNameParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.ListVdcGroup()
			return []cav.EndpointRequestOption{
				cav.WithQueryParam(ep.QueryParams[0], "name=="+p.OwnerName),
			}, nil
		},
		Extract: func(resp *cav.Response, _ listVDCGroupByOwnerNameParams) (*itypes.ApiResponseListVdcGroup, error) {
			list, ok := resp.Result().(*itypes.ApiResponseListVdcGroup)
			if !ok || list == nil {
				return nil, fmt.Errorf("unexpected list VDC group response type %T", resp.Result())
			}

			return list, nil
		},
	}
)

// CreateEdgeGateway creates an edge gateway for a VDC or VDC group owner.
func (c *Client) CreateEdgeGateway(ctx context.Context, params types.ParamsCreateEdgeGateway) (*types.ModelEdgeGateway, error) {
	if params.OwnerName == "" {
		return nil, fmt.Errorf("%s: validate: owner name is required", opCreateEdgeGateway)
	}

	dependencies, err := c.resolveCreateEdgeGatewayDependencies(ctx, params.OwnerName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve dependencies: %w", opCreateEdgeGateway, err)
	}

	t0, bandwidth, err := resolveEdgeGatewayCreateT0(dependencies.t0s, params)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve t0: %w", opCreateEdgeGateway, err)
	}

	if len(t0.EdgeGateways) >= t0.MaxEdgeGateways {
		return nil, fmt.Errorf("%s: %w", opCreateEdgeGateway, errors.New("Maximum number of edge gateways reached for T0: "+t0.Name))
	}

	if !t0.Bandwidth.AllowUnlimited && !slices.Contains(t0.Bandwidth.AllowedBandwidthValues, bandwidth) {
		return nil, fmt.Errorf("%s: %w", opCreateEdgeGateway, errors.New("Invalid bandwidth value for SHARED T0"))
	}

	ownerName, ownerType, err := resolveEdgeGatewayCreateOwner(dependencies.vdcs, dependencies.vdcGroups, params.OwnerName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve owner: %w", opCreateEdgeGateway, err)
	}

	jobID, err := cav.Execute(ctx, c.c, createEdgeGatewaySubmitOp, createEdgeGatewaySubmitParams{
		OwnerType: ownerType,
		OwnerName: ownerName,
		T0Name:    t0.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: create: %w", opCreateEdgeGateway, err)
	}

	edgeGatewayCreated, err := cav.AwaitJob(ctx, c.c, jobID, cav.JobPollOptions{
		Timeout:         30 * time.Second,
		PollingInterval: 1 * time.Second,
	}, extractCreatedEdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: await job: %w", opCreateEdgeGateway, err)
	}

	edgeCreated, err := c.GetEdgeGateway(ctx, types.ParamsEdgeGateway{Name: edgeGatewayCreated})
	if err != nil {
		return nil, fmt.Errorf("%s: get created: %w", opCreateEdgeGateway, err)
	}

	if bandwidth > 5 {
		if _, err := cav.Execute(ctx, c.c, updateCreatedEdgeGatewayBandwidthOp, updateCreatedEdgeGatewayBandwidthParams{
			EdgeGatewayID: edgeCreated.ID,
			Bandwidth:     bandwidth,
		}); err != nil {
			return edgeCreated, &cav.PartialSuccessError[*types.ModelEdgeGateway]{
				Result: edgeCreated,
				Err:    fmt.Errorf("%s: update bandwidth: %w", opCreateEdgeGateway, err),
			}
		}
	}

	return edgeCreated, nil
}

type createEdgeGatewayDependencies struct {
	vdcs      *itypes.ApiResponseListVDC
	vdcGroups *itypes.ApiResponseListVdcGroup
	t0s       *types.ModelT0s
}

func (c *Client) resolveCreateEdgeGatewayDependencies(ctx context.Context, ownerName string) (*createEdgeGatewayDependencies, error) {
	dependencies := &createEdgeGatewayDependencies{}
	errGroup, errCtx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		resp, err := cav.Execute(errCtx, c.c, listVDCGroupByOwnerNameOp, listVDCGroupByOwnerNameParams{OwnerName: ownerName})
		if err != nil {
			return fmt.Errorf("failed to list VDC Groups: %w", err)
		}

		dependencies.vdcGroups = resp
		return nil
	})

	errGroup.Go(func() error {
		resp, err := cav.Execute(errCtx, c.c, listVDCByOwnerNameOp, listVDCByOwnerNameParams{OwnerName: ownerName})
		if err != nil {
			return fmt.Errorf("failed to list VDCs: %w", err)
		}

		dependencies.vdcs = resp
		return nil
	})

	errGroup.Go(func() error {
		var err error
		dependencies.t0s, err = c.ListT0(errCtx)
		if err != nil {
			return fmt.Errorf("failed to list T0 routers: %w", err)
		}

		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	return dependencies, nil
}

func extractCreatedEdgeGatewayName(job *resty.Response) (string, error) {
	if job == nil {
		return "", errors.New("created edge gateway name not found in job response")
	}

	r, ok := job.Result().(*cav.CerberusJobAPIResponse)
	if !ok {
		return "", errors.New("created edge gateway name not found in job response")
	}
	if len(*r) == 0 {
		return "", errors.New("no job information returned")
	}

	var fallback string
	for _, action := range (*r)[0].Actions {
		if fallback == "" && action.Details != "" {
			fallback = action.Details
		}
		if err := validators.New().Var(action.Details, "resource_name=edgegateway"); err == nil {
			return action.Details, nil
		}
	}

	if fallback != "" {
		return fallback, nil
	}

	return "", errors.New("created edge gateway name not found in job response")
}

func resolveEdgeGatewayCreateOwner(vdcs *itypes.ApiResponseListVDC, vdcGroups *itypes.ApiResponseListVdcGroup, ownerName string) (string, string, error) {
	matchingVDCs := make([]itypes.ApiResponseListVDCRecord, 0)
	if vdcs != nil {
		for _, record := range vdcs.Records {
			if record.Name == ownerName {
				matchingVDCs = append(matchingVDCs, record)
			}
		}
	}

	matchingVDCGroups := make([]itypes.ApiResponseListVdcGroupDetails, 0)
	if vdcGroups != nil {
		for _, group := range vdcGroups.Values {
			if group.Name == ownerName {
				matchingVDCGroups = append(matchingVDCGroups, group)
			}
		}
	}

	switch {
	case len(matchingVDCs) == 0 && len(matchingVDCGroups) == 0:
		return "", "", errors.New("no VDCs or VDC groups found for owner: " + ownerName)
	case len(matchingVDCs) >= 1 && len(matchingVDCGroups) >= 1:
		return "", "", errors.New("both VDCs and VDC groups found for owner: " + ownerName)
	case len(matchingVDCs) == 1:
		return matchingVDCs[0].Name, "vdc", nil
	case len(matchingVDCGroups) == 1:
		return matchingVDCGroups[0].Name, "vdcgroup", nil
	default:
		return "", "", errors.New("ambiguous owner: " + ownerName)
	}
}

func resolveEdgeGatewayCreateT0(t0s *types.ModelT0s, params types.ParamsCreateEdgeGateway) (types.ModelT0, int, error) {
	if t0s.Count == 0 {
		return types.ModelT0{}, 0, errors.New("no T0 routers available to connect edge gateway")
	}

	if params.T0Name == "" {
		if t0s.Count > 1 {
			return types.ModelT0{}, 0, errors.New("multiple T0 routers found, please specify T0Name")
		}

		selected := t0s.T0s[0]
		bandwidth := params.Bandwidth
		if !selected.Bandwidth.AllowUnlimited && bandwidth <= 0 {
			bandwidth = 5
		}
		return selected, bandwidth, nil
	}

	for _, t0Model := range t0s.T0s {
		if t0Model.Name != params.T0Name {
			continue
		}

		bandwidth := params.Bandwidth
		if !t0Model.Bandwidth.AllowUnlimited && bandwidth <= 0 {
			bandwidth = 5
		}
		return t0Model, bandwidth, nil
	}

	return types.ModelT0{}, 0, errors.New("T0 router not found: " + params.T0Name)
}
