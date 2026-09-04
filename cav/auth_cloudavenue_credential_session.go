/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import "errors"

const sessionKeyOrganizationID = "organizationID"

// getSession returns session data for cache persistence.
func (c *cloudavenueCredential) getSession() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]string{
		"organization":           c.organization,
		sessionKeyOrganizationID: c.organizationID,
		"siteID":                 c.siteID,
		"bearer":                 c.bearer,
	}
}

// restoreSession loads cached session data.
func (c *cloudavenueCredential) restoreSession(data map[string]string) error {
	if data == nil {
		return errors.New("invalid session data")
	}

	xlogger.Debug("Restoring session from cache", "data", data)

	c.mu.Lock()
	defer c.mu.Unlock()
	c.organization = data["organization"]
	c.bearer = data["bearer"]
	c.organizationID = data[sessionKeyOrganizationID]
	c.siteID = data["siteID"]

	return nil
}

// getExtraData returns backend metadata derived from session state.
func (c *cloudavenueCredential) getExtraData() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]string{
		sessionKeyOrganizationID: c.organizationID,
		"siteID":                 c.siteID,
	}
}
