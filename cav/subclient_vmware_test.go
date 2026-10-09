/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

func doVmwareRequest(t *testing.T, status int, contentType, body string) *resty.Response {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	c := resty.New().SetResultError(vmwareError{})
	t.Cleanup(func() { _ = c.Close() })

	resp, err := c.R().Post(srv.URL + "/test")
	require.NoError(t, err)

	return resp
}

func TestVmwareParseAPIErrorWAFRejection(t *testing.T) {
	const body = `<html><head><title>Request Rejected</title></head>` +
		`<body>The requested URL was rejected. Your support ID is: 123456789</body></html>`

	resp := doVmwareRequest(t, http.StatusOK, "text/html; charset=UTF-8", body)

	apiErr := (&vmware{}).parseAPIError("op", resp)
	require.NotNil(t, apiErr, "WAF rejection must not be reported as success")
	require.Equal(t, errors.CustomerAPIWAFRejectedCode(), apiErr.StatusMessage)
	require.Contains(t, apiErr.Message, "123456789")
	require.Equal(t, resp.Duration(), apiErr.Duration)
	require.Equal(t, resp.Request.URL, apiErr.Endpoint)
	require.Equal(t, resp.Request.Method, apiErr.Method)
}

func TestVmwareParseAPIErrorSuccess(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
	}{
		{name: "JSON", contentType: "application/json", body: `{"ok":true}`},
		{name: "HTML without WAF marker", contentType: "text/html", body: `<html><body>hello</body></html>`},
		{name: "non-HTML rejection marker", contentType: "application/json", body: `{"message":"The requested URL was rejected"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doVmwareRequest(t, http.StatusOK, tt.contentType, tt.body)
			require.Nil(t, (&vmware{}).parseAPIError("op", resp))
		})
	}
}

func TestVmwareParseAPIErrorMetadata(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "typed error", body: `{"message":"bad request","majorErrorCode":400,"minorErrorCode":"BAD_REQUEST"}`},
		{name: "fallback error", body: `not JSON`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doVmwareRequest(t, http.StatusBadRequest, "application/json", tt.body)
			apiErr := (&vmware{}).parseAPIError("op", resp)
			require.NotNil(t, apiErr)
			require.Equal(t, resp.Duration(), apiErr.Duration)
			require.Equal(t, resp.Request.URL, apiErr.Endpoint)
			require.Equal(t, resp.Request.Method, apiErr.Method)
		})
	}
}

func TestVmwareContextData(t *testing.T) {
	v := &vmware{subclient: subclient{credential: &cloudavenueCredential{
		organizationID: "org-1",
		siteID:         "site-1",
	}}}

	data := v.ContextData(t.Context())
	require.Equal(t, "org-1", data.OrganizationID)
	require.Equal(t, "site-1", data.SiteID)
}
