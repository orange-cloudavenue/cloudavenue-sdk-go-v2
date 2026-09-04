/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

type (
	// requestOption stores request-scoped configuration.
	requestOption struct{}

	// RequestOption configures request construction.
	//
	// It is kept for API compatibility and is currently unused by the runtime.
	RequestOption func(*requestOption) error
)
