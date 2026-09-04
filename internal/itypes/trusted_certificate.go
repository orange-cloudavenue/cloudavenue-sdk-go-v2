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
	APIResponseListTrustedCertificate struct {
		Values []APIResponseTrustedCertificate `json:"values" fakesize:"2"`
	}

	APIResponseTrustedCertificate struct {
		ID          string `json:"id,omitempty" fake:"{urn:trustedCertificate}"`
		Alias       string `json:"alias,omitempty" fake:"mock-trusted-cert-{word}"`
		Certificate string `json:"certificate,omitempty" fake:"-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----"`
	}

	APIRequestTrustedCertificate struct {
		ID          string `json:"id,omitempty" fake:"{urn:trustedCertificate}"`
		Alias       string `json:"alias,omitempty" fake:"mock-trusted-cert-{word}"`
		Certificate string `json:"certificate,omitempty" fake:"-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----"`
	}
)

func (r *APIResponseListTrustedCertificate) ToModel() *types.ModelListTrustedCertificate {
	model := &types.ModelListTrustedCertificate{Certificates: make([]types.ModelGetTrustedCertificate, 0, len(r.Values))}
	for _, certificate := range r.Values {
		model.Certificates = append(model.Certificates, certificate.ToModel())
	}
	return model
}

func (r *APIResponseTrustedCertificate) ToModel() types.ModelGetTrustedCertificate {
	return types.ModelGetTrustedCertificate{ID: r.ID, Name: r.Alias, Certificate: r.Certificate}
}

func (r *APIRequestTrustedCertificate) ToModel() types.ModelGetTrustedCertificate {
	return types.ModelGetTrustedCertificate{ID: r.ID, Name: r.Alias, Certificate: r.Certificate}
}
