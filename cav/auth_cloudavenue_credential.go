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
	"errors"
	"log/slog"
	"sync"

	"github.com/orange-cloudavenue/common-go/validators"
	"resty.dev/v3"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/http-client"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

var _ auth = (*cloudavenueCredential)(nil)

// cloudavenueCredential implements auth for username/password login.
type cloudavenueCredential struct {
	logger         *slog.Logger
	httpC          *resty.Client
	username       string `validate:"required"`
	password       string `validate:"required"`
	bearer         string
	organization   string `validate:"required"`
	organizationID string
	siteID         string
	console        consoles.ConsoleName
	mu             sync.RWMutex
}

// cloudavenueCredentialXVmwareAccessToken carries VMware access token.
const cloudavenueCredentialXVmwareAccessToken = "X-VMWARE-VCLOUD-ACCESS-TOKEN" // #nosec G101

func newCloudavenueCredential(c consoles.ConsoleName, organization, username, password string) (auth, error) {
	cc := &cloudavenueCredential{
		logger:       xlogger.WithGroup("auth"),
		console:      c,
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

	cc.logger = cc.logger.With("organization", cc.organization)
	cc.httpC = httpclient.NewHTTPClient().SetBaseURL(c.GetAPIVCDEndpoint())

	return cc, nil
}

// Headers returns authorization headers for authenticated requests.
func (c *cloudavenueCredential) Headers() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.bearer
	return headers
}

// IsInitialized reports whether bearer token is available.
func (c *cloudavenueCredential) IsInitialized() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bearer != ""
}
