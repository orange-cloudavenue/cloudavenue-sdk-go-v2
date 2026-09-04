/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package consoles provides CloudAvenue console metadata and lookup helpers.
//
// Example:
//
//	console, ok := consoles.FindByOrganizationName("cav01ev01ocb1234567")
//	if ok {
//		endpoint := console.GetAPIVCDEndpoint()
//		_ = endpoint
//	}
//
// All exported functions and methods are safe for concurrent use.

package consoles
