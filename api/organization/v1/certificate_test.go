/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package organization

import (
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestListCertificate(t *testing.T) {
	certificateID := generator.MustGenerate("{urn:certificateLibraryItem}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListCertificate())
	ms.SetResponse(endpoints.ListCertificate(), &itypes.ApiResponseListCertificate{Values: []itypes.ApiResponseCertificate{{
		ID:          certificateID,
		Alias:       "cert-1",
		Description: "desc",
		Certificate: "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----",
	}}}, nil)

	resp, err := client.ListCertificate(t.Context(), types.ParamsListCertificate{Name: "cert-1"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Certificates, 1)
	assert.Equal(t, certificateID, resp.Certificates[0].ID)
	assert.Equal(t, "cert-1", resp.Certificates[0].Name)

	ms.CleanResponse(endpoints.ListCertificate())
}

func TestGetCertificateByName(t *testing.T) {
	certificateID := generator.MustGenerate("{urn:certificateLibraryItem}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListCertificate())
	ms.SetResponse(endpoints.ListCertificate(), &itypes.ApiResponseListCertificate{Values: []itypes.ApiResponseCertificate{{
		ID:          certificateID,
		Alias:       "cert-1",
		Description: "desc",
		Certificate: "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----",
	}}}, nil)

	resp, err := client.GetCertificate(t.Context(), types.ParamsGetCertificate{Name: "cert-1"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, certificateID, resp.ID)
	assert.Equal(t, "cert-1", resp.Name)
	assert.Equal(t, "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----", resp.Certificate)

	ms.CleanResponse(endpoints.ListCertificate())
}

func TestGetCertificateByNameNotFound(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListCertificate())
	ms.SetResponse(endpoints.ListCertificate(), &itypes.ApiResponseListCertificate{Values: []itypes.ApiResponseCertificate{}}, nil)

	resp, err := client.GetCertificate(t.Context(), types.ParamsGetCertificate{Name: "missing-cert"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "not found")

	ms.CleanResponse(endpoints.ListCertificate())
}

func TestGetCertificateByNameMultipleMatches(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListCertificate())
	ms.SetResponse(endpoints.ListCertificate(), &itypes.ApiResponseListCertificate{Values: []itypes.ApiResponseCertificate{
		{ID: generator.MustGenerate("{urn:certificateLibraryItem}"), Alias: "cert-1", Certificate: "-----BEGIN CERTIFICATE-----one-----END CERTIFICATE-----"},
		{ID: generator.MustGenerate("{urn:certificateLibraryItem}"), Alias: "cert-1", Certificate: "-----BEGIN CERTIFICATE-----two-----END CERTIFICATE-----"},
	}}, nil)

	resp, err := client.GetCertificate(t.Context(), types.ParamsGetCertificate{Name: "cert-1"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "multiple certificates found")

	ms.CleanResponse(endpoints.ListCertificate())
}

func TestCreateCertificate(t *testing.T) {
	certificateID := generator.MustGenerate("{urn:certificateLibraryItem}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.CreateCertificate())
	ms.SetResponse(endpoints.CreateCertificate(), &itypes.ApiResponseCertificate{
		ID:          certificateID,
		Alias:       "cert-1",
		Description: "desc",
		Certificate: "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----",
	}, nil)

	resp, err := client.CreateCertificate(t.Context(), types.ParamsCreateCertificate{
		Name:        "cert-1",
		Description: "desc",
		Certificate: "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----",
		PrivateKey:  "-----BEGIN PRIVATE KEY-----mock-----END PRIVATE KEY-----",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, certificateID, resp.ID)
	assert.Equal(t, "cert-1", resp.Name)

	ms.CleanResponse(endpoints.CreateCertificate())
}

func TestCreateCertificateRejectsDuplicateAlias(t *testing.T) {
	certificateID := generator.MustGenerate("{urn:certificateLibraryItem}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListCertificate())
	ms.SetResponse(endpoints.ListCertificate(), &itypes.ApiResponseListCertificate{Values: []itypes.ApiResponseCertificate{{
		ID:          certificateID,
		Alias:       "cert-1",
		Description: "desc",
		Certificate: "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----",
	}}}, nil)

	resp, err := client.CreateCertificate(t.Context(), types.ParamsCreateCertificate{
		Name:        "cert-1",
		Description: "desc",
		Certificate: "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "already exists")

	ms.CleanResponse(endpoints.ListCertificate())
}

func TestCreateCertificateRequiresPrivateKeyForPassphrase(t *testing.T) {
	client, _ := newClient(t)

	resp, err := client.CreateCertificate(t.Context(), types.ParamsCreateCertificate{
		Name:                 "cert-1",
		Certificate:          "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----",
		PrivateKeyPassphrase: "secret",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "private_key is required")
}

func TestUpdateCertificate(t *testing.T) {
	certificateID := generator.MustGenerate("{urn:certificateLibraryItem}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetCertificate())
	ms.SetResponse(endpoints.GetCertificate(), &itypes.ApiResponseCertificate{
		ID:          certificateID,
		Alias:       "cert-1",
		Description: "old-desc",
		Certificate: "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----",
	}, nil)
	ms.CleanResponse(endpoints.UpdateCertificate())

	resp, err := client.UpdateCertificate(t.Context(), types.ParamsUpdateCertificate{
		ID:          certificateID,
		Name:        "cert-2",
		Description: "new-desc",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, certificateID, resp.ID)
	assert.Equal(t, "cert-2", resp.Name)
	assert.Equal(t, "new-desc", resp.Description)
	assert.Equal(t, "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----", resp.Certificate)

	ms.CleanResponse(endpoints.GetCertificate())
	ms.CleanResponse(endpoints.UpdateCertificate())
}

func TestUpdateCertificateKeepsImmutableContent(t *testing.T) {
	certificateID := generator.MustGenerate("{urn:certificateLibraryItem}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetCertificate())
	ms.SetResponse(endpoints.GetCertificate(), &itypes.ApiResponseCertificate{
		ID:          certificateID,
		Alias:       "cert-1",
		Description: "old-desc",
		Certificate: "-----BEGIN CERTIFICATE-----immutable-----END CERTIFICATE-----",
	}, nil)
	ms.CleanResponse(endpoints.UpdateCertificate())

	resp, err := client.UpdateCertificate(t.Context(), types.ParamsUpdateCertificate{
		ID:          certificateID,
		Description: "new-desc",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "-----BEGIN CERTIFICATE-----immutable-----END CERTIFICATE-----", resp.Certificate)

	ms.CleanResponse(endpoints.GetCertificate())
	ms.CleanResponse(endpoints.UpdateCertificate())
}

func TestDeleteCertificate(t *testing.T) {
	certificateID := generator.MustGenerate("{urn:certificateLibraryItem}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetCertificate())
	ms.SetResponse(endpoints.GetCertificate(), &itypes.ApiResponseCertificate{
		ID:          certificateID,
		Alias:       "cert-1",
		Description: "desc",
		Certificate: "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----",
	}, nil)
	ms.CleanResponse(endpoints.DeleteCertificate())

	err := client.DeleteCertificate(t.Context(), types.ParamsDeleteCertificate{ID: certificateID})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetCertificate())
	ms.CleanResponse(endpoints.DeleteCertificate())
}

func TestDeleteCertificateByNameNotFound(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListCertificate())
	ms.SetResponse(endpoints.ListCertificate(), &itypes.ApiResponseListCertificate{Values: []itypes.ApiResponseCertificate{}}, nil)

	err := client.DeleteCertificate(t.Context(), types.ParamsDeleteCertificate{Name: "missing-cert"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	ms.CleanResponse(endpoints.ListCertificate())
}
