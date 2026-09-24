/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

type APIResponseAccessControlGrants struct {
	*Page
	Values []AccessControlGrant `json:"values"`
}

type Page struct {
	ResultTotal  int64         `json:"resultTotal"`
	PageCount    int64         `json:"pageCount"`
	Page         int64         `json:"page"`
	PageSize     int64         `json:"pageSize"`
	Associations []Association `json:"associations"`
}

type Association struct {
	EntityID      string `json:"entityId"`
	AssociationID string `json:"associationId"`
}

type AccessControlGrant struct {
	ID            string `json:"id"`
	Tenant        Tenant `json:"tenant"`
	GrantType     string `json:"grantType"`
	ObjectID      string `json:"objectId"`
	AccessLevelID string `json:"accessLevelId"`
}

type Tenant struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
