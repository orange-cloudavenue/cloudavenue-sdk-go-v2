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

func TestListCertificateConsumers(t *testing.T) {
	certificateID := generator.MustGenerate("{urn:certificateLibraryItem}")
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetCertificate())
	ms.SetResponse(endpoints.GetCertificate(), &itypes.APIResponseCertificate{ID: certificateID, Alias: "cert-1", Certificate: "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----"}, nil)
	ms.CleanResponse(endpoints.ListCertificateConsumers())
	ms.SetResponse(endpoints.ListCertificateConsumers(), &itypes.APIEntityReferences{Values: []itypes.APIEntityReference{{ID: generator.MustGenerate("{urn:vdcGroup}"), Name: "vdcgroup-1"}}}, nil)

	resp, err := client.ListCertificateConsumers(t.Context(), types.ParamsListCertificateConsumers{CertificateID: certificateID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.References, 1)
	assert.Equal(t, "vdcgroup-1", resp.References[0].Name)

	ms.CleanResponse(endpoints.GetCertificate())
	ms.CleanResponse(endpoints.ListCertificateConsumers())
}

func TestAddCertificateConsumer(t *testing.T) {
	certificateID := generator.MustGenerate("{urn:certificateLibraryItem}")
	consumerID := generator.MustGenerate("{urn:vdcGroup}")
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetCertificate())
	ms.SetResponse(endpoints.GetCertificate(), &itypes.APIResponseCertificate{ID: certificateID, Alias: "cert-1", Certificate: "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----"}, nil)
	ms.CleanResponse(endpoints.AddCertificateConsumer())
	ms.SetResponse(endpoints.AddCertificateConsumer(), &itypes.APIEntityReference{ID: consumerID, Name: "vdcgroup-1"}, nil)

	resp, err := client.AddCertificateConsumer(t.Context(), types.ParamsAddCertificateConsumer{CertificateID: certificateID, ConsumerID: consumerID, ConsumerName: "vdcgroup-1"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, consumerID, resp.ID)
	assert.Equal(t, "vdcgroup-1", resp.Name)

	ms.CleanResponse(endpoints.GetCertificate())
	ms.CleanResponse(endpoints.AddCertificateConsumer())
}

func TestSetCertificateConsumers(t *testing.T) {
	certificateID := generator.MustGenerate("{urn:certificateLibraryItem}")
	consumerID1 := generator.MustGenerate("{urn:vdcGroup}")
	consumerID2 := generator.MustGenerate("{urn:vdc}")
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetCertificate())
	ms.SetResponse(endpoints.GetCertificate(), &itypes.APIResponseCertificate{ID: certificateID, Alias: "cert-1", Certificate: "-----BEGIN CERTIFICATE-----mock-----END CERTIFICATE-----"}, nil)
	ms.CleanResponse(endpoints.SetCertificateConsumers())
	ms.SetResponse(endpoints.SetCertificateConsumers(), &itypes.APIEntityReferences{Values: []itypes.APIEntityReference{{ID: consumerID1, Name: "vdcgroup-1"}, {ID: consumerID2, Name: "edge-1"}}}, nil)

	resp, err := client.SetCertificateConsumers(t.Context(), types.ParamsSetCertificateConsumers{
		CertificateID: certificateID,
		Consumers:     []types.ModelEntityReference{{ID: consumerID1, Name: "vdcgroup-1"}, {ID: consumerID2, Name: "edge-1"}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.References, 2)
	assert.Equal(t, "edge-1", resp.References[1].Name)

	ms.CleanResponse(endpoints.GetCertificate())
	ms.CleanResponse(endpoints.SetCertificateConsumers())
}
