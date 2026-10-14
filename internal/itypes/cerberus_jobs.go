/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

type APIResponseJobs []struct {
	JobID string                        `json:"jobId"`
	Links map[string]APIResponseJobLink `json:"links"`
}

type APIResponseJobLink struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}
