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
	"fmt"
	"net/url"
	"path"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

const (
	ovfEnvelopeXMLNS    = "http://schemas.dmtf.org/ovf/envelope/1"
	vcloudXMLNS         = "http://www.vmware.com/vcloud/v1.5"
	xmlSchemaInstanceNS = "http://www.w3.org/2001/XMLSchema-instance"
)

func extractVMwareJob(resp *cav.Response, action string) (cav.Job, error) {
	if job, ok := resp.Result().(*cav.Job); ok && job != nil && job.ID != "" {
		return *job, nil
	}

	if resp == nil || resp.Raw == nil {
		return cav.Job{}, fmt.Errorf("unexpected %s response: missing raw response", action)
	}

	location := resp.Raw.Header().Get("Location")
	if location == "" {
		return cav.Job{}, fmt.Errorf("unexpected %s response: missing job location", action)
	}

	parsed, err := url.Parse(location)
	if err != nil {
		return cav.Job{}, fmt.Errorf("unexpected %s response: parse location: %w", action, err)
	}

	jobID := path.Base(parsed.Path)
	if jobID == "" || jobID == "." || jobID == "/" {
		return cav.Job{}, fmt.Errorf("unexpected %s response: missing job id in location %q", action, location)
	}

	return cav.Job{ID: jobID}, nil
}
