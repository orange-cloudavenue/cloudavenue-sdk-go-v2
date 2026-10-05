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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

// doCerberusRequest serves given response and returns resty response,
// configured like Cerberus-routed customer API client (JSON result error).
func doCerberusRequest(t *testing.T, status int, contentType, body string) *resty.Response {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	c := resty.New().SetResultError(errors.CustomerAPIErrorBody{})
	t.Cleanup(func() { _ = c.Close() })

	resp, err := c.R().Get(srv.URL + "/test")
	require.NoError(t, err)

	return resp
}

func TestCerberusParseAPIError(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		contentType    string
		body           string
		wantCode       string
		wantMessage    string
		wantSentinel   error
		wantIsNotFound bool
	}{
		{
			name:         "BAD_REQUEST",
			status:       http.StatusBadRequest,
			contentType:  "application/json",
			body:         `{"minorErrorCode":"BAD_REQUEST","message":"invalid name","stackTrace":null}`,
			wantCode:     "BAD_REQUEST",
			wantMessage:  "BAD_REQUEST: invalid name",
			wantSentinel: errors.ErrBadRequest,
		},
		{
			name:           "RESOURCE_NOT_FOUND",
			status:         http.StatusNotFound,
			contentType:    "application/json",
			body:           `{"minorErrorCode":"RESOURCE_NOT_FOUND","message":"user not found","stackTrace":null}`,
			wantCode:       "RESOURCE_NOT_FOUND",
			wantMessage:    "RESOURCE_NOT_FOUND: user not found",
			wantSentinel:   errors.ErrNotFound,
			wantIsNotFound: true,
		},
		{
			name:         "ACCESS_TO_RESOURCE_IS_FORBIDDEN",
			status:       http.StatusForbidden,
			contentType:  "application/json",
			body:         `{"minorErrorCode":"ACCESS_TO_RESOURCE_IS_FORBIDDEN","message":"nope","stackTrace":null}`,
			wantCode:     "ACCESS_TO_RESOURCE_IS_FORBIDDEN",
			wantMessage:  "ACCESS_TO_RESOURCE_IS_FORBIDDEN: nope",
			wantSentinel: errors.ErrForbidden,
		},
		{
			name:         "INTERNAL_SERVER_ERROR",
			status:       http.StatusInternalServerError,
			contentType:  "application/json",
			body:         `{"minorErrorCode":"INTERNAL_SERVER_ERROR","message":"boom","stackTrace":null}`,
			wantCode:     "INTERNAL_SERVER_ERROR",
			wantMessage:  "INTERNAL_SERVER_ERROR: boom",
			wantSentinel: errors.ErrInternalServerError,
		},
		{
			name:         "legacy code field still honoured",
			status:       http.StatusBadRequest,
			contentType:  "application/json",
			body:         `{"code":"LEGACY","message":"old shape"}`,
			wantCode:     "LEGACY",
			wantMessage:  "LEGACY: old shape",
			wantSentinel: errors.ErrBadRequest,
		},
		{
			name:         "no code no reason: no leading separator",
			status:       http.StatusBadRequest,
			contentType:  "application/json",
			body:         `{"message":"only a message"}`,
			wantCode:     "",
			wantMessage:  "only a message",
			wantSentinel: errors.ErrBadRequest,
		},
		{
			name:         "empty message keeps code",
			status:       http.StatusBadRequest,
			contentType:  "application/json",
			body:         `{"minorErrorCode":"BAD_REQUEST","message":""}`,
			wantCode:     "BAD_REQUEST",
			wantMessage:  "BAD_REQUEST",
			wantSentinel: errors.ErrBadRequest,
		},
		{
			name:         "empty code and message falls back to unknown",
			status:       http.StatusBadRequest,
			contentType:  "application/json",
			body:         `{}`,
			wantCode:     "",
			wantMessage:  unknownErrorMessage,
			wantSentinel: errors.ErrBadRequest,
		},
		{
			name:         "non-JSON body falls back to unknown",
			status:       http.StatusBadGateway,
			contentType:  "text/html",
			body:         `<html><body>Bad gateway</body></html>`,
			wantCode:     "",
			wantMessage:  unknownErrorMessage,
			wantSentinel: errors.ErrBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doCerberusRequest(t, tt.status, tt.contentType, tt.body)

			apiErr := (&cerberus{}).parseAPIError("op", resp)
			require.NotNil(t, apiErr)

			assert.Equal(t, tt.wantCode, apiErr.StatusMessage)
			assert.Equal(t, tt.wantMessage, apiErr.Message)
			assert.Equal(t, tt.status, apiErr.StatusCode)
			assert.NotContains(t, apiErr.Message[:1], ":", "leading colon artifact")
			require.ErrorIs(t, apiErr, tt.wantSentinel)
			assert.Equal(t, tt.wantIsNotFound, apiErr.IsNotFound())
		})
	}
}

func TestCerberusParseAPIErrorSuccess(t *testing.T) {
	t.Run("nil response", func(t *testing.T) {
		require.Nil(t, (&cerberus{}).parseAPIError("op", nil))
	})

	tests := []struct {
		name        string
		status      int
		contentType string
		body        string
	}{
		{"200 JSON", http.StatusOK, "application/json", `{"ok":true}`},
		{"204 no content", http.StatusNoContent, "", ``},
		{"200 XML", http.StatusOK, "application/xml", `<a/>`},
		{"200 legit HTML without WAF marker", http.StatusOK, "text/html", `<html><body>hello</body></html>`},
		{"200 non-HTML mentioning rejection", http.StatusOK, "application/json", `{"message":"The requested URL was rejected"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doCerberusRequest(t, tt.status, tt.contentType, tt.body)
			require.Nil(t, (&cerberus{}).parseAPIError("op", resp))
		})
	}
}

func TestCerberusParseAPIErrorWAFRejection(t *testing.T) {
	const body = `<html><head><title>Request Rejected</title></head>` +
		`<body>The requested URL was rejected. Please consult with your administrator.` +
		`<br><br>Your support ID is: 3850989646093916041<br><br>` +
		`<a href='javascript:history.back();'>[Go Back]</a></body></html>`

	resp := doCerberusRequest(t, http.StatusOK, "text/html; charset=UTF-8", body)

	apiErr := (&cerberus{}).parseAPIError("op", resp)
	require.NotNil(t, apiErr, "WAF rejection must not be reported as success")

	assert.Equal(t, http.StatusOK, apiErr.StatusCode)
	assert.Equal(t, errors.CustomerAPIWAFRejectedCode(), apiErr.StatusMessage)
	assert.Contains(t, apiErr.Message, "3850989646093916041")
}

func TestCerberusContextData(t *testing.T) {
	data := (&cerberus{}).ContextData(t.Context())
	require.Equal(t, ContextData{}, data)
}
