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
	"strings"
	"time"

	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/validators"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// CreatePublicIP allocates a public IP for an edge gateway.
func (c *Client) CreatePublicIP(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGatewayPublicIP, error) {
	if params.ID == "" && params.Name == "" {
		return nil, fmt.Errorf("id or name is required")
	}

	ep := endpoints.CreatePublicIP()

	if params.ID == "" {
		var err error
		params.ID, err = c.retrieveEdgeGatewayIDByName(ctx, params.Name)
		if err != nil {
			return nil, err
		}
	}

	edgeID, err := extractor.ExtractUUID(params.ID)
	if err != nil {
		return nil, err
	}

	body := itypes.APIRequestEdgegatewayPublicIP{
		NetworkType:   "internet",
		EdgeGatewayID: edgeID,
		Properties: itypes.APIRequestEdgegatewayPublicIPProperties{
			Announced: true,
		},
	}

	resp, err := c.c.Do(
		ctx,
		ep,
		cav.SetBody(body),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create public IP: %w", err)
	}

	// Parse jobID from the create response (HTTP 201 with {"jobId":"...","message":"..."})
	jobID := resp.Result().(*cav.CerberusJobCreatedAPIResponse).ID
	if jobID == "" {
		return nil, fmt.Errorf("Failed to create public IP: %w", errors.New("job id not found in create response"))
	}

	// Poll job completion and extract the created public IP from the job response
	publicipCreated, err := cav.AwaitJob(ctx, c.c, jobID, cav.JobPollOptions{
		Timeout:         30 * time.Second,
		PollingInterval: 1 * time.Second,
	}, func(resp *resty.Response) (string, error) {
		// The job status response may be either Cerberus or VMware format depending on which endpoint succeeded
		if r, ok := resp.Result().(*cav.CerberusJobAPIResponse); ok {
			if len(*r) == 0 {
				return "", errors.New("no job information returned")
			}

			job := (*r)[0]
			for _, j := range job.Actions {
				if err := validators.New().Var(j.Details, "ip4_addr"); err == nil {
					return j.Details, nil
				}
			}
		}

		return "", errors.New("public IP not found in job response")
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create public IP: %w", err)
	}

	return c.GetPublicIP(ctx, types.ParamsGetEdgeGatewayPublicIP{
		ID:   params.ID,
		Name: params.Name,
		IP:   publicipCreated,
	})
}

// ListPublicIP lists public IPs assigned to an edge gateway.
func (c *Client) ListPublicIP(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGatewayPublicIPs, error) {
	if params.ID == "" && params.Name == "" {
		return nil, fmt.Errorf("id or name is required")
	}

	services, err := c.GetServices(ctx, params)
	if err != nil {
		return nil, err
	}

	ips := &types.ModelEdgeGatewayPublicIPs{
		EdgegatewayID:   services.ID,
		EdgegatewayName: services.Name,
		PublicIPs:       make([]types.ModelEdgeGatewayPublicIP, 0, len(services.PublicIP)),
	}

	for _, publicip := range services.PublicIP {
		ips.PublicIPs = append(ips.PublicIPs, types.ModelEdgeGatewayPublicIP{
			EdgegatewayID:                    services.ID,
			EdgegatewayName:                  services.Name,
			ModelEdgeGatewayServicesPublicIP: *publicip,
		})
	}

	return ips, nil
}

// GetPublicIP returns a public IP assigned to an edge gateway.
func (c *Client) GetPublicIP(ctx context.Context, params types.ParamsGetEdgeGatewayPublicIP) (*types.ModelEdgeGatewayPublicIP, error) {
	if params.IP == "" {
		return nil, fmt.Errorf("ip is required")
	}
	if err := validators.New().Var(params.IP, "ip4_addr"); err != nil {
		return nil, fmt.Errorf("invalid IP address: %w", err)
	}

	if params.ID == "" {
		var err error
		params.ID, err = c.retrieveEdgeGatewayIDByName(ctx, params.Name)
		if err != nil {
			return nil, err
		}
	}

	ep := endpoints.GetEdgeGatewayServices()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], params.ID),
		cav.WithQueryParam(ep.QueryParams[1], params.Name),
		cav.WithQueryParam(ep.QueryParams[2], params.IP),
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving network services for edge gateway %s: %w", params.ID, err)
	}

	data := resp.Result().(*itypes.APIResponseNetworkServices).ToModel(types.ParamsEdgeGateway{
		ID:   params.ID,
		Name: params.Name,
	})
	if data == nil || len(data.PublicIP) == 0 {
		return nil, fmt.Errorf("no public IPs found for edge gateway %s", params.ID)
	}

	for _, publicip := range data.PublicIP {
		if publicip.IP == params.IP {
			return &types.ModelEdgeGatewayPublicIP{
				EdgegatewayID:                    data.ID,
				EdgegatewayName:                  data.Name,
				ModelEdgeGatewayServicesPublicIP: *publicip,
			}, nil
		}
	}

	return nil, fmt.Errorf("public IP %s not found in edge gateway %s", params.IP, params.ID)
}

// DeletePublicIP releases a public IP from an edge gateway.
func (c *Client) DeletePublicIP(ctx context.Context, params types.ParamsDeleteEdgeGatewayPublicIP) error {
	if params.IP == "" {
		return fmt.Errorf("ip is required")
	}
	if err := validators.New().Var(params.IP, "ip4_addr"); err != nil {
		return fmt.Errorf("invalid IP address: %w", err)
	}

	ep := endpoints.DisableCloudavenueServices()
	ipID := fmt.Sprintf("ip-%s", strings.ReplaceAll(params.IP, ".", "-"))

	_, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], ipID),
	)

	return err
}
