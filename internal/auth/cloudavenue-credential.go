/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package auth

import (
	"context"
	"errors"
	"log"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

var _ Auth = (*CloudavenueCredential)(nil)

// CloudavenueCredential is the pair of a username and password
// used to authenticate with the CloudAvenue service.
// This type is public because it's an internal package.
type CloudavenueCredential struct {
	username     string
	password     string
	bearer       string
	httpC        *resty.Client
	organization string
	console      consoles.Console
}

// XVmwareAccessToken is the header used to retrieve the Bearer token in the authentication process.
const XVmwareAccessToken = "X-VMWARE-VCLOUD-ACCESS-TOKEN" // #nosec G101

const VDCVersion = "38.1"

// NewCloudavenueCredential creates a new CloudavenueCredential
// with the given username and password.
var NewCloudavenueCredential = newCloudavenueCredential

func newCloudavenueCredential(httpC *resty.Client, c consoles.Console, organization, username, password string) (Auth, error) {
	if httpC == nil {
		return nil, errors.New("http client cannot be nil")
	}
	if organization == "" {
		return nil, errors.New("organization cannot be empty")
	}
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}
	return &CloudavenueCredential{
		httpC:        httpC,
		console:      c,
		organization: organization,
		username:     username,
		password:     password,
	}, nil
}

// Headers returns the HTTP headers required for authentication
// using the CloudavenueCredential.
func (c *CloudavenueCredential) Headers() map[string]string {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.bearer
	return headers
}

// Refresh is a placeholder method for refreshing the authentication token.
// It currently does nothing but should be implemented to refresh the token
// when needed.
func (c *CloudavenueCredential) Refresh(ctx context.Context) error {
	r := c.httpC.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json;version="+VDCVersion)
	if c.bearer != "" {
		r.SetAuthToken(c.bearer)
	} else {
		log.Default().Println("No bearer token found, using basic auth for authentication.")
		log.Default().Printf("Using username: %s, organization: %s\n", c.username, c.organization)
		r.SetBasicAuth(c.username+"@"+c.organization, c.password)
	}

	resp, err := r.Post(c.console.GetAPIVCDEndpoint() + "/1.0.0/sessions")
	if err != nil {
		c.bearer = ""
		return err
	}

	if resp.IsError() {
		c.bearer = ""
		return errors.New("failed to authenticate: " + resp.String())
	}

	c.bearer = resp.Header().Get(XVmwareAccessToken)

	return nil
}

// IsInitialized checks if the CloudavenueCredential is initialized.
func (c *CloudavenueCredential) IsInitialized() bool {
	return c.bearer != ""
}
