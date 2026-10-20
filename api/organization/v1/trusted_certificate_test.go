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

func TestListTrustedCertificate(t *testing.T) {
	certificateID := "urn:vcloud:trustedCertificate:" + generator.MustGenerate("{uuid}")
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListTrustedCertificate())
	ms.SetResponse(endpoints.ListTrustedCertificate(), &itypes.APIResponseListTrustedCertificate{Values: []itypes.APIResponseTrustedCertificate{{
		ID:          certificateID,
		Alias:       "trusted-1",
		Certificate: "-----BEGIN CERTIFICATE-----trusted-----END CERTIFICATE-----",
	}}}, nil)

	resp, err := client.ListTrustedCertificate(t.Context(), types.ParamsListTrustedCertificate{Name: "trusted-1"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Certificates, 1)
	assert.Equal(t, certificateID, resp.Certificates[0].ID)

	ms.CleanResponse(endpoints.ListTrustedCertificate())
}

func TestGetTrustedCertificateByName(t *testing.T) {
	certificateID := "urn:vcloud:trustedCertificate:" + generator.MustGenerate("{uuid}")
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListTrustedCertificate())
	ms.SetResponse(endpoints.ListTrustedCertificate(), &itypes.APIResponseListTrustedCertificate{Values: []itypes.APIResponseTrustedCertificate{{
		ID:          certificateID,
		Alias:       "trusted-1",
		Certificate: "-----BEGIN CERTIFICATE-----trusted-----END CERTIFICATE-----",
	}}}, nil)

	resp, err := client.GetTrustedCertificate(t.Context(), types.ParamsGetTrustedCertificate{Name: "trusted-1"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, certificateID, resp.ID)
	assert.Equal(t, "trusted-1", resp.Name)

	ms.CleanResponse(endpoints.ListTrustedCertificate())
}

func TestCreateTrustedCertificate(t *testing.T) {
	certificateID := "urn:vcloud:trustedCertificate:" + generator.MustGenerate("{uuid}")
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListTrustedCertificate())
	ms.SetResponse(endpoints.ListTrustedCertificate(), &itypes.APIResponseListTrustedCertificate{Values: []itypes.APIResponseTrustedCertificate{}}, nil)
	ms.CleanResponse(endpoints.CreateTrustedCertificate())
	ms.SetResponse(endpoints.CreateTrustedCertificate(), &itypes.APIResponseTrustedCertificate{
		ID:          certificateID,
		Alias:       "trusted-1",
		Certificate: "-----BEGIN CERTIFICATE-----trusted-----END CERTIFICATE-----",
	}, nil)

	resp, err := client.CreateTrustedCertificate(t.Context(), types.ParamsCreateTrustedCertificate{
		Name:        "trusted-1",
		Certificate: "-----BEGIN CERTIFICATE-----trusted-----END CERTIFICATE-----",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, certificateID, resp.ID)

	ms.CleanResponse(endpoints.ListTrustedCertificate())
	ms.CleanResponse(endpoints.CreateTrustedCertificate())
}

func TestUpdateTrustedCertificate(t *testing.T) {
	certificateID := "urn:vcloud:trustedCertificate:" + generator.MustGenerate("{uuid}")
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetTrustedCertificate())
	ms.SetResponse(endpoints.GetTrustedCertificate(), &itypes.APIResponseTrustedCertificate{
		ID:          certificateID,
		Alias:       "trusted-1",
		Certificate: "-----BEGIN CERTIFICATE-----old-----END CERTIFICATE-----",
	}, nil)
	ms.CleanResponse(endpoints.UpdateTrustedCertificate())
	ms.SetResponse(endpoints.UpdateTrustedCertificate(), &itypes.APIResponseTrustedCertificate{
		ID:          certificateID,
		Alias:       "trusted-2",
		Certificate: "-----BEGIN CERTIFICATE-----new-----END CERTIFICATE-----",
	}, nil)

	resp, err := client.UpdateTrustedCertificate(t.Context(), types.ParamsUpdateTrustedCertificate{
		ID:          certificateID,
		Name:        "trusted-2",
		Certificate: "-----BEGIN CERTIFICATE-----new-----END CERTIFICATE-----",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "trusted-2", resp.Name)
	assert.Equal(t, "-----BEGIN CERTIFICATE-----new-----END CERTIFICATE-----", resp.Certificate)

	ms.CleanResponse(endpoints.GetTrustedCertificate())
	ms.CleanResponse(endpoints.UpdateTrustedCertificate())
}

func TestDeleteTrustedCertificate(t *testing.T) {
	certificateID := "urn:vcloud:trustedCertificate:" + generator.MustGenerate("{uuid}")
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetTrustedCertificate())
	ms.SetResponse(endpoints.GetTrustedCertificate(), &itypes.APIResponseTrustedCertificate{
		ID:          certificateID,
		Alias:       "trusted-1",
		Certificate: "-----BEGIN CERTIFICATE-----trusted-----END CERTIFICATE-----",
	}, nil)
	ms.CleanResponse(endpoints.DeleteTrustedCertificate())

	err := client.DeleteTrustedCertificate(t.Context(), types.ParamsDeleteTrustedCertificate{ID: certificateID})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetTrustedCertificate())
	ms.CleanResponse(endpoints.DeleteTrustedCertificate())
}

func TestGetTrustedCertificateByNameNotFound(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListTrustedCertificate())
	ms.SetResponse(endpoints.ListTrustedCertificate(), &itypes.APIResponseListTrustedCertificate{Values: []itypes.APIResponseTrustedCertificate{}}, nil)

	resp, err := client.GetTrustedCertificate(t.Context(), types.ParamsGetTrustedCertificate{Name: "missing-trusted"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "not found")

	ms.CleanResponse(endpoints.ListTrustedCertificate())
}

func TestGetTrustedCertificateByNameMultipleMatches(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListTrustedCertificate())
	ms.SetResponse(endpoints.ListTrustedCertificate(), &itypes.APIResponseListTrustedCertificate{Values: []itypes.APIResponseTrustedCertificate{
		{ID: "urn:vcloud:trustedCertificate:" + generator.MustGenerate("{uuid}"), Alias: "trusted-1", Certificate: "-----BEGIN CERTIFICATE-----one-----END CERTIFICATE-----"},
		{ID: "urn:vcloud:trustedCertificate:" + generator.MustGenerate("{uuid}"), Alias: "trusted-1", Certificate: "-----BEGIN CERTIFICATE-----two-----END CERTIFICATE-----"},
	}}, nil)

	resp, err := client.GetTrustedCertificate(t.Context(), types.ParamsGetTrustedCertificate{Name: "trusted-1"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "multiple trusted certificates found")

	ms.CleanResponse(endpoints.ListTrustedCertificate())
}

func TestCreateTrustedCertificateRejectsDuplicateAlias(t *testing.T) {
	certificateID := "urn:vcloud:trustedCertificate:" + generator.MustGenerate("{uuid}")
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListTrustedCertificate())
	ms.SetResponse(endpoints.ListTrustedCertificate(), &itypes.APIResponseListTrustedCertificate{Values: []itypes.APIResponseTrustedCertificate{{
		ID:          certificateID,
		Alias:       "trusted-1",
		Certificate: "-----BEGIN CERTIFICATE-----trusted-----END CERTIFICATE-----",
	}}}, nil)

	resp, err := client.CreateTrustedCertificate(t.Context(), types.ParamsCreateTrustedCertificate{Name: "trusted-1", Certificate: "-----BEGIN CERTIFICATE-----trusted-----END CERTIFICATE-----"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "already exists")

	ms.CleanResponse(endpoints.ListTrustedCertificate())
}

func TestDeleteTrustedCertificateByNameNotFound(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListTrustedCertificate())
	ms.SetResponse(endpoints.ListTrustedCertificate(), &itypes.APIResponseListTrustedCertificate{Values: []itypes.APIResponseTrustedCertificate{}}, nil)

	err := client.DeleteTrustedCertificate(t.Context(), types.ParamsDeleteTrustedCertificate{Name: "missing-trusted"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	ms.CleanResponse(endpoints.ListTrustedCertificate())
}
