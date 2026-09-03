/*
* SPDX-FileCopyrightText: Copyright (c) 2025 Orange
* SPDX-License-Identifier: Mozilla Public License 2.0
*
* This software is distributed under the MPL-2.0 license.
* the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
* or see the "LICENSE" file for more details.
 */

package cav

import "resty.dev/v3"

type (
	jobsInterface interface {
		// JobRefresh refreshes job state.
		JobRefresh(httpC *resty.Client, resp *resty.Response, reqOpts []EndpointRequestOption) (job *Job, err error)

		// JobParser parses job data from response.
		JobParser(resp *resty.Response) (job *Job, err error)

		// JobStatusParser maps backend status strings to JobStatus.
		JobStatusParser(status string) (s JobStatus, err error)

		idempotentRetryCondition() resty.RetryConditionFunc
	}
)
