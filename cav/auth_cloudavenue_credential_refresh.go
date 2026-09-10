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
	"errors"

	"resty.dev/v3"
)

// Refresh refreshes authentication token and session metadata.
func (c *cloudavenueCredential) Refresh(ctx context.Context) error {
	logger := c.logger.WithGroup("refresh")
	ep, err := GetEndpoint("SessionVmware")
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get endpoint for CreateSessionVmware", "error", err)
		return errors.New("failed to get endpoint for CreateSessionVmware: " + err.Error())
	}

	resp, err := c.executeRefreshRequest(ctx, ep)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to refresh session", "error", err)
		return err
	}

	if err := (&vmware{}).parseAPIError("SessionVmware", resp); err != nil {
		c.clearBearer()
		logger.ErrorContext(ctx, "Failed to refresh session", "error", err)
		return err
	}

	logger.DebugContext(
		ctx, "Successfully refreshed session",
		"endpoint", ep.PathTemplate,
		"method", ep.Method,
		"status", resp.StatusCode(),
	)

	c.storeRefreshSession(resp)

	return nil
}

// executeRefreshRequest sends VMware session refresh request.
func (c *cloudavenueCredential) executeRefreshRequest(ctx context.Context, ep *Endpoint) (*resty.Response, error) {
	request := c.httpC.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json;version="+vmwareVCDVersion).
		SetResult(&apiResponseSessionVmware{}).
		SetResultError(&vmwareError{}).
		SetURL(c.console.GetAPIVCDEndpoint())

	if err := c.applyRefreshAuth(request); err != nil {
		return nil, err
	}

	return request.Execute("POST", ep.PathTemplate)
}

// applyRefreshAuth applies bearer or basic authentication to request.
func (c *cloudavenueCredential) applyRefreshAuth(request *resty.Request) error {
	bearer := c.currentBearer()
	if bearer != "" {
		request.SetAuthToken(bearer)
		return nil
	}

	request.SetBasicAuth(c.username+"@"+c.organization, c.password)
	return nil
}

// currentBearer returns current bearer token.
func (c *cloudavenueCredential) currentBearer() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bearer
}

// clearBearer removes cached bearer token.
func (c *cloudavenueCredential) clearBearer() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bearer = ""
}

// storeRefreshSession updates cached session data from refresh response.
func (c *cloudavenueCredential) storeRefreshSession(resp *resty.Response) {
	session := resp.Result().(*apiResponseSessionVmware)

	c.mu.Lock()
	defer c.mu.Unlock()
	c.bearer = resp.Header().Get(cloudavenueCredentialXVmwareAccessToken)
	c.organizationID = session.Org.ID
	c.siteID = session.Site.ID
}
