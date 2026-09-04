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
	"log/slog"
	"time"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

// Client is SDK runtime client.
type Client interface {
	NewRequest(ctx context.Context, endpoint *Endpoint, opts ...RequestOption) (req *resty.Request, err error)
	NewRequestWithBackend(ctx context.Context, backend BackendTarget, endpoint *Endpoint, opts ...RequestOption) (req *resty.Request, err error)
	Logger() *slog.Logger
	Do(ctx context.Context, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error)
	DoWithBackend(ctx context.Context, backend BackendTarget, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error)
	GetConsole() consoles.ConsoleName
	Close() error
}

type client struct {
	logger             *slog.Logger
	console            consoles.ConsoleName
	clientsInitialized map[subClientName]subClientInterface

	cachePassphrase, cachePath string
}

// NewClient creates client bound to organization.
func NewClient(organization string, opts ...ClientOption) (Client, error) {
	settings := newSettings(organization)

	if err := withConsole()(settings); err != nil {
		return nil, err
	}

	client := &client{
		console: settings.Console,
	}

	for _, opt := range opts {
		if err := opt(settings); err != nil {
			return nil, err
		}
	}

	client.logger = xlogger.WithGroup("client").With("organization", settings.Organization)
	client.clientsInitialized = settings.SubClients

	if settings.CachePassphrase != "" && settings.CachePath != "" {
		if err := client.restoreSessionsFromCache(settings.CachePassphrase, settings.CachePath); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// Logger returns client logger.
func (c *client) Logger() *slog.Logger {
	return c.logger
}

// GetConsole returns client console.
func (c *client) GetConsole() consoles.ConsoleName {
	return c.console
}

// Close closes client and persists cache when configured.
func (c *client) Close() error {
	errGroup := []error{}

	for _, subClient := range c.clientsInitialized {
		if err := subClient.close(); err != nil {
			errGroup = append(errGroup, fmt.Errorf("failed to close subclient: %w", err))
		}
	}

	if len(errGroup) > 0 {
		c.logger.Error("failed to close some subclients", "errors", errGroup)
		return fmt.Errorf("failed to close some subclients: %v", errGroup)
	}

	if c.cachePassphrase != "" && c.cachePath != "" {
		if err := c.storeSessionsToCache(c.cachePassphrase, c.cachePath); err != nil {
			c.logger.Error("failed to store sessions to cache", "error", err)
			return err
		}
	}

	c.logger.Debug("Closing client", "console", c.console)
	return nil
}

// identifyClient returns subclient backing backend.
func (c *client) identifyClient(_ context.Context, backend BackendTarget) (subClientInterface, error) {
	clientName, err := backendToSubClientName(backend)
	if err != nil {
		return nil, err
	}
	if c.clientsInitialized[clientName] == nil {
		return nil, fmt.Errorf("invalid client %s", clientName)
	}
	return c.clientsInitialized[clientName], nil
}

// backendToSubClientName maps backends to runtime subclients.
func backendToSubClientName(backend BackendTarget) (subClientName, error) {
	switch backend {
	case BackendInfrapi:
		return ClientCerberus, nil
	case BackendVMware:
		return ClientVmware, nil
	case BackendOSE:
		return ClientOSE, nil
	case BackendNetBackup:
		return ClientNetbackup, nil
	default:
		return "", fmt.Errorf("unknown backend %d", backend)
	}
}

// NewRequestWithBackend creates request configured for endpoint and explicit backend.
func (c *client) NewRequestWithBackend(ctx context.Context, backend BackendTarget, endpoint *Endpoint, _ ...RequestOption) (req *resty.Request, err error) {
	sc, err := c.identifyClient(ctx, backend)
	if err != nil {
		return nil, err
	}

	clientName, _ := backendToSubClientName(backend)
	ctxv := context.WithValue(ctx, contextKeyClientName, clientName)
	hC, err := sc.newHTTPClient(ctxv)
	if err != nil {
		return nil, err
	}

	contextData := sc.ContextData(ctxv)
	ctxv = storeExtraDataInContext(ctxv, contextData)

	var (
		retryCount       = 5
		retryWaitTime    = 60 * time.Second
		retryMaxWaitTime = 5 * time.Second
		retryConditions  = make([]resty.RetryConditionFunc, 0)
		retryIdempotent  = false
	)

	switch endpoint.Method {
	case MethodPOST, MethodPUT, MethodDELETE:
		var conflictRetry resty.RetryConditionFunc = func(resp *resty.Response, _ error) bool {
			if sc.idempotentRetryCondition()(resp, nil) {
				// Extend retries for BUSY_ENTITY responses with unknown server-side resolution time.
				resp.Request.RetryCount++
				return true
			}

			return false
		}

		retryConditions = append(retryConditions, conflictRetry)
		retryIdempotent = true
	}

	hR := hC.NewRequest().
		SetContext(ctxv).
		SetRetryDefaultConditions(true).
		SetRetryCount(retryCount).
		SetRetryMaxWaitTime(retryMaxWaitTime).
		SetRetryWaitTime(retryWaitTime).
		AddRetryConditions(retryConditions...).
		SetRetryAllowNonIdempotent(retryIdempotent)

	for _, q := range endpoint.QueryParams {
		if q.Value != "" {
			hR.SetQueryParam(q.Name, q.Value)
		}
	}

	for _, p := range endpoint.PathParams {
		if p.Value != "" {
			hR.SetPathParam(p.Name, p.Value)
		}
	}

	return hR, nil
}

// NewRequest creates request configured for endpoint.
func (c *client) NewRequest(ctx context.Context, endpoint *Endpoint, _ ...RequestOption) (req *resty.Request, err error) {
	return c.NewRequestWithBackend(ctx, endpoint.Backend, endpoint)
}

// DoWithBackend executes endpoint request with explicit backend.
func (c *client) DoWithBackend(ctx context.Context, backend BackendTarget, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error) {
	req, err := c.NewRequestWithBackend(ctx, backend, endpoint)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err := opt(endpoint, req); err != nil {
			return nil, err
		}
	}

	resp, err := req.SetResult(endpoint.ResponseType).Execute(endpoint.Method.String(), endpoint.PathTemplate)
	if err != nil {
		return nil, err
	}

	sc, err := c.identifyClient(ctx, backend)
	if err != nil {
		return nil, err
	}

	if errAPI := sc.parseAPIError(endpoint.Description, resp); errAPI != nil {
		return nil, errAPI
	}

	return resp, nil
}

// Do executes endpoint request.
func (c *client) Do(ctx context.Context, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error) {
	return c.DoWithBackend(ctx, endpoint.Backend, endpoint, opts...)
}
