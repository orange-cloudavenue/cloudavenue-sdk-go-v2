/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

type (
	// ModelObjectReference represents a reference to an object in the API.
	// It contains the ID and name of the object.
	ModelObjectReference struct {
		ID   string `json:"id" fake:"{urn:vdc}"`
		Name string `json:"name"`
	}
)
