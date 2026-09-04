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
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"resty.dev/v3"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/http-client"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

var _ subClientInterface = &netbackup{}

type netbackup struct {
	subclient

	mu                 sync.RWMutex
	accessToken        string
	storedRefreshToken string
	tokenExpiresAt     time.Time
	baseURL            string
}

func newNetbackupClient() subClientInterface {
	return &netbackup{}
}

func (n *netbackup) getID() string {
	return string(ClientNetbackup)
}

func (n *netbackup) newHTTPClient(ctx context.Context) (*resty.Client, error) {
	baseURL := n.getBaseURL()
	if baseURL == "" {
		return nil, fmt.Errorf("netbackup base URL not configured")
	}

	hC := httpclient.NewHTTPClient().
		SetBaseURL(baseURL).
		SetHeader("Accept", "application/json").
		SetResultError(netbackupError{})

	if err := n.ensureToken(ctx); err != nil {
		return nil, fmt.Errorf("failed to ensure netbackup token: %w", err)
	}

	hC.SetHeader("Authorization", "Bearer "+n.getAccessToken())

	return hC, nil
}

func (n *netbackup) parseAPIError(operation string, resp *resty.Response) *errors.APIError {
	if resp == nil || resp.StatusCode() < http.StatusBadRequest {
		return nil
	}

	if err, ok := resp.ResultError().(*netbackupError); ok {
		return &errors.APIError{
			Operation:  operation,
			StatusCode: resp.StatusCode(),
			Message:    err.Message,
			Duration:   resp.Duration(),
			Endpoint:   resp.Request.URL,
			Method:     resp.Request.Method,
			Err:        classifyStatusCode(resp.StatusCode()),
		}
	}

	return &errors.APIError{
		Operation:  operation,
		StatusCode: resp.StatusCode(),
		Message:    unknownErrorMessage,
		Duration:   resp.Duration(),
		Endpoint:   resp.Request.URL,
		Method:     resp.Request.Method,
		Err:        classifyStatusCode(resp.StatusCode()),
	}
}

func (n *netbackup) idempotentRetryCondition() resty.RetryConditionFunc {
	return func(resp *resty.Response, err error) bool {
		return false
	}
}

func (n *netbackup) ContextData(_ context.Context) ContextData {
	if n == nil {
		return ContextData{}
	}

	return ContextData{
		OrganizationID: n.getAccessToken(),
		SiteID:         n.getBaseURL(),
	}
}

func (n *netbackup) getBaseURL() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.baseURL != "" {
		return n.baseURL
	}
	return n.console.GetNetbackupEndpoint()
}

func (n *netbackup) getAccessToken() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.accessToken
}

func (n *netbackup) getRefreshToken() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.storedRefreshToken
}

func (n *netbackup) setTokens(accessToken, refreshToken string, expiresIn int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.accessToken = accessToken
	n.storedRefreshToken = refreshToken
	if expiresIn > 0 {
		n.tokenExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	} else {
		n.tokenExpiresAt = time.Time{}
	}
}

func (n *netbackup) isTokenExpired() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.tokenExpiresAt.IsZero() {
		return true
	}
	// Refresh 5 minutes before actual expiry
	return time.Now().Add(5 * time.Minute).After(n.tokenExpiresAt)
}

func (n *netbackup) ensureToken(ctx context.Context) error {
	if !n.isTokenExpired() && n.getAccessToken() != "" {
		return nil
	}

	if n.getRefreshToken() != "" {
		if err := n.refreshToken(ctx); err == nil {
			return nil
		}
		// If refresh fails, fall through to password grant
	}

	return n.getToken(ctx)
}

func (n *netbackup) getToken(ctx context.Context) error {
	username := os.Getenv("NETBACKUP_USERNAME")
	password := os.Getenv("NETBACKUP_PASSWORD")
	if username == "" || password == "" {
		return fmt.Errorf("NETBACKUP_USERNAME and NETBACKUP_PASSWORD environment variables must be set")
	}

	baseURL := n.getBaseURL()
	if baseURL == "" {
		return fmt.Errorf("netbackup base URL not configured")
	}

	httpC := httpclient.NewHTTPClient().SetBaseURL(baseURL)

	resp, err := httpC.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type": "password",
			"username":   username,
			"password":   password,
		}).
		SetResult(&netbackupTokenResponse{}).
		SetResultError(&netbackupError{}).
		Post("/NetBackupSelfService/Api/auth/token")
	if err != nil {
		return fmt.Errorf("failed to get netbackup token: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to get netbackup token: %s", resp.Status())
	}

	tokenResp, ok := resp.Result().(*netbackupTokenResponse)
	if !ok || tokenResp == nil {
		return fmt.Errorf("invalid netbackup token response")
	}

	n.setTokens(tokenResp.AccessToken, tokenResp.RefreshToken, tokenResp.ExpiresIn)
	return nil
}

func (n *netbackup) refreshToken(ctx context.Context) error {
	refreshToken := n.getRefreshToken()
	if refreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	baseURL := n.getBaseURL()
	if baseURL == "" {
		return fmt.Errorf("netbackup base URL not configured")
	}

	httpC := httpclient.NewHTTPClient().SetBaseURL(baseURL)

	resp, err := httpC.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
		}).
		SetResult(&netbackupTokenResponse{}).
		SetResultError(&netbackupError{}).
		Post("/NetBackupSelfService/Api/auth/token")
	if err != nil {
		return fmt.Errorf("failed to refresh netbackup token: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to refresh netbackup token: %s", resp.Status())
	}

	tokenResp, ok := resp.Result().(*netbackupTokenResponse)
	if !ok || tokenResp == nil {
		return fmt.Errorf("invalid netbackup token response")
	}

	n.setTokens(tokenResp.AccessToken, tokenResp.RefreshToken, tokenResp.ExpiresIn)
	return nil
}

func (n *netbackup) getSession() map[string]string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return map[string]string{
		"accessToken":  n.accessToken,
		"refreshToken": n.storedRefreshToken,
		"tokenExpires": n.tokenExpiresAt.Format(time.RFC3339),
		"baseURL":      n.baseURL,
	}
}

func (n *netbackup) restoreSession(data map[string]string) error {
	if data == nil {
		return fmt.Errorf("invalid session data")
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	n.accessToken = data["accessToken"]
	n.storedRefreshToken = data["refreshToken"]

	if expiresStr := data["tokenExpires"]; expiresStr != "" {
		if t, err := time.Parse(time.RFC3339, expiresStr); err == nil {
			n.tokenExpiresAt = t
		}
	}

	n.baseURL = data["baseURL"]

	return nil
}

func (n *netbackup) close() error {
	return nil
}

type netbackupError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type netbackupTokenResponse struct {
	AccessToken  string `json:"access_token"`  //nolint:tagliatelle
	TokenType    string `json:"token_type"`    //nolint:tagliatelle
	ExpiresIn    int    `json:"expires_in"`    //nolint:tagliatelle
	RefreshToken string `json:"refresh_token"` //nolint:tagliatelle
}

// WithNetbackupCredential configures NetBackup credentials from environment variables.
func WithNetbackupCredential() ClientOption {
	return func(s *settings) error {
		logger := xlogger.WithGroup("client").WithGroup("options").WithGroup("WithNetbackupCredential")

		username := os.Getenv("NETBACKUP_USERNAME")
		password := os.Getenv("NETBACKUP_PASSWORD")
		if username == "" || password == "" {
			logger.Error("NETBACKUP_USERNAME and NETBACKUP_PASSWORD environment variables must be set")
			return fmt.Errorf("NETBACKUP_USERNAME and NETBACKUP_PASSWORD environment variables must be set")
		}

		if _, ok := s.SubClients[ClientNetbackup]; !ok {
			s.SubClients[ClientNetbackup] = newNetbackupClient()
		}

		nb := s.SubClients[ClientNetbackup]
		nb.setConsole(s.Console)

		logger.Debug("NetBackup credential configured", "console", s.Console)
		return nil
	}
}
