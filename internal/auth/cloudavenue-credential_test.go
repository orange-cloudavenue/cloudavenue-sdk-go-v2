/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package auth

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

var getMockConsole = func() consoles.Console {
	c, _ := consoles.FindBySiteID("mock")
	return c
}

func getMockRestyClient() *resty.Client {
	client := resty.New()
	httpmock.ActivateNonDefault(client.Client())
	return client
}

func TestCloudavenueCredential_Headers_WithBearer(t *testing.T) {
	httpC := resty.New()
	console := getMockConsole()
	cred := &CloudavenueCredential{
		httpC:        httpC,
		console:      console,
		organization: "org",
		username:     "user",
		password:     "pass",
		bearer:       "mockBearerToken",
	}

	headers := cred.Headers()
	assert.Equal(t, "Bearer mockBearerToken", headers["Authorization"])
}

func TestCloudavenueCredential_Headers_WithoutBearer(t *testing.T) {
	httpC := resty.New()
	console := getMockConsole()
	cred := &CloudavenueCredential{
		httpC:        httpC,
		console:      console,
		organization: "org",
		username:     "user",
		password:     "pass",
		bearer:       "",
	}

	headers := cred.Headers()
	assert.Equal(t, "Bearer ", headers["Authorization"])
}

func TestCloudavenueCredential_Refresh_WithBearer(t *testing.T) {
	client := getMockRestyClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"http://mock.api/cloudapi/1.0.0/sessions",
		httpmock.NewStringResponder(200, "ok").
			HeaderSet(http.Header{
				XVmwareAccessToken: []string{"token-from-header"},
			}),
	)

	cred := &CloudavenueCredential{
		httpC:        client,
		console:      getMockConsole(),
		organization: "org",
		username:     "user",
		password:     "pass",
		bearer:       "existingBearer",
	}

	err := cred.Refresh(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, "token-from-header", cred.bearer)
}

func TestCloudavenueCredential_Refresh_WithoutBearer(t *testing.T) {
	client := getMockRestyClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"http://mock.api/cloudapi/1.0.0/sessions",
		httpmock.NewStringResponder(200, "ok").
			HeaderSet(http.Header{
				XVmwareAccessToken: []string{"token-basic-auth"},
			}),
	)

	cred := &CloudavenueCredential{
		httpC:        client,
		console:      getMockConsole(),
		organization: "org",
		username:     "user",
		password:     "pass",
		bearer:       "",
	}

	err := cred.Refresh(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, "token-basic-auth", cred.bearer)
}

func TestCloudavenueCredential_Refresh_Error(t *testing.T) {
	client := getMockRestyClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"http://mock.api/cloudapi/1.0.0/sessions",
		httpmock.NewErrorResponder(errors.New("network error")),
	)

	cred := &CloudavenueCredential{
		httpC:        client,
		console:      getMockConsole(),
		organization: "org",
		username:     "user",
		password:     "pass",
		bearer:       "",
	}

	err := cred.Refresh(t.Context())
	assert.Error(t, err)
	assert.Empty(t, cred.bearer)
}

func TestCloudavenueCredential_Refresh_AuthFailed(t *testing.T) {
	client := getMockRestyClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"http://mock.api/cloudapi/1.0.0/sessions",
		httpmock.NewStringResponder(401, "unauthorized"),
	)

	cred := &CloudavenueCredential{
		httpC:        client,
		console:      getMockConsole(),
		organization: "org",
		username:     "user",
		password:     "pass",
		bearer:       "",
	}

	err := cred.Refresh(t.Context())
	assert.Error(t, err)
	assert.Empty(t, cred.bearer)
}

func TestCloudavenueCredential_IsInitialized(t *testing.T) {
	credWithBearer := &CloudavenueCredential{
		bearer: "sometoken",
	}
	assert.True(t, credWithBearer.IsInitialized(), "Should be initialized when bearer is set")

	credWithoutBearer := &CloudavenueCredential{
		bearer: "",
	}
	assert.False(t, credWithoutBearer.IsInitialized(), "Should not be initialized when bearer is empty")
}

func TestCloudavenueCredential_NewCloudavenueCredential(t *testing.T) {
	cred, err := newCloudavenueCredential(getMockRestyClient(), getMockConsole(), "org", "user", "pass")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.NotNil(t, cred)
}

func TestCloudavenueCredential_NewCloudavenueCredential_NilClient(t *testing.T) {
	_, err := newCloudavenueCredential(nil, getMockConsole(), "org", "user", "pass")
	if err == nil {
		t.Fatal("expected error for nil http client, got nil")
	}
	assert.Equal(t, "http client cannot be nil", err.Error())
}

func TestCloudavenueCredential_NewCloudavenueCredential_EmptyOrganization(t *testing.T) {
	_, err := newCloudavenueCredential(getMockRestyClient(), getMockConsole(), "", "user", "pass")
	if err == nil {
		t.Fatal("expected error for empty organization, got nil")
	}
	assert.Equal(t, "organization cannot be empty", err.Error())
}

func TestCloudavenueCredential_NewCloudavenueCredential_EmptyUsername(t *testing.T) {
	_, err := newCloudavenueCredential(getMockRestyClient(), getMockConsole(), "org", "", "pass")
	if err == nil {
		t.Fatal("expected error for empty username, got nil")
	}
	assert.Equal(t, "username cannot be empty", err.Error())
}

func TestCloudavenueCredential_NewCloudavenueCredential_EmptyPassword(t *testing.T) {
	_, err := newCloudavenueCredential(getMockRestyClient(), getMockConsole(), "org", "user", "")
	if err == nil {
		t.Fatal("expected error for empty password, got nil")
	}
	assert.Equal(t, "password cannot be empty", err.Error())
}
