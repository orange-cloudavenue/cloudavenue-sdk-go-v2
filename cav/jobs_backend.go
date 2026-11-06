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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

const (
	getJobCerberusEndpointName = "GetJobCerberus"
	getJobVMwareEndpointName   = "GetJobVmware"
)

// fetchJob reads the current state of a job from its backend and parses it.
func fetchJob(ctx context.Context, c Client, backend BackendTarget, jobID string) (*Job, error) {
	endpointName, err := getJobEndpointName(backend)
	if err != nil {
		return nil, err
	}

	ep, err := GetEndpoint(endpointName)
	if err != nil {
		return nil, err
	}

	resp, err := c.DoWithBackend(ctx, backend, ep, buildJobRequestOptions(backend, ep, jobID)...)
	if err != nil {
		return nil, err
	}

	job, err := parseJobResponse(resp, backend)
	if err != nil {
		return job, fmt.Errorf("parse job status: %w", err)
	}

	return job, nil
}

func getJobEndpointName(backend BackendTarget) (string, error) {
	switch backend {
	case BackendInfrapi:
		return getJobCerberusEndpointName, nil
	case BackendVMware:
		return getJobVMwareEndpointName, nil
	default:
		return "", fmt.Errorf("backend %d does not support jobs", backend)
	}
}

func buildJobRequestOptions(backend BackendTarget, endpoint *Endpoint, jobID string) []EndpointRequestOption {
	return []EndpointRequestOption{
		WithPathParam(endpoint.PathParams[0], jobID),
		SetCustomRestyOption(func(r *resty.Request) {
			switch backend {
			case BackendVMware:
				r.SetResult(&vmwareJobAPIResponse{})
				r.SetResultError(&vmwareError{})
			default:
				r.SetResult(&CerberusJobAPIResponse{})
				r.SetResultError(&errors.CustomerAPIErrorBody{})
			}
		}),
	}
}

func parseJobResponse(resp *resty.Response, backend BackendTarget) (*Job, error) {
	if resp == nil {
		return nil, fmt.Errorf("job response is nil")
	}

	switch backend {
	case BackendInfrapi:
		return parseCerberusJob(resp)
	case BackendVMware:
		return parseVMwareJob(resp)
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
