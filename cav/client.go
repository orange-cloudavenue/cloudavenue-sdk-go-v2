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

	"resty.dev/v3"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

type client struct {
	logger             *slog.Logger
	httpClient         *resty.Client
	console            consoles.Console
	clientsInitialized map[SubClientName]SubClient
}

type Client interface {
	NewRequest(ctx context.Context, endpoint *Endpoint, opts ...RequestOption) (req *resty.Request, err error)
	ParseAPIError(action string, resp *resty.Response) *errors.APIError
	Logger() *slog.Logger
}

// NewClient creates a new client object
//
// Zero or more ClientOption object can be passed as a parameter.
// These options will then be applied to the client.
func NewClient(organization string, opts ...ClientOption) (Client, error) {
	// organization format validation are done in the withConsole option.

	settings := newSettings(organization)

	// Load the console based on the organization name.
	// This will set the Console property in the settings.
	if err := withConsole()(settings); err != nil {
		return nil, err
	}

	if settings.httpClient == nil {
		settings.httpClient = httpclient.NewHTTPClient()
	}

	client := &client{
		httpClient: settings.httpClient,
		console:    settings.Console,
	}

	for _, opt := range opts {
		if err := opt(settings); err != nil {
			return nil, err
		}
	}

	client.logger = xlogger.WithGroup("client").With("organization", settings.Organization)
	client.clientsInitialized = settings.SubClients

	return client, nil
}

// ParseAPIError parses the API error response from the subclient.
func (c *client) ParseAPIError(action string, resp *resty.Response) *errors.APIError {
	if resp == nil {
		return nil
	}

	clientName, ok := resp.Request.Context().Value(contextKeyClientName).(SubClientName)
	if !ok {
		return &errors.APIError{
			StatusCode: resp.StatusCode(),
			Message:    "unknown client",
			Duration:   resp.Duration(),
			Endpoint:   resp.Request.URL,
		}
	}
	if v, ok := c.clientsInitialized[clientName]; ok {
		return v.ParseAPIError(action, resp)
	}
	return &errors.APIError{
		StatusCode: resp.StatusCode(),
		Message:    "unknown client",
		Duration:   resp.Duration(),
		Endpoint:   resp.Request.URL,
	}
}

// Logger returns the logger for the client.
func (c *client) Logger() *slog.Logger {
	return c.logger
}

// identifyClient identifies the client type.
func (c *client) identifyClient(_ context.Context, cN SubClientName) (SubClient, error) {
	if c.clientsInitialized[cN] == nil {
		return nil, fmt.Errorf("invalid client %s", cN)
	}
	return c.clientsInitialized[cN], nil
}
