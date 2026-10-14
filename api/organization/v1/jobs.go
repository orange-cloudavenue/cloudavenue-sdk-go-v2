/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is licensed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package organization

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// GetJobs returns current organization root jobs through Cerberus.
func (c *Client) GetJobs(ctx context.Context) ([]types.ModelOrganizationJob, error) {
	resp, err := c.c.Do(ctx, endpoints.GetJobs())
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}
	result, ok := resp.Result().(*itypes.APIResponseJobs)
	if !ok || result == nil {
		return nil, fmt.Errorf("unexpected jobs response type %T", resp.Result())
	}
	jobs := make([]types.ModelOrganizationJob, len(*result))
	for i, job := range *result {
		links := make(map[string]types.ModelJobLink, len(job.Links))
		for name, link := range job.Links {
			links[name] = types.ModelJobLink{Href: link.Href, Method: link.Method}
		}
		jobs[i] = types.ModelOrganizationJob{JobID: job.JobID, Links: links}
	}
	return jobs, nil
}
