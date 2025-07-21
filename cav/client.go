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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

// isMockClient is a boolean flag to indicate if the client is a mock client.
var isMockClient bool

type client struct {
	logger             *slog.Logger
	console            consoles.Console
	clientsInitialized map[SubClientName]SubClient
}

type Client interface {
	NewRequest(ctx context.Context, endpoint *Endpoint, opts ...RequestOption) (req *resty.Request, err error)
	Logger() *slog.Logger
	Do(ctx context.Context, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error)
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

	// Detect if the client is a mock client based on the organization name.
	// This is a simple heuristic to determine if the client is a mock client.
	if organization == "cav01ev01ocb0001234" {
		isMockClient = true
	}

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
			Method:     resp.Request.Method,
		}
	}
	if v, ok := c.clientsInitialized[clientName]; ok {
		return v.parseAPIError(action, resp)
	}
	return &errors.APIError{
		StatusCode: resp.StatusCode(),
		Message:    "unknown client",
		Duration:   resp.Duration(),
		Endpoint:   resp.Request.URL,
		Method:     resp.Request.Method,
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
