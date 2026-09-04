/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"

type (
	ApiEntityReferences struct {
		Values []ApiEntityReference `json:"values" fakesize:"2"`
	}

	ApiEntityReference struct {
		Name string `json:"name,omitempty" fake:"mock-entity-{word}"`
		ID   string `json:"id,omitempty" fake:"{urn}"`
	}

	ApiResponseListCertificate struct {
		Values []ApiResponseCertificate `json:"values" fakesize:"2"`
	}

	ApiResponseCertificate struct {
		ID                   string `json:"id,omitempty" fake:"{urn:certificateLibraryItem}"`
		Alias                string `json:"alias,omitempty" fake:"mock-certificate-{word}"`
		PrivateKey           string `json:"privateKey,omitempty" fake:"-----BEGIN PRIVATE KEY-----mock-----END PRIVATE KEY-----"`
		PrivateKeyPassphrase string `json:"privateKeyPassphrase,omitempty" fake:"{password}"`
		Certificate          string `json:"certificate,omitempty" fake:"-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----"`
		Description          string `json:"description,omitempty" fake:"{sentence}"`
		ConsumerCount        int    `json:"consumerCount,omitempty" fake:"{number:0,5}"`
	}

	ApiRequestCertificate struct {
		ID                   string `json:"id,omitempty" fake:"{urn:certificateLibraryItem}"`
		Alias                string `json:"alias,omitempty" fake:"mock-certificate-{word}"`
		PrivateKey           string `json:"privateKey,omitempty" fake:"-----BEGIN PRIVATE KEY-----mock-----END PRIVATE KEY-----"`
		PrivateKeyPassphrase string `json:"privateKeyPassphrase,omitempty" fake:"{password}"`
		Certificate          string `json:"certificate,omitempty" fake:"-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----"`
		Description          string `json:"description,omitempty" fake:"{sentence}"`
		ConsumerCount        int    `json:"consumerCount,omitempty" fake:"{number:0,5}"`
	}
)

func (r *ApiResponseListCertificate) ToModel() *types.ModelListCertificate {
	model := &types.ModelListCertificate{Certificates: make([]types.ModelGetCertificate, 0, len(r.Values))}
	for _, certificate := range r.Values {
		model.Certificates = append(model.Certificates, certificate.ToModel())
	}
	return model
}

func (r *ApiEntityReferences) ToModel() *types.ModelListEntityReference {
	model := &types.ModelListEntityReference{References: make([]types.ModelEntityReference, 0, len(r.Values))}
	for _, ref := range r.Values {
		model.References = append(model.References, ref.ToModel())
	}
	return model
}

func (r *ApiEntityReference) ToModel() types.ModelEntityReference {
	return types.ModelEntityReference{ID: r.ID, Name: r.Name}
}

func (r *ApiResponseCertificate) ToModel() types.ModelGetCertificate {
	return types.ModelGetCertificate{
		ID:            r.ID,
		Name:          r.Alias,
		Description:   r.Description,
		Certificate:   r.Certificate,
		ConsumerCount: r.ConsumerCount,
	}
}

func (r *ApiRequestCertificate) ToModel() types.ModelGetCertificate {
	return types.ModelGetCertificate{
		ID:            r.ID,
		Name:          r.Alias,
		Description:   r.Description,
		Certificate:   r.Certificate,
		ConsumerCount: r.ConsumerCount,
	}
}
