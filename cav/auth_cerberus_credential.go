/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/orange-cloudavenue/common-go/validators"
	"resty.dev/v3"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/http-client"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

const (
	cerberusTokenTypeBearer = "Bearer"
	cerberusAuthPath        = "/auth/v1/user/token"
	cerberusExpiryBuffer    = 30 * time.Second
	formGrantTypeKey        = "grant_type"
)

type cerberusAuthProvider struct {
	logger       *slog.Logger
	httpC        *resty.Client
	organization string `validate:"required"`
	username     string `validate:"required"`
	password     string `validate:"required"`

	mu          sync.RWMutex
	accessToken string
	tokenType   string
	expiresAt   time.Time
}

type cerberusSession struct {
	organization string
	accessToken  string
	tokenType    string
	expiresAt    time.Time
}

type cerberusAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type cerberusErrorResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
}

func newCerberusAuthProvider(c consoles.ConsoleName, organization, username, password string) (authProvider, error) {
	cc := &cerberusAuthProvider{
		logger:       xlogger.WithGroup("auth").With("backend", ClientCerberus),
		organization: organization,
		username:     username,
		password:     password,
	}

	if err := validators.New().Var(cc.username, "required"); err != nil {
		cc.logger.Error("Failed to validate username", "error", err)
		return nil, err
	}
	if err := validators.New().Var(cc.password, "required"); err != nil {
		cc.logger.Error("Failed to validate password", "error", err)
		return nil, err
	}
	if err := validators.New().Var(cc.organization, "required"); err != nil {
		cc.logger.Error("Failed to validate organization", "error", err)
		return nil, err
	}
	if ok := consoles.IsValidOrganizationName(organization); !ok {
		cc.logger.Error("Invalid organization name", "organization", organization)
		return nil, errors.New("invalid organization name")
	}

	cc.httpC = httpclient.NewHTTPClient().SetBaseURL(c.GetAPICerberusEndpoint())

	return cc, nil
}

func (c *cerberusAuthProvider) Headers() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.accessToken == "" {
		return map[string]string{}
	}

	tokenType := c.tokenType
	if tokenType == "" {
		tokenType = cerberusTokenTypeBearer
	}

	return map[string]string{"Authorization": tokenType + " " + c.accessToken}
}

func (c *cerberusAuthProvider) CurrentAuthorization(ctx context.Context) (string, error) {
	if !c.IsInitialized() {
		if err := c.Refresh(ctx); err != nil {
			return "", err
		}
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.accessToken == "" {
		return "", errors.New("authorization not initialized")
	}

	return c.accessToken, nil
}

func (c *cerberusAuthProvider) Refresh(ctx context.Context) error {
	c.mu.RLock()
	if c.isInitializedLocked() && !c.isExpiredLocked() {
		c.mu.RUnlock()
		return nil
	}
	c.mu.RUnlock()

	resp, err := c.httpC.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		SetFormData(map[string]string{
			formGrantTypeKey: "client_credentials",
			"client_id":      c.username,
			"client_secret":  c.password,
			"scope":          "tenant:" + c.organization,
		}).
		SetResult(&cerberusAuthResponse{}).
		SetResultError(&cerberusErrorResponse{}).
		Post(cerberusAuthPath)
	if err != nil {
		return fmt.Errorf("authentication request failed: %w", err)
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusMultipleChoices {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.clearSessionLocked()
		return fmt.Errorf("authentication failed: %w", toCerberusError(resp))
	}

	authResp, ok := resp.Result().(*cerberusAuthResponse)
	if !ok || authResp == nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.clearSessionLocked()
		return errors.New("authentication failed: invalid response format")
	}
	if authResp.AccessToken == "" {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.clearSessionLocked()
		return errors.New("authentication failed: empty access token received")
	}

	tokenType := authResp.TokenType
	if tokenType == "" {
		tokenType = cerberusTokenTypeBearer
	}

	nextSession := cerberusSession{
		organization: c.organization,
		accessToken:  authResp.AccessToken,
		tokenType:    tokenType,
		expiresAt:    time.Now().Add(time.Duration(authResp.ExpiresIn) * time.Second),
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isInitializedLocked() && !c.isExpiredLocked() {
		return nil
	}
	c.applySessionLocked(nextSession)

	return nil
}

func (c *cerberusAuthProvider) IsInitialized() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isInitializedLocked() && !c.isExpiredLocked()
}

func (c *cerberusAuthProvider) getSession() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]string{
		"organization": c.organization,
		"accessToken":  c.accessToken,
		"tokenType":    c.tokenType,
		"expiresAt":    c.expiresAt.UTC().Format(time.RFC3339Nano),
	}
}

func (c *cerberusAuthProvider) restoreSession(data map[string]string) error {
	if data == nil {
		return errors.New("invalid session data")
	}

	var parsedExpiresAt time.Time
	setExpiresAt := false
	if rawExpiresAt := data["expiresAt"]; rawExpiresAt != "" {
		parsed, err := time.Parse(time.RFC3339Nano, rawExpiresAt)
		if err != nil {
			return fmt.Errorf("restore session: parse expiresAt: %w", err)
		}
		parsedExpiresAt = parsed
		setExpiresAt = true
	}

	nextSession := cerberusSession{
		organization: data["organization"],
		accessToken:  data["accessToken"],
		tokenType:    data["tokenType"],
		expiresAt:    time.Time{},
	}
	if setExpiresAt {
		nextSession.expiresAt = parsedExpiresAt
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.applySessionLocked(nextSession)

	return nil
}

func (c *cerberusAuthProvider) getExtraData() map[string]string {
	return map[string]string{}
}

func (c *cerberusAuthProvider) isInitializedLocked() bool {
	return c.accessToken != ""
}

func (c *cerberusAuthProvider) isExpiredLocked() bool {
	if c.expiresAt.IsZero() {
		return true
	}
	return c.expiresAt.Add(-cerberusExpiryBuffer).Before(time.Now())
}

func (c *cerberusAuthProvider) applySessionLocked(session cerberusSession) {
	c.organization = session.organization
	c.accessToken = session.accessToken
	c.tokenType = session.tokenType
	c.expiresAt = session.expiresAt
}

func (c *cerberusAuthProvider) clearSessionLocked() {
	c.accessToken = ""
	c.tokenType = ""
	c.expiresAt = time.Time{}
}

func (e *cerberusErrorResponse) formatError() string {
	var parts []string
	if e.Code != 0 {
		parts = append(parts, fmt.Sprintf("code:%d", e.Code))
	}
	if e.Message != "" {
		parts = append(parts, fmt.Sprintf("message:%s", e.Message))
	}
	if e.Description != "" {
		parts = append(parts, fmt.Sprintf("description:%s", e.Description))
	}

	return strings.Join(parts, ": ")
}

func toCerberusError(resp *resty.Response) error {
	statusCode := resp.StatusCode()

	cerberusErr, _ := resp.ResultError().(*cerberusErrorResponse)
	if cerberusErr != nil {
		if formatted := cerberusErr.formatError(); formatted != "" {
			return fmt.Errorf("HTTPCode:%d - %s", statusCode, formatted)
		}
	}

	body := strings.TrimSpace(resp.String())
	if body == "" {
		return fmt.Errorf("HTTPCode:%s", resp.Status())
	}

	return fmt.Errorf("HTTPCode:%s - body: %s", resp.Status(), body)
}
