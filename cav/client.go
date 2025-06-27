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

	"resty.dev/v3"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	subclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/subClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

type client struct {
	httpClient         *resty.Client
	console            consoles.Console
	clientsInitialized map[subclient.Name]subclient.Client
}

const (
	Vmware    = subclient.Vmware
	Cerberus  = subclient.Cerberus
	Netbackup = subclient.Netbackup
	mock      = subclient.Name("mock")     // For testing purposes
	mockJob   = subclient.Name("mock-job") // For testing purposes with jobs
)

type Client interface {
	NewRequest(ctx context.Context, client subclient.Name, reqOpt ...RequestOption) (req *resty.Request, err error)
	ParseAPIError(resp *resty.Response) *errors.APIError
}

// NewClient creates a new client object
//
// Zero or more ClientOption object can be passed as a parameter.
// These options will then be applied to the client.
func NewClient(organization string, opts ...ClientOption) (Client, error) {
	settings := newSettings(organization)

	// Load the console based on the organization name.
	// This will set the Console property in the settings.
	if err := withConsole()(settings); err != nil {
		return nil, err
	}

	if settings.httpClient == nil {
		settings.httpClient = httpclient.NewHTTPClient()
	}

	for _, opt := range opts {
		if err := opt(settings); err != nil {
			return nil, err
		}
	}

	return &client{
		httpClient:         settings.httpClient,
		console:            settings.Console,
		clientsInitialized: settings.SubClients,
	}, nil
}

// ParseAPIError parses the API error response from the subclient.
func (c *client) ParseAPIError(resp *resty.Response) *errors.APIError {
	if resp == nil {
		return nil
	}

	clientName, ok := resp.Request.Context().Value(subclient.ContextKeyClientName).(subclient.Name)
	if !ok {
		return &errors.APIError{
			StatusCode: resp.StatusCode(),
			Message:    "unknown client",
			Duration:   resp.Duration(),
			Endpoint:   resp.Request.URL,
		}
	}
	if v, ok := c.clientsInitialized[clientName]; ok {
		return v.ParseAPIError(resp)
	}
	return &errors.APIError{
		StatusCode: resp.StatusCode(),
		Message:    "unknown client",
		Duration:   resp.Duration(),
		Endpoint:   resp.Request.URL,
	}
}

// identifyClient identifies the client type.
func (c *client) identifyClient(_ context.Context, cN subclient.Name) (subclient.Client, error) {
	if c.clientsInitialized[cN] == nil {
		return nil, fmt.Errorf("invalid client %s", cN)
	}
	return c.clientsInitialized[cN], nil
}
