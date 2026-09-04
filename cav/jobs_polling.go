/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"resty.dev/v3"
)

func getJobStatus(ctx context.Context, c Client, jobID string) (*resty.Response, BackendTarget, error) {
	endpointNames := []string{"GetJobCerberus", "GetJobVmware"}

	var lastErr error
	for _, endpointName := range endpointNames {
		ep, err := GetEndpoint(endpointName)
		if err != nil {
			lastErr = err
			continue
		}

		resp, err := c.Do(ctx, ep, buildJobRequestOptions(ep.Backend, ep, jobID)...)
		if err != nil {
			lastErr = err
			continue
		}

		return resp, ep.Backend, nil
	}

	return nil, 0, fmt.Errorf("no job endpoint found: %w", lastErr)
}

func buildJobRequestOptions(backend BackendTarget, endpoint *Endpoint, jobID string) []EndpointRequestOption {
	return []EndpointRequestOption{
		WithPathParam(endpoint.PathParams[0], jobID),
		SetCustomRestyOption(func(r *resty.Request) {
			switch backend {
			case BackendInfrapi:
				r.SetResult(&CerberusJobAPIResponse{})
				r.SetResultError(&cerberusError{})
			case BackendVMware:
				r.SetResult(&vmwareJobAPIResponse{})
				r.SetResultError(&vmwareError{})
			default:
				r.SetResult(&CerberusJobAPIResponse{})
				r.SetResultError(&cerberusError{})
			}
		}),
	}
}

func parseJobResponse(resp *Response, backend BackendTarget) (*Job, error) {
	if resp == nil || resp.Raw == nil {
		return nil, fmt.Errorf("job response is nil")
	}

	switch backend {
	case BackendInfrapi:
		return (&cerberus{}).JobParser(resp.Raw)
	case BackendVMware:
		return (&vmware{}).JobParser(resp.Raw)
	default:
		return nil, fmt.Errorf("backend %d does not support jobs", backend)
	}
}

func withJitter(interval, jitter time.Duration) time.Duration {
	if jitter <= 0 {
		return interval
	}

	maxJitter := int64(jitter)*2 + 1
	n, err := rand.Int(rand.Reader, big.NewInt(maxJitter))
	if err != nil {
		return interval
	}

	return interval + time.Duration(n.Int64()) - jitter
}
