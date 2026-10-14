/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is licensed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

// ModelOrganizationJob describes an organization root job reference.
type ModelOrganizationJob struct {
	JobID string                  `documentation:"Organization root job identifier"`
	Links map[string]ModelJobLink `documentation:"Available job operations"`
}

// ModelJobLink describes an operation link returned by Cerberus.
type ModelJobLink struct {
	Href   string `documentation:"Job operation URL"`
	Method string `documentation:"HTTP method"`
}
