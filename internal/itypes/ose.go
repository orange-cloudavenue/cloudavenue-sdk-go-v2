/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

// APIResponseAssociatedTenant represents a tenant associated with the current user.
type APIResponseAssociatedTenant struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// APIResponseListAssociatedTenants represents the response for listing associated tenants.
type APIResponseListAssociatedTenants []APIResponseAssociatedTenant

// APIResponseS3Credentials represents S3 access credentials returned by OSE.
type APIResponseS3Credentials struct {
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	SessionToken string `json:"sessionToken,omitempty"`
	Expiration   string `json:"expiration,omitempty"`
}
